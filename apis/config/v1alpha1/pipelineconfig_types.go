package v1alpha1

import (
	"reflect"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// A PipelineConfigSpec defines the desired state of a PipelineConfig.
type PipelineConfigSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PipelineConfigForProvider `json:"forProvider"`
}

// A PipelineConfigStatus represents the observed state of a PipelineConfig.
type PipelineConfigStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          *runtime.RawExtension `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// nolint:staticcheck
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,gocd}
type PipelineConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineConfigSpec   `json:"spec"`
	Status PipelineConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineConfigList contains a list of PipelineConfig
type PipelineConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineConfig `json:"items"`
}

// PipelineConfig type metadata.
var (
	PipelineConfigKind             = reflect.TypeOf(PipelineConfig{}).Name()
	PipelineConfigGroupKind        = schema.GroupKind{Group: Group, Kind: PipelineConfigKind}.String()
	PipelineConfigKindAPIVersion   = PipelineConfigKind + "." + SchemeGroupVersion.String()
	PipelineConfigGroupVersionKind = SchemeGroupVersion.WithKind(PipelineConfigKind)
)

func init() {
	SchemeBuilder.Register(&PipelineConfig{}, &PipelineConfigList{})
}
