// Package v1alpha1 /*
package v1alpha1

import (
  "reflect"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime/schema"

  xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// AuthorizationConfigurationParameters are the configurable fields of an AuthorizationConfiguration.
type AuthorizationConfigurationParameters struct {
  // The identifier of the authorization configuration.
  ID string `json:"id"`
  // The plugin identifier of the authorization plugin.
  PluginID string `json:"pluginId"`
  // Allow only those users to login who have explicitly been added by an administrator.
  AllowOnlyKnowUsersToLogin bool `json:"allowOnlyKnowUsersToLogin"`
  // The list of configuration properties that represent the configuration of this authorization configuration.
  Properties []KeyValue `json:"properties"`
}

// AuthorizationConfigurationObservation are the observable fields of an AuthorizationConfiguration.
type AuthorizationConfigurationObservation struct {
  // The identifier of the authorization configuration.
  ID string `json:"id,omitempty"`
  // The plugin identifier of the authorization plugin.
  PluginID string `json:"pluginId,omitempty"`
  // Allow only those users to login who have explicitly been added by an administrator.
  AllowOnlyKnowUsersToLogin bool `json:"allowOnlyKnowUsersToLogin,omitempty"`
  // The list of configuration properties that represent the configuration of this authorization configuration.
  Properties    []KeyValue  `json:"properties,omitempty"`
  Links         EntityLinks `json:"links"`
  TransactionID string      `json:"transactionId,omitempty"`
}

// An AuthorizationConfigurationSpec defines the desired state of an AuthorizationConfiguration.
type AuthorizationConfigurationSpec struct {
  xpv1.ResourceSpec `json:",inline"`
  ForProvider       AuthorizationConfigurationParameters `json:"forProvider"`
}

// An AuthorizationConfigurationStatus represents the observed state of an AuthorizationConfiguration.
type AuthorizationConfigurationStatus struct {
  xpv1.ResourceStatus `json:",inline"`
  AtProvider          AuthorizationConfigurationObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// An AuthorizationConfiguration is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Cluster",categories={crossplane,managed,gocd}
type AuthorizationConfiguration struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  Spec   AuthorizationConfigurationSpec   `json:"spec"`
  Status AuthorizationConfigurationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AuthorizationConfigurationList contains a list of AuthorizationConfiguration
type AuthorizationConfigurationList struct {
  metav1.TypeMeta `json:",inline"`
  metav1.ListMeta `json:"metadata,omitempty"`
  Items           []AuthorizationConfiguration `json:"items"`
}

// AuthorizationConfiguration type metadata.
var (
  AuthorizationConfigurationKind             = reflect.TypeOf(AuthorizationConfiguration{}).Name()
  AuthorizationConfigurationGroupKind        = schema.GroupKind{Group: Group, Kind: AuthorizationConfigurationKind}.String()
  AuthorizationConfigurationKindAPIVersion   = AuthorizationConfigurationKind + "." + SchemeGroupVersion.String()
  AuthorizationConfigurationGroupVersionKind = SchemeGroupVersion.WithKind(AuthorizationConfigurationKind)
)

func init() {
  SchemeBuilder.Register(&AuthorizationConfiguration{}, &AuthorizationConfigurationList{})
}
