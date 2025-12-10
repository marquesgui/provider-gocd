package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type KeyValue struct {
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

func KeyValuesEqual(a, b []KeyValue) bool {
	if len(a) != len(b) {
		return false
	}
	mapA := make(map[string]string, len(a))
	for _, v := range a {
		mapA[v.Key] = v.Value
	}
	for _, v := range b {
		if mapA[v.Key] != v.Value {
			return false
		}
	}
	return true
}

type Link struct {
	Href string `json:"href"`
}

type EntityLinks struct {
	Self Link `json:"self"`
	Doc  Link `json:"doc"`
	Find Link `json:"find"`
}

type EnvironmentVariable struct { //nolint:recvcheck
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	Value string `json:"value,omitempty"`
	// +kubebuilder:validation:Optional
	ValueFrom *ValueSource `json:"valueFrom,omitempty"`
}

type ValueSource struct {
	// +kubebuilder:validation:Optional
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// +kubebuilder:validation:Optional
	SecretKeyRef *xpv1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

type ConfigMapKeySelector struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}
