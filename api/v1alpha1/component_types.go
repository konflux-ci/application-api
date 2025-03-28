/*
Copyright 2021-2023 Red Hat, Inc.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ComponentSrcType describes the type of
// the src for the Component.
// Only one of the following location type may be specified.
// +kubebuilder:validation:Enum=Git;Image
type ComponentSrcType string

const (
	GitComponentSrcType   ComponentSrcType = "Git"
	ImageComponentSrcType ComponentSrcType = "Image"
)

type GitSource struct {
	// An HTTPS URL representing the git repository to create the component from.
	URL string `json:"url"`

	// Specify a branch/tag/commit id. If not specified, default is `main`/`master`.
	// Example: devel.
	// Optional.
	Revision string `json:"revision,omitempty"`

	// A relative path inside the git repo containing the component
	// Example: folderA/folderB/gitops.
	// Optional.
	Context string `json:"context,omitempty"`

	// If specified, the devfile at the URI will be used for the component. Can be a local path inside the repository, or an external URL.
	// Example: https://raw.githubusercontent.com/devfile-samples/devfile-sample-java-springboot-basic/main/devfile.yaml.
	// Optional.
	DevfileURL string `json:"devfileUrl,omitempty"`

	// If specified, the dockerfile at the URI will be used for the component. Can be a local path inside the repository, or an external URL.
	// Optional.
	DockerfileURL string `json:"dockerfileUrl,omitempty"`
}

// ComponentSource describes the Component source
type ComponentSource struct {
	ComponentSourceUnion `json:",inline"`
}

// +union
type ComponentSourceUnion struct {
	// Git Source for a Component.
	// Optional.
	GitSource *GitSource `json:"git,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ComponentSpec defines the desired state of Component
type ComponentSpec struct {

	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:validation:MaxLength=63
	// ComponentName is name of the component to be added to the Application. The name must adhere to DNS-1123 validation.
	// Required.
	// +required
	ComponentName string `json:"componentName"`

	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// Application is the name of the application resource that the component belongs to.
	// Required.
	// +required
	Application string `json:"application"`

	// Secret describes the name of a Kubernetes secret containing either:
	// 1. A Personal Access Token to access the Component's git repostiory (if using a Git-source component) or
	// 2. An Image Pull Secret to access the Component's container image (if using an Image-source component).
	// Optional.
	// +optional
	Secret string `json:"secret,omitempty"`

	// Source describes the Component source.
	// Optional.
	// +optional
	Source ComponentSource `json:"source,omitempty"`

	// Compute Resources required by this component.
	// Optional.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// The number of replicas to deploy the component with.
	// Optional.
	// +optional
	Replicas *int `json:"replicas,omitempty"`

	// The port to expose the component over.
	// Optional.
	// +optional
	TargetPort int `json:"targetPort,omitempty"`

	// The route to expose the component with.
	// Optional.
	// +optional
	Route string `json:"route,omitempty"`

	// An array of environment variables to add to the component (ValueFrom not currently supported)
	// Optional
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// The container image to build or create the component from
	// Example: quay.io/someorg/somerepository:latest.
	// Optional.
	// +optional
	ContainerImage string `json:"containerImage,omitempty"`

	// Whether or not to bypass the generation of GitOps resources for the Component. Defaults to false.
	// Optional.
	// +optional
	SkipGitOpsResourceGeneration bool `json:"skipGitOpsResourceGeneration,omitempty"`

	// The list of components to be nudged by this components build upon a successful result.
	// Optional.
	// +optional
	BuildNudgesRef []string `json:"build-nudges-ref,omitempty"`
}

// ComponentStatus defines the observed state of Component
type ComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions is an array of the Component's status conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Webhook URL generated by Builds
	Webhook string `json:"webhook,omitempty"`

	// ContainerImage stores the associated built container image for the component
	ContainerImage string `json:"containerImage,omitempty"`

	// The devfile model for the Component CR
	Devfile string `json:"devfile,omitempty"`

	// GitOps specific status for the Component CR
	GitOps GitOpsStatus `json:"gitops,omitempty"`

	// The last built commit id (SHA-1 checksum) from the latest component build.
	// Example: 41fbdb124775323f58fd5ce93c70bb7d79c20650.
	LastBuiltCommit string `json:"lastBuiltCommit,omitempty"`

	// The last digest image component promoted with.
	// Example: quay.io/someorg/somerepository@sha256:5ca85b7f7b9da18a9c4101e81ee1d9bac35ac2b0b0221908ff7389204660a262.
	LastPromotedImage string `json:"lastPromotedImage,omitempty"`

	// The list of names of Components whose builds nudge this resource (their spec.build-nudges-ref[] references this component)
	BuildNudgedBy []string `json:"build-nudged-by,omitempty"`
}

// GitOpsStatus contains GitOps repository-specific status for the component
type GitOpsStatus struct {
	// RepositoryURL is the gitops repository URL for the component
	RepositoryURL string `json:"repositoryURL,omitempty"`

	// Branch is the git branch used for the gitops repository
	Branch string `json:"branch,omitempty"`

	// Context is the path within the gitops repository used for the gitops resources
	Context string `json:"context,omitempty"`

	// ResourceGenerationSkipped is whether or not GitOps resource generation was skipped for the component
	ResourceGenerationSkipped bool `json:"resourceGenerationSkipped,omitempty"`

	// CommitID is the most recent commit ID in the GitOps repository for this component
	CommitID string `json:"commitID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Component is the Schema for the components API.    For description, refer to <a href="https://konflux-ci.dev/docs/reference/kube-apis/application-api/"> Hybrid Application Service Kube API </a>
// +kubebuilder:resource:path=components,shortName=hascmp;hc;comp
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[-1].status"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[-1].reason"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.conditions[-1].type"
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec"`
	Status ComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComponentList contains a list of Component
type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Component `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Component{}, &ComponentList{})
}
