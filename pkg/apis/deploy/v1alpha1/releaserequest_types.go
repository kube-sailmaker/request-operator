package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReleaseRequestSpec defines the desired state of ReleaseRequest
type ReleaseRequestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Namespace string    `json:"release-namespace"`
	Apps      []AppItem `json:"apps"`
}

// ReleaseRequestStatus defines the observed state of ReleaseRequest
type ReleaseRequestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AppItem struct {
	Name     string            `json:"name"`
	Alias    string            `json:"alias"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata""`
	Status   string            `json:"status"`
	Message  string            `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleaseRequest is the Schema for the releaserequests API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=releaserequests,scope=Namespaced
type ReleaseRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReleaseRequestSpec   `json:"spec,omitempty"`
	Status ReleaseRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleaseRequestList contains a list of ReleaseRequest
type ReleaseRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReleaseRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReleaseRequest{}, &ReleaseRequestList{})
}
