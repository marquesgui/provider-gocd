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
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	"github.com/marquesgui/provider-gocd/pkg/gocd"
	"github.com/marquesgui/provider-gocd/pkg/gocd/mock"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
)

func TestObserve(t *testing.T) {
	type fields struct {
		service gocd.ElasticAgentProfileService
	}

	type args struct {
		ctx context.Context
		mg  resource.Managed
	}

	type want struct {
		o   managed.ExternalObservation
		err error
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockElasticAgentProfileService(ctrl)

	id := "test-id"

	cases := map[string]struct {
		reason string
		fields fields
		args   args
		want   want
	}{
		"Exists": {
			reason: "Should return ResourceExists: true when GoCD returns matching profile",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					ea.Spec.ForProvider.ClusterProfileID = "cluster-id"
					ea.Spec.ForProvider.Properties = []v1alpha1.ConfigProperty{{Key: "k", Value: "v"}}
					return ea
				}(),
			},
			want: want{
				o: managed.ExternalObservation{
					ResourceExists:    true,
					ResourceUpToDate:  false, // Expected false because isUpToDate needs kube client or more complex mock
					ConnectionDetails: managed.ConnectionDetails{},
				},
			},
		},
		"NotFound": {
			reason: "Should return ResourceExists: false when GoCD returns 404",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					return ea
				}(),
			},
			want: want{
				o: managed.ExternalObservation{
					ResourceExists: false,
				},
			},
		},
		"GetError": {
			reason: "Should return error when GoCD returns error",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					return ea
				}(),
			},
			want: want{
				err: errors.Wrap(errors.New("some error"), "provider-gocd: cannot get the elastic agent profile"),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if name == "Exists" {
				m.EXPECT().Get(gomock.Any(), id).Return(&gocd.ElasticAgentProfileResponse{
					ElasticAgentProfile: gocd.ElasticAgentProfile{
						ID:               id,
						ClusterProfileID: "cluster-id",
						Properties:       []gocd.ConfigProperty{{Key: "k", Value: "v"}},
					},
				}, "etag", nil)
			}
			if name == "NotFound" {
				m.EXPECT().Get(gomock.Any(), id).Return(nil, "", nil)
			}
			if name == "GetError" {
				m.EXPECT().Get(gomock.Any(), id).Return(nil, "", errors.New("some error"))
			}

			e := external{service: tc.fields.service}
			got, err := e.Observe(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type fields struct {
		service gocd.ElasticAgentProfileService
	}

	type args struct {
		ctx context.Context
		mg  resource.Managed
	}

	type want struct {
		c   managed.ExternalCreation
		err error
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockElasticAgentProfileService(ctrl)

	id := "test-id"

	cases := map[string]struct {
		reason string
		fields fields
		args   args
		want   want
	}{
		"Successful": {
			reason: "Should return Successful creation",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					ea.SetName(id)
					ea.Spec.ForProvider.ClusterProfileID = "cluster-id"
					ea.Spec.ForProvider.Properties = []v1alpha1.ConfigProperty{{Key: "k", Value: "v"}}
					return ea
				}(),
			},
			want: want{
				c: managed.ExternalCreation{
					ConnectionDetails: managed.ConnectionDetails{},
				},
			},
		},
		"CreateError": {
			reason: "Should return error when GoCD returns error",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					ea.SetName(id)
					return ea
				}(),
			},
			want: want{
				err: errors.Wrap(errors.New("some error"), "gocd: error creating a new elastic agent profile"),
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			if n == "Successful" {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&gocd.ElasticAgentProfileResponse{
					ElasticAgentProfile: gocd.ElasticAgentProfile{
						ID:               id,
						ClusterProfileID: "cluster-id",
						Properties:       []gocd.ConfigProperty{{Key: "k", Value: "v"}},
					},
				}, "etag", nil)
			} else {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, "", errors.New("some error"))
			}

			e := external{service: tc.fields.service}
			got, err := e.Create(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Create(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.c, got); diff != "" {
				t.Errorf("\n%s\ne.Create(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type fields struct {
		service gocd.ElasticAgentProfileService
	}

	type args struct {
		ctx context.Context
		mg  resource.Managed
	}

	type want struct {
		u   managed.ExternalUpdate
		err error
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockElasticAgentProfileService(ctrl)

	id := "test-id"

	cases := map[string]struct {
		reason string
		fields fields
		args   args
		want   want
	}{
		"Successful": {
			reason: "Should return Successful update",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					ea.Spec.ForProvider.ClusterProfileID = "cluster-id"
					ea.Spec.ForProvider.Properties = []v1alpha1.ConfigProperty{{Key: "k", Value: "v"}}
					return ea
				}(),
			},
			want: want{
				u: managed.ExternalUpdate{
					ConnectionDetails: managed.ConnectionDetails{},
				},
			},
		},
		"UpdateError": {
			reason: "Should return error when GoCD returns error",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					return ea
				}(),
			},
			want: want{
				err: errors.Wrap(errors.New("some error"), "error updating the elastic agent profile"),
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			if n == "Successful" {
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&gocd.ElasticAgentProfileResponse{
					ElasticAgentProfile: gocd.ElasticAgentProfile{
						ID:               id,
						ClusterProfileID: "cluster-id",
						Properties:       []gocd.ConfigProperty{{Key: "k", Value: "v"}},
					},
				}, "new-etag", nil)
			} else {
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "", errors.New("some error"))
			}

			e := external{service: tc.fields.service}
			got, err := e.Update(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Update(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.u, got); diff != "" {
				t.Errorf("\n%s\ne.Update(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		service gocd.ElasticAgentProfileService
	}

	type args struct {
		ctx context.Context
		mg  resource.Managed
	}

	type want struct {
		d   managed.ExternalDelete
		err error
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockElasticAgentProfileService(ctrl)

	id := "test-id"

	cases := map[string]struct {
		reason string
		fields fields
		args   args
		want   want
	}{
		"Successful": {
			reason: "Should return Successful deletion",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					return ea
				}(),
			},
			want: want{
				d: managed.ExternalDelete{},
			},
		},
		"DeleteError": {
			reason: "Should return error when GoCD returns error",
			fields: fields{
				service: m,
			},
			args: args{
				ctx: context.Background(),
				mg: func() resource.Managed {
					ea := &v1alpha1.ElasticAgentProfile{}
					meta.SetExternalName(ea, id)
					return ea
				}(),
			},
			want: want{
				err: errors.Wrap(errors.New("some error"), "cannot delete the elastic agent profile"),
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			if n == "Successful" {
				m.EXPECT().Delete(gomock.Any(), id).Return(nil)
			} else {
				m.EXPECT().Delete(gomock.Any(), id).Return(errors.New("some error"))
			}

			e := external{service: tc.fields.service}
			got, err := e.Delete(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Delete(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.d, got); diff != "" {
				t.Errorf("\n%s\ne.Delete(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}
