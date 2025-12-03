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

package pipelineconfig

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ToSha256 returns the sha256 hash of a given string
func ToSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func calculateHashes(ctx context.Context, kube client.Client, pc v1alpha1.PipelineConfigForProvider) (map[string]string, error) {
	hashes := make(map[string]string)

	// Helper function to process environment variables
	processEnvVars := func(envVars []v1alpha1.EnvironmentVariable, prefix string) error {
		for _, v := range envVars {
			var value string
			if v.Value != "" {
				value = v.Value
			} else if v.ValueFrom != nil {
				var err error
				value, _, err = GetValueFrom(ctx, kube, v.ValueFrom)
				if err != nil {
					return errors.Wrapf(err, "failed to get value for environment variable %s", v.Name)
				}
			}
			key := fmt.Sprintf("%s.%s", prefix, v.Name)
			hashes[key] = ToSha256(value)
		}
		return nil
	}

	// Pipeline-level environment variables
	if err := processEnvVars(pc.EnvironmentVariables, "pipeline"); err != nil {
		return nil, err
	}

	// Stage-level environment variables
	for _, s := range pc.Stages {
		stagePrefix := fmt.Sprintf("stage.%s", s.Name)
		if err := processEnvVars(s.EnvironmentVariables, stagePrefix); err != nil {
			return nil, err
		}

		// Job-level environment variables
		for _, j := range s.Jobs {
			jobPrefix := fmt.Sprintf("job.%s.%s", s.Name, j.Name)
			if err := processEnvVars(j.EnvironmentVariables, jobPrefix); err != nil {
				return nil, err
			}
		}
	}

	return hashes, nil
}
