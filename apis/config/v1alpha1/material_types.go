package v1alpha1

type MaterialType string

const (
	MaterialTypeGit        MaterialType = "git"
	MaterialTypeSvn        MaterialType = "svn"
	MaterialTypeHg         MaterialType = "hg"
	MaterialTypeP4         MaterialType = "p4"
	MaterialTypeTfs        MaterialType = "tfs"
	MaterialTypePackage    MaterialType = "package"
	MaterialTypePlugin     MaterialType = "plugin"
	MaterialTypeDependency MaterialType = "dependency"
)

func (m MaterialType) String() string {
	return string(m)
}

type MaterialAttributesGit struct { //nolint:recvcheck

	Name   string `json:"name,omitempty"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
	// +kubebuilder:validation:Optional
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`
	// +kubebuilder:validation:Optional
	Destination string `json:"destination,omitempty"`
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate,omitempty"`
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter,omitempty"`
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter,omitempty"`
	// +kubebuilder:validation:Optional
	SubmoduleFolder string `json:"submoduleFolder,omitempty"`
	// +kubebuilder:validation:Optional
	ShallowClone bool `json:"shallowClone,omitempty"`
}

type MaterialAttributesSvn struct { //nolint:recvcheck

	Name string `json:"name"`
	URL  string `json:"url"`
	// +kubebuilder:validation:Optional
	Username string `json:"username"`
	// +kubebuilder:validation:Optional
	Password string `json:"password"`
	// +kubebuilder:validation:Optional
	EncryptedPassword string `json:"encryptedPassword"`
	// +kubebuilder:validation:Optional
	Destination string `json:"destination"`
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter"`
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate"`
	// +kubebuilder:validation:Optional
	CheckExternals bool `json:"checkExternals"`
}

type MaterialAttributesHg struct { //nolint:recvcheck

	Name   string `json:"name"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
	// +kubebuilder:validation:Optional
	Username string `json:"username"`
	// +kubebuilder:validation:Optional
	Password string `json:"password"`
	// +kubebuilder:validation:Optional
	EncryptedPassword string `json:"encryptedPassword"`
	// +kubebuilder:validation:Optional
	Destination string `json:"destination"`
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter"`
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate"`
}

type MaterialAttributesP4 struct { //nolint:recvcheck

	Name string `json:"name"`
	Port string `json:"port"`
	// +kubebuilder:validation:Optional
	UseTickets bool `json:"useTickets"`
	// +kubebuilder:validation:Optional
	View string `json:"view"`
	// +kubebuilder:validation:Optional
	Username string `json:"username"`
	// +kubebuilder:validation:Optional
	Password string `json:"password"`
	// +kubebuilder:validation:Optional
	EncryptedPassword string `json:"encryptedPassword"`
	// +kubebuilder:validation:Optional
	Destination string `json:"destination"`
	// +kubebuilder:validation:Optional
	Filter Filter `json:"filter"`
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
	// +kubebuilder:validation:Optional
	AutoUpdate bool `json:"autoUpdate"`
}

type MaterialAttributesTfs struct { //nolint:recvcheck

	Name              string `json:"name"`
	URL               string `json:"url"`
	ProjectPath       string `json:"projectPath"`
	Domain            string `json:"domain"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"encryptedPassword"`
	Destination       string `json:"destination"`
	AutoUpdate        bool   `json:"autoUpdate"`
	Filter            Filter `json:"filter"`
	InvertFilter      bool   `json:"invertFilter"`
}

type MaterialAttributesDependency struct { //nolint:recvcheck

	Name                string `json:"name"`
	Pipeline            string `json:"pipeline"`
	Stage               string `json:"stage"`
	AutoUpdate          bool   `json:"autoUpdate"`
	IgnoreForScheduling bool   `json:"ignoreForScheduling"`
}

type MaterialAttributesPackage struct {
	Ref *string `json:"ref"`
}

type MaterialAttributesPlugin struct { //nolint:recvcheck

	Ref         string `json:"ref"`
	Destination string `json:"destination"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Filter Filter `json:"filter"`
	// +kubebuilder:validation:Optional
	InvertFilter bool `json:"invertFilter"`
}

type Filter struct { //nolint:recvcheck
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Ignore []string `json:"ignore"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:default:={}
	Includes []string `json:"includes"`
}

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
	// +kubebuilder:validation:Optional
	GitAttributes *MaterialAttributesGit `json:"gitAttributes,omitempty"`
	// +kubebuilder:validation:Optional
	SvnAttributes *MaterialAttributesSvn `json:"svnAttributes,omitempty"`
	// +kubebuilder:validation:Optional
	HgAttributes *MaterialAttributesHg `json:"hgAttributes,omitempty"`
	// +kubebuilder:validation:Optional
	P4Attributes *MaterialAttributesP4 `json:"p4Attributes,omitempty"`
	// +kubebuilder:validation:Optional
	TfsAttributes *MaterialAttributesTfs `json:"tfsAttributes,omitempty"`
	// +kubebuilder:validation:Optional
	DependencyAttributes *MaterialAttributesDependency `json:"dependencyAttributes,omitempty"`
	// +kubebuilder:validation:Optional
	PackageAttributes *MaterialAttributesPackage `json:"packageAttributes,omitempty"`
	// +kubebuilder:validation:Optional
	PluginAttributes *MaterialAttributesPlugin `json:"pluginAttributes,omitempty"`
}
