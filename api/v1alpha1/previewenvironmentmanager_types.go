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
	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PreviewEnvironmentManagerKind = "PreviewEnvironmentManager"

// PreviewEnvironmentManagerSpec defines the desired state of PreviewEnvironmentManager
type PreviewEnvironmentManagerSpec struct {
	// +required
	Watch WatchObject `json:"watch"`

	// +required
	Rules Rules `json:"rules"`

	// +required
	Template TemplateSpec `json:"template"`

	// +required
	Interval metav1.Duration `json:"interval"`

	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// +optional
	Limit int `json:"limit,omitempty"`
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
}

type KustomizationSpec struct {
	// +required
	Path string `json:"path"`

	// +required
	Interval metav1.Duration `json:"interval"`

	// +kubebuilder:default:=true
	// +optional
	Prune bool `json:"prune"`
}

type SourceSpec struct {
	// +kubebuilder:default:='1m'
	// +optional
	Interval metav1.Duration `json:"interval"`

	// +optional
	SecretRef *meta.LocalObjectReference `json:"secretRef"`
}

// Rules define the rules that determine which branches will
// create a PreviewEnvironment
type Rules struct {
	// MatchBranch is be a regex string that will be used to spawn a new
	// preview environment. `branch-name` will be substituted when this
	// rule is used
	// +optional
	MatchBranch string `json:"matchBranch,omitempty"`
}

// PreviewEnvironmentManagerStatus defines the observed state of PreviewEnvironmentManager
type PreviewEnvironmentManagerStatus struct {
	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions holds the conditions for the GitRepository.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// GetLimitl Returns the reconcilation interval
func (in *PreviewEnvironmentManager) GetLimit() int {
	return in.Spec.Limit
}

// GetInterval Returns the reconcilation interval
func (in *PreviewEnvironmentManager) GetInterval() metav1.Duration {
	return in.Spec.Interval
}

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
