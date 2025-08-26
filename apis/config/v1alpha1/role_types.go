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
)

// RoleParametersAttributes are the attributes of a role.
type RoleParametersAttributes struct {
  // The list of users belongs to the role.
  Users []string `json:"users,omitempty"`
  // The authorization configuration identifier.
  AuthConfigID string `json:"authConfigId,omitempty"`
  // The list of configuration properties that represent the configuration of this plugin role.
  Properties []KeyValue `json:"properties,omitempty"`
}

type RoleParametersPolicy struct {
  // The type of permission which can be either allow or deny.
  // +kubebuilder:validation:Enum=allow;deny
  Permission string `json:"permission"`
  // The action that is being controlled via this rule. Can be one of *, view, administer
  // +kubebuilder:validation:Enum=view;administer
  Action string `json:"action"`
  // The type of entity that the rule is applied on. Can be one of *, environment.
  // +kubebuilder:validation:Enum=environment;"*"
  Type string `json:"type"`
  // The actual entity on which the rule is applied. Resource should be the name of the entity or a wildcard which matches one or more entities.
  Resource string `json:"resource"`
}

// RoleParameters are the configurable fields of a role.
type RoleParameters struct {
  // The name of the role.
  Name string `json:"name"`
  // The type of the role.
  // +kubebuilder:validation:Enum=gocd;plugin
  Type string `json:"type"`
  // The attributes of the role.
  Attributes RoleParametersAttributes `json:"attributes"`
  Policy     []RoleParametersPolicy   `json:"policy"`
}

// RoleObservation represents the observed state of a role.
type RoleObservation struct {
  Name       string                   `json:"name,omitempty"`
  Type       string                   `json:"type,omitempty"`
  Attributes RoleParametersAttributes `json:"attributes,omitempty"`
  Policy     []RoleParametersPolicy   `json:"policy,omitempty"`
  Links      EntityLinks              `json:"links"`
}

// A RoleSpec defines the desired state of a role.
type RoleSpec struct {
  xpv1.ResourceSpec `json:",inline"`
  ForProvider       RoleParameters `json:"forProvider"`
}

// A RoleStatus represents the observed state of a role.
type RoleStatus struct {
  xpv1.ResourceStatus `json:",inline"`
  AtProvider          RoleObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// nolint:staticcheck
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,gocd}
type Role struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  Spec   RoleSpec   `json:"spec"`
  Status RoleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type RoleList struct {
  metav1.TypeMeta `json:",inline"`
  metav1.ListMeta `json:"metadata,omitempty"`
  Items           []Role `json:"items"`
}

// role type metadata.
var (
  RoleKind             = reflect.TypeOf(Role{}).Name()
  RoleGroupKind        = schema.GroupKind{Group: Group, Kind: RoleKind}.String()
  RoleKindAPIVersion   = RoleKind + "." + SchemeGroupVersion.String()
  RoleGroupVersionKind = SchemeGroupVersion.WithKind(RoleKind)
)

func init() {
  SchemeBuilder.Register(&Role{}, &RoleList{})
}
