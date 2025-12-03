/*
Copyright 2025 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
ou may not use this file except in compliance with the License.
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
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCalculateHashes(t *testing.T) {
	type args struct {
		kube client.Client
		pc   v1alpha1.PipelineConfigForProvider
	}
	type want struct {
		hashes map[string]string
		err    error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"SimplePipeline": {
			reason: "Should correctly calculate hashes for a simple pipeline with only pipeline-level environment variables.",
			args: args{
				kube: fake.NewClientBuilder().Build(),
				pc: v1alpha1.PipelineConfigForProvider{
					EnvironmentVariables: []v1alpha1.EnvironmentVariable{
						{Name: "VAR1", Value: "VALUE1"},
						{Name: "VAR2", Value: "VALUE2"},
					},
				},
			},
			want: want{
				hashes: map[string]string{
					"pipeline.VAR1": ToSha256("VALUE1"),
					"pipeline.VAR2": ToSha256("VALUE2"),
				},
			},
		},
		"ComplexPipeline": {
			reason: "Should correctly calculate hashes for a complex pipeline with environment variables at all levels.",
			args: args{
				kube: fake.NewClientBuilder().WithRuntimeObjects(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{Name: "secret1", Namespace: "default"},
						Data:       map[string][]byte{"key1": []byte("secretValue")},
					},
				).Build(),
				pc: v1alpha1.PipelineConfigForProvider{
					EnvironmentVariables: []v1alpha1.EnvironmentVariable{
						{Name: "PIPELINE_VAR", Value: "pipelineValue"},
					},
					Stages: []v1alpha1.Stage{
						{
							Name: "stage1",
							EnvironmentVariables: []v1alpha1.EnvironmentVariable{
								{Name: "STAGE_VAR", Value: "stageValue"},
							},
							Jobs: []v1alpha1.Job{
								{
									Name: "job1",
									EnvironmentVariables: []v1alpha1.EnvironmentVariable{
										{Name: "JOB_VAR", ValueFrom: &v1alpha1.EnvVarSource{
											SecretKeyRef: &xpv1.SecretKeySelector{
												SecretReference: xpv1.SecretReference{
													Name:      "secret1",
													Namespace: "default",
												},
												Key: "key1",
											},
										}},
									},
								},
							},
						},
					},
				},
			},
			want: want{
				hashes: map[string]string{
					"pipeline.PIPELINE_VAR":   ToSha256("pipelineValue"),
					"stage.stage1.STAGE_VAR":  ToSha256("stageValue"),
					"job.stage1.job1.JOB_VAR": ToSha256("secretValue"),
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			hashes, err := calculateHashes(context.Background(), tc.args.kube, tc.args.pc)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ncalculateHashes(...): -want error, +got error:\n%s", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.hashes, hashes); diff != "" {
				t.Errorf("\n%s\ncalculateHashes(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}
