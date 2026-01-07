package v1alpha1

import (
	"reflect"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// TaskPluggableAttributesPluginConfiguration defines the plugin configuration for a pluggable task in a GoCD pipeline job.
type TaskPluggableAttributesPluginConfiguration struct {
	// ID is the identifier of the plugin.
	ID string `json:"id"`
	// Version specifies the version of the plugin to use.
	Version string `json:"version"`
}

// TaskPluggableAttributes defines the attributes for a pluggable task in a GoCD pipeline job.
type TaskPluggableAttributes struct { //nolint:recvcheck
	// RunIf specifies the conditions under which the pluggable task should run.
	RunIf []TaskAttributesRunIfTypes `json:"runIf"`
	// PluginConfiguration contains the plugin configuration details.
	PluginConfiguration TaskPluggableAttributesPluginConfiguration `json:"pluginConfiguration"`
	// Configuration is an optional list of key-value pairs for additional plugin configuration.
	// +kubebuilder:validation:Optional
	Configuration []KeyValue `json:"configuration"`
}

// TaskRakeAttributes defines the attributes for a rake task in a GoCD pipeline job.
type TaskRakeAttributes struct { //nolint:recvcheck
	// RunIf specifies the conditions under which the rake task should run.
	RunIf []TaskAttributesRunIfTypes `json:"runIf"`
	// BuildFile is the path to the Rake build file.
	BuildFile string `json:"buildFile"`
	// Target is the Rake target to execute.
	Target string `json:"target"`
	// WorkingDirectory specifies the directory in which to run the Rake task.
	WorkingDirectory string `json:"workingDirectory"`
}

// TaskNantAttributes defines the attributes for a nant task in a GoCD pipeline job.
type TaskNantAttributes struct { //nolint:recvcheck
	// RunIf specifies the conditions under which the nant task should run.
	RunIf []TaskAttributesRunIfTypes `json:"runIf"`
	// BuildFile is the path to the Nant build file.
	BuildFile string `json:"buildFile"`
	// Target is the Nant target to execute.
	Target string `json:"target"`
	// NantPath specifies the path to the Nant executable.
	NantPath string `json:"nantPath"`
	// WorkingDirectory specifies the directory in which to run the Nant task.
	WorkingDirectory string `json:"workingDirectory"`
}

// TaskAntAttributes defines the attributes for an ant task in a GoCD pipeline job.
type TaskAntAttributes struct { //nolint:recvcheck
	// RunIf specifies the conditions under which the ant task should run.
	RunIf []TaskAttributesRunIfTypes `json:"runIf"`
	// BuildFile is the path to the Ant build file.
	BuildFile string `json:"buildFile"`
	// Target is the Ant target to execute.
	Target string `json:"target"`
	// WorkingDirectory specifies the directory in which to run the Ant task.
	WorkingDirectory string `json:"workingDirectory"`
}

// TaskFetchAttributesArtifactOrigin represents the origin of an artifact for a fetch task.
// It specifies whether the artifact is from GoCD or an external source.
type TaskFetchAttributesArtifactOrigin string

const (
	// TaskFetchAttributesArtifactOriginGoCD indicates the artifact originates from GoCD.
	TaskFetchAttributesArtifactOriginGoCD TaskFetchAttributesArtifactOrigin = "gocd"
	// TaskFetchAttributesArtifactOriginExternal indicates the artifact originates from an external source.
	TaskFetchAttributesArtifactOriginExternal TaskFetchAttributesArtifactOrigin = "external"
)

// TaskFetchAttributes defines the attributes for a fetch task in a GoCD pipeline job.
type TaskFetchAttributes struct { //nolint:recvcheck
	// ArtifactOrigin specifies the origin of the artifact to fetch (e.g., gocd, external).
	// +kubebuilder:validation:Enum=gocd;external
	ArtifactOrigin TaskFetchAttributesArtifactOrigin `json:"artifactOrigin"`
	// RunIf specifies the conditions under which the fetch task should run.
	RunIf []TaskAttributesRunIfTypes `json:"runIf"`
	// Pipeline is the name of the pipeline from which to fetch the artifact.
	Pipeline string `json:"pipeline"`
	// Stage is the name of the stage from which to fetch the artifact.
	Stage string `json:"stage"`
	// Job is the name of the job from which to fetch the artifact.
	Job string `json:"job"`
	// Source is the source path of the artifact to fetch.
	Source string `json:"source"`
	// IsSourceAFile indicates if the source is a file.
	IsSourceAFile bool `json:"isSourceAFile"`
	// Destination is the destination path where the artifact will be placed.
	Destination string `json:"destination"`
	// ArtifactID is the identifier of the artifact to fetch.
	ArtifactID string `json:"artifactId"`
	// Configuration is an optional list of key-value pairs for additional configuration.
	// +kubebuilder:validation:Optional
	Configuration []KeyValue `json:"configuration"`
}

// TaskExecAttributes defines the attributes for an exec task in a GoCD pipeline job.
type TaskExecAttributes struct { //nolint:recvcheck
	// RunIf specifies the conditions under which the task should run.
	RunIf []TaskAttributesRunIfTypes `json:"runIf"`
	// Command is the executable to run.
	Command string `json:"command"`
	// Arguments is a list of arguments to pass to the command.
	// +kubebuilder:validation:Optional
	Arguments []string `json:"arguments,omitempty"`
	// WorkingDirectory specifies the directory in which to run the command.
	WorkingDirectory *string `json:"workingDirectory,omitempty"`
}

// TaskType defines the supported types of tasks for GoCD pipeline jobs.
type TaskType string

const (
	// TaskTypeExec represents an exec task.
	TaskTypeExec TaskType = "exec"
	// TaskTypeRake represents a rake task.
	TaskTypeRake TaskType = "rake"
	// TaskTypeFetch represents a fetch task.
	TaskTypeFetch TaskType = "fetch"
	// TaskTypeAnt represents an ant task.
	TaskTypeAnt TaskType = "ant"
	// TaskTypeNant represents a nant task.
	TaskTypeNant TaskType = "nant"
	// TaskTypePluggable represents a pluggable task.
	TaskTypePluggable TaskType = "pluggable_task"
)

// Task represents a task in a GoCD pipeline job, supporting multiple task types and their attributes.
// nolint:staticcheck
// +kubebuilder:validation:XValidation:rule="self.type == 'fetch' ? has(self.fetchAttributes) : !has(self.fetchAttributes)",message="fetchAttributes must be set only when type is 'fetch'"
// +kubebuilder:validation:XValidation:rule="self.type == 'exec' ? has(self.execAttributes) : !has(self.execAttributes)",message="execAttributes must be set only when type is 'exec'"
// +kubebuilder:validation:XValidation:rule="self.type == 'ant' ? has(self.antAttributes) : !has(self.antAttributes)",message="antAttributes must be set only when type is 'ant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'nant' ? has(self.nantAttributes) : !has(self.nantAttributes)",message="nantAttributes must be set only when type is 'nant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'rake' ? has(self.rakeAttributes) : !has(self.rakeAttributes)",message="rakeAttributes must be set only when type is 'rake'"
// +kubebuilder:validation:XValidation:rule="self.type == 'pluggable_task' ? has(self.pluggableAttributes) : !has(self.pluggableAttributes)",message="pluggableAttributes must be set only when type is 'pluggable_task'"
type Task struct {
	// Type specifies the kind of task (e.g., exec, rake, fetch, ant, nant, pluggable_task).
	// Only the attributes corresponding to the selected type should be set; others must be nil.
	// +kubebuilder:validation:Enum=exec;rake;fetch;ant;nant;pluggable_task
	Type TaskType `json:"type"`
	// ExecAttributes contains configuration for exec tasks.
	// +kubebuilder:validation:Optional
	ExecAttributes *TaskExecAttributes `json:"execAttributes"`
	// FetchAttributes contains configuration for fetch tasks.
	// +kubebuilder:validation:Optional
	FetchAttributes *TaskFetchAttributes `json:"fetchAttributes"`
	// AntAttributes contains configuration for ant tasks.
	// +kubebuilder:validation:Optional
	AntAttributes *TaskAntAttributes `json:"antAttributes"`
	// NantAttributes contains configuration for nant tasks.
	// +kubebuilder:validation:Optional
	NantAttributes *TaskNantAttributes `json:"nantAttributes"`
	// RakeAttributes contains configuration for rake tasks.
	// +kubebuilder:validation:Optional
	RakeAttributes *TaskRakeAttributes `json:"rakeAttributes"`
	// PluggableAttributes contains configuration for pluggable tasks.
	// +kubebuilder:validation:Optional
	PluggableAttributes *TaskPluggableAttributes `json:"pluggableAttributes"`
}

// TaskAttributesRunIfTypes defines the possible conditions for running a task in a GoCD pipeline job.
type TaskAttributesRunIfTypes string

const (
	// TaskExecAttributesRunIfTypesPassed runs the task if the previous task passed.
	TaskExecAttributesRunIfTypesPassed = "passed"
	// TaskExecAttributesRunIfTypesFailed runs the task if the previous task failed.
	TaskExecAttributesRunIfTypesFailed = "failed"
	// TaskExecAttributesRunIfTypesAny runs the task regardless of the previous task's result.
	TaskExecAttributesRunIfTypesAny = "any"
)

// TaskPluggableAttributesWithCancel defines the attributes for a pluggable task in a GoCD pipeline job, including an optional onCancel task.
type TaskPluggableAttributesWithCancel struct { //nolint:recvcheck
	// TaskPluggableAttributes contains the standard pluggable task configuration.
	TaskPluggableAttributes `json:"-,inline"`
	// OnCancel specifies a task to execute if the main pluggable task is canceled.
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

// TaskRakeAttributesWithCancel defines the attributes for a rake task in a GoCD pipeline job, including an optional onCancel task.
type TaskRakeAttributesWithCancel struct { //nolint:recvcheck
	// TaskRakeAttributes contains the standard rake task configuration.
	TaskRakeAttributes `json:",inline"`
	// OnCancel specifies a task to execute if the main rake task is canceled.
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

// TaskNantAttributesWithCancel defines the attributes for a nant task in a GoCD pipeline job, including an optional onCancel task.
type TaskNantAttributesWithCancel struct { //nolint:recvcheck
	// TaskNantAttributes contains the standard nant task configuration.
	TaskNantAttributes `json:",inline"`
	// OnCancel specifies a task to execute if the main nant task is canceled.
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

// TaskAntAttributesWithCancel defines the attributes for an ant task in a GoCD pipeline job, including an optional onCancel task.
type TaskAntAttributesWithCancel struct { //nolint:recvcheck
	// TaskAntAttributes contains the standard ant task configuration.
	TaskAntAttributes `json:",inline"`
	// OnCancel specifies a task to execute if the main ant task is canceled.
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

// TaskFetchAttributesWithCancel defines the attributes for a fetch task in a GoCD pipeline job, including an optional onCancel task.
type TaskFetchAttributesWithCancel struct { //nolint:recvcheck
	// TaskFetchAttributes contains the standard fetch task configuration.
	TaskFetchAttributes `json:"-,inline"`
	// OnCancel specifies a task to execute if the main fetch task is canceled.
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel"`
}

// TaskExecAttributesWithCancel defines the attributes for an exec task in a GoCD pipeline job, including an optional onCancel task.
type TaskExecAttributesWithCancel struct { //nolint:recvcheck
	// TaskExecAttributes contains the standard exec task configuration.
	TaskExecAttributes `json:",inline"`
	// OnCancel specifies a task to execute if the main exec task is canceled.
	// +kubebuilder:validation:Optional
	OnCancel *Task `json:"onCancel,omitempty"`
}

// TaskWithCancel represents a task in a GoCD pipeline job, supporting multiple task types and their attributes.
// nolint:staticcheck
// +kubebuilder:validation:XValidation:rule="self.type == 'exec' ? has(self.execAttributes) : !has(self.execAttributes)",message="execAttributes must be set only when type is 'exec'"
// +kubebuilder:validation:XValidation:rule="self.type == 'fetch' ? has(self.fetchAttributes) : !has(self.fetchAttributes)",message="fetchAttributes must be set only when type is 'fetch'"
// +kubebuilder:validation:XValidation:rule="self.type == 'ant' ? has(self.antAttributes) : !has(self.antAttributes)",message="antAttributes must be set only when type is 'ant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'nant' ? has(self.nantAttributes) : !has(self.nantAttributes)",message="nantAttributes must be set only when type is 'nant'"
// +kubebuilder:validation:XValidation:rule="self.type == 'rake' ? has(self.rakeAttributes) : !has(self.rakeAttributes)",message="rakeAttributes must be set only when type is 'rake'"
// +kubebuilder:validation:XValidation:rule="self.type == 'pluggable_task' ? has(self.pluggableAttributes) : !has(self.pluggableAttributes)",message="pluggableAttributes must be set only when type is 'pluggable_task'"
type TaskWithCancel struct { //nolint:recvcheck
	// Type specifies the kind of task (e.g., exec, rake, fetch, ant, nant, pluggable_task).
	// Only the attributes corresponding to the selected type should be set; others must be nil.
	// +kubebuilder:validation:Enum=exec;rake;fetch;ant;nant;pluggable_task
	Type TaskType `json:"type"`
	// ExecAttributes contains configuration for exec tasks.
	// +kubebuilder:validation:Optional
	ExecAttributes *TaskExecAttributesWithCancel `json:"execAttributes"`
	// FetchAttributes contains configuration for fetch tasks.
	// +kubebuilder:validation:Optional
	FetchAttributes *TaskFetchAttributesWithCancel `json:"fetchAttributes"`
	// AntAttributes contains configuration for ant tasks.
	// +kubebuilder:validation:Optional
	AntAttributes *TaskAntAttributesWithCancel `json:"antAttributes"`
	// NantAttributes contains configuration for nant tasks.
	// +kubebuilder:validation:Optional
	NantAttributes *TaskNantAttributesWithCancel `json:"nantAttributes"`
	// RakeAttributes contains configuration for rake tasks.
	// +kubebuilder:validation:Optional
	RakeAttributes *TaskRakeAttributesWithCancel `json:"rakeAttributes"`
	// PluggableAttributes contains configuration for pluggable tasks.
	// +kubebuilder:validation:Optional
	PluggableAttributes *TaskPluggableAttributesWithCancel `json:"pluggableAttributes"`
}

// JobTab represents a tab in the job UI, allowing custom display of files or reports.
type JobTab struct {
	// Name is the display name of the tab.
	Name string `json:"name"`
	// Path is the file path or pattern to display in the tab.
	Path string `json:"path"`
}

// JobArtifactType is a string type for specifying the type of a job artifact.
type JobArtifactType string

const (
	// JobArtifactTypeTest represents a test artifact.
	JobArtifactTypeTest = "test"
	// JobArtifactTypeBuild represents a build artifact.
	JobArtifactTypeBuild = "build"
	// JobArtifactTypeExternal represents an external artifact.
	JobArtifactTypeExternal = "external"
)

// JobArtifactTypeFromString returns a JobArtifactType from a string value.
// If the value is not recognized, it defaults to JobArtifactTypeBuild.
func JobArtifactTypeFromString(s string) JobArtifactType {
	switch JobArtifactType(s) {
	case JobArtifactTypeTest, JobArtifactTypeBuild, JobArtifactTypeExternal:
		return JobArtifactType(s)
	default:
		return JobArtifactTypeBuild
	}
}

// JobArtifact represents an artifact produced or consumed by a job.
// The fields required depend on the artifact type.
// For 'test' and 'build', source and destination are required.
// For 'external', id and storeId are required.
// nolint:staticcheck
// +kubebuilder:validation:XValidation:rule="(self.type in ['test', 'build'])? has(self.source) && self.source != \"\" : true",message="source must be set only when type is 'test' or 'build'"
// +kubebuilder:validation:XValidation:rule="(self.type in ['test', 'build'])? has(self.destination) && self.destination != \"\" : true",message="destination must be set only when type is 'test' or 'build'"
// +kubebuilder:validation:XValidation:rule="(self.type == 'external')? has(self.id) && self.id != \"\" : true",message="id must be set when type is external"
// +kubebuilder:validation:XValidation:rule="(self.type == 'external')? has(self.storeId) && self.storeId != \"\" : true",message="storeId must be set when type is external"
type JobArtifact struct { //nolint:recvcheck
	// Type specifies the type of artifact. Allowed values: test, build, external.
	// +kubebuilder:validation:Enum=test;build;external
	Type JobArtifactType `json:"type"`
	// Source is the source path of the artifact (required for test/build).
	// +kubebuilder:validation:Optional
	Source string `json:"source"`
	// Destination is the destination path for the artifact (required for test/build).
	// +kubebuilder:validation:Optional
	Destination *string `json:"destination"`
	// ID is the identifier for the external artifact (required for external).
	// +kubebuilder:validation:Optional
	ID string `json:"id"`
	// StoreID is the identifier of the external artifact store (required for external).
	// +kubebuilder:validation:Optional
	StoreID *string `json:"storeId"`
	// Configuration is an optional list of key-value pairs for artifact configuration.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=100
	Configuration []KeyValue `json:"configuration"`
}

// Job represents a job within a stage in a GoCD pipeline.
type Job struct { //nolint:recvcheck
	// Name is the identifier for the job.
	Name string `json:"name"`
	// RunInstanceCount specifies the number of instances to run, as an integer or string.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=""
	RunInstanceCount intstr.IntOrString `json:"runInstanceCount"`
	// Timeout specifies the job timeout, as an integer or string.
	// +kubebuilder:default:="never"
	Timeout intstr.IntOrString `json:"timeout,omitempty"`
	// EnvironmentVariables is a list of environment variables available to the job.
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
	// Resources is a list of resources required by the job.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:default:={}
	Resources []string `json:"resources"`
	// Tasks is a list of tasks to be executed by the job.
	// +kubebuilder:validation:MaxItems=50
	Tasks []TaskWithCancel `json:"tasks"`
	// Tabs defines custom tabs to display job information in the GoCD UI.
	Tabs []JobTab `json:"tabs,omitempty"`
	// Artifacts is a list of artifacts produced by the job.
	// +kubebuilder:validation:Optional
	Artifacts []JobArtifact `json:"artifacts"`
	// ElasticProfileID specifies the elastic agent profile to use for the job.
	ElasticProfileID string `json:"elasticProfileID,omitempty"`
}

// StageApprovalAuthorization defines the authorization configuration for stage approval in a GoCD pipeline.
type StageApprovalAuthorization struct {
	// Users is a list of users authorized to approve the stage.
	Users []string `json:"users"`
	// Roles is a list of roles authorized to approve the stage.
	Roles []string `json:"roles"`
}

// StageApprovalType defines the type of approval required for a stage in a GoCD pipeline.
type StageApprovalType string

const (
	// StageApprovalTypeManual indicates manual approval is required.
	StageApprovalTypeManual StageApprovalType = "manual"
	// StageApprovalTypeSuccess indicates automatic approval upon successful completion of the previous stage.
	StageApprovalTypeSuccess StageApprovalType = "success"
)

// StageApproval defines the approval configuration for a stage in a GoCD pipeline.
type StageApproval struct {
	// Type specifies the approval type, such as manual or success.
	// +kubebuilder:validation:Enum=manual;success
	Type StageApprovalType `json:"type"`
	// AllowOnlyOnSuccess indicates if approval is allowed only when the previous stage succeeds.
	AllowOnlyOnSuccess bool `json:"allowOnlyOnSuccess,omitempty"`
	// Authorization defines the authorization configuration for stage approval.
	Authorization StageApprovalAuthorization `json:"authorization"`
}

// Stage represents a stage in a GoCD pipeline.
type Stage struct {
	// Name is the identifier for the stage.
	Name string `json:"name"`
	// FetchMaterials determines if materials should be fetched before running the stage.
	FetchMaterials bool `json:"fetchMaterials"`
	// CleanWorkingDir specifies if the working directory should be cleaned before the stage runs.
	CleanWorkingDir bool `json:"cleanWorkingDir"`
	// NeverCleanupArtifacts indicates if artifacts should never be cleaned up for this stage.
	NeverCleanupArtifacts bool `json:"neverCleanupArtifacts"`
	// Approval defines the approval configuration for the stage.
	Approval StageApproval `json:"approval,omitempty"`
	// EnvironmentVariables is a list of environment variables available to the stage.
	// +kubebuilder:validation:MaxItems=50
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
	// Jobs is a list of jobs that are part of this stage.
	// +kubebuilder:validation:MaxItems=50
	Jobs []Job `json:"jobs"`
}

// MaterialAttributesPlugin defines the configuration attributes for a plugin material in a GoCD pipeline.
type MaterialAttributesPlugin struct { //nolint:recvcheck
	// Ref specifies the reference to the plugin.
	Ref string `json:"ref"`
	// Destination is the folder where the plugin material will be placed.
	Destination string `json:"destination"`
	// Filter specifies file patterns to include or exclude.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Filter Filter `json:"filter"`
	// InvertFilter indicates if the filter should be inverted.
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
}

// MaterialAttributesPackage defines the configuration attributes for a package material in a GoCD pipeline.
type MaterialAttributesPackage struct {
	// Ref specifies the reference to the package repository.
	Ref *string `json:"ref"`
}

// MaterialAttributesDependency defines the configuration attributes for a dependency material in a GoCD pipeline.
type MaterialAttributesDependency struct { //nolint:recvcheck
	// Name is the identifier for the dependency material.
	Name string `json:"name"`
	// Pipeline specifies the name of the upstream pipeline.
	Pipeline string `json:"pipeline"`
	// Stage specifies the stage in the upstream pipeline to depend on.
	Stage string `json:"stage"`
	// AutoUpdate determines if GoCD should automatically poll for changes in the dependency.
	AutoUpdate bool `json:"autoUpdate"`
	// IgnoreForScheduling indicates if this dependency should be ignored for scheduling purposes.
	IgnoreForScheduling bool `json:"ignoreForScheduling"`
}

// MaterialAttributesTfs defines the configuration attributes for a Team Foundation Server (TFS) material in a GoCD pipeline.
type MaterialAttributesTfs struct { //nolint:recvcheck
	// Name is an optional identifier for the material.
	Name string `json:"name"`
	// URL specifies the TFS server URL.
	URL string `json:"url"`
	// ProjectPath specifies the path to the TFS project.
	ProjectPath string `json:"projectPath"`
	// Domain specifies the domain for TFS authentication.
	Domain string `json:"domain"`
	// Username is the username for TFS authentication.
	Username string `json:"username"`
	// Password is the password for TFS authentication.
	Password string `json:"password"`
	// EncryptedPassword is an optional encrypted password for repository authentication.
	EncryptedPassword string `json:"encryptedPassword"`
	// Destination is an optional folder where the repository will be checked out.
	Destination string `json:"destination"`
	// AutoUpdate determines if GoCD should automatically poll for changes.
	AutoUpdate bool `json:"autoUpdate"`
	// Filter specifies file patterns to include or exclude.
	Filter Filter `json:"filter"`
	// InvertFilter indicates if the filter should be inverted.
	InvertFilter bool `json:"invertFilter"`
}

// MaterialAttributesP4 defines the configuration attributes for a Perforce (P4) material in a GoCD pipeline.
type MaterialAttributesP4 struct { //nolint:recvcheck
	// Name is an optional identifier for the material.
	Name string `json:"name"`
	// Port specifies the Perforce server address.
	Port string `json:"port"`
	// UseTickets indicates whether to use Perforce tickets for authentication.
	// +kubebuilder:validation:Optional
	UseTickets bool `json:"useTickets"`
	// View specifies the Perforce view mapping.
	// +kubebuilder:validation:Optional
	View string `json:"view"`
	// Username is an optional username for repository authentication.
	// +kubebuilder:validation:Optional
	Username string `json:"username"`
	// Password is an optional password for repository authentication.
	// +kubebuilder:validation:Optional
	Password string `json:"password"`
	// EncryptedPassword is an optional encrypted password for repository authentication.
	// +kubebuilder:validation:Optional
	EncryptedPassword string `json:"encryptedPassword"`
	// Destination is an optional folder where the repository will be checked out.
	// +kubebuilder:validation:Optional
	Destination string `json:"destination"`
	// Filter specifies file patterns to include or exclude.
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter"`
	// InvertFilter indicates if the filter should be inverted.
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
	// AutoUpdate determines if GoCD should automatically poll for changes.
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate"`
}

// MaterialAttributesHg defines the configuration attributes for a Mercurial (Hg) material in a GoCD pipeline.
type MaterialAttributesHg struct { //nolint:recvcheck
	// Name is an optional identifier for the material.
	Name string `json:"name"`
	// URL specifies the Mercurial repository location.
	URL string `json:"url"`
	// Branch specifies the branch to track.
	Branch string `json:"branch"`
	// Username is an optional username for repository authentication.
	// +kubebuilder:validation:Optional
	Username string `json:"username"`
	// Password is an optional password for repository authentication.
	// +kubebuilder:validation:Optional
	Password string `json:"password"`
	// EncryptedPassword is an optional encrypted password for repository authentication.
	// +kubebuilder:validation:Optional
	EncryptedPassword string `json:"encryptedPassword"`
	// Destination is an optional folder where the repository will be checked out.
	// +kubebuilder:validation:Optional
	Destination string `json:"destination"`
	// Filter specifies file patterns to include or exclude.
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter"`
	// InvertFilter indicates if the filter should be inverted.
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
	// AutoUpdate determines if GoCD should automatically poll for changes.
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate"`
}

// MaterialAttributesSvn defines the configuration attributes for an SVN material in a GoCD pipeline.
type MaterialAttributesSvn struct { //nolint:recvcheck
	// Name is an optional identifier for the material.
	Name string `json:"name"`
	// URL specifies the SVN repository location.
	URL string `json:"url"`
	// Username is an optional username for repository authentication.
	// +kubebuilder:validation:Optional
	Username string `json:"username"`
	// Password is an optional password for repository authentication.
	// +kubebuilder:validation:Optional
	Password string `json:"password"`
	// EncryptedPassword is an optional encrypted password for repository authentication.
	// +kubebuilder:validation:Optional
	EncryptedPassword string `json:"encryptedPassword"`
	// Destination is an optional folder where the repository will be checked out.
	// +kubebuilder:validation:Optional
	Destination string `json:"destination"`
	// Filter specifies file patterns to include or exclude.
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter"`
	// InvertFilter indicates if the filter should be inverted.
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
	// AutoUpdate determines if GoCD should automatically poll for changes.
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate"`
	// CheckExternals specifies whether to check SVN externals.
	// +kubebuilder:validation:Optional
	CheckExternals bool `json:"checkExternals"`
}

// Filter defines file patterns to include or exclude for a material in a GoCD pipeline.
type Filter struct { //nolint:recvcheck
	// Ignore specifies file patterns to be excluded from triggering pipeline runs.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Ignore []string `json:"ignore"`
	// Includes specifies file patterns to be included for triggering pipeline runs.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:default:={}
	Includes []string `json:"includes"`
}

// MaterialAttributesGit defines the configuration attributes for a Git material in a GoCD pipeline.
type MaterialAttributesGit struct { //nolint:recvcheck
	// Name is an optional identifier for the material.
	Name string `json:"name,omitempty"`
	// URL specifies the Git repository location.
	URL string `json:"url"`
	// Branch specifies the branch to track.
	Branch string `json:"branch"`
	// +kubebuilder:validation:Optional
	// Username is an optional username for repository authentication.
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Optional
	// Password is an optional password for repository authentication.
	Password string `json:"password,omitempty"`
	// Destination is an optional folder where the repository will be cloned.
	// +kubebuilder:validation:Optional
	Destination string `json:"destination,omitempty"`
	// AutoUpdate determines if GoCD should automatically poll for changes.
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate,omitempty"`
	// Filter specifies file patterns to include or exclude.
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter,omitempty"`
	// InvertFilter indicates if the filter should be inverted.
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter,omitempty"`
	// SubmoduleFolder specifies a subfolder for Git submodules.
	// +kubebuilder:validation:Optional
	SubmoduleFolder string `json:"submoduleFolder,omitempty"`
	// ShallowClone determines if a shallow clone should be performed.
	// +kubebuilder:validation:Optional
	ShallowClone bool `json:"shallowClone,omitempty"`
}

// MaterialType defines the supported types of materials for GoCD pipelines.
type MaterialType string

const (
	// MaterialTypeGit represents a Git material.
	MaterialTypeGit MaterialType = "git"
	// MaterialTypeSvn represents a Subversion (SVN) material.
	MaterialTypeSvn MaterialType = "svn"
	// MaterialTypeHg represents a Mercurial (Hg) material.
	MaterialTypeHg MaterialType = "hg"
	// MaterialTypeP4 represents a Perforce (P4) material.
	MaterialTypeP4 MaterialType = "p4"
	// MaterialTypeTfs represents a Team Foundation Server (TFS) material.
	MaterialTypeTfs MaterialType = "tfs"
	// MaterialTypePackage represents a package material.
	MaterialTypePackage MaterialType = "package"
	// MaterialTypePlugin represents a plugin material.
	MaterialTypePlugin MaterialType = "plugin"
	// MaterialTypeDependency represents a dependency material.
	MaterialTypeDependency MaterialType = "dependency"
)

// String returns the string representation of the MaterialType.
func (m MaterialType) String() string {
	return string(m)
}

// Material represents a material source for a GoCD pipeline.
// Type specifies the kind of material (e.g., git, svn, hg, p4, tfs, dependency, package, plugin).
// Each material type has its own set of attributes, only one of which should be set at a time.
// // nolint:staticcheck
// // +kubebuilder:validation:XValidation:rule="self.type == 'git' ? has(self.gitAttributes) : !has(self.gitAttributes)",message="gitAttributes must be set only when type is 'git'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'hg' ? has(self.hgAttributes) : !has(self.hgAttributes)",message="hgAttributes must be set only when type is 'hg'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'svn' ? has(self.svnAttributes) : !has(self.svnAttributes)",message="svnAttributes must be set only when type is 'svn'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'p4' ? has(self.p4Attributes) : !has(self.p4Attributes)",message="p4Attributes must be set only when type is 'p4'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'tfs' ? has(self.tfsAttributes) : !has(self.tfsAttributes)",message="tfsAttributes must be set only when type is 'tfs'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'dependency' ? has(self.dependencyAttributes) : !has(self.dependencyAttributes)",message="dependencyAttributes must be set only when type is 'dependency'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'package' ? has(self.packageAttributes) : !has(self.packageAttributes)",message="packageAttributes must be set only when type is 'package'"
// // +kubebuilder:validation:XValidation:rule="self.type == 'plugin' ? has(self.pluginAttributes) : !has(self.pluginAttributes)",message="pluginAttributes must be set only when type is 'plugin'"
type Material struct { //nolint:recvcheck
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=git;svn;hg;p4;tfs;dependency;package;plugin
	Type MaterialType `json:"type"`
	// GitAttributes contains configuration for git materials.
	// +kubebuilder:validation:Optional
	GitAttributes *MaterialAttributesGit `json:"gitAttributes,omitempty"`
	// SvnAttributes contains configuration for svn materials.
	// +kubebuilder:validation:Optional
	SvnAttributes *MaterialAttributesSvn `json:"svnAttributes,omitempty"`
	// HgAttributes contains configuration for mercurial materials.
	// +kubebuilder:validation:Optional
	HgAttributes *MaterialAttributesHg `json:"hgAttributes,omitempty"`
	// P4Attributes contains configuration for perforce materials.
	// +kubebuilder:validation:Optional
	P4Attributes *MaterialAttributesP4 `json:"p4Attributes,omitempty"`
	// TfsAttributes contains configuration for TFS materials.
	// +kubebuilder:validation:Optional
	TfsAttributes *MaterialAttributesTfs `json:"tfsAttributes,omitempty"`
	// DependencyAttributes contains configuration for dependency materials.
	// +kubebuilder:validation:Optional
	DependencyAttributes *MaterialAttributesDependency `json:"dependencyAttributes,omitempty"`
	// PackageAttributes contains configuration for package materials.
	// +kubebuilder:validation:Optional
	PackageAttributes *MaterialAttributesPackage `json:"packageAttributes,omitempty"`
	// PluginAttributes contains configuration for plugin materials.
	// +kubebuilder:validation:Optional
	PluginAttributes *MaterialAttributesPlugin `json:"pluginAttributes,omitempty"`
}

// Parameter represents a key-value pair used for parameterizing GoCD pipeline configurations.
type Parameter struct {
	// Name is the identifier for the parameter.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	// Value is the value assigned to the parameter.
	Value string `json:"value"`
}

// OriginType is a string type that specifies the origin of the pipeline configuration.
type OriginType string

const (
	// OriginTypeGoCD indicates the pipeline is defined natively in GoCD.
	OriginTypeGoCD OriginType = "gocd"
	// OriginTypeConfigRepo indicates the pipeline is defined in a configuration repository.
	OriginTypeConfigRepo OriginType = "config_repo"
)

// // nolint:staticcheck
// // +kubebuilder:validation:XValidation:rule="self.type == 'config_repo' ? has(self.id) : !has(self.id)",message="id must be set only when the type is 'config_repo'"
type Origin struct {
	// +kubebuilder:validation:Enum=gocd;config_repo
	Type OriginType `json:"type"`
	// +kubebuilder:validation:Optional
	ID string `json:"id"`
}

// LockBehavior defines the locking behavior for a GoCD pipeline.
type LockBehavior string

const (
	// LockBehaviorLockOnFailure locks the pipeline only if a stage fails.
	LockBehaviorLockOnFailure LockBehavior = "lockOnFailure"
	// LockBehaviorUnlockWhenFinished unlocks the pipeline when the stage finishes, regardless of result.
	LockBehaviorUnlockWhenFinished LockBehavior = "unlockWhenFinished"
	// LockBehaviorNone means no locking is applied to the pipeline.
	LockBehaviorNone LockBehavior = "none"
)

// LockBehaviorFromString converts a string to a LockBehavior pointer.
// If the input string matches a known LockBehavior value, it returns a pointer to that value.
// Otherwise, it defaults to LockBehaviorNone.
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

// String returns the string representation of the LockBehavior value.
func (lb LockBehavior) String() string {
	return string(lb)
}

// PipelineConfigForProvider defines the configuration for a GoCD pipeline as required by the provider.
type PipelineConfigForProvider struct { //nolint:recvcheck
	// Group is the GoCD pipeline group to which this pipeline belongs.
	// +kubebuilder:validation:Optional
	Group string `json:"group,omitempty"`
	// LabelTemplate is the template used for pipeline labeling.
	// +kubebuilder:validation:Optional
	LabelTemplate string `json:"labelTemplate,omitempty"`
	// LockBehavior defines the locking behavior for the pipeline.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="none"
	LockBehavior LockBehavior `json:"lockBehavior,omitempty"`
	// Name is the name of the pipeline.
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// Template is the name of the pipeline template to use.
	// +kubebuilder:validation:Optional
	Template string `json:"template,omitempty"`
	// Origin specifies the origin of the pipeline configuration.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={"type":"gocd"}
	Origin Origin `json:"origin,omitempty"`
	// Parameters is a list of parameters for the pipeline.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=50
	Parameters []Parameter `json:"parameters,omitempty"`
	// EnvironmentVariables is a list of environment variables for the pipeline.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:Optional
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
	// Materials is a list of materials (sources) for the pipeline.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	Materials []Material `json:"materials,omitempty"`
	// Stages is a list of stages for the pipeline.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	Stages []Stage `json:"stages,omitempty"`
	// TrackingTool specifies the tracking tool configuration for the pipeline.
	// +kubebuilder:validation:Optional
	TrackingTool TrackingTool `json:"trackingTool,omitempty"`
	// Timer specifies the cron time when the pipeline should be triggered.
	// +kubebuilder:validation:Optional
	Timer Timer `json:"timer,omitempty"`
}

// TrackingToolAttributes defines the attributes for a tracking tool used in a GoCD pipeline.
type TrackingToolAttributes struct {
	// URLPattern specifies the URL pattern for linking to issues or tickets.
	URLPattern string `json:"urlPattern"`
	// Regex defines the regular expression used to identify references in commit messages.
	Regex string `json:"regex"`
}

// TrackingTool represents a tool used for tracking issues or tickets in a GoCD pipeline configuration.
type TrackingTool struct {
	// Type specifies the kind of tracking tool (e.g., JIRA, Mingle).
	Type string `json:"type"`
	// Attributes contains the configuration details specific to the tracking tool.
	Attributes TrackingToolAttributes `json:"attributes"`
}

// Timer defines the scheduling configuration for a GoCD pipeline.
type Timer struct {
	// Spec specifies the cron-like schedule for triggering the pipeline.
	Spec string `json:"spec"`
	// OnlyOnChanges indicates whether the pipeline should only be triggered when there are material changes.
	OnlyOnChanges bool `json:"onlyOnChanges"`
}

// A PipelineConfigSpec defines the desired state of a PipelineConfig.
type PipelineConfigSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PipelineConfigForProvider `json:"forProvider"`
}

// A PipelineConfigStatus represents the observed state of a PipelineConfig.
type PipelineConfigStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          *runtime.RawExtension `json:"atProvider,omitempty"`
	// EnvironmentVariableHashes stores the hashes of the environment variables
	// to detect changes in secure variables.
	// +optional
	EnvironmentVariableHashes map[string]string `json:"environmentVariableHashes,omitempty"`
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
