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

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	previewv1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
	v1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
)

// PreviewEnvironmentReconciler reconciles a PreviewEnvironment object
type PreviewEnvironmentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironments/finalizers,verbs=update
//+kubebuilder:rbac:groups=kustomize.toolkit.fluxcd.io,resources=kustomizations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=buckets;gitrepositories,verbs=create;update;patch;delete;get;list;watch

func (r *PreviewEnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	obj := &v1alpha1.PreviewEnvironment{}
	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Info("object not found", "name", req.NamespacedName.Name, "namespace", req.NamespacedName.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcile(ctx, obj)
}

func (r *PreviewEnvironmentReconciler) reconcile(ctx context.Context, obj *v1alpha1.PreviewEnvironment) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	req := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetName(),
	}

	obj.SetKustomizeSubstitution("branch", obj.Spec.Branch)
	obj.SetKustomizeSubstitution("commit", obj.Spec.Commit)

	kustomization := &kustomizev1.Kustomization{}
	if err := r.Client.Get(ctx, req, kustomization); err != nil {
		if apierrors.IsNotFound(err) {
			kustomization = &kustomizev1.Kustomization{
				ObjectMeta: metav1.ObjectMeta{
					Name:      obj.GetName(),
					Namespace: obj.GetNamespace(),
				},
				Spec: obj.GetKustomizationSpec(),
			}
			if err := r.Client.Create(ctx, kustomization); err != nil {
				log.Error(err, "error creating Kustomization")
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, err
		}
	}

	// kustomization.Spec = obj.GetKustomizationSpec()
	//
	// if err := r.Client.Update(ctx, kustomization); err != nil {
	//     log.Error(err, "error updating Kustomization")
	//     return ctrl.Result{}, err
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PreviewEnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&previewv1alpha1.PreviewEnvironment{}).
		Owns(&kustomizev1.Kustomization{}).
		Complete(r)
}
