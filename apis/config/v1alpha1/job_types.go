package v1alpha1

import (
	"k8s.io/apimachinery/pkg/util/intstr"
)

type JobTab struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

const (
	JobArtifactTypeTest     = "test"
	JobArtifactTypeBuild    = "build"
	JobArtifactTypeExternal = "external"
)

type JobArtifactType string

func JobArtifactTypeFromString(s string) JobArtifactType {
	switch JobArtifactType(s) {
	case JobArtifactTypeTest, JobArtifactTypeBuild, JobArtifactTypeExternal:
		return JobArtifactType(s)
	default:
		return JobArtifactTypeBuild
	}
}

// nolint:staticcheck
// +kubebuilder:validation:XValidation:rule="(self.type in ['test', 'build'])? has(self.source) && self.source != \"\" : true",message="source must be set only when type is 'test' or 'build'"
// +kubebuilder:validation:XValidation:rule="(self.type in ['test', 'build'])? has(self.destination) && self.destination != \"\" : true",message="destination must be set only when type is 'test' or 'build'"
// +kubebuilder:validation:XValidation:rule="(self.type == 'external')? has(self.id) && self.id != \"\" : true",message="id must be set when type is external"
// +kubebuilder:validation:XValidation:rule="(self.type == 'external')? has(self.storeId) && self.storeId != \"\" : true",message="storeId must be set when type is external"
type JobArtifact struct { //nolint:recvcheck
	// +kubebuilder:validation:Enum=test;build;external
	Type JobArtifactType `json:"type"`
	// +kubebuilder:validation:Optional
	Source string `json:"source"`
	// +kubebuilder:validation:Optional
	Destination *string `json:"destination"`
	// +kubebuilder:validation:Optional
	ID string `json:"id"`
	// +kubebuilder:validation:Optional
	StoreID *string `json:"storeId"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=100
	Configuration []KeyValue `json:"configuration"`
}

type Job struct { //nolint:recvcheck
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=""
	RunInstanceCount intstr.IntOrString `json:"runInstanceCount"`
	// +kubebuilder:default:="never"
	Timeout              intstr.IntOrString    `json:"timeout,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:default:={}
	Resources []string `json:"resources"`
	// +kubebuilder:validation:MaxItems=50
	Tasks []TaskWithCancel `json:"tasks"`
	Tabs  []JobTab         `json:"tabs,omitempty"`
	// +kubebuilder:validation:Optional
	Artifacts        []JobArtifact `json:"artifacts"`
	ElasticProfileID string        `json:"elasticProfileID,omitempty"`
}
