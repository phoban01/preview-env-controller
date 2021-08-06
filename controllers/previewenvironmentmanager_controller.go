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
	"time"

	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/google/go-github/github"
	previewv1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
	v1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// PreviewEnvironmentManagerReconciler reconciles a PreviewEnvironmentManager object
type PreviewEnvironmentManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironmentmanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironmentmanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironmentmanagers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// SetupWithManager sets up the controller with the Manager.
func (r *PreviewEnvironmentManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&previewv1alpha1.PreviewEnvironmentManager{}).
		Owns(&previewv1alpha1.PreviewEnvironment{}).
		Complete(r)
}

func (r *PreviewEnvironmentManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	obj := &v1alpha1.PreviewEnvironmentManager{}

	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Info("object not found", "name", req.NamespacedName.Name, "namespace", req.NamespacedName.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcile(ctx, obj)
}

func (r *PreviewEnvironmentManagerReconciler) reconcile(ctx context.Context, obj *v1alpha1.PreviewEnvironmentManager) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// fetch source branches
	secret := &corev1.Secret{}
	if err := r.Get(ctx,
		types.NamespacedName{
			Name:      obj.Spec.Watch.CredentialsRef.Name,
			Namespace: obj.Spec.Watch.CredentialsRef.Namespace,
		}, secret); err != nil {
		return ctrl.Result{RequeueAfter: time.Minute * 1}, err
	}

	if _, ok := secret.Data["password"]; !ok {
		return ctrl.Result{RequeueAfter: time.Minute * 1}, nil
	}

	token := secret.Data["password"]

	branches, err := getBranches(ctx, obj.Spec.Watch.URL, token)
	if err != nil {
		log.Info("rate limited by github")
		return ctrl.Result{RequeueAfter: time.Minute * 1}, nil
	}

	// iterate existing previewenvironment
	// in the current namespace
	prEnvs := &v1alpha1.PreviewEnvironmentList{}
	if err := r.Client.List(ctx, prEnvs, &client.ListOptions{Namespace: obj.GetNamespace()}); err != nil {
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

	re, err := regexp.Compile(obj.Spec.Rules.MatchBranch)
	if err != nil {
		return ctrl.Result{}, err
	}

	// iterate branches
	for branch, sha := range branches {
		// check for matching rules
		if ok := re.MatchString(branch); !ok {
			log.Info("branch doesn't match rule",
				"branch", branch,
				"rule", obj.Spec.Rules.MatchBranch)
			continue
		}

		log.Info("branch matches rule",
			"branch", branch,
			"rule", obj.Spec.Rules.MatchBranch)

		// create previewenvironment
		if len(prEnvs.Items)-len(gcPrEnvs.Items) > obj.GetLimit() {
			log.Info("preview environment limit reached")
			return ctrl.Result{RequeueAfter: obj.GetInterval().Duration}, nil
		}
		penvName := fmt.Sprintf("%s-%s", strings.TrimSuffix(obj.Spec.Template.Prefix, "-"), branch)
		penv := &v1alpha1.PreviewEnvironment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      penvName,
				Namespace: obj.GetNamespace(),
			},
			Spec: v1alpha1.PreviewEnvironmentSpec{
				Branch:        branch,
				Commit:        sha,
				Kustomization: obj.Spec.Template.Kustomization,
			},
		}

		if err := r.Client.Create(ctx, penv); err != nil {
			if apierrors.IsAlreadyExists(err) {
				// if err := r.Update(ctx, penv); err != nil {
				return ctrl.Result{}, nil
				// }
			}
			return ctrl.Result{}, err
		}
		log.Info("preview environment created", "name", penv.GetName(), "namespace", penv.GetNamespace())
	}

	return ctrl.Result{RequeueAfter: obj.GetInterval().Duration}, nil
}

func getBranches(ctx context.Context, sourceURL string, repoToken []byte) (map[string]string, error) {
	httpClient := http.DefaultClient

	if repoToken != nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: string(repoToken)},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	c := github.NewClient(httpClient)

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
	p := strings.Split(strings.TrimLeft(u.Path, "/"), "/")

	owner := p[0]

	repo := strings.TrimSuffix(p[1], ".git")

	return owner, repo, nil

}
