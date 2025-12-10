/*
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

package elasticagentprofile

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/feature"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"

	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	apisv1alpha1 "github.com/marquesgui/provider-gocd/apis/v1alpha1"
	"github.com/marquesgui/provider-gocd/internal/features"
	"github.com/marquesgui/provider-gocd/pkg/gocd"
)

const (
	errNotElasticAgentProfile = "managed resource is not a ElasticAgentProfile custom resource"
	errTrackPCUsage           = "cannot track ProviderConfig usage"
	errGetPC                  = "cannot get ProviderConfig"
	errGetCreds               = "cannot get credentials"
	errNewClient              = "cannot create new Service"
	etagAnnotationKey         = "gocd.crossplane.io/etag"
)

var newService = func(creds []byte) (any, error) {
	cfg, err := apisv1alpha1.ParseGocdProviderConfig(creds)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse gocd provider config")
	}
	c, err := gocd.New(gocd.Config{
		BaseURL:  cfg.BaseURL,
		Username: cfg.Username,
		Password: cfg.Password,
		Token:    cfg.Token,
		Insecure: cfg.Insecure,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new GoCD client")
	}
	return c.ElasticAgentProfile(), nil
}

// Setup adds a controller that reconciles ElasticAgentProfile managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ElasticAgentProfileGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	opts := []managed.ReconcilerOption{
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newService,
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
			mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.ElasticAgentProfileList{}, o.MetricOptions.PollStateMetricInterval,
		)
		if err := mgr.Add(stateMetricsRecorder); err != nil {
			return errors.Wrap(err, "cannot register MR state metrics recorder for kind v1alpha1.ElasticAgentProfileList")
		}
	}

	r := managed.NewReconciler(mgr, resource.ManagedKind(v1alpha1.ElasticAgentProfileGroupVersionKind), opts...)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.ElasticAgentProfile{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (interface{}, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.ElasticAgentProfile)
	if !ok {
		return nil, errors.New(errNotElasticAgentProfile)
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

	s, ok := svc.(gocd.ElasticAgentProfileService)
	if !ok {
		return nil, errors.New("returned service does not implement gocd.ElasticAgentProfileService")
	}

	return &external{service: s, kube: c.kube}, nil
}

type external struct {
	service gocd.ElasticAgentProfileService
	kube    client.Client
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	ea, ok := mg.(*v1alpha1.ElasticAgentProfile)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotElasticAgentProfile)
	}

	got, etag, err := c.service.Get(ctx, ea.Spec.ForProvider.ID)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "provider-gocd: cannot get the elastic agent profile")
	}
	if got != nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	if err := updateStatus(ea, got); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "could not update the satus")
	}

	keepETag(ea, etag)
	upToDate := isUpToDate(ctx, c.kube, ea, got)

	if upToDate {
		ea.SetConditions(xpv1.Available())
	} else {
		ea.SetConditions(xpv1.Unavailable())
	}

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  upToDate,
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.ElasticAgentProfile)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotElasticAgentProfile)
	}

	elasticAgentProfile := gocd.ElasticAgentProfile{
		ID:               cr.Spec.ForProvider.ID,
		ClusterProfileID: cr.Spec.ForProvider.ClusterProfileID,
		Properties:       mapProperties(cr.Spec.ForProvider.Properties),
	}

	got, etag, err := c.service.Create(ctx, elasticAgentProfile)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "gocd: error creating a new elastic agent profile")
	}

	if got != nil {
		if err := updateStatus(cr, got); err != nil {
			return managed.ExternalCreation{}, errors.Wrap(err, "gocd: error while updating elastic agent profile status")
		}
	}
	keepETag(cr, etag)

	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.ElasticAgentProfile)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotElasticAgentProfile)
	}

	rb := gocd.ElasticAgentProfile{
		ID:               cr.Spec.ForProvider.ID,
		ClusterProfileID: cr.Spec.ForProvider.ClusterProfileID,
		Properties:       mapProperties(cr.Spec.ForProvider.Properties),
	}

	etag := cr.Annotations[etagAnnotationKey]
	got, newEtag, err := c.service.Update(ctx, rb, etag)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "error updating the elastic agent profile")
	}

	if err := updateStatus(cr, got); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "could not update the elastic agent status")
	}
	keepETag(cr, newEtag)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.ElasticAgentProfile)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotElasticAgentProfile)
	}

	if err := c.service.Delete(ctx, cr.Spec.ForProvider.ID); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "cannot delete the elastic agent profile")
	}

	return managed.ExternalDelete{}, nil
}

func (c *external) Disconnect(ctx context.Context) error {
	return nil
}

func updateStatus(ea *v1alpha1.ElasticAgentProfile, got *gocd.ElasticAgentProfileResponse) error {
	b, err := json.Marshal(got)
	if err != nil {
		return errors.Wrap(err, "provider-gocd: error marshalling elastic agent profile")
	}
	ea.Status.AtProvider = &runtime.RawExtension{Raw: b}
	return nil
}

func keepETag(ea *v1alpha1.ElasticAgentProfile, etag string) {
	if ea.Annotations == nil {
		ea.Annotations = map[string]string{}
	}
	if etag != "" {
		ea.Annotations[etagAnnotationKey] = etag
	}
}

func isUpToDate(ctx context.Context, kube client.Client, ea *v1alpha1.ElasticAgentProfile, got *gocd.ElasticAgentProfileResponse) bool {
	propertiesAreEqual := func(eaProperties []v1alpha1.ConfigProperty, gotProperties []gocd.ConfigProperty) bool {
		if len(eaProperties) != len(gotProperties) {
			return false
		}

		propertiesValue := make(map[string]string)
		propertiesCount := make(map[string]int)
		for _, p := range eaProperties {
			propertiesValue[p.Key] = p.Value
			propertiesCount[p.Key]++
		}

		for _, p := range gotProperties {
			v, ok := propertiesValue[p.Key]
			if !ok {
				return false
			}

			if v != p.Value {
				return false
			}

			propertiesCount[p.Key]--
			if propertiesCount[p.Key] <= 0 {
				return false
			}
		}

		return true
	}

	return ea.Spec.ForProvider.ID == got.ID &&
		ea.Spec.ForProvider.ClusterProfileID == got.ClusterProfileID &&
		propertiesAreEqual(ea.Spec.ForProvider.Properties, got.Properties)
}

func mapProperties(p []v1alpha1.ConfigProperty) []gocd.ConfigProperty {
	prts := make([]gocd.ConfigProperty, 0, 0)
	for _, v := range p {
		prts = append(prts, gocd.ConfigProperty{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return prts
}
