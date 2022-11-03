/*
Copyright 2022.

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
)

// ApplicationSnapshotEnvironmentBindingSpec defines the desired state of ApplicationSnapshotEnvironmentBinding
type ApplicationSnapshotEnvironmentBindingSpec struct {

	// Application is a reference to the Application resource (defined in the namespace) involved in the binding
	Application string `json:"application"`

	// Environment is the Environment resource (defined in the namespace) that the binding will deploy to
	Environment string `json:"environment"`

	// Snapshot is the Snapshot resource (defined in the namespace) that contains the container image versions
	// for the components of the Application
	Snapshot string `json:"snapshot"`

	// Components contains individual component data
	Components []BindingComponent `json:"components"`
}

// ApplicationSnapshotEnvironmentBindingStatus defines the observed state of ApplicationSnapshotEnvironmentBinding
type ApplicationSnapshotEnvironmentBindingStatus struct {

	// GitOpsDeployments describes the set of GitOpsDeployment resources that correspond to the binding.
	// To determine the health/sync status of a binding, you can look at the GitOpsDeployments described here.
	GitOpsDeployments []BindingStatusGitOpsDeployment `json:"gitopsDeployments,omitempty"`

	// Components describes a component's GitOps repository information.
	// This status is updated by the Application Service controller.
	Components []ComponentStatus `json:"components,omitempty"`

	// Condition describes operations on the GitOps repository, for example, if there were issues with generating/processing the repository.
	// This status is updated by the Application Service controller.
	GitOpsRepoConditions []metav1.Condition `json:"gitopsRepoConditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApplicationSnapshotEnvironmentBinding is the Schema for the applicationsnapshotenvironmentbindings API
// +kubebuilder:resource:path=applicationsnapshotenvironmentbindings,shortName=aseb;binding
type ApplicationSnapshotEnvironmentBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSnapshotEnvironmentBindingSpec   `json:"spec"`
	Status ApplicationSnapshotEnvironmentBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationSnapshotEnvironmentBindingList contains a list of ApplicationSnapshotEnvironmentBinding
type ApplicationSnapshotEnvironmentBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationSnapshotEnvironmentBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationSnapshotEnvironmentBinding{}, &ApplicationSnapshotEnvironmentBindingList{})
}
