package v1alpha1

const (
	OriginTypeGoCD       OriginType = "gocd"
	OriginTypeConfigRepo OriginType = "config_repo"
)

type OriginType string

// // nolint:staticcheck
// // +kubebuilder:validation:XValidation:rule="self.type == 'config_repo' ? has(self.id) : !has(self.id)",message="id must be set only when the type is 'config_repo'"
type Origin struct {
	// +kubebuilder:validation:Enum=gocd;config_repo
	Type OriginType `json:"type"`
	// +kubebuilder:validation:Optional
	ID string `json:"id"`
}
