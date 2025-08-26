/*
Package pipelineconfig
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
package pipelineconfig

import (
	"context"
	"encoding/json"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/marquesgui/provider-gocd/pkg/gocd"

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
)

const (
	errNotPipelineConfig = "managed resource is not a PipelineConfig custom resource"
	errTrackPCUsage      = "cannot track ProviderConfig usage"
	errGetPC             = "cannot get ProviderConfig"
	errGetCreds          = "cannot get credentials"
	errNewClient         = "cannot create new Service"
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
	return c.PipelineConfigs(), nil
}

// Setup adds a controller that reconciles PipelineConfig managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.PipelineConfigGroupKind)

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
			mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.PipelineConfigList{}, o.MetricOptions.PollStateMetricInterval,
		)
		if err := mgr.Add(stateMetricsRecorder); err != nil {
			return errors.Wrap(err, "cannot register MR state metrics recorder for kind v1alpha1.PipelineConfigList")
		}
	}

	r := managed.NewReconciler(mgr, resource.ManagedKind(v1alpha1.PipelineConfigGroupVersionKind), opts...)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.PipelineConfig{}).
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
	cr, ok := mg.(*v1alpha1.PipelineConfig)
	if !ok {
		return nil, errors.New(errNotPipelineConfig)
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

	s, ok := svc.(gocd.PipelineConfigsService)
	if !ok {
		return nil, errors.New("returned service does not implement gocd.PipelineConfigsService")
	}
	return &external{service: s}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service gocd.PipelineConfigsService
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	pc, ok := mg.(*v1alpha1.PipelineConfig)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotPipelineConfig)
	}

	got, etag, err := c.service.Get(ctx, pc.Spec.ForProvider.Name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "cannot get pipeline config")
	}
	if got == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	err = updateStatus(pc, got)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "cannot update status")
	}

	keepETag(pc, etag)
	upToDate := isUpToDate(pc.Spec.ForProvider, got)

	if upToDate {
		pc.SetConditions(xpv1.Available())
	} else {
		pc.SetConditions(xpv1.Unavailable())
	}

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  upToDate,
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.PipelineConfig)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotPipelineConfig)
	}

	requestBody := mapAPIToDtoPipelineConfig(cr.Spec.ForProvider)
	out, etag, err := c.service.Create(ctx, requestBody)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "cannot create pipeline config")
	}
	if out != nil {
		err = updateStatus(cr, out)
		if err != nil {
			return managed.ExternalCreation{}, errors.Wrap(err, "cannot update status")
		}
	}
	keepETag(cr, etag)

	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.PipelineConfig)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotPipelineConfig)
	}

	requestBody := mapAPIToDtoPipelineConfig(cr.Spec.ForProvider)
	etag := ""
	if cr.Annotations != nil {
		etag = cr.Annotations["gocd.crossplane.io/etag"]
	}
	out, etag, err := c.service.Update(ctx, etag, requestBody)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "cannot update pipeline config")
	}
	keepETag(cr, etag)
	err = updateStatus(cr, out)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "cannot update status")
	}
	return managed.ExternalUpdate{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.PipelineConfig)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotPipelineConfig)
	}

	err := c.service.Delete(ctx, cr.Spec.ForProvider.Name)
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "cannot delete pipeline config")
	}
	return managed.ExternalDelete{}, nil
}

func (c *external) Disconnect(context.Context) error {
	return nil
}

func keepETag(pc *v1alpha1.PipelineConfig, etag string) {
	if pc.Annotations == nil {
		pc.Annotations = map[string]string{}
	}
	if etag != "" {
		pc.Annotations["gocd.crossplane.io/etag"] = etag
	}
}

func updateStatus(pc *v1alpha1.PipelineConfig, got *gocd.PipelineConfig) error {
	b, err := json.Marshal(got)
	if err != nil {
		return errors.Wrap(err, "error marshalling pipeline config parameters")
	}
	pc.Status.AtProvider = &runtime.RawExtension{Raw: b}
	return nil
}

func isUpToDate(pc v1alpha1.PipelineConfigForProvider, got *gocd.PipelineConfig) bool {
	desired := mapAPIToDtoPipelineConfig(pc)
	isEqual := desired.Equal(got)
	return isEqual
}
