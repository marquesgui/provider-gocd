package v1alpha1

type TaskType string

const (
	TaskTypeExec      TaskType = "exec"
	TaskTypeRake      TaskType = "rake"
	TaskTypeFetch     TaskType = "fetch"
	TaskTypeAnt       TaskType = "ant"
	TaskTypeNant      TaskType = "nant"
	TaskTypePluggable TaskType = "pluggable_task"
)

type TaskAttributesRunIfTypes string

const (
	TaskExecAttributesRunIfTypesPassed = "passed"
	TaskExecAttributesRunIfTypesFailed = "failed"
	TaskExecAttributesRunIfTypesAny    = "any"
)

type TaskExecAttributes struct { //nolint:recvcheck
	RunIf   []TaskAttributesRunIfTypes `json:"runIf"`
	Command string                     `json:"command"`
	// +kubebuilder:validation:Optional
	Arguments        []string `json:"arguments,omitempty"`
	WorkingDirectory *string  `json:"workingDirectory,omitempty"`
}

type TaskFetchAttributesArtifactOrigin string

const (
	TaskFetchAttributesArtifactOriginGoCD     TaskFetchAttributesArtifactOrigin = "gocd"
	TaskFetchAttributesArtifactOriginExternal TaskFetchAttributesArtifactOrigin = "external"
)

type TaskAntAttributes struct { //nolint:recvcheck
	RunIf            []TaskAttributesRunIfTypes `json:"runIf"`
	BuildFile        string                     `json:"buildFile"`
	Target           string                     `json:"target"`
	WorkingDirectory string                     `json:"workingDirectory"`
}

type TaskNantAttributes struct { //nolint:recvcheck
	RunIf            []TaskAttributesRunIfTypes `json:"runIf"`
	BuildFile        string                     `json:"buildFile"`
	Target           string                     `json:"target"`
	NantPath         string                     `json:"nantPath"`
	WorkingDirectory string                     `json:"workingDirectory"`
}

type TaskRakeAttributes struct { //nolint:recvcheck
	RunIf            []TaskAttributesRunIfTypes `json:"runIf"`
	BuildFile        string                     `json:"buildFile"`
	Target           string                     `json:"target"`
	WorkingDirectory string                     `json:"workingDirectory"`
}

type TaskPluggableAttributesPluginConfiguration struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type TaskPluggableAttributes struct { //nolint:recvcheck
	RunIf               []TaskAttributesRunIfTypes                 `json:"runIf"`
	PluginConfiguration TaskPluggableAttributesPluginConfiguration `json:"pluginConfiguration"`
	// +kubebuilder:validation:Optional
	Configuration []KeyValue `json:"configuration"`
}

// nolint:staticcheck
// +kubebuilder:validation:XValidation:rule="self.type == 'fetch' ? has(self.fetchAttributes) : !has(self.fetchAttributes)",message="fetchAttributes must be set only when type is 'fetch'"
// +kubebuilder:validation:XValidation:rule="self.type == 'exec' ? has(self.execAttributes) : !has(self.execAttributes)",message="execAttributes must be set only when type is 'exec'"
// +kubebuilder:validation:XValidation:rule="self.type == 'ant' ? has(self.antAttributes) : !has(self.antAttributes)",message="antAttributes must be set only when type is 'ant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'nant' ? has(self.nantAttributes) : !has(self.nantAttributes)",message="nantAttributes must be set only when type is 'nant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'rake' ? has(self.rakeAttributes) : !has(self.rakeAttributes)",message="rakeAttributes must be set only when type is 'rake'"
// +kubebuilder:validation:XValidation:rule="self.type == 'pluggable_task' ? has(self.pluggableAttributes) : !has(self.pluggableAttributes)",message="pluggableAttributes must be set only when type is 'pluggable_task'"
type Task struct {
	// +kubebuilder:validation:Enum=exec;rake;fetch;ant;nant;pluggable_task
	Type TaskType `json:"type"`
	// +kubebuilder:validation:Optional
	ExecAttributes *TaskExecAttributes `json:"execAttributes"`
	// +kubebuilder:validation:Optional
	FetchAttributes *TaskFetchAttributes `json:"fetchAttributes"`
	// +kubebuilder:validation:Optional
	AntAttributes *TaskAntAttributes `json:"antAttributes"`
	// +kubebuilder:validation:Optional
	NantAttributes *TaskNantAttributes `json:"nantAttributes"`
	// +kubebuilder:validation:Optional
	RakeAttributes *TaskRakeAttributes `json:"rakeAttributes"`
	// +kubebuilder:validation:Optional
	PluggableAttributes *TaskPluggableAttributes `json:"pluggableAttributes"`
}

type TaskExecAttributesWithCancel struct { //nolint:recvcheck
	TaskExecAttributes `json:",inline"`
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel,omitempty"`
}

type TaskFetchAttributes struct { //nolint:recvcheck
	// +kubebuilder:validation:Enum=gocd;external
	ArtifactOrigin TaskFetchAttributesArtifactOrigin `json:"artifactOrigin"`
	RunIf          []TaskAttributesRunIfTypes        `json:"runIf"`
	Pipeline       string                            `json:"pipeline"`
	Stage          string                            `json:"stage"`
	Job            string                            `json:"job"`
	Source         string                            `json:"source"`
	IsSourceAFile  bool                              `json:"isSourceAFile"`
	Destination    string                            `json:"destination"`
	ArtifactID     string                            `json:"artifactId"`
	// +kubebuilder:validation:Optional
	Configuration []KeyValue `json:"configuration"`
}

type TaskFetchAttributesWithCancel struct { //nolint:recvcheck
	TaskFetchAttributes `json:"-,inline"`
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

type TaskAntAttributesWithCancel struct { //nolint:recvcheck
	TaskAntAttributes `json:",inline"`
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

type TaskNantAttributesWithCancel struct { //nolint:recvcheck
	TaskNantAttributes `json:",inline"`
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

type TaskRakeAttributesWithCancel struct { //nolint:recvcheck
	TaskRakeAttributes `json:",inline"`
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

type TaskPluggableAttributesWithCancel struct { //nolint:recvcheck
	TaskPluggableAttributes `json:"-,inline"`
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

// nolint:staticcheck
// +kubebuilder:validation:XValidation:rule="self.type == 'exec' ? has(self.execAttributes) : !has(self.execAttributes)",message="execAttributes must be set only when type is 'exec'"
// +kubebuilder:validation:XValidation:rule="self.type == 'fetch' ? has(self.fetchAttributes) : !has(self.fetchAttributes)",message="fetchAttributes must be set only when type is 'fetch'"
// +kubebuilder:validation:XValidation:rule="self.type == 'ant' ? has(self.antAttributes) : !has(self.antAttributes)",message="antAttributes must be set only when type is 'ant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'nant' ? has(self.nantAttributes) : !has(self.nantAttributes)",message="nantAttributes must be set only when type is 'nant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'rake' ? has(self.rakeAttributes) : !has(self.rakeAttributes)",message="rakeAttributes must be set only when type is 'rake'"
// +kubebuilder:validation:XValidation:rule="self.type == 'pluggable_task' ? has(self.pluggableAttributes) : !has(self.pluggableAttributes)",message="pluggableAttributes must be set only when type is 'pluggable_task'"
type TaskWithCancel struct { //nolint:recvcheck
	// +kubebuilder:validation:Enum=exec;rake;fetch;ant;nant;pluggable_task
	Type TaskType `json:"type"`
	// +kubebuilder:validation:Optional
	ExecAttributes *TaskExecAttributesWithCancel `json:"execAttributes"`
	// +kubebuilder:validation:Optional
	FetchAttributes *TaskFetchAttributesWithCancel `json:"fetchAttributes"`
	// +kubebuilder:validation:Optional
	AntAttributes *TaskAntAttributesWithCancel `json:"antAttributes"`
	// +kubebuilder:validation:Optional
	NantAttributes *TaskNantAttributesWithCancel `json:"nantAttributes"`
	// +kubebuilder:validation:Optional
	RakeAttributes *TaskRakeAttributesWithCancel `json:"rakeAttributes"`
	// +kubebuilder:validation:Optional
	PluggableAttributes *TaskPluggableAttributesWithCancel `json:"pluggableAttributes"`
}
