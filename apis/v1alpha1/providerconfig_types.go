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
  "encoding/json"
  "reflect"

  "github.com/pkg/errors"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime/schema"

  xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// A ProviderConfigSpec defines the desired state of a ProviderConfig.
type ProviderConfigSpec struct {
  // Credentials required to authenticate to this provider.
  Credentials ProviderCredentials `json:"credentials"`
}

// ProviderCredentials required to authenticate.
type ProviderCredentials struct {
  // Source of the provider credentials.
  // +kubebuilder:validation:Enum=None;Secret;InjectedIdentity;Environment;Filesystem
  Source xpv1.CredentialsSource `json:"source"`

  xpv1.CommonCredentialSelectors `json:",inline"`
}

// A ProviderConfigStatus reflects the observed state of a ProviderConfig.
type ProviderConfigStatus struct {
  xpv1.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// A ProviderConfig configures a GoCD provider.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Cluster
type ProviderConfig struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  Spec   ProviderConfigSpec   `json:"spec"`
  Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
  metav1.TypeMeta `json:",inline"`
  metav1.ListMeta `json:"metadata,omitempty"`
  Items           []ProviderConfig `json:"items"`
}

// ProviderConfig type metadata.
var (
  ProviderConfigKind             = reflect.TypeOf(ProviderConfig{}).Name()
  ProviderConfigGroupKind        = schema.GroupKind{Group: Group, Kind: ProviderConfigKind}.String()
  ProviderConfigKindAPIVersion   = ProviderConfigKind + "." + SchemeGroupVersion.String()
  ProviderConfigGroupVersionKind = SchemeGroupVersion.WithKind(ProviderConfigKind)
)

type GocdProviderConfig struct {
  BaseURL  string `json:"baseURL"`
  Username string `json:"username"`
  Password string `json:"password"`
  Token    string `json:"token"`
  Insecure bool   `json:"insecure"`
}

func ParseGocdProviderConfig(cfg []byte) (*GocdProviderConfig, error) {
  var gpcfg GocdProviderConfig
  err := json.Unmarshal(cfg, &gpcfg)
  if err != nil {
    return nil, errors.New("cannot parse gocd provider config")
  }
  return &gpcfg, nil
}

func init() {
  SchemeBuilder.Register(&ProviderConfig{}, &ProviderConfigList{})
}
