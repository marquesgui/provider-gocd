package pipelineconfig

import (
	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	"github.com/marquesgui/provider-gocd/pkg/gocd"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func stringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intstrOrNil(i intstr.IntOrString) *intstr.IntOrString {
	switch i.Type {
	case intstr.Int:
		if i.IntVal == 0 {
			return nil
		}
	case intstr.String:
		if i.StrVal == "" {
			return nil
		}
	}
	return &i
}

func mapAPIToDtoPipelineConfig(cr v1alpha1.PipelineConfigForProvider) *gocd.PipelineConfig {
	return &gocd.PipelineConfig{
		Group:                stringOrNil(cr.Group),
		LabelTemplate:        stringOrNil(cr.LabelTemplate),
		LockBehavior:         gocd.PipelineConfigLockBehaviorFromString(cr.LockBehavior.String()),
		Name:                 stringOrNil(cr.Name),
		Template:             stringOrNil(cr.Template),
		Origin:               mapAPIOriginToDTO(cr.Origin),
		Parameters:           mapAPIParametersToDTO(cr.Parameters),
		EnvironmentVariables: mapAPIEnvironmentVariablesToDTO(cr.EnvironmentVariables),
		Materials:            mapAPIMaterialsToDTO(cr.Materials),
		Stages:               mapAPIStagesToDto(cr.Stages),
		TrackingTool:         mapAPITrackingToolToDTO(cr.TrackingTool),
		Timer:                mapAPITimerToDTO(cr.Timer),
	}
}

func mapAPITimerToDTO(timer v1alpha1.Timer) *gocd.PipelineConfigTimer {
	if timer.Spec == "" {
		return nil
	}

	return &gocd.PipelineConfigTimer{
		Spec:          timer.Spec,
		OnlyOnChanges: timer.OnlyOnChanges,
	}
}

func mapAPITrackingToolToDTO(tt v1alpha1.TrackingTool) *gocd.PipelineConfigTrackingTool {
	if tt.Type == "" {
		return nil
	}

	return &gocd.PipelineConfigTrackingTool{
		Type: tt.Type,
		Attributes: gocd.PipelineConfigTrackingToolAttributes{
			URLPattern: tt.Attributes.URLPattern,
			Regex:      tt.Attributes.Regex,
		},
	}
}

func mapAPIOriginToDTO(o v1alpha1.Origin) *gocd.PipelineConfigOrigin {
	return &gocd.PipelineConfigOrigin{
		Type: gocd.PipelineConfigOriginTypeFromString(string(o.Type)),
		ID:   stringOrNil(o.ID),
	}
}

func mapAPIStagesToDto(stages []v1alpha1.Stage) []gocd.PipelineConfigStage {
	out := make([]gocd.PipelineConfigStage, 0, len(stages))
	for _, v := range stages {
		out = append(out, gocd.PipelineConfigStage{
			Name:                  v.Name,
			FetchMaterials:        v.FetchMaterials,
			CleanWorkingDirectory: v.CleanWorkingDir,
			NeverCleanupArtifacts: v.NeverCleanupArtifacts,
			Approval:              mapAPIStageApprovalToDTO(v.Approval),
			EnvironmentVariables:  mapAPIEnvironmentVariablesToDTO(v.EnvironmentVariables),
			Jobs:                  mapAPIStageJobsToDto(v.Jobs),
		})
	}
	return out
}

func mapAPIStageJobsToDto(jobs []v1alpha1.Job) []gocd.PipelineConfigStageJobs {
	out := make([]gocd.PipelineConfigStageJobs, 0, len(jobs))
	for _, v := range jobs {
		out = append(out, gocd.PipelineConfigStageJobs{
			Name:                 v.Name,
			RunInstanceCount:     intstrOrNil(v.RunInstanceCount),
			Timeout:              v.Timeout,
			EnvironmentVariables: mapAPIEnvironmentVariablesToDTO(v.EnvironmentVariables),
			Resources:            v.Resources,
			Tasks:                mapAPIJobTasksToDTO(v.Tasks),
			Tabs:                 mapAPIJobTabsToDTO(v.Tabs),
			Artifacts:            mapAPIJobArtifactsToDTO(v.Artifacts),
			ElasticProfileID:     v.ElasticProfileID,
		})
	}
	return out
}

func mapAPIJobArtifactsToDTO(artifacts []v1alpha1.JobArtifact) []gocd.PipelineConfigStageJobsArtifact {
	out := make([]gocd.PipelineConfigStageJobsArtifact, 0, len(artifacts))
	for _, v := range artifacts {
		out = append(out, gocd.PipelineConfigStageJobsArtifact{
			Type:          gocd.PipelineConfigStageJobsArtifactTypeFromString(string(v.Type)),
			Source:        stringOrNil(v.Source),
			Destination:   v.Destination,
			ArtifactID:    stringOrNil(v.ID),
			StoreID:       v.StoreID,
			Configuration: mapJobTaskConfiguration(v.Configuration),
		})
	}
	return out
}

func mapAPIJobTabsToDTO(tabs []v1alpha1.JobTab) []gocd.PipelineConfigStageJobsTab {
	out := make([]gocd.PipelineConfigStageJobsTab, 0, len(tabs))
	for _, v := range tabs {
		out = append(out, gocd.PipelineConfigStageJobsTab{
			Name: v.Name,
			Path: v.Path,
		})
	}
	return out
}

func mapAPIJobTasksToDTO(tasks []v1alpha1.TaskWithCancel) []gocd.PipelineConfigStageJobsTask {
	out := make([]gocd.PipelineConfigStageJobsTask, 0, len(tasks))
	for _, v := range tasks {
		t := gocd.PipelineConfigStageJobsTask{
			Type: gocd.PipelineConfigStageJobsTaskTypeFromString(string(v.Type)),
		}

		switch t.Type {
		case gocd.PipelineConfigStageJobsTaskTypeExec:
			t.Attributes = mapAPIJobTaskExecWithCancelToDTO(v.ExecAttributes)
		case gocd.PipelineConfigStageJobsTaskTypeAnt:
			t.Attributes = mapAPIJobTaskAntWithCancelToDTO(v.AntAttributes)
		case gocd.PipelineConfigStageJobsTaskTypeNant:
			t.Attributes = mapAPIJobTaskNantWithCancelToDTO(v.NantAttributes)
		case gocd.PipelineConfigStageJobsTaskTypeRake:
			t.Attributes = mapAPIJobTaskRakeWithCancelToDTO(v.RakeAttributes)
		case gocd.PipelineConfigStageJobsTaskTypeFetch:
			t.Attributes = mapAPIJobTaskFetchWithCancelToDTO(v.FetchAttributes)
		case gocd.PipelineConfigStageJobsTaskTypePluggableTask:
			t.Attributes = mapAPIJobTaskPluggableWithCancelToDTO(v.PluggableAttributes)
		}
		out = append(out, t)
	}

	return out
}

func mapAPIJobTaskPluggableWithCancelToDTO(attributes *v1alpha1.TaskPluggableAttributesWithCancel) gocd.PipelineConfigStageJobsTaskAttributes {
	if attributes == nil {
		return nil
	}

	return &gocd.PipelineConfigStageJobsTaskAttributesPluggable{
		RunIf: mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		PluginConfiguration: gocd.PipelineConfigStageJobsTaskAttributesPluggableTaskPluginConfiguration{
			ID:      attributes.PluginConfiguration.ID,
			Version: attributes.PluginConfiguration.Version,
		},
		OnCancel:      mapAPIJobTaskToDTO(attributes.OnCancel),
		Configuration: mapJobTaskConfiguration(attributes.Configuration),
	}
}

func mapAPIJobTaskFetchWithCancelToDTO(attributes *v1alpha1.TaskFetchAttributesWithCancel) gocd.PipelineConfigStageJobsTaskAttributes {
	if attributes == nil {
		return nil
	}

	return &gocd.PipelineConfigStageJobsTaskAttributesFetch{
		ArtifactOrigin: gocd.PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeFromString(string(attributes.ArtifactOrigin)),
		RunIf:          mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		Pipeline:       attributes.Pipeline,
		Stage:          attributes.Stage,
		Job:            attributes.Job,
		Source:         attributes.Source,
		IsSourceAFile:  attributes.IsSourceAFile,
		Destination:    attributes.Destination,
		OnCancel:       mapAPIJobTaskToDTO(attributes.OnCancel),
		ArtifactID:     attributes.ArtifactID,
		Configuration:  mapJobTaskConfiguration(attributes.Configuration),
	}
}

func mapAPIJobTaskRakeWithCancelToDTO(attributes *v1alpha1.TaskRakeAttributesWithCancel) gocd.PipelineConfigStageJobsTaskAttributes {
	if attributes == nil {
		return nil
	}

	return &gocd.PipelineConfigStageJobsTaskAttributesRake{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		BuildFile:        attributes.BuildFile,
		Target:           attributes.Target,
		WorkingDirectory: attributes.WorkingDirectory,
		OnCancel:         mapAPIJobTaskToDTO(attributes.OnCancel),
	}
}

func mapAPIJobTaskNantWithCancelToDTO(attributes *v1alpha1.TaskNantAttributesWithCancel) gocd.PipelineConfigStageJobsTaskAttributes {
	if attributes == nil {
		return nil
	}

	return &gocd.PipelineConfigStageJobsTaskAttributesNant{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		BuildFile:        attributes.BuildFile,
		Target:           attributes.Target,
		NantPath:         attributes.NantPath,
		WorkingDirectory: attributes.WorkingDirectory,
		OnCancel:         mapAPIJobTaskToDTO(attributes.OnCancel),
	}
}

func mapAPIJobTaskAntWithCancelToDTO(attributes *v1alpha1.TaskAntAttributesWithCancel) gocd.PipelineConfigStageJobsTaskAttributes {
	if attributes == nil {
		return nil
	}

	return &gocd.PipelineConfigStageJobsTaskAttributesAnt{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		BuildFile:        attributes.BuildFile,
		Target:           attributes.Target,
		WorkingDirectory: attributes.WorkingDirectory,
		OnCancel:         mapAPIJobTaskToDTO(attributes.OnCancel),
	}
}

func mapAPIJobTaskExecWithCancelToDTO(attributes *v1alpha1.TaskExecAttributesWithCancel) gocd.PipelineConfigStageJobsTaskAttributes {
	if attributes == nil {
		return nil
	}

	return &gocd.PipelineConfigStageJobsTaskAttributesExec{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		Command:          attributes.Command,
		Arguments:        attributes.Arguments,
		WorkingDirectory: attributes.WorkingDirectory,
		OnCancel:         mapAPIJobTaskToDTO(attributes.OnCancel),
	}
}

func mapAPIJobTaskToDTO(cancel *v1alpha1.Task) *gocd.PipelineConfigStageJobsTask {
	if cancel == nil {
		return nil
	}

	out := &gocd.PipelineConfigStageJobsTask{
		Type: gocd.PipelineConfigStageJobsTaskTypeFromString(string(cancel.Type)),
	}

	switch out.Type {
	case gocd.PipelineConfigStageJobsTaskTypeExec:
		out.Attributes = mapAPIJobTaskExecToDTO(cancel.ExecAttributes)
	case gocd.PipelineConfigStageJobsTaskTypeAnt:
		out.Attributes = mapAPIJobTaskAntToDTO(cancel.AntAttributes)
	case gocd.PipelineConfigStageJobsTaskTypeNant:
		out.Attributes = mapAPIJobTaskNantToDTO(cancel.NantAttributes)
	case gocd.PipelineConfigStageJobsTaskTypeRake:
		out.Attributes = mapAPIJobTaskRakeToDTO(cancel.RakeAttributes)
	case gocd.PipelineConfigStageJobsTaskTypeFetch:
		out.Attributes = mapAPIJobTaskFetchToDTO(cancel.FetchAttributes)
	case gocd.PipelineConfigStageJobsTaskTypePluggableTask:
		out.Attributes = mapAPIJobTaskPluggableToDTO(cancel.PluggableAttributes)
	}

	return out
}

func mapAPIJobTaskPluggableToDTO(attributes *v1alpha1.TaskPluggableAttributes) gocd.PipelineConfigStageJobsTaskAttributes {
	return &gocd.PipelineConfigStageJobsTaskAttributesPluggable{
		RunIf: mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		PluginConfiguration: gocd.PipelineConfigStageJobsTaskAttributesPluggableTaskPluginConfiguration{
			ID:      attributes.PluginConfiguration.ID,
			Version: attributes.PluginConfiguration.Version,
		},
		Configuration: mapJobTaskConfiguration(attributes.Configuration),
	}
}

func mapAPIJobTaskFetchToDTO(attributes *v1alpha1.TaskFetchAttributes) gocd.PipelineConfigStageJobsTaskAttributes {
	return &gocd.PipelineConfigStageJobsTaskAttributesFetch{
		ArtifactOrigin: gocd.PipelineConfigStageJobsTaskAttributesFetchArtifactOriginTypeFromString(string(attributes.ArtifactOrigin)),
		RunIf:          mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		Pipeline:       attributes.Pipeline,
		Stage:          attributes.Stage,
		Job:            attributes.Job,
		Source:         attributes.Source,
		IsSourceAFile:  attributes.IsSourceAFile,
		Destination:    attributes.Destination,
		ArtifactID:     attributes.ArtifactID,
		Configuration:  mapJobTaskConfiguration(attributes.Configuration),
	}
}

func mapJobTaskConfiguration(configuration []v1alpha1.KeyValue) []gocd.ConfigProperty {
	out := make([]gocd.ConfigProperty, 0, len(configuration))
	for _, v := range configuration {
		out = append(out, gocd.ConfigProperty{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return out
}

func mapAPIJobTaskRakeToDTO(attributes *v1alpha1.TaskRakeAttributes) gocd.PipelineConfigStageJobsTaskAttributes {
	return &gocd.PipelineConfigStageJobsTaskAttributesRake{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		BuildFile:        attributes.BuildFile,
		Target:           attributes.Target,
		WorkingDirectory: attributes.WorkingDirectory,
	}
}

func mapAPIJobTaskNantToDTO(attributes *v1alpha1.TaskNantAttributes) gocd.PipelineConfigStageJobsTaskAttributes {
	return &gocd.PipelineConfigStageJobsTaskAttributesNant{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		BuildFile:        attributes.BuildFile,
		Target:           attributes.Target,
		NantPath:         attributes.NantPath,
		WorkingDirectory: attributes.WorkingDirectory,
	}
}

func mapAPIJobTaskAntToDTO(attributes *v1alpha1.TaskAntAttributes) gocd.PipelineConfigStageJobsTaskAttributes {
	return &gocd.PipelineConfigStageJobsTaskAttributesAnt{
		RunIf:            mapAPIJobTaskRunIfToDTO(attributes.RunIf),
		BuildFile:        attributes.BuildFile,
		Target:           attributes.Target,
		WorkingDirectory: attributes.WorkingDirectory,
	}
}

func mapAPIJobTaskExecToDTO(cancel *v1alpha1.TaskExecAttributes) gocd.PipelineConfigStageJobsTaskAttributes {
	return &gocd.PipelineConfigStageJobsTaskAttributesExec{
		RunIf:            mapAPIJobTaskRunIfToDTO(cancel.RunIf),
		Command:          cancel.Command,
		Arguments:        cancel.Arguments,
		WorkingDirectory: cancel.WorkingDirectory,
	}
}

func mapAPIJobTaskRunIfToDTO(runIf []v1alpha1.TaskAttributesRunIfTypes) []gocd.RunIfType {
	out := make([]gocd.RunIfType, 0, len(runIf))
	for _, v := range runIf {
		out = append(out, gocd.RunIfTypeFromString(string(v)))
	}
	return out
}

func mapAPIStageApprovalToDTO(approval v1alpha1.StageApproval) gocd.PipelineConfigApproval {
	return gocd.PipelineConfigApproval{
		Type:               gocd.PipelineConfigApprovalTypeFromString(string(approval.Type)),
		AllowOnlyOnSuccess: approval.AllowOnlyOnSuccess,
		Authorization: gocd.PipelineConfigApprovalAuthorization{
			Users: approval.Authorization.Users,
			Roles: approval.Authorization.Roles,
		},
	}
}

func mapAPIMaterialsToDTO(materials []v1alpha1.Material) []gocd.PipelineConfigMaterial {
	m := make([]gocd.PipelineConfigMaterial, 0, len(materials))
	for _, v := range materials {
		mat := gocd.PipelineConfigMaterial{
			Type: gocd.PipelineConfigMaterialTypeFromString(v.Type.String()),
		}
		var attr gocd.PipelineConfigMaterialAttributes
		switch mat.Type {
		case gocd.PipelineConfigMaterialTypeGit:
			attr = mapAPIMaterialGitAttributesToDTO(v.GitAttributes)
		case gocd.PipelineConfigMaterialTypeSvn:
			attr = mapAPIMaterialSvnAttributesToDTO(v.SvnAttributes)
		case gocd.PipelineConfigMaterialTypeP4:
			attr = mapAPIMaterialP4AttributesToDTO(v.P4Attributes)
		case gocd.PipelineConfigMaterialTypeHg:
			attr = mapAPIMaterialHgAttributesToDTO(v.HgAttributes)
		case gocd.PipelineConfigMaterialTypeTfs:
			attr = mapAPIMaterialTfsAttributesToDTO(v.TfsAttributes)
		case gocd.PipelineConfigMaterialTypeDependency:
			attr = mapAPIMaterialDependencyToDTO(v.DependencyAttributes)
		case gocd.PipelineConfigMaterialTypePackage:
			attr = mapAPIMaterialPackageToDTO(v.PackageAttributes)
		case gocd.PipelineConfigMaterialTypePlugin:
			attr = mapAPIMaterialPluginToDTO(v.PluginAttributes)
		}
		mat.Attributes = attr
		m = append(m, mat)
	}

	return m
}

func mapAPIMaterialPluginToDTO(attributes *v1alpha1.MaterialAttributesPlugin) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesPlugin{
		Ref:         stringOrNil(attributes.Ref),
		Destination: stringOrNil(attributes.Destination),
		Filter: &gocd.PipelineConfigMaterialFilter{
			Ignore:   attributes.Filter.Ignore,
			Includes: attributes.Filter.Includes,
		},
		InvertFilter: attributes.InvertFilter,
	}
}

func mapAPIMaterialPackageToDTO(attributes *v1alpha1.MaterialAttributesPackage) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesPackage{
		Ref: attributes.Ref,
	}
}

func mapAPIMaterialDependencyToDTO(attributes *v1alpha1.MaterialAttributesDependency) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesDependency{
		Name:                stringOrNil(attributes.Name),
		Pipeline:            stringOrNil(attributes.Pipeline),
		Stage:               stringOrNil(attributes.Stage),
		AutoUpdate:          attributes.AutoUpdate,
		IgnoreForScheduling: attributes.IgnoreForScheduling,
	}
}

func mapAPIMaterialTfsAttributesToDTO(attributes *v1alpha1.MaterialAttributesTfs) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesTfs{
		Name:              stringOrNil(attributes.Name),
		URL:               stringOrNil(attributes.URL),
		ProjectPath:       stringOrNil(attributes.ProjectPath),
		Domain:            stringOrNil(attributes.Domain),
		Username:          stringOrNil(attributes.Username),
		Password:          stringOrNil(attributes.Password),
		EncryptedPassword: stringOrNil(attributes.EncryptedPassword),
		Destination:       stringOrNil(attributes.Destination),
		AutoUpdate:        attributes.AutoUpdate,
		Filter:            mapAPIFilterToDTO(attributes.Filter),
		InvertFilter:      attributes.InvertFilter,
	}
}

func mapAPIMaterialHgAttributesToDTO(attributes *v1alpha1.MaterialAttributesHg) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesHg{
		Name:              stringOrNil(attributes.Name),
		URL:               stringOrNil(attributes.URL),
		Username:          stringOrNil(attributes.Username),
		Password:          stringOrNil(attributes.Password),
		EncryptedPassword: stringOrNil(attributes.EncryptedPassword),
		Branch:            stringOrNil(attributes.Branch),
		Destination:       stringOrNil(attributes.Destination),
		Filter:            mapAPIFilterToDTO(attributes.Filter),
		InvertFilter:      attributes.InvertFilter,
		AutoUpdate:        attributes.AutoUpdate,
	}
}

func mapAPIMaterialP4AttributesToDTO(attributes *v1alpha1.MaterialAttributesP4) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesP4{
		Name:              stringOrNil(attributes.Name),
		Port:              stringOrNil(attributes.Port),
		UseTickets:        attributes.UseTickets,
		View:              stringOrNil(attributes.View),
		Username:          stringOrNil(attributes.Username),
		Password:          stringOrNil(attributes.Password),
		EncryptedPassword: stringOrNil(attributes.EncryptedPassword),
		Destination:       stringOrNil(attributes.Destination),
		Filter:            mapAPIFilterToDTO(attributes.Filter),
		InvertFilter:      attributes.InvertFilter,
		AutoUpdate:        attributes.AutoUpdate,
	}
}

func mapAPIMaterialSvnAttributesToDTO(attributes *v1alpha1.MaterialAttributesSvn) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesSvn{
		Name:              stringOrNil(attributes.Name),
		URL:               stringOrNil(attributes.URL),
		Username:          stringOrNil(attributes.Username),
		Password:          stringOrNil(attributes.Password),
		EncryptedPassword: stringOrNil(attributes.EncryptedPassword),
		Destination:       stringOrNil(attributes.Destination),
		Filter:            mapAPIFilterToDTO(attributes.Filter),
		InvertFilter:      attributes.InvertFilter,
		AutoUpdate:        attributes.AutoUpdate,
		CheckExternals:    attributes.CheckExternals,
	}
}

func mapAPIMaterialGitAttributesToDTO(attributes *v1alpha1.MaterialAttributesGit) gocd.PipelineConfigMaterialAttributes {
	return &gocd.PipelineConfigMaterialAttributesGit{
		Name:            stringOrNil(attributes.Name),
		URL:             stringOrNil(attributes.URL),
		Username:        stringOrNil(attributes.Username),
		Password:        stringOrNil(attributes.Password),
		Branch:          stringOrNil(attributes.Branch),
		Destination:     stringOrNil(attributes.Destination),
		AutoUpdate:      attributes.AutoUpdate,
		Filter:          mapAPIFilterToDTO(attributes.Filter),
		InvertFilter:    attributes.InvertFilter,
		SubmoduleFolder: stringOrNil(attributes.SubmoduleFolder),
		ShallowClone:    attributes.ShallowClone,
	}
}

func mapAPIFilterToDTO(f v1alpha1.Filter) *gocd.PipelineConfigMaterialFilter {
	if len(f.Ignore) == 0 {
		return nil
	}

	return &gocd.PipelineConfigMaterialFilter{
		Ignore:   f.Ignore,
		Includes: f.Includes,
	}
}

func mapAPIEnvironmentVariablesToDTO(variables []v1alpha1.EnvironmentVariable) []gocd.EnvironmentVariable {
	out := make([]gocd.EnvironmentVariable, 0, len(variables))
	for _, v := range variables {
		out = append(out, gocd.EnvironmentVariable{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return out
}

func mapAPIParametersToDTO(parameters []v1alpha1.Parameter) []gocd.PipelineConfigParameter {
	out := make([]gocd.PipelineConfigParameter, 0, len(parameters))
	for _, p := range parameters {
		out = append(out, gocd.PipelineConfigParameter{
			Name:  p.Name,
			Value: p.Value,
		})
	}
	return out
}
