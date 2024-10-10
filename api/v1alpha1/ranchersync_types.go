/*
Copyright 2024.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RancherSyncSpec defines the desired state of RancherSync
type RancherSyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Api   string `json:"api"`
	Token string `json:"token"`
}

// RancherSyncStatus defines the observed state of RancherSync
type RancherSyncStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Other status fields can be added here, like progress, state, etc.
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RancherSync is the Schema for the ranchersyncs API
type RancherSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RancherSyncSpec   `json:"spec,omitempty"`
	Status RancherSyncStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RancherSyncList contains a list of RancherSync
type RancherSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RancherSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RancherSync{}, &RancherSyncList{})
}
