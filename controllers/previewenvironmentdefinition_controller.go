/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/google/go-github/github"
	v1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PreviewEnvironmentDefinitionReconciler reconciles a PreviewEnvironmentDefinition object
type PreviewEnvironmentDefinitionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironmentdefinitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironmentdefinitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironmentdefinitions/finalizers,verbs=update

// Reconcile is the main sync loop for our controller
func (r *PreviewEnvironmentDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	obj := &v1alpha1.PreviewEnvironmentDefinition{}

	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Info("object not found", "name", req.NamespacedName.Name, "namespace", req.NamespacedName.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *PreviewEnvironmentDefinitionReconciler) reconcile(ctx context.Context, obj *v1alpha1.PreviewEnvironmentDefinition) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// fetch source branches
	source := &sourcev1.GitRepository{}
	if err := r.Client.Get(ctx, obj.Spec.SourceRef, source); err != nil {
		log.Info("source not found", "name", source.Name, "namespace", source.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	branches, err := getBranches(ctx, source.Spec.URL)
	if err != nil {
		return ctrl.Result{}, err
	}

	// iterate existing previewenvironment
	// in the current namespace
	prEnvs := &v1alpha1.PreviewEnvironmentList{}
	if err := r.Client.List(ctx, prEnvs, &client.ListOptions{Namespace: obj.GetNamespace()}); err != nil {
		log.Info("source not found", "name", source.Name, "namespace", source.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if branch is missing then delete previewenvironment
	gcPrEnvs := &v1alpha1.PreviewEnvironmentList{}
	for _, p := range prEnvs.Items {
		if _, ok := branches[p.Spec.Branch]; !ok {
			gcPrEnvs.Items = append(gcPrEnvs.Items, p)
		}
	}

	for _, p := range gcPrEnvs.Items {
		if err := r.Client.Delete(ctx, &p); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	re, err := regexp.Compile(obj.Spec.SpawnRules.MatchBranch)
	if err != nil {
		return ctrl.Result{}, err
	}

	// iterate branches
	for b := range branches {
		// check for matching rules
		if ok := re.MatchString(b); !ok {
			continue
		}
		// create previewenvironment
		penvName := fmt.Sprintf("%s-%s", obj.Spec.Template.Prefix, b)
		penv := newPreviewEnvironment(penvName, obj.GetNamespace(), b)
		if err := r.Client.Create(ctx, penv); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PreviewEnvironmentDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PreviewEnvironmentDefinition{}).
		Owns(&v1alpha1.PreviewEnvironment{}).
		Complete(r)
}

func getBranches(ctx context.Context, sourceURL string) (map[string]string, error) {
	c := github.NewClient(http.DefaultClient)

	owner, repo, err := parseURL(sourceURL)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.Repositories.ListBranches(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	branches := make(map[string]string, len(resp))

	for _, b := range resp {
		branches[*b.Name] = *b.Commit.SHA
	}

	return branches, nil
}
func parseURL(repoURL string) (string, string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}
	p := strings.Split(u.Path, "/")

	owner := p[0]

	repo := strings.TrimSuffix(p[1], ".git")

	return owner, repo, nil

}

func newPreviewEnvironment(name, namespace, branch string) *v1alpha1.PreviewEnvironment {
	return &v1alpha1.PreviewEnvironment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.PreviewEnvironmentSpec{
			Branch: branch,
		},
	}
}
