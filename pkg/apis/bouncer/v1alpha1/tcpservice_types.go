package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TcpServiceSpec defines the desired state of TcpService
// +k8s:openapi-gen=true
type TcpServiceSpec struct {
	NodePort int `json:"nodePort"`
}

// TcpServiceStatus defines the observed state of TcpService
// +k8s:openapi-gen=true
type TcpServiceStatus struct {
	Exposed bool `json:"exposed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TcpService is the Schema for the tcpservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type TcpService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TcpServiceSpec   `json:"spec,omitempty"`
	Status TcpServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TcpServiceList contains a list of TcpService
type TcpServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TcpService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TcpService{}, &TcpServiceList{})
}
