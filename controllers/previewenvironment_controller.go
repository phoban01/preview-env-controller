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
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	previewv1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
	v1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
)

var (
	ownerKey      = ".metadata.controller"
	finalizerName = "preview.gitops.phoban.io/finalizer"
)

// PreviewEnvironmentReconciler reconciles a PreviewEnvironment object
type PreviewEnvironmentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=preview.gitops.phoban.io,resources=previewenvironments/finalizers,verbs=update
//+kubebuilder:rbac:groups=kustomize.toolkit.fluxcd.io,resources=kustomizations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories,verbs=create;update;patch;delete;get;list;watch
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories/finalizers,verbs=update;patch;delete;get;list

func (r *PreviewEnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	obj := &v1alpha1.PreviewEnvironment{}
	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Info("object not found", "name", req.NamespacedName.Name, "namespace", req.NamespacedName.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if obj.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(obj, finalizerName) {
			controllerutil.AddFinalizer(obj, finalizerName)
			if err := r.Update(ctx, obj); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(obj, finalizerName) {
			if obj.Spec.CreateNamespace {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: obj.Status.EnvNamespace,
					},
				}
				if err := r.Client.Delete(ctx, ns); err != nil {
					return ctrl.Result{}, err
				}
			}

			controllerutil.RemoveFinalizer(obj, finalizerName)
			if err := r.Update(ctx, obj); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	return r.reconcile(ctx, obj)
}

func (r *PreviewEnvironmentReconciler) reconcile(ctx context.Context, obj *v1alpha1.PreviewEnvironment) (ctrl.Result, error) {
	manager := &v1alpha1.PreviewEnvironmentManager{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name:      metav1.GetControllerOf(obj).Name,
		Namespace: obj.GetNamespace(),
	}, manager); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if obj.Spec.CreateNamespace {
		if err := r.reconcileNamespace(ctx, obj, manager); err != nil {
			return ctrl.Result{}, nil
		}
	}

	if err := r.reconcileGitRepository(ctx, obj, manager); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileKustomization(ctx, obj, manager); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *PreviewEnvironmentReconciler) reconcileNamespace(ctx context.Context, obj *v1alpha1.PreviewEnvironment, manager *v1alpha1.PreviewEnvironmentManager) error {
	log := ctrl.LoggerFrom(ctx)

	if obj.Status.EnvNamespace == "" {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: obj.Name,
			},
		}

		if err := r.Client.Create(ctx, ns); err != nil {
			log.Error(err, "error creating namespace")
			return err
		}

		obj.Status.EnvNamespace = obj.Name

		if err := r.Client.Status().Update(ctx, obj); err != nil {
			log.Error(err, "error updating status")
			return err
		}
	}

	return nil
}

func (r *PreviewEnvironmentReconciler) reconcileGitRepository(ctx context.Context, obj *v1alpha1.PreviewEnvironment, manager *v1alpha1.PreviewEnvironmentManager) error {
	log := ctrl.LoggerFrom(ctx)

	defer func() {
		if err := r.Client.Status().Update(ctx, obj); err != nil {
			log.Error(err, "error updating status")
		}
	}()

	if obj.Status.GitRepository == nil {
		gitRepo := &sourcev1.GitRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name:      obj.Name,
				Namespace: obj.Namespace,
			},
			Spec: sourcev1.GitRepositorySpec{
				//TODO: both of these could be fetched via helper methods
				//on manager
				URL:       manager.Spec.Watch.URL,
				SecretRef: &manager.Spec.Watch.CredentialsRef,
				Interval:  manager.Spec.Template.Interval,
				Reference: &sourcev1.GitRepositoryRef{
					Branch: obj.Spec.Branch,
				},
			},
		}
		if err := ctrl.SetControllerReference(obj, gitRepo, r.Scheme); err != nil {
			return err
		}
		if err := r.Client.Create(ctx, gitRepo); err != nil {
			log.Error(err, "error creating git repo")
			return err
		}
		obj.Status.GitRepository = &meta.NamespacedObjectReference{
			Name:      gitRepo.Name,
			Namespace: gitRepo.Namespace,
		}

		return nil
	}

	curRepo, err := r.getGitRepo(ctx, obj)
	if err != nil {
		return err
	}

	if curRepo.Status.Artifact == nil {
		return nil
	}

	if curRepo.Spec.URL == manager.Spec.Watch.URL &&
		curRepo.Spec.SecretRef.Name == manager.Spec.Watch.CredentialsRef.Name &&
		curRepo.Spec.Interval == manager.Spec.Template.Interval &&
		curRepo.Status.Artifact.Revision == obj.Status.Commit {
		return nil
	}

	newRepo := curRepo.DeepCopy()
	newRepo.Spec.URL = manager.Spec.Watch.URL
	newRepo.Spec.SecretRef = manager.Spec.Watch.CredentialsRef.DeepCopy()
	newRepo.Spec.Interval = manager.Spec.Template.Interval

	obj.Status.Commit = curRepo.Status.Artifact.Revision

	return r.Client.Patch(ctx, curRepo, client.MergeFrom(newRepo))
}

func (r *PreviewEnvironmentReconciler) getGitRepo(ctx context.Context, obj *v1alpha1.PreviewEnvironment) (*sourcev1.GitRepository, error) {
	repo := &sourcev1.GitRepository{}
	req := types.NamespacedName{
		Namespace: obj.Status.GitRepository.Namespace,
		Name:      obj.Status.GitRepository.Name,
	}
	if err := r.Client.Get(ctx, req, repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *PreviewEnvironmentReconciler) reconcileKustomization(ctx context.Context, obj *v1alpha1.PreviewEnvironment, manager *v1alpha1.PreviewEnvironmentManager) error {
	log := ctrl.LoggerFrom(ctx)

	if obj.Status.Commit == "" {
		return nil
	}

	if obj.Status.Kustomization == nil {
		var sourceRef *kustomizev1.CrossNamespaceSourceReference
		if manager.Spec.Template.KustomizationSpec.SourceRef == nil {
			sourceRef = &kustomizev1.CrossNamespaceSourceReference{
				Kind:      "GitRepository",
				Namespace: obj.Namespace,
				Name:      obj.Status.GitRepository.Name,
			}
		} else {
			sourceRef = manager.Spec.Template.KustomizationSpec.SourceRef
		}
		kustomization := &kustomizev1.Kustomization{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: obj.Namespace,
				Name:      obj.Name,
			},
			Spec: kustomizev1.KustomizationSpec{
				SourceRef: *sourceRef,
				Path:      manager.Spec.Template.KustomizationSpec.Path,
				Prune:     manager.Spec.Template.KustomizationSpec.Prune,
				Interval:  manager.Spec.Template.Interval,
				PostBuild: &kustomizev1.PostBuild{
					Substitute: map[string]string{
						"branch": obj.Spec.Branch,
						"commit": obj.Status.Commit,
					},
				},
			},
		}

		if obj.Spec.CreateNamespace {
			kustomization.Spec.TargetNamespace = obj.Status.EnvNamespace
			kustomization.Spec.PostBuild.Substitute["namespace"] = obj.Status.EnvNamespace
		}

		if err := ctrl.SetControllerReference(obj, kustomization, r.Scheme); err != nil {
			return err
		}

		if err := r.Client.Create(ctx, kustomization); err != nil {
			log.Error(err, "error creating kustomization")
			return err
		}

		obj.Status.Kustomization = &meta.NamespacedObjectReference{
			Namespace: kustomization.Namespace,
			Name:      kustomization.Name,
		}

		if err := r.Client.Status().Update(ctx, obj); err != nil {
			log.Error(err, "error updating status")
		}

		return nil
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PreviewEnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &sourcev1.GitRepository{}, ownerKey, func(obj client.Object) []string {
		repo := obj.(*sourcev1.GitRepository)
		owner := metav1.GetControllerOf(repo)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != v1alpha1.GroupVersion.Version || owner.Kind != v1alpha1.PreviewEnvironmentKind {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&previewv1alpha1.PreviewEnvironment{}).
		Owns(&sourcev1.GitRepository{}).
		Owns(&kustomizev1.Kustomization{}).
		Complete(r)
}
