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

// Package apis contains Kubernetes API for the GoCD provider.
package apis

// samplev1alpha1 "github.com/crossplane/provider-gocd/apis/sample/v1alpha1"
import (
  configgocdv1alpha1 "github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
  gocdv1alpha1 "github.com/marquesgui/provider-gocd/apis/v1alpha1"
  "k8s.io/apimachinery/pkg/runtime"
)

func init() {
  // Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
  // samplev1alpha1.SchemeBuilder.AddToScheme,
  AddToSchemes = append(AddToSchemes,
    gocdv1alpha1.SchemeBuilder.AddToScheme,
    configgocdv1alpha1.SchemeBuilder.AddToScheme,
  )
}

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
  return AddToSchemes.AddToScheme(s)
}
