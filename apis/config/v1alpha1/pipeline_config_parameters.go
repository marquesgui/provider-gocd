// Package pipelineconfig provides types and utilities for configuring GoCD pipeline resources,
// including pipeline parameters, environment variables, and related configuration options.
package v1alpha1

// PipelineConfigForProvider are the configurable fields of a PipelineConfig.
type PipelineConfigForProvider struct { //nolint:recvcheck
	// +kubebuilder:validation:Optional
	Group string `json:"group,omitempty"`
	// +kubebuilder:validation:Optional
	LabelTemplate string `json:"labelTemplate,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="none"
	LockBehavior LockBehavior `json:"lockBehavior,omitempty"`
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	Template string `json:"template,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={"type":"gocd"}
	Origin Origin `json:"origin,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=50
	Parameters []Parameter `json:"parameters,omitempty"`
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:Optional
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	Materials []Material `json:"materials,omitempty"`
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	Stages []Stage `json:"stages,omitempty"`
	// +kubebuilder:validation:Optional
	TrackingTool TrackingTool `json:"trackingTool,omitempty"`
	// +kubebuilder:validation:Optional
	Timer Timer `json:"timer,omitempty"`
}

type TrackingToolAttributes struct {
	URLPattern string `json:"urlPattern"`
	Regex      string `json:"regex"`
}

type TrackingTool struct {
	Type       string                 `json:"type"`
	Attributes TrackingToolAttributes `json:"attributes"`
}

type Timer struct {
	Spec          string `json:"spec"`
	OnlyOnChanges bool   `json:"onlyOnChanges"`
}

type LockBehavior string

const (
	LockBehaviorLockOnFailure      LockBehavior = "lockOnFailure"
	LockBehaviorUnlockWhenFinished LockBehavior = "unlockWhenFinished"
	LockBehaviorNone               LockBehavior = "none"
)

func LockBehaviorFromString(s string) *LockBehavior {
	var out LockBehavior
	switch LockBehavior(s) {
	case LockBehaviorLockOnFailure, LockBehaviorNone, LockBehaviorUnlockWhenFinished:
		out = LockBehavior(s)
	default:
		out = LockBehaviorNone
	}
	return &out
}

func (lb LockBehavior) String() string {
	return string(lb)
}

type Parameter struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

func SortParameters(a, b Parameter) bool {
	return a.Name+a.Value < b.Name+b.Value
}
