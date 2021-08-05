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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PreviewEnvironmentManagerSpec defines the desired state of PreviewEnvironmentManager
type PreviewEnvironmentManagerSpec struct {
	SourceRef types.NamespacedName `json:"sourceRef"`

	// +required
	Template PreviewEnvironmentTemplateSpec `json:"template"`

	// +required
	SpawnRules Rules `json:"spawnRules"`

	// +required
	Interval metav1.Duration `json:"interval"`

	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

// PreviewEnvironmentTemplateSpec defines the type of PreviewEnvironments that will be created
type PreviewEnvironmentTemplateSpec struct {
	// Only supports Kustomization resources
	// at present
	// +required
	TemplateRef types.NamespacedName `json:"templateRef"`

	// +optional
	Limit int `json:"limit,omitempty"`

	// +optional
	Prefix string `json:"prefix,omitempty"`
}

// Rules define the rules that determine which branches will
// create a PreviewEnvironment
type Rules struct {
	// MatchBranch is be a regex string that will be used to spawn a new
	// preview environment.
	// `branch-name` will be substituted when this
	// rule is used
	// +optional
	MatchBranch string `json:"matchBranch,omitempty"`

	// +optional
	// GitHub GitHubRules `json:"github,omitempty"`

	// +optional
	// GitLab GitLabRules `json:"gitlab,omitempty"`
}

// type GitHubRules struct {
//     PullRequest bool     `json:"pullRequest"`
//     WithLabel   []string `json:"withLabel"`
// }
//
// type GitLabRules struct {
//     MergeRequest bool `json:"mergeRequest"`
// }

// PreviewEnvironmentManagerStatus defines the observed state of PreviewEnvironmentManager
type PreviewEnvironmentManagerStatus struct {
	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions holds the conditions for the GitRepository.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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
