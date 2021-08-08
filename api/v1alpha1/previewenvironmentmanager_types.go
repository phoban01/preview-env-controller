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

package v1alpha1

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PreviewEnvironmentManagerStrategy string

const (
	PreviewEnvironmentManagerKind                                   = "PreviewEnvironmentManager"
	PullRequestStrategy           PreviewEnvironmentManagerStrategy = "PullRequest"
	BranchStrategy                PreviewEnvironmentManagerStrategy = "Branch"
)

// PreviewEnvironmentManagerSpec defines the desired state of PreviewEnvironmentManager
type PreviewEnvironmentManagerSpec struct {
	// +required
	Watch WatchObject `json:"watch"`

	// +required
	Strategy Strategy `json:"strategy"`

	// +required
	Template TemplateSpec `json:"template"`

	// +required
	Interval metav1.Duration `json:"interval"`

	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// +optional
	Limit int `json:"limit,omitempty"`

	// +kubebuilder:default:=true
	// +optional
	Prune bool `json:"prune"`
}

// WatchObject defines a repository to watch
// for branches
type WatchObject struct {
	// +required
	URL string `json:"url"`

	// +required
	Ref Ref `json:"ref"`

	// +required
	CredentialsRef meta.LocalObjectReference `json:"credentialsRef"`
}

type Ref struct {
	// +required
	Branch string `json:"branch"`
}

type Strategy struct {
	// =kubebuilder:validation:Enum=PullRequest,Branch
	// +optional
	Type PreviewEnvironmentManagerStrategy `json:"type"`

	// +optional
	Rules Rules `json:"rules"`
}

// Rules define the rules that determine which branches will
// create a PreviewEnvironment
type Rules struct {
	// Match is be a regex string that will be used to spawn a new
	// preview environment. `branch-name` will be substituted when this
	// rule is used
	// +optional
	Match string `json:"match,omitempty"`

	// +optional
	Labels []string `json:"label"`

	// +optional
	Draft bool `json:"draft"`
}

// TemplateSpec defines the type of PreviewEnvironments that will be created
type TemplateSpec struct {
	// KustomizationSpec
	// +required
	KustomizationSpec KustomizationSpec `json:"kustomizationSpec"`

	// SourceSpec
	// +optional
	SourceSpec SourceSpec `json:"sourceSpec"`

	// +optional
	Prefix string `json:"prefix,omitempty"`

	// +optional
	CreateNamespace bool `json:"createNamespace"`

	// +optional
	TargetNamespace string `json:"targetNamespace"`
}

type KustomizationSpec struct {
	// +required
	Path string `json:"path"`

	// +required
	Interval metav1.Duration `json:"interval"`

	// +kubebuilder:default:=true
	// +optional
	Prune bool `json:"prune"`

	// +optional
	SourceRef *kustomizev1.CrossNamespaceSourceReference `json:"sourceRef"`
}

type SourceSpec struct {
	// +kubebuilder:default:='1m'
	// +optional
	Interval metav1.Duration `json:"interval"`

	// +optional
	SecretRef *meta.LocalObjectReference `json:"secretRef"`
}

// PreviewEnvironmentManagerStatus defines the observed state of PreviewEnvironmentManager
type PreviewEnvironmentManagerStatus struct {
	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions holds the conditions for the GitRepository.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	EnvironmentCount int `json:"environmentCount"`
}

// GetLimitl Returns the reconcilation interval
func (in *PreviewEnvironmentManager) GetLimit() int {
	return in.Spec.Limit
}

// GetInterval Returns the reconcilation interval
func (in *PreviewEnvironmentManager) GetInterval() metav1.Duration {
	return in.Spec.Interval
}

// GetEnvironmentCount returns the number of active preview environments
func (in *PreviewEnvironmentManager) GetEnvironmentCount() int {
	return in.Status.EnvironmentCount
}

//+kubebuilder:printcolumn:name="Count",type=integer,JSONPath=`.status.environmentCount`
//+kubebuilder:resource:shortName=pman
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PreviewEnvironmentManager is the Schema for the previewenvironmentmanagers API
type PreviewEnvironmentManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PreviewEnvironmentManagerSpec   `json:"spec,omitempty"`
	Status PreviewEnvironmentManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PreviewEnvironmentManagerList contains a list of PreviewEnvironmentManager
type PreviewEnvironmentManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PreviewEnvironmentManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PreviewEnvironmentManager{}, &PreviewEnvironmentManagerList{})
}
