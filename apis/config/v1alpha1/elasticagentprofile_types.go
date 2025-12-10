/*
Copyright 2025 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ConfigProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ElasticAgentProfileParameters struct {
	ID               string           `json:"id"`
	ClusterProfileID string           `json:"clusterProfileID"`
	Properties       []ConfigProperty `json:"properties"`
}

// ElasticAgentProfileObservation are the observable fields of a ElasticAgentProfile.
type ElasticAgentProfileObservation struct {
	ConfigurableField string `json:"configurableField"`
	ObservableField   string `json:"observableField,omitempty"`
}

// A ElasticAgentProfileSpec defines the desired state of a ElasticAgentProfile.
type ElasticAgentProfileSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ElasticAgentProfileParameters `json:"forProvider"`
}

// A ElasticAgentProfileStatus represents the observed state of a ElasticAgentProfile.
type ElasticAgentProfileStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          *runtime.RawExtension `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ElasticAgentProfile is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,gocd}
type ElasticAgentProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticAgentProfileSpec   `json:"spec"`
	Status ElasticAgentProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ElasticAgentProfileList contains a list of ElasticAgentProfile
type ElasticAgentProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticAgentProfile `json:"items"`
}

// ElasticAgentProfile type metadata.
var (
	ElasticAgentProfileKind             = reflect.TypeOf(ElasticAgentProfile{}).Name()
	ElasticAgentProfileGroupKind        = schema.GroupKind{Group: Group, Kind: ElasticAgentProfileKind}.String()
	ElasticAgentProfileKindAPIVersion   = ElasticAgentProfileKind + "." + SchemeGroupVersion.String()
	ElasticAgentProfileGroupVersionKind = SchemeGroupVersion.WithKind(ElasticAgentProfileKind)
)

func init() {
	SchemeBuilder.Register(&ElasticAgentProfile{}, &ElasticAgentProfileList{})
}
