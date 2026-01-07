/*
Package role
Copyright 2025 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package role

import (
	"context"
	"encoding/json"
	"maps"
	"slices"
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"
	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	apisv1alpha1 "github.com/marquesgui/provider-gocd/apis/v1alpha1"
	"github.com/marquesgui/provider-gocd/internal/controller/helper"
	"github.com/marquesgui/provider-gocd/internal/features"
	"github.com/marquesgui/provider-gocd/pkg/gocd"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errNotRole      = "managed resource is not a role custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"
	errNewClient    = "cannot create new Service"
)

type gocdRoleService interface {
	Get(ctx context.Context, id string) (*gocd.Role, string, error)
	Create(ctx context.Context, cfg gocd.Role) (*gocd.Role, string, error)
	Update(ctx context.Context, id string, cfg gocd.Role, etag string) (*gocd.Role, string, error)
	Delete(ctx context.Context, id string, etag string) error
}

// newServiceFn builds the real GoCD Role service from credentials bytes.
var newServiceFn = func(creds []byte) (any, error) {
	type c struct {
		BaseURL  string `json:"baseURL"`
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token"`
		Insecure bool   `json:"insecure"`
	}
	var cfg c
	_ = json.Unmarshal(creds, &cfg)
	gc, err := gocd.New(gocd.Config{
		BaseURL:  cfg.BaseURL,
		Username: cfg.Username,
		Password: cfg.Password,
		Token:    cfg.Token,
		Insecure: cfg.Insecure,
	})
	if err != nil {
		return nil, err
	}
	return gc.Roles(), nil
}

// Setup adds a controller that reconciles role managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RoleGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	opts := []managed.ReconcilerOption{
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newServiceFn,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...),
		managed.WithManagementPolicies(),
	}

	if o.Features.Enabled(feature.EnableAlphaChangeLogs) {
		opts = append(opts, managed.WithChangeLogger(o.ChangeLogOptions.ChangeLogger))
	}

	if o.MetricOptions != nil {
		opts = append(opts, managed.WithMetricRecorder(o.MetricOptions.MRMetrics))
	}

	if o.MetricOptions != nil && o.MetricOptions.MRStateMetrics != nil {
		stateMetricsRecorder := statemetrics.NewMRStateRecorder(
			mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.RoleList{}, o.MetricOptions.PollStateMetricInterval,
		)
		if err := mgr.Add(stateMetricsRecorder); err != nil {
			return errors.Wrap(err, "cannot register MR state metrics recorder for kind v1alpha1.roleList")
		}
	}

	r := managed.NewReconciler(mgr, resource.ManagedKind(v1alpha1.RoleGroupVersionKind), opts...)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.Role{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (any, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Role)
	if !ok {
		return nil, errors.New(errNotRole)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	rs, ok := svc.(gocdRoleService)
	if !ok {
		return nil, errors.New("returned service does not implement gocdRoleService")
	}

	return &external{service: rs}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	service gocdRoleService
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	r, ok := mg.(*v1alpha1.Role)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRole)
	}

	name := meta.GetExternalName(r)
	if name == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	got, etag, err := c.service.Get(ctx, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "cannot get role")
	}
	if got == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	updateStatus(r, got)
	helper.KeepETag(r, etag)

	upToDate := isUpToDate(r, got)
	if upToDate {
		r.SetConditions(xpv1.Available())
	} else {
		r.SetConditions(xpv1.Unavailable())
	}

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  upToDate,
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Role)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRole)
	}

	name := helper.GetID(cr, cr.Spec.ForProvider.Name)

	in := createRoleRequest(name, cr)
	out, etag, err := c.service.Create(ctx, in)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "cannot create role")
	}
	if out != nil {
		meta.SetExternalName(cr, out.Name)
		cr.Status.AtProvider.Name = out.Name
	}
	helper.KeepETag(cr, etag)

	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Role)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRole)
	}

	name := meta.GetExternalName(cr)
	in := createRoleRequest(name, cr)

	etag := helper.GetETag(cr)
	out, newETag, err := c.service.Update(ctx, name, in, etag)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "cannot update role")
	}
	if out != nil {
		updateStatus(cr, out)
	}
	helper.KeepETag(cr, newETag)

	return managed.ExternalUpdate{ConnectionDetails: managed.ConnectionDetails{}}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Role)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotRole)
	}

	name := meta.GetExternalName(cr)
	etag := helper.GetETag(cr)

	if err := c.service.Delete(ctx, name, etag); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "cannot delete role")
	}
	return managed.ExternalDelete{}, nil
}

func (c *external) Disconnect(_ context.Context) error {
	return nil
}

func isUpToDate(cr *v1alpha1.Role, got *gocd.Role) bool { //nolint:gocyclo
	current := &cr.Spec.ForProvider
	if meta.GetExternalName(cr) != got.Name || current.Type != got.Type {
		return false
	}

	if current.Attributes.AuthConfigID != got.Attributes.AuthConfigID {
		return false
	}

	propertiesIsUpToDate := func(current []v1alpha1.KeyValue, got []gocd.ConfigProperty) bool {
		if len(current) != len(got) {
			return false
		}
		mapCurrent := make(map[string]string)
		for _, v := range current {
			mapCurrent[v.Key] = v.Value
		}
		mapGot := make(map[string]string)
		for _, v := range got {
			mapGot[v.Key] = v.Value
		}
		return maps.Equal(mapCurrent, mapGot)
	}(current.Attributes.Properties, got.Attributes.Properties)

	userIsUpToDate := func(current []string, got []string) bool {
		if len(current) != len(got) {
			return false
		}
		ac := slices.Clone(current)
		bc := slices.Clone(got)
		slices.Sort(ac)
		slices.Sort(bc)
		return slices.Equal(ac, bc)
	}(current.Attributes.Users, got.Attributes.Users)

	policiesIsUpToDate := func(c []v1alpha1.RoleParametersPolicy, g []gocd.Policy) bool {
		if len(c) != len(g) {
			return false
		}

		type sortablePolicy struct {
			Permission string
			Action     string
			Type       string
			Resource   string
		}

		toPolicies := func(in any) []sortablePolicy {
			var result []sortablePolicy
			switch v := in.(type) {
			case []v1alpha1.RoleParametersPolicy:
				for _, p := range v {
					result = append(result, sortablePolicy{
						Permission: p.Permission,
						Action:     p.Action,
						Type:       p.Type,
						Resource:   p.Resource,
					})
				}
			case []gocd.Policy:
				for _, p := range v {
					result = append(result, sortablePolicy{
						Permission: p.Permission,
						Action:     p.Action,
						Type:       p.Type,
						Resource:   p.Resource,
					})
				}
			}
			return result
		}

		cp := toPolicies(c)
		gp := toPolicies(g)

		sortFunc := func(a, b sortablePolicy) int {
			if a.Permission != b.Permission {
				return strings.Compare(a.Permission, b.Permission)
			}
			if a.Action != b.Action {
				return strings.Compare(a.Action, b.Action)
			}
			if a.Type != b.Type {
				return strings.Compare(a.Type, b.Type)
			}
			return strings.Compare(a.Resource, b.Resource)
		}
		slices.SortFunc(cp, sortFunc)
		slices.SortFunc(gp, sortFunc)

		return slices.Equal(cp, gp)
	}(current.Policy, got.Policy)

	return propertiesIsUpToDate && userIsUpToDate && policiesIsUpToDate
}

func updateStatus(cr *v1alpha1.Role, got *gocd.Role) {
	cr.Status.AtProvider.Name = got.Name
	cr.Status.AtProvider.Type = got.Type

	if atts := got.Attributes; atts != nil {
		cr.Status.AtProvider.Attributes.Users = atts.Users
		if atts.AuthConfigID != "" {
			cr.Status.AtProvider.Attributes.AuthConfigID = got.Attributes.AuthConfigID
			cr.Status.AtProvider.Attributes.Properties = make([]v1alpha1.KeyValue, 0)
			for _, p := range got.Attributes.Properties {
				cr.Status.AtProvider.Attributes.Properties = append(cr.Status.AtProvider.Attributes.Properties, v1alpha1.KeyValue{
					Key:   p.Key,
					Value: p.Value,
				})
			}
		}
	}
	if got.Links != nil {
		if got.Links.Self != nil {
			cr.Status.AtProvider.Links.Self.Href = got.Links.Self.Href
		}
		if got.Links.Doc != nil {
			cr.Status.AtProvider.Links.Doc.Href = got.Links.Doc.Href
		}
		if got.Links.Find != nil {
			cr.Status.AtProvider.Links.Find.Href = got.Links.Find.Href
		}
	}

	if policies := got.Policy; policies != nil {
		cr.Status.AtProvider.Policy = make([]v1alpha1.RoleParametersPolicy, 0)
		for _, p := range policies {
			cr.Status.AtProvider.Policy = append(cr.Status.AtProvider.Policy, v1alpha1.RoleParametersPolicy{
				Permission: p.Permission,
				Action:     p.Action,
				Type:       p.Type,
				Resource:   p.Resource,
			})
		}
	}
}

func createRoleRequest(name string, cr *v1alpha1.Role) gocd.Role {
	prop := make([]gocd.ConfigProperty, 0)
	for _, v := range cr.Spec.ForProvider.Attributes.Properties {
		prop = append(prop, gocd.ConfigProperty{Key: v.Key, Value: v.Value})
	}

	poly := make([]gocd.Policy, 0)
	for _, v := range cr.Spec.ForProvider.Policy {
		poly = append(poly, gocd.Policy{
			Permission: v.Permission,
			Action:     v.Action,
			Type:       v.Type,
			Resource:   v.Resource,
		})
	}
	return gocd.Role{
		Name: name,
		Type: cr.Spec.ForProvider.Type,
		Attributes: &gocd.RoleAttributes{
			AuthConfigID: cr.Spec.ForProvider.Attributes.AuthConfigID,
			Users:        cr.Spec.ForProvider.Attributes.Users,
			Properties:   prop,
		},
		Policy: poly,
	}
}
