package v1alpha1

type StageApprovalAuthorization struct {
	Users []string `json:"users"`
	Roles []string `json:"roles"`
}

type StageApprovalType string

const (
	StageApprovalTypeManual  StageApprovalType = "manual"
	StageApprovalTypeSuccess StageApprovalType = "success"
)

type StageApproval struct {
	// +kubebuilder:validation:Enum=manual;success
	Type               StageApprovalType          `json:"type"`
	AllowOnlyOnSuccess bool                       `json:"allowOnlyOnSuccess,omitempty"`
	Authorization      StageApprovalAuthorization `json:"authorization"`
}

type Stage struct {
	Name                  string        `json:"name"`
	FetchMaterials        bool          `json:"fetchMaterials"`
	CleanWorkingDir       bool          `json:"cleanWorkingDir"`
	NeverCleanupArtifacts bool          `json:"neverCleanupArtifacts"`
	Approval              StageApproval `json:"approval,omitempty"`
	// +kubebuilder:validation:MaxItems=50
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
	// +kubebuilder:validation:MaxItems=50
	Jobs []Job `json:"jobs"`
}
