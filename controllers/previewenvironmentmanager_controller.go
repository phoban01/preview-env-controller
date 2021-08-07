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
		log.Info("object not found",
			"namespace", req.NamespacedName.Namespace,
			"name", req.NamespacedName.Name)
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
			Namespace: obj.Namespace,
		}, secret); err != nil {
		return ctrl.Result{RequeueAfter: time.Minute * 1}, err
	}

	if _, ok := secret.Data["password"]; !ok {
		return ctrl.Result{RequeueAfter: time.Minute * 1}, nil
	}

	log.Info("Found token", "token", secret.Data["password"])
	branches, err := getBranches(ctx, obj.Spec.Watch.URL, secret.Data["password"])
	if err != nil {
		log.Info("rate limited by github")
		return ctrl.Result{RequeueAfter: time.Minute * 1}, nil
	}

	// get previewenvironments managed by thie manager
	// should use controllerRef
	prEnvs := &v1alpha1.PreviewEnvironmentList{}
	if err := r.Client.List(ctx, prEnvs, &client.ListOptions{Namespace: obj.GetNamespace()}); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	re, err := regexp.Compile(obj.Spec.Rules.MatchBranch)
	if err != nil {
		return ctrl.Result{}, err
	}

	// if branch is missing then delete previewenvironment
	// or doesn't match current rule
	gcPrEnvs := &v1alpha1.PreviewEnvironmentList{}
	for _, p := range prEnvs.Items {
		if _, ok := branches[p.Spec.Branch]; !ok {
			gcPrEnvs.Items = append(gcPrEnvs.Items, p)
		}

		// TODO: make configurable via Prune field
		if ok := re.MatchString(p.Spec.Branch); !ok {
			log.Info("Existing environment branch no longer matches rule",
				"branch", p.Spec.Branch,
				"rule", obj.Spec.Rules.MatchBranch)
			gcPrEnvs.Items = append(gcPrEnvs.Items, p)
		}
	}

	for _, p := range gcPrEnvs.Items {
		if err := r.Client.Delete(ctx, &p); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
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

		if err := r.reconcilePreviewEnv(ctx, obj, branch, sha); err != nil {
			log.Info("error reconciling preview env")
			return ctrl.Result{}, err
		}
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

func (r *PreviewEnvironmentManagerReconciler) reconcilePreviewEnv(ctx context.Context, obj *v1alpha1.PreviewEnvironmentManager, branch, sha string) error {
	log := ctrl.LoggerFrom(ctx)

	name := fmt.Sprintf("%s-%s", strings.TrimSuffix(obj.Spec.Template.Prefix, "-"), branch)

	//check if exists
	//TODO: refactor this to use cache
	if r.previewEnvExists(ctx, name, obj.Namespace) != true {
		log.Info("creating preview environment", "name", name)
		return r.createPreviewEnv(ctx, name, branch, obj)
	}

	log.Info("preview environment exists", "name", name)

	return nil
}

//TODO: refactor this to use cache
func (r *PreviewEnvironmentManagerReconciler) previewEnvExists(ctx context.Context, name, namespace string) bool {
	err := r.Client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &v1alpha1.PreviewEnvironment{})

	return apierrors.IsNotFound(err) != true
}

func (r *PreviewEnvironmentManagerReconciler) createPreviewEnv(
	ctx context.Context,
	name,
	branch string,
	obj *v1alpha1.PreviewEnvironmentManager) error {
	log := ctrl.LoggerFrom(ctx)

	newEnv := &v1alpha1.PreviewEnvironment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: obj.Namespace,
		},
		Spec: v1alpha1.PreviewEnvironmentSpec{
			Branch: branch,
		},
	}

	if err := ctrl.SetControllerReference(obj, newEnv, r.Scheme); err != nil {
		return err
	}

	if err := r.Client.Create(ctx, newEnv); err != nil {
		log.Info("error creating preview environment", "name", name)
		return err
	}

	log.Info("preview environment created", "name", name)

	return nil
}
