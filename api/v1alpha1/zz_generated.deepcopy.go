// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KustomizationSpec) DeepCopyInto(out *KustomizationSpec) {
	*out = *in
	if in.SourceRef != nil {
		in, out := &in.SourceRef, &out.SourceRef
		*out = new(v1beta1.CrossNamespaceSourceReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KustomizationSpec.
func (in *KustomizationSpec) DeepCopy() *KustomizationSpec {
	if in == nil {
		return nil
	}
	out := new(KustomizationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironment) DeepCopyInto(out *PreviewEnvironment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironment.
func (in *PreviewEnvironment) DeepCopy() *PreviewEnvironment {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PreviewEnvironment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentList) DeepCopyInto(out *PreviewEnvironmentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PreviewEnvironment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentList.
func (in *PreviewEnvironmentList) DeepCopy() *PreviewEnvironmentList {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PreviewEnvironmentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentManager) DeepCopyInto(out *PreviewEnvironmentManager) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentManager.
func (in *PreviewEnvironmentManager) DeepCopy() *PreviewEnvironmentManager {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PreviewEnvironmentManager) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentManagerList) DeepCopyInto(out *PreviewEnvironmentManagerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PreviewEnvironmentManager, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentManagerList.
func (in *PreviewEnvironmentManagerList) DeepCopy() *PreviewEnvironmentManagerList {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentManagerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PreviewEnvironmentManagerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentManagerSpec) DeepCopyInto(out *PreviewEnvironmentManagerSpec) {
	*out = *in
	out.Watch = in.Watch
	in.Strategy.DeepCopyInto(&out.Strategy)
	in.Template.DeepCopyInto(&out.Template)
	out.Interval = in.Interval
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentManagerSpec.
func (in *PreviewEnvironmentManagerSpec) DeepCopy() *PreviewEnvironmentManagerSpec {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentManagerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentManagerStatus) DeepCopyInto(out *PreviewEnvironmentManagerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentManagerStatus.
func (in *PreviewEnvironmentManagerStatus) DeepCopy() *PreviewEnvironmentManagerStatus {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentManagerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentSpec) DeepCopyInto(out *PreviewEnvironmentSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentSpec.
func (in *PreviewEnvironmentSpec) DeepCopy() *PreviewEnvironmentSpec {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviewEnvironmentStatus) DeepCopyInto(out *PreviewEnvironmentStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.GitRepository != nil {
		in, out := &in.GitRepository, &out.GitRepository
		*out = new(meta.NamespacedObjectReference)
		**out = **in
	}
	if in.Kustomization != nil {
		in, out := &in.Kustomization, &out.Kustomization
		*out = new(meta.NamespacedObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewEnvironmentStatus.
func (in *PreviewEnvironmentStatus) DeepCopy() *PreviewEnvironmentStatus {
	if in == nil {
		return nil
	}
	out := new(PreviewEnvironmentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ref) DeepCopyInto(out *Ref) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ref.
func (in *Ref) DeepCopy() *Ref {
	if in == nil {
		return nil
	}
	out := new(Ref)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rules) DeepCopyInto(out *Rules) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rules.
func (in *Rules) DeepCopy() *Rules {
	if in == nil {
		return nil
	}
	out := new(Rules)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SourceSpec) DeepCopyInto(out *SourceSpec) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(meta.LocalObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SourceSpec.
func (in *SourceSpec) DeepCopy() *SourceSpec {
	if in == nil {
		return nil
	}
	out := new(SourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Strategy) DeepCopyInto(out *Strategy) {
	*out = *in
	in.Rules.DeepCopyInto(&out.Rules)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Strategy.
func (in *Strategy) DeepCopy() *Strategy {
	if in == nil {
		return nil
	}
	out := new(Strategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemplateSpec) DeepCopyInto(out *TemplateSpec) {
	*out = *in
	out.Interval = in.Interval
	in.KustomizationSpec.DeepCopyInto(&out.KustomizationSpec)
	in.SourceSpec.DeepCopyInto(&out.SourceSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemplateSpec.
func (in *TemplateSpec) DeepCopy() *TemplateSpec {
	if in == nil {
		return nil
	}
	out := new(TemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WatchObject) DeepCopyInto(out *WatchObject) {
	*out = *in
	out.Ref = in.Ref
	out.CredentialsRef = in.CredentialsRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WatchObject.
func (in *WatchObject) DeepCopy() *WatchObject {
	if in == nil {
		return nil
	}
	out := new(WatchObject)
	in.DeepCopyInto(out)
	return out
}
