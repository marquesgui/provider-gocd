package gocd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/marquesgui/provider-gocd/pkg/cmp"
)

const (
	acceptPipelineConfigs      = "application/vnd.go.cd.v11+json"
	pipelineConfigsServicePath = "/go/api/admin/pipelines"
	nilStr                     = "<nil>"
)

type PipelineConfigLockBehavior string

func PipelineConfigLockBehaviorFromString(s string) *PipelineConfigLockBehavior {
	var out PipelineConfigLockBehavior
	switch PipelineConfigLockBehavior(s) {
	case PipelineConfigLockBehaviorUnlockWhenFinished,
		PipelineConfigLockBehaviorLockOnFailure,
		PipelineConfigLockBehaviorNone:
		out = PipelineConfigLockBehavior(s)
		return &out
	default:
		out = PipelineConfigLockBehaviorNone
	}
	return &out
}

func (p PipelineConfigLockBehavior) String() string {
	return string(p)
}

func (p PipelineConfigLockBehavior) Equal(o PipelineConfigLockBehavior) bool {
	return p == o
}

const (
	PipelineConfigLockBehaviorUnlockWhenFinished PipelineConfigLockBehavior = "unlockWhenFinished"
	PipelineConfigLockBehaviorLockOnFailure      PipelineConfigLockBehavior = "lockOnFailure"
	PipelineConfigLockBehaviorNone               PipelineConfigLockBehavior = "none"
)

const (
	PipelineConfigOriginTypeGoCD       PipelineConfigOriginType = "gocd"
	PipelineConfigOriginTypeConfigRepo PipelineConfigOriginType = "config_repo"
)

type PipelineConfigOriginType string

func (t PipelineConfigOriginType) String() string {
	return string(t)
}

func PipelineConfigOriginTypeFromString(s string) *PipelineConfigOriginType {
	var out PipelineConfigOriginType
	switch PipelineConfigOriginType(s) {
	case PipelineConfigOriginTypeGoCD, PipelineConfigOriginTypeConfigRepo:
		out = PipelineConfigOriginType(s)
	default:
		out = PipelineConfigOriginTypeGoCD
	}
	return &out
}

type PipelineConfigOrigin struct {
	Type  *PipelineConfigOriginType `json:"type"`
	ID    *string                   `json:"id"`
	Links *HALLinks                 `json:"_links,omitempty"`
}

func (p *PipelineConfigOrigin) Equal(other *PipelineConfigOrigin) bool {
	if p == nil || other == nil {
		return p == other
	}

	typeIsEqual := cmp.PtrEqual(p.Type, other.Type)
	idIsEqual := cmp.PtrEqual(p.ID, other.ID)

	return typeIsEqual && idIsEqual
}

type PipelineConfigParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (p PipelineConfigParameter) Equal(o PipelineConfigParameter) bool {
	nameIsEqual := p.Name == o.Name
	valueIsEqual := p.Value == o.Value
	return nameIsEqual && valueIsEqual
}

type EnvironmentVariable struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	EncryptedValue string `json:"encrypted_value"`
	Secure         bool   `json:"secure"`
}

func (e EnvironmentVariable) Equal(other EnvironmentVariable) bool {
	return e.Name == other.Name
}

type PipelineConfigMaterialType string

const (
	PipelineConfigMaterialTypeGit        PipelineConfigMaterialType = "git"
	PipelineConfigMaterialTypeSvn        PipelineConfigMaterialType = "svn"
	PipelineConfigMaterialTypeHg         PipelineConfigMaterialType = "hg"
	PipelineConfigMaterialTypeP4         PipelineConfigMaterialType = "p4"
	PipelineConfigMaterialTypeTfs        PipelineConfigMaterialType = "tfs"
	PipelineConfigMaterialTypeDependency PipelineConfigMaterialType = "dependency"
	PipelineConfigMaterialTypePackage    PipelineConfigMaterialType = "package"
	PipelineConfigMaterialTypePlugin     PipelineConfigMaterialType = "plugin"
)

func PipelineConfigMaterialTypeFromString(s string) PipelineConfigMaterialType {
	switch PipelineConfigMaterialType(s) {
	case PipelineConfigMaterialTypeGit, PipelineConfigMaterialTypeSvn, PipelineConfigMaterialTypeHg, PipelineConfigMaterialTypeP4, PipelineConfigMaterialTypeTfs, PipelineConfigMaterialTypeDependency, PipelineConfigMaterialTypePackage, PipelineConfigMaterialTypePlugin:
		return PipelineConfigMaterialType(s)
	default:
		return PipelineConfigMaterialTypeGit
	}
}

type PipelineConfigMaterialAttributes interface {
	getMaterialAttrID() string
	Equal(PipelineConfigMaterialAttributes) bool
}

type PipelineConfigMaterialFilter struct {
	Ignore   []string `json:"ignore"`
	Includes []string `json:"includes"`
}

func (p *PipelineConfigMaterialFilter) Equal(other *PipelineConfigMaterialFilter) bool {
	if p == nil || other == nil {
		return p == other
	}

	ingnoreIsEqual := func() bool {
		pIgnoreCopy := slices.Clone(p.Ignore)
		otherIgnoreCopy := slices.Clone(other.Ignore)

		slices.Sort(pIgnoreCopy)
		slices.Sort(otherIgnoreCopy)

		return slices.Equal(pIgnoreCopy, otherIgnoreCopy)
	}()

	includesIsEqual := func() bool {
		pIncludesCopy := slices.Clone(p.Includes)
		otherIncludesCopy := slices.Clone(other.Includes)

		slices.Sort(pIncludesCopy)
		slices.Sort(otherIncludesCopy)

		return slices.Equal(pIncludesCopy, otherIncludesCopy)
	}()

	return includesIsEqual && ingnoreIsEqual
}

type PipelineConfigMaterialAttributesGit struct {
	Name            *string                       `json:"name"`
	URL             *string                       `json:"url"`
	Username        *string                       `json:"username,omitempty"`
	Password        *string                       `json:"password,omitempty"`
	Branch          *string                       `json:"branch"`
	Destination     *string                       `json:"destination"`
	AutoUpdate      bool                          `json:"auto_update"`
	Filter          *PipelineConfigMaterialFilter `json:"filter"`
	InvertFilter    bool                          `json:"invert_filter"`
	SubmoduleFolder *string                       `json:"submodule_folder"`
	ShallowClone    bool                          `json:"shallow_clone"`
}

func (p *PipelineConfigMaterialAttributesGit) Equal(other PipelineConfigMaterialAttributes) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}
	o, ok := other.(*PipelineConfigMaterialAttributesGit)
	if !ok {
		return false
	}

	nameIsEqual := cmp.PtrEqual(p.Name, o.Name)
	urlIsEqual := cmp.PtrEqual(p.URL, o.URL)
	usernameIsEqual := cmp.PtrEqual(p.Username, o.Username)
	branchIsEqual := cmp.PtrEqual(p.Branch, o.Branch)
	destinationIsEqual := cmp.PtrEqual(p.Destination, o.Destination)
	autoUpdateIsEqual := p.AutoUpdate == o.AutoUpdate
	filterAreEqual := p.Filter.Equal(o.Filter)
	invertFilterIsEqual := p.InvertFilter == o.InvertFilter
	submoduleFolderIsEqual := cmp.PtrEqual(p.SubmoduleFolder, o.SubmoduleFolder)
	shallowCloneIsEqual := p.ShallowClone == o.ShallowClone

	return nameIsEqual &&
		urlIsEqual &&
		usernameIsEqual &&
		branchIsEqual &&
		destinationIsEqual &&
		autoUpdateIsEqual &&
		filterAreEqual &&
		invertFilterIsEqual &&
		submoduleFolderIsEqual &&
		shallowCloneIsEqual
}

func (p *PipelineConfigMaterialAttributesGit) getMaterialAttrID() string {
	if p == nil {
		return nilStr
	}
	return *p.URL + *p.Branch
}

type PipelineConfigMaterialAttributesSvn struct {
	Name              *string                       `json:"name,omitempty"`
	URL               *string                       `json:"url,omitempty"`
	Username          *string                       `json:"username,omitempty"`
	Password          *string                       `json:"password,omitempty"`
	EncryptedPassword *string                       `json:"encrypted_password,omitempty"`
	Destination       *string                       `json:"destination,omitempty"`
	Filter            *PipelineConfigMaterialFilter `json:"filter,omitempty"`
	InvertFilter      bool                          `json:"invert_filter,omitempty"`
	AutoUpdate        bool                          `json:"auto_update,omitempty"`
	CheckExternals    bool                          `json:"check_externals,omitempty"`
}

func (p *PipelineConfigMaterialAttributesSvn) Equal(other PipelineConfigMaterialAttributes) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesSvn)
	if !ok {
		return false
	}

	return *p.Name == *o.Name &&
		*p.URL == *o.URL &&
		*p.Username == *o.Username &&
		*p.Destination == *o.Destination &&
		p.Filter.Equal(o.Filter) &&
		p.InvertFilter == o.InvertFilter &&
		p.AutoUpdate == o.AutoUpdate &&
		p.CheckExternals == o.CheckExternals
}

func (p *PipelineConfigMaterialAttributesSvn) getMaterialAttrID() string {
	if p == nil {
		return nilStr
	}

	return *p.URL
}

type PipelineConfigMaterialAttributesHg struct {
	Name              *string                       `json:"name,omitempty"`
	URL               *string                       `json:"url,omitempty"`
	Username          *string                       `json:"username,omitempty"`
	Password          *string                       `json:"password,omitempty"`
	EncryptedPassword *string                       `json:"encrypted_password,omitempty"`
	Branch            *string                       `json:"branch,omitempty"`
	Destination       *string                       `json:"destination,omitempty"`
	Filter            *PipelineConfigMaterialFilter `json:"filter,omitempty"`
	InvertFilter      bool                          `json:"invert_filter,omitempty"`
	AutoUpdate        bool                          `json:"auto_update,omitempty"`
}

func (p *PipelineConfigMaterialAttributesHg) Equal(other PipelineConfigMaterialAttributes) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesHg)
	if !ok {
		return false
	}

	nameIsEqual := *p.Name == *o.Name
	urlIsEqual := *p.URL == *o.URL
	usernameIsEqual := *p.Username == *o.Username
	branchIsEqual := *p.Branch == *o.Branch
	destinationIsEqual := *p.Destination == *o.Destination
	filterAreEqual := p.Filter.Equal(o.Filter)
	InvertFilterIsEqual := p.InvertFilter == o.InvertFilter
	autoUpdateIsEqual := p.AutoUpdate == o.AutoUpdate

	return nameIsEqual &&
		urlIsEqual &&
		usernameIsEqual &&
		branchIsEqual &&
		destinationIsEqual &&
		filterAreEqual &&
		InvertFilterIsEqual &&
		autoUpdateIsEqual
}

func (p *PipelineConfigMaterialAttributesHg) getMaterialAttrID() string {
	if p == nil {
		return nilStr
	}

	return *p.URL + *p.Branch
}

type PipelineConfigMaterialAttributesP4 struct {
	Name              *string                       `json:"name,omitempty"`
	Port              *string                       `json:"port,omitempty"`
	UseTickets        bool                          `json:"use_tickets,omitempty"`
	View              *string                       `json:"view,omitempty"`
	Username          *string                       `json:"username,omitempty"`
	Password          *string                       `json:"password,omitempty"`
	EncryptedPassword *string                       `json:"encrypted_password,omitempty"`
	Destination       *string                       `json:"destination,omitempty"`
	Filter            *PipelineConfigMaterialFilter `json:"filter,omitempty"`
	InvertFilter      bool                          `json:"invert_filter,omitempty"`
	AutoUpdate        bool                          `json:"auto_update,omitempty"`
}

func (p *PipelineConfigMaterialAttributesP4) Equal(other PipelineConfigMaterialAttributes) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesP4)
	if !ok {
		return false
	}

	nameIsEqual := *p.Name == *o.Name
	portIsEqual := *p.Port == *o.Port
	UseTicketsIsEqual := p.UseTickets == o.UseTickets
	viewIsEqual := *p.View == *o.View
	usernameIsEqual := *p.Username == *o.Username
	destinationIsEqual := *p.Destination == *o.Destination
	filterIsEqual := p.Filter.Equal(o.Filter)
	invertFilterIsEqual := p.InvertFilter == o.InvertFilter
	autoUpdateIsEqual := p.AutoUpdate == o.AutoUpdate

	return nameIsEqual &&
		portIsEqual &&
		UseTicketsIsEqual &&
		viewIsEqual &&
		usernameIsEqual &&
		destinationIsEqual &&
		filterIsEqual &&
		invertFilterIsEqual &&
		autoUpdateIsEqual
}

func (p *PipelineConfigMaterialAttributesP4) getMaterialAttrID() string {
	if p == nil {
		return nilStr
	}

	return *p.Port
}

type PipelineConfigMaterialAttributesTfs struct {
	Name              *string                       `json:"name,omitempty"`
	URL               *string                       `json:"url,omitempty"`
	ProjectPath       *string                       `json:"project_path,omitempty"`
	Domain            *string                       `json:"domain,omitempty"`
	Username          *string                       `json:"username,omitempty"`
	Password          *string                       `json:"password,omitempty"`
	EncryptedPassword *string                       `json:"encrypted_password,omitempty"`
	Destination       *string                       `json:"destination,omitempty"`
	AutoUpdate        bool                          `json:"auto_update,omitempty"`
	Filter            *PipelineConfigMaterialFilter `json:"filter,omitempty"`
	InvertFilter      bool                          `json:"invert_filter,omitempty"`
}

func (p *PipelineConfigMaterialAttributesTfs) Equal(other PipelineConfigMaterialAttributes) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesTfs)
	if !ok {
		return false
	}

	nameIsEqual := *p.Name == *o.Name
	urlIsEqual := *p.URL == *o.URL
	projectPathIsEqual := *p.ProjectPath == *o.ProjectPath
	domainIsEqual := *p.Domain == *o.Domain
	usernameIsEqual := *p.Username == *o.Username
	destinationIsEqual := *p.Destination == *o.Destination
	autoUpdateIsEqual := p.AutoUpdate == o.AutoUpdate
	filterAreEqual := p.Filter.Equal(o.Filter)
	invertFilterIsEqual := p.InvertFilter == o.InvertFilter

	return nameIsEqual &&
		urlIsEqual &&
		projectPathIsEqual &&
		domainIsEqual &&
		usernameIsEqual &&
		destinationIsEqual &&
		autoUpdateIsEqual &&
		filterAreEqual &&
		invertFilterIsEqual
}

func (p *PipelineConfigMaterialAttributesTfs) getMaterialAttrID() string {
	if p == nil {
		return nilStr
	}
	return *p.URL + *p.ProjectPath
}

type PipelineConfigMaterialAttributesDependency struct {
	Name                *string `json:"name,omitempty"`
	Pipeline            *string `json:"pipeline,omitempty"`
	Stage               *string `json:"stage,omitempty"`
	AutoUpdate          bool    `json:"auto_update,omitempty"`
	IgnoreForScheduling bool    `json:"ignore_for_scheduling,omitempty"`
}

func (p *PipelineConfigMaterialAttributesDependency) Equal(other PipelineConfigMaterialAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesDependency)
	if !ok {
		return false
	}

	nameIsEqual := *p.Name == *o.Name
	pipelineIsEqual := *p.Pipeline == *o.Pipeline
	stageIsEqual := *p.Stage == *o.Stage
	autoUpdateIsEqual := p.AutoUpdate == o.AutoUpdate
	ignoreForSchedulingIsEqual := p.IgnoreForScheduling == o.IgnoreForScheduling

	return nameIsEqual &&
		pipelineIsEqual &&
		stageIsEqual &&
		autoUpdateIsEqual &&
		ignoreForSchedulingIsEqual
}

func (p *PipelineConfigMaterialAttributesDependency) getMaterialAttrID() string {
	if p == nil {
		return nilStr
	}
	return *p.Pipeline + *p.Stage
}

type PipelineConfigMaterialAttributesPackage struct {
	Ref *string `json:"ref,omitempty"`
}

func (p *PipelineConfigMaterialAttributesPackage) Equal(other PipelineConfigMaterialAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesPackage)
	if !ok {
		return false
	}

	return *p.Ref == *o.Ref
}

func (p *PipelineConfigMaterialAttributesPackage) getMaterialAttrID() string {
	return *p.Ref
}

type PipelineConfigMaterialAttributesPlugin struct {
	Ref          *string                       `json:"ref,omitempty"`
	Destination  *string                       `json:"destination,omitempty"`
	Filter       *PipelineConfigMaterialFilter `json:"filter,omitempty"`
	InvertFilter bool                          `json:"invert_filter,omitempty"`
}

func (p *PipelineConfigMaterialAttributesPlugin) Equal(other PipelineConfigMaterialAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigMaterialAttributesPlugin)
	if !ok {
		return false
	}

	return *p.Ref == *o.Ref &&
		*p.Destination == *o.Destination &&
		p.Filter.Equal(o.Filter) &&
		p.InvertFilter == o.InvertFilter
}

func (p *PipelineConfigMaterialAttributesPlugin) getMaterialAttrID() string {
	return *p.Ref
}

type PipelineConfigMaterial struct { //nolint:recvcheck
	Type       PipelineConfigMaterialType       `json:"type"`
	Attributes PipelineConfigMaterialAttributes `json:"attributes"`
}

func (m PipelineConfigMaterial) Equal(o PipelineConfigMaterial) bool {
	return m.Type == o.Type && m.Attributes.Equal(o.Attributes)
}

func (m PipelineConfigMaterial) GetID() string {
	return string(m.Type) + ":" + m.Attributes.getMaterialAttrID()
}

func (m *PipelineConfigMaterial) UnmarshalJSON(b []byte) error { //nolint:gocyclo
	var aux struct {
		Type       PipelineConfigMaterialType `json:"type"`
		Attributes json.RawMessage            `json:"attributes"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	m.Type = aux.Type

	switch aux.Type {
	case PipelineConfigMaterialTypeGit:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesGit](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypeSvn:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesSvn](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypeHg:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesHg](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypeP4:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesP4](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypeTfs:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesTfs](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypeDependency:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesDependency](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypePackage:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesPackage](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	case PipelineConfigMaterialTypePlugin:
		a, err := unmarshallAttrs[*PipelineConfigMaterialAttributesPlugin](aux.Attributes)
		if err != nil {
			return err
		}
		m.Attributes = a
	}
	return nil
}

func unmarshallAttrs[T PipelineConfigMaterialAttributes](b []byte) (T, error) {
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return v, err
	}
	return v, nil
}

func PipelineConfigMaterialLessFunc(a, b PipelineConfigMaterial) bool {
	return a.Attributes.getMaterialAttrID() < b.Attributes.getMaterialAttrID()
}

type PipelineConfigApprovalAuthorization struct {
	Users []string `json:"users"`
	Roles []string `json:"roles"`
}

func (p PipelineConfigApprovalAuthorization) Equal(other PipelineConfigApprovalAuthorization) bool {
	usersAreEqual := func() bool {
		pUsersClone := slices.Clone(p.Users)
		otherUsersClone := slices.Clone(other.Users)

		slices.Sort(pUsersClone)
		slices.Sort(otherUsersClone)

		return slices.Equal(pUsersClone, otherUsersClone)
	}()

	rolesAreEqual := func() bool {
		pRolesClone := slices.Clone(p.Roles)
		otherRolesClone := slices.Clone(other.Roles)

		slices.Sort(pRolesClone)
		slices.Sort(otherRolesClone)

		return slices.Equal(pRolesClone, otherRolesClone)
	}()

	return usersAreEqual && rolesAreEqual
}

type PipelineConfigApprovalType string

const (
	PipelineConfigApprovalTypeSuccess PipelineConfigApprovalType = "success"
	PipelineConfigApprovalTypeManual  PipelineConfigApprovalType = "manual"
)

func (pc PipelineConfigApprovalType) String() string {
	return string(pc)
}

func PipelineConfigApprovalTypeFromString(s string) PipelineConfigApprovalType {
	switch PipelineConfigApprovalType(s) {
	case PipelineConfigApprovalTypeSuccess, PipelineConfigApprovalTypeManual:
		return PipelineConfigApprovalType(s)
	default:
		return PipelineConfigApprovalTypeSuccess
	}
}

type PipelineConfigApproval struct {
	Type               PipelineConfigApprovalType          `json:"type"`
	AllowOnlyOnSuccess bool                                `json:"allow_only_on_success,omitempty"`
	Authorization      PipelineConfigApprovalAuthorization `json:"authorization"`
}

func (p PipelineConfigApproval) Equal(other PipelineConfigApproval) bool {
	typeIsEqual := p.Type == other.Type
	allowOnlyOnSuccessIsEqual := p.AllowOnlyOnSuccess == other.AllowOnlyOnSuccess
	authorizationIsEqual := p.Authorization.Equal(other.Authorization)

	return typeIsEqual && allowOnlyOnSuccessIsEqual && authorizationIsEqual
}

type RunIfType string

const (
	RunIfTypePassed RunIfType = "passed"
	RunIfTypeFailed RunIfType = "failed"
	RunIfTypeAny    RunIfType = "any"
)

func (t RunIfType) String() string {
	return string(t)
}

func RunIfTypeFromString(s string) RunIfType {
	switch RunIfType(s) {
	case RunIfTypePassed, RunIfTypeFailed, RunIfTypeAny:
		return RunIfType(s)
	default:
		return RunIfTypePassed
	}
}

type PipelineConfigStageJobsTaskType string

const (
	PipelineConfigStageJobsTaskTypeExec          PipelineConfigStageJobsTaskType = "exec"
	PipelineConfigStageJobsTaskTypeAnt           PipelineConfigStageJobsTaskType = "ant"
	PipelineConfigStageJobsTaskTypeNant          PipelineConfigStageJobsTaskType = "nant"
	PipelineConfigStageJobsTaskTypeRake          PipelineConfigStageJobsTaskType = "rake"
	PipelineConfigStageJobsTaskTypeFetch         PipelineConfigStageJobsTaskType = "fetch"
	PipelineConfigStageJobsTaskTypePluggableTask PipelineConfigStageJobsTaskType = "pluggable_task"
)

func (tt PipelineConfigStageJobsTaskType) String() string {
	return string(tt)
}

func PipelineConfigStageJobsTaskTypeFromString(s string) PipelineConfigStageJobsTaskType {
	switch PipelineConfigStageJobsTaskType(s) {
	case PipelineConfigStageJobsTaskTypeExec, PipelineConfigStageJobsTaskTypeAnt, PipelineConfigStageJobsTaskTypeNant,
		PipelineConfigStageJobsTaskTypeRake, PipelineConfigStageJobsTaskTypeFetch,
		PipelineConfigStageJobsTaskTypePluggableTask:

		return PipelineConfigStageJobsTaskType(s)
	default:
		return PipelineConfigStageJobsTaskTypeExec
	}
}

type PipelineConfigStageJobsTaskAttributes interface {
	isPipelineConfigStageJobsTaskAttributes()
	Equal(other PipelineConfigStageJobsTaskAttributes) bool
}
type PipelineConfigStageJobsTaskAttributesExec struct {
	RunIf            []RunIfType                  `json:"run_if"`
	Command          string                       `json:"command"`
	Arguments        []string                     `json:"arguments,omitempty"`
	WorkingDirectory *string                      `json:"working_directory"`
	OnCancel         *PipelineConfigStageJobsTask `json:"on_cancel,omitempty"`
}

func (*PipelineConfigStageJobsTaskAttributesExec) isPipelineConfigStageJobsTaskAttributes() {}
func (p *PipelineConfigStageJobsTaskAttributesExec) Equal(other PipelineConfigStageJobsTaskAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}
	o, ok := other.(*PipelineConfigStageJobsTaskAttributesExec)
	if !ok {
		return false
	}

	runIfIsEqual := func() bool {
		pRunIfClone := append(make([]RunIfType, 0, len(p.RunIf)), p.RunIf...)
		oRunIfClone := append(make([]RunIfType, 0, len(o.RunIf)), o.RunIf...)

		slices.Sort(pRunIfClone)
		slices.Sort(oRunIfClone)

		return slices.Equal(pRunIfClone, oRunIfClone)
	}()

	commandIsEqual := p.Command == o.Command

	argumentsAreEqual := slices.Equal(p.Arguments, o.Arguments)

	workingDirectoryIsEqual := cmp.PtrEqual(p.WorkingDirectory, o.WorkingDirectory)
	onCancelIsEqual := p.OnCancel.Equal(o.OnCancel)
	return runIfIsEqual && commandIsEqual && argumentsAreEqual && workingDirectoryIsEqual && onCancelIsEqual
}

type PipelineConfigStageJobsTaskAttributesAnt struct {
	RunIf            []RunIfType                  `json:"run_if"`
	BuildFile        string                       `json:"build_file"`
	Target           string                       `json:"target"`
	WorkingDirectory string                       `json:"working_directory"`
	OnCancel         *PipelineConfigStageJobsTask `json:"on_cancel,omitempty"`
}

func (*PipelineConfigStageJobsTaskAttributesAnt) isPipelineConfigStageJobsTaskAttributes() {}
func (p *PipelineConfigStageJobsTaskAttributesAnt) Equal(other PipelineConfigStageJobsTaskAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigStageJobsTaskAttributesAnt)
	if !ok {
		return false
	}

	runIfIsEqual := func() bool {
		pRunIfClone := append(make([]RunIfType, 0, len(p.RunIf)), p.RunIf...)
		oRunIfClone := append(make([]RunIfType, 0, len(o.RunIf)), o.RunIf...)

		slices.Sort(pRunIfClone)
		slices.Sort(oRunIfClone)

		return slices.Equal(pRunIfClone, oRunIfClone)
	}()

	buildFileIsEqual := p.BuildFile == o.BuildFile
	targetIsEqual := p.Target == o.Target
	workingDirectoryIsEqual := p.WorkingDirectory == o.WorkingDirectory
	onCancelIsEqual := p.OnCancel.Equal(o.OnCancel)

	return runIfIsEqual && buildFileIsEqual && targetIsEqual && workingDirectoryIsEqual && onCancelIsEqual
}

type PipelineConfigStageJobsTaskAttributesNant struct {
	RunIf            []RunIfType                  `json:"run_if"`
	BuildFile        string                       `json:"build_file"`
	Target           string                       `json:"target"`
	NantPath         string                       `json:"nant_path"`
	WorkingDirectory string                       `json:"working_directory"`
	OnCancel         *PipelineConfigStageJobsTask `json:"on_cancel,omitempty"`
}

func (*PipelineConfigStageJobsTaskAttributesNant) isPipelineConfigStageJobsTaskAttributes() {}
func (p *PipelineConfigStageJobsTaskAttributesNant) Equal(other PipelineConfigStageJobsTaskAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigStageJobsTaskAttributesNant)
	if !ok {
		return false
	}

	runIfIsEqual := func() bool {
		pRunIfClone := append(make([]RunIfType, 0, len(p.RunIf)), p.RunIf...)
		oRunIfClone := append(make([]RunIfType, 0, len(o.RunIf)), o.RunIf...)

		slices.Sort(pRunIfClone)
		slices.Sort(oRunIfClone)

		return slices.Equal(pRunIfClone, oRunIfClone)
	}()

	buildFileIsEqual := p.BuildFile == o.BuildFile
	targetIsEqual := p.Target == o.Target
	nantPathIsEqual := p.NantPath == o.NantPath
	workingDirectoryIsEqual := p.WorkingDirectory == o.WorkingDirectory
	onCancelIsEqual := p.OnCancel.Equal(o.OnCancel)

	return runIfIsEqual && buildFileIsEqual && targetIsEqual && nantPathIsEqual && workingDirectoryIsEqual && onCancelIsEqual
}

type PipelineConfigStageJobsTaskAttributesRake struct {
	RunIf            []RunIfType                  `json:"run_if"`
	BuildFile        string                       `json:"build_file"`
	Target           string                       `json:"target"`
	WorkingDirectory string                       `json:"working_directory"`
	OnCancel         *PipelineConfigStageJobsTask `json:"on_cancel,omitempty"`
}

func (*PipelineConfigStageJobsTaskAttributesRake) isPipelineConfigStageJobsTaskAttributes() {}
func (p *PipelineConfigStageJobsTaskAttributesRake) Equal(other PipelineConfigStageJobsTaskAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigStageJobsTaskAttributesRake)
	if !ok {
		return false
	}

	runIfIsEqual := func() bool {
		pRunIfClone := append(make([]RunIfType, 0, len(p.RunIf)), p.RunIf...)
		oRunIfClone := append(make([]RunIfType, 0, len(o.RunIf)), o.RunIf...)

		slices.Sort(pRunIfClone)
		slices.Sort(oRunIfClone)

		return slices.Equal(pRunIfClone, oRunIfClone)
	}()

	buildFileIsEqual := p.BuildFile == o.BuildFile
	targetIsEqual := p.Target == o.Target
	workingDirectoryIsEqual := p.WorkingDirectory == o.WorkingDirectory
	onCancelIsEqual := p.OnCancel.Equal(o.OnCancel)

	return runIfIsEqual && buildFileIsEqual && targetIsEqual && workingDirectoryIsEqual && onCancelIsEqual
}

type PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType string

const (
	PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeGoCD     PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType = "gocd"
	PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeExternal PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType = "external"
)

func (t PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType) String() string {
	return string(t)
}

func PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeFromString(s string) PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType {
	switch PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType(s) {
	case PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeGoCD, PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeExternal:
		return PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType(s)
	default:
		return PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeGoCD
	}
}

type PipelineConfigStageJobsTaskAttributesFetch struct {
	ArtifactOrigin PipelineConfigStageJobsTaskAttributesFetchArtifactOriginType `json:"artifact_origin"`
	RunIf          []RunIfType                                                  `json:"run_if"`
	Pipeline       string                                                       `json:"pipeline"`
	Stage          string                                                       `json:"stage"`
	Job            string                                                       `json:"job"`
	Source         string                                                       `json:"source"`
	IsSourceAFile  bool                                                         `json:"is_source_a_file"`
	Destination    string                                                       `json:"destination"`
	OnCancel       *PipelineConfigStageJobsTask                                 `json:"on_cancel,omitempty"`
	ArtifactID     string                                                       `json:"artifact_id"`
	Configuration  []ConfigProperty                                             `json:"configuration"`
}

func (*PipelineConfigStageJobsTaskAttributesFetch) isPipelineConfigStageJobsTaskAttributes() {}
func (p *PipelineConfigStageJobsTaskAttributesFetch) Equal(other PipelineConfigStageJobsTaskAttributes) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigStageJobsTaskAttributesFetch)
	if !ok {
		return false
	}

	artifactOriginIsEqual := p.ArtifactID == o.ArtifactID
	runIfIsEqual := func() bool {
		pRunIfClone := append(make([]RunIfType, 0, len(p.RunIf)), p.RunIf...)
		oRunIfClone := append(make([]RunIfType, 0, len(o.RunIf)), o.RunIf...)

		slices.Sort(pRunIfClone)
		slices.Sort(oRunIfClone)

		return slices.Equal(pRunIfClone, oRunIfClone)
	}()
	pipelineIsEqual := p.Pipeline == o.Pipeline
	stageIsEqual := p.Stage == o.Stage
	jobIsEqual := p.Job == o.Job
	sourceIsEqual := p.Source == o.Source
	isSourceAFileIsEqual := p.IsSourceAFile == o.IsSourceAFile
	destinationIsEqual := p.Destination == o.Destination
	onCancelIsEqual := p.OnCancel.Equal(o.OnCancel)
	artifactIDIsEqual := p.ArtifactID == o.ArtifactID
	configurationIsEqual := cmp.SlicesEqualUnordered(p.Configuration, o.Configuration, func(v ConfigProperty) string {
		return v.Key
	})

	return artifactOriginIsEqual && runIfIsEqual && pipelineIsEqual && stageIsEqual &&
		jobIsEqual && sourceIsEqual && isSourceAFileIsEqual && destinationIsEqual && onCancelIsEqual &&
		artifactIDIsEqual && configurationIsEqual
}

type PipelineConfigStageJobsTaskAttributesPluggableTaskPluginConfiguration struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

func (p PipelineConfigStageJobsTaskAttributesPluggableTaskPluginConfiguration) Equal(c PipelineConfigStageJobsTaskAttributesPluggableTaskPluginConfiguration) bool {
	return p.ID == c.ID && p.Version == c.Version
}

type PipelineConfigStageJobsTaskAttributesPluggable struct {
	RunIf               []RunIfType                                                           `json:"run_if"`
	PluginConfiguration PipelineConfigStageJobsTaskAttributesPluggableTaskPluginConfiguration `json:"plugin_configuration"`
	Configuration       []ConfigProperty                                                      `json:"configuration"`
	OnCancel            *PipelineConfigStageJobsTask                                          `json:"on_cancel,omitempty"`
}

func (*PipelineConfigStageJobsTaskAttributesPluggable) isPipelineConfigStageJobsTaskAttributes() {}
func (p *PipelineConfigStageJobsTaskAttributesPluggable) Equal(other PipelineConfigStageJobsTaskAttributes) bool {
	if p == nil || other == nil {
		return p == other
	}

	o, ok := other.(*PipelineConfigStageJobsTaskAttributesPluggable)
	if !ok {
		return false
	}

	runIfIsEqual := func() bool {
		pRunIfClone := append(make([]RunIfType, 0, len(p.RunIf)), p.RunIf...)
		oRunIfClone := append(make([]RunIfType, 0, len(o.RunIf)), o.RunIf...)

		slices.Sort(pRunIfClone)
		slices.Sort(oRunIfClone)

		return slices.Equal(pRunIfClone, oRunIfClone)
	}()
	puglinConfigurationIsEqual := p.PluginConfiguration.Equal(o.PluginConfiguration)
	configurationIsEqual := cmp.SlicesEqualUnordered(p.Configuration, o.Configuration, func(c ConfigProperty) string {
		return c.Key
	})
	onCancelIsEqual := p.OnCancel.Equal(o.OnCancel)

	return runIfIsEqual && puglinConfigurationIsEqual && configurationIsEqual && onCancelIsEqual
}

type PipelineConfigStageJobsTask struct { //nolint:recvcheck
	Type       PipelineConfigStageJobsTaskType       `json:"type"`
	Attributes PipelineConfigStageJobsTaskAttributes `json:"attributes"`
}

func (t *PipelineConfigStageJobsTask) Equal(o *PipelineConfigStageJobsTask) bool {
	if t == nil || o == nil {
		return t == o
	}
	return t.Type == o.Type && t.Attributes.Equal(o.Attributes)
}

func (t *PipelineConfigStageJobsTask) UnmarshalJSON(b []byte) error { //nolint:gocyclo
	var aux struct {
		Type       PipelineConfigStageJobsTaskType `json:"type"`
		Attributes json.RawMessage                 `json:"attributes"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	t.Type = aux.Type

	switch t.Type {
	case PipelineConfigStageJobsTaskTypeExec:
		var v PipelineConfigStageJobsTaskAttributesExec
		err := json.Unmarshal(aux.Attributes, &v)
		if err != nil {
			return err
		}
		t.Attributes = &v
	case PipelineConfigStageJobsTaskTypeAnt:
		var v PipelineConfigStageJobsTaskAttributesAnt
		err := json.Unmarshal(aux.Attributes, &v)
		if err != nil {
			return err
		}
		t.Attributes = &v
	case PipelineConfigStageJobsTaskTypeNant:
		var v PipelineConfigStageJobsTaskAttributesNant
		err := json.Unmarshal(aux.Attributes, &v)
		if err != nil {
			return err
		}
		t.Attributes = &v
	case PipelineConfigStageJobsTaskTypeRake:
		var v PipelineConfigStageJobsTaskAttributesRake
		err := json.Unmarshal(aux.Attributes, &v)
		if err != nil {
			return err
		}
		t.Attributes = &v
	case PipelineConfigStageJobsTaskTypeFetch:
		var v PipelineConfigStageJobsTaskAttributesFetch
		err := json.Unmarshal(aux.Attributes, &v)
		if err != nil {
			return err
		}
		t.Attributes = &v
	case PipelineConfigStageJobsTaskTypePluggableTask:
		var v PipelineConfigStageJobsTaskAttributesPluggable
		err := json.Unmarshal(aux.Attributes, &v)
		if err != nil {
			return err
		}
		t.Attributes = &v
	default:
		return fmt.Errorf("unknown task type: %s", t.Type)
	}
	return nil
}

type PipelineConfigStageJobsTab struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (p PipelineConfigStageJobsTab) Equal(other PipelineConfigStageJobsTab) bool {
	return p.Name == other.Name && p.Path == other.Path
}

type PipelineConfigStageJobsArtifactType string

const (
	PipelineConfigStageJobsArtifactTypeTest     = "test"
	PipelineConfigStageJobsArtifactTypeBuild    = "build"
	PipelineConfigStageJobsArtifactTypeExternal = "external"
)

func (receiver PipelineConfigStageJobsArtifactType) String() string {
	return string(receiver)
}

func PipelineConfigStageJobsArtifactTypeFromString(s string) PipelineConfigStageJobsArtifactType {
	switch PipelineConfigStageJobsArtifactType(s) {
	case PipelineConfigStageJobsArtifactTypeTest, PipelineConfigStageJobsArtifactTypeBuild, PipelineConfigStageJobsArtifactTypeExternal:
		return PipelineConfigStageJobsArtifactType(s)
	default:
		return PipelineConfigStageJobsArtifactTypeBuild
	}
}

type PipelineConfigStageJobsArtifact struct {
	Type          PipelineConfigStageJobsArtifactType `json:"type"`
	Source        *string                             `json:"source,omitempty"`
	Destination   *string                             `json:"destination,omitempty"`
	ArtifactID    *string                             `json:"artifact_id,omitempty"`
	StoreID       *string                             `json:"store_id,omitempty"`
	Configuration []ConfigProperty                    `json:"configuration,omitempty"`
}

func (p PipelineConfigStageJobsArtifact) Equal(other PipelineConfigStageJobsArtifact) bool {
	typeIsEqual := p.Type == other.Type
	sourceIsEqual := *p.Source == *other.Source
	destinationIsEqual := *p.Destination == *other.Destination
	artifactIDIsEqual := *p.ArtifactID == *other.ArtifactID
	storeIDIsEqual := *p.StoreID == *other.StoreID
	configurationIsEqual := cmp.SlicesEqualUnordered(p.Configuration, other.Configuration, func(c ConfigProperty) string {
		return c.Key
	})

	return typeIsEqual && sourceIsEqual && destinationIsEqual && artifactIDIsEqual && storeIDIsEqual && configurationIsEqual
}

type PipelineConfigStageJobs struct {
	Name                 string                            `json:"name"`
	RunInstanceCount     *intstr.IntOrString               `json:"run_instance_count"`
	Timeout              intstr.IntOrString                `json:"timeout"`
	EnvironmentVariables []EnvironmentVariable             `json:"environment_variables"`
	Resources            []string                          `json:"resources"`
	Tasks                []PipelineConfigStageJobsTask     `json:"tasks"`
	Tabs                 []PipelineConfigStageJobsTab      `json:"tabs,omitempty"`
	Artifacts            []PipelineConfigStageJobsArtifact `json:"artifacts"`
	ElasticProfileID     string                            `json:"elastic_profile_id,omitempty"`
}

func (j PipelineConfigStageJobs) Equal(o PipelineConfigStageJobs) bool { //nolint:gocyclo
	nameIsEqual := j.Name == o.Name
	runInstanceCountIsEqual := func() bool {
		if j.RunInstanceCount == nil || o.RunInstanceCount == nil {
			return j.RunInstanceCount == o.RunInstanceCount
		}
		return j.RunInstanceCount.IntVal == o.RunInstanceCount.IntVal &&
			j.RunInstanceCount.StrVal == o.RunInstanceCount.StrVal &&
			j.RunInstanceCount.Type == o.RunInstanceCount.Type
	}()
	timeoutIsEqual := func() bool {
		return j.Timeout.IntVal == o.Timeout.IntVal &&
			j.Timeout.StrVal == o.Timeout.StrVal &&
			j.Timeout.Type == o.Timeout.Type
	}()
	evironmentVariablesAreEqual := cmp.SlicesEqualUnordered(
		j.EnvironmentVariables,
		o.EnvironmentVariables,
		func(e EnvironmentVariable) string {
			return e.Name + e.Value
		},
	)
	resourcesAreEqual := func() bool {
		jResourceClone := append(make([]string, 0, len(j.Resources)), j.Resources...)
		oResourceClone := append(make([]string, 0, len(o.Resources)), o.Resources...)

		return slices.Equal(jResourceClone, oResourceClone)
	}()
	tasksAreEqual := slices.EqualFunc(j.Tasks, o.Tasks, func(a, b PipelineConfigStageJobsTask) bool {
		return a.Equal(&b)
	})
	tabsAreEqual := cmp.SlicesEqualUnordered(j.Tabs, o.Tabs, func(t PipelineConfigStageJobsTab) string {
		return t.Name
	})
	artifactsAreEqual := cmp.SlicesEqualUnordered(j.Artifacts, o.Artifacts, func(p PipelineConfigStageJobsArtifact) string {
		return *p.ArtifactID
	})
	elasticProfileIDIsEqual := j.ElasticProfileID == o.ElasticProfileID

	return nameIsEqual &&
		runInstanceCountIsEqual &&
		timeoutIsEqual &&
		evironmentVariablesAreEqual &&
		resourcesAreEqual &&
		tasksAreEqual &&
		tabsAreEqual &&
		artifactsAreEqual &&
		elasticProfileIDIsEqual
}

type PipelineConfigStage struct {
	Name                  string                    `json:"name"`
	FetchMaterials        bool                      `json:"fetch_materials"`
	CleanWorkingDirectory bool                      `json:"clean_working_directory"`
	NeverCleanupArtifacts bool                      `json:"never_cleanup_artifacts"`
	Approval              PipelineConfigApproval    `json:"approval"`
	EnvironmentVariables  []EnvironmentVariable     `json:"environment_variables"`
	Jobs                  []PipelineConfigStageJobs `json:"jobs"`
}

func (s PipelineConfigStage) Equal(o PipelineConfigStage) bool {
	nameIsEqual := s.Name == o.Name
	fetchMaterialsIsEqual := s.FetchMaterials == o.FetchMaterials
	cleanWorkDirectryIsEqual := s.CleanWorkingDirectory == o.CleanWorkingDirectory
	neverCleanupArtifactIsEqual := s.NeverCleanupArtifacts == o.NeverCleanupArtifacts
	approvalIsEqual := s.Approval.Equal(o.Approval)
	environmentVariablesAreEqual := cmp.SlicesEqualUnordered(
		s.EnvironmentVariables,
		o.EnvironmentVariables,
		func(e EnvironmentVariable) string { return e.Name })
	jobsAreEqual := cmp.SlicesEqualUnordered(
		s.Jobs,
		o.Jobs,
		func(j PipelineConfigStageJobs) string {
			return j.Name
		},
	)

	return nameIsEqual && fetchMaterialsIsEqual && cleanWorkDirectryIsEqual && neverCleanupArtifactIsEqual && approvalIsEqual && environmentVariablesAreEqual && jobsAreEqual
}

type PipelineConfigTrackingToolAttributes struct {
	URLPattern string `json:"url_pattern"`
	Regex      string `json:"regex"`
}

func (p PipelineConfigTrackingToolAttributes) Equal(other PipelineConfigTrackingToolAttributes) bool {
	return p.URLPattern == other.URLPattern && p.Regex == other.Regex
}

type PipelineConfigTrackingTool struct {
	Type       string                               `json:"type"`
	Attributes PipelineConfigTrackingToolAttributes `json:"attributes"`
}

func (p *PipelineConfigTrackingTool) Equal(other *PipelineConfigTrackingTool) bool {
	if p == nil || other == nil {
		return p == other
	}

	typeIsEqual := p.Type == other.Type
	attributesAreEqual := p.Attributes.Equal(other.Attributes)

	return typeIsEqual && attributesAreEqual
}

type PipelineConfigTimer struct {
	Spec          string `json:"spec"`
	OnlyOnChanges bool   `json:"only_on_changes"`
}

func (p *PipelineConfigTimer) Equal(other *PipelineConfigTimer) bool {
	if p == nil || other == nil {
		return p == other
	}

	return p.Spec == other.Spec && p.OnlyOnChanges == other.OnlyOnChanges
}

type PipelineConfig struct {
	Group                *string                     `json:"group,omitempty"`
	LabelTemplate        *string                     `json:"label_template,omitempty"`
	LockBehavior         *PipelineConfigLockBehavior `json:"lock_behavior,omitempty"`
	Name                 *string                     `json:"name,omitempty"`
	Template             *string                     `json:"template"`
	Origin               *PipelineConfigOrigin       `json:"origin,omitempty"`
	Parameters           []PipelineConfigParameter   `json:"parameters,omitempty"`
	EnvironmentVariables []EnvironmentVariable       `json:"environment_variables,omitempty"`
	Materials            []PipelineConfigMaterial    `json:"materials,omitempty"`
	Stages               []PipelineConfigStage       `json:"stages,omitempty"`
	TrackingTool         *PipelineConfigTrackingTool `json:"tracking_tool,omitempty"`
	Timer                *PipelineConfigTimer        `json:"timer,omitempty"`
	Links                *HALLinks                   `json:"_links,omitempty"`
}

func (p *PipelineConfig) Equal(other *PipelineConfig) bool { //nolint:gocyclo
	if p == nil || other == nil {
		return p == other
	}

	groupIsEqual := cmp.PtrEqual(p.Group, other.Group)
	labelTemplateIsEqual := cmp.PtrEqual(p.LabelTemplate, other.LabelTemplate)
	lockBehaviorIsEqual := cmp.PtrEqual(p.LockBehavior, other.LockBehavior)
	nameIsEqual := cmp.PtrEqual(p.Name, other.Name)
	templateIsEqual := cmp.PtrEqual(p.Template, other.Template)
	originIsEqual := p.Origin.Equal(other.Origin)
	parametersAreEqual := cmp.SlicesEqualUnordered(
		p.Parameters,
		other.Parameters,
		func(p PipelineConfigParameter) string {
			return p.Name + p.Value
		})
	environmentVariablesAreEqual := cmp.SlicesEqualUnordered(
		p.EnvironmentVariables,
		other.EnvironmentVariables,
		func(e EnvironmentVariable) string {
			return e.Name
		})
	materialsAreEqual := cmp.SlicesEqualUnordered(
		p.Materials,
		other.Materials,
		func(m PipelineConfigMaterial) string {
			return m.GetID()
		})
	stagesAreEqual := cmp.SlicesEqualUnordered(
		p.Stages,
		other.Stages,
		func(s PipelineConfigStage) string {
			return s.Name
		})
	trackingToolsAreEquals := p.TrackingTool.Equal(other.TrackingTool)
	timerIsEqual := p.Timer.Equal(other.Timer)

	return groupIsEqual &&
		labelTemplateIsEqual &&
		lockBehaviorIsEqual &&
		nameIsEqual &&
		templateIsEqual &&
		originIsEqual &&
		parametersAreEqual &&
		environmentVariablesAreEqual &&
		materialsAreEqual &&
		stagesAreEqual &&
		trackingToolsAreEquals &&
		timerIsEqual
}

type pipelineConfigCreateRequest struct {
	Group    *string         `json:"group"`
	Pipeline *PipelineConfig `json:"pipeline"`
}

type PipelineConfigsService interface {
	Get(ctx context.Context, name string) (*PipelineConfig, string, error)
	Create(ctx context.Context, body *PipelineConfig) (*PipelineConfig, string, error)
	Update(ctx context.Context, etag string, body *PipelineConfig) (*PipelineConfig, string, error)
	Delete(ctx context.Context, name string) error
}

func (c *client) PipelineConfigs() PipelineConfigsService {
	return &pipelineConfigsService{c: c}
}

type pipelineConfigsService struct {
	c *client
}

func (p *pipelineConfigsService) Get(ctx context.Context, name string) (*PipelineConfig, string, error) {
	path := fmt.Sprintf("%s/%s", pipelineConfigsServicePath, url.PathEscape(name))
	resp, err := p.c.do(ctx, http.MethodGet, path, acceptPipelineConfigs, nil, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to get pipeline config")
	}
	if resp.StatusCode == http.StatusNotFound {
		_ = resp.Body.Close()
		return nil, "", nil
	}
	var result PipelineConfig
	err = decodeJSON(resp, &result)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to decode response")
	}
	return &result, resp.Header.Get("ETag"), nil
}

func (p *pipelineConfigsService) Create(ctx context.Context, pc *PipelineConfig) (*PipelineConfig, string, error) {
	b := newPipelineConfigCreateRequest(pc)

	resp, err := p.c.do(ctx, http.MethodPost, pipelineConfigsServicePath, acceptPipelineConfigs, nil, b)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to create pipeline config")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, "", errors.New("gocd: pipeline group not found")
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, "", fmt.Errorf("gocd: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var result PipelineConfig
	err = decodeJSON(resp, &result)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to decode response")
	}
	return &result, resp.Header.Get("ETag"), nil
}

func newPipelineConfigCreateRequest(pc *PipelineConfig) *pipelineConfigCreateRequest {
	return &pipelineConfigCreateRequest{
		Group: pc.Group,
		Pipeline: &PipelineConfig{
			LabelTemplate:        pc.LabelTemplate,
			LockBehavior:         pc.LockBehavior,
			Name:                 pc.Name,
			Template:             pc.Template,
			Origin:               pc.Origin,
			Parameters:           pc.Parameters,
			EnvironmentVariables: pc.EnvironmentVariables,
			Materials:            pc.Materials,
			Stages:               pc.Stages,
			TrackingTool:         pc.TrackingTool,
			Timer:                pc.Timer,
		},
	}
}

func (p *pipelineConfigsService) Update(ctx context.Context, etag string, body *PipelineConfig) (*PipelineConfig, string, error) {
	path := fmt.Sprintf("%s/%s", pipelineConfigsServicePath, url.PathEscape(*body.Name))
	headers := map[string]string{
		"If-Match": etag,
	}
	resp, err := p.c.do(ctx, http.MethodPut, path, acceptPipelineConfigs, headers, body)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to update pipeline config")
	}
	var result PipelineConfig
	err = decodeJSON(resp, &result)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to decode response")
	}
	return &result, resp.Header.Get("ETag"), nil
}

func (p *pipelineConfigsService) Delete(ctx context.Context, name string) error {
	resp, err := p.c.do(ctx, http.MethodDelete, pipelineConfigsServicePath+"/"+url.PathEscape(name), acceptPipelineConfigs, nil, nil)
	if err != nil {
		return errors.Wrap(err, "gocd: failed to delete pipeline config")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("gocd: failed to delete pipeline config: %s", resp.Status)
	}
	_ = resp.Body.Close()
	return nil
}
