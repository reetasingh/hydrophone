/*
Copyright 2024 The Kubernetes Authors.

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

package common

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestValidateArgs(t *testing.T) {
	// Create a temporary output directory
	tempDir := t.TempDir()
	viper.Set("output-dir", tempDir)

	// Set up the test cases
	testCases := []struct {
		name          string
		extraArgs     []string
		expectedArgs  []string
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "With extra args",
			extraArgs:     []string{"--key1=value1", "--key2=value2"},
			expectedArgs:  []string{"--key1=value1", "--key2=value2"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name:          "Invalid extra args format",
			extraArgs:     []string{"invalid-arg"},
			expectedArgs:  []string{},
			wantErr:       true,
			expectedErr:   "expected [[invalid-arg]] in [[invalid-arg]] to be of --key=value format",
		},
		{
			name:          "Extra args with missing values",
			extraArgs:     []string{"--key1=value1", "--key2"},
			expectedArgs:  []string{},
			wantErr:       true,
			expectedErr:   "expected [[--key2]] in [[--key1=value1 --key2]] to be of --key=value format",
		},
		{
			name:          "Extra args with invalid key format",
			extraArgs:     []string{"key1=value1", "--key2=value2"},
			expectedArgs:  []string{},
			wantErr:       true,
			expectedErr:   "expected key [key1] in [[key1=value1 --key2=value2]] to start with prefix --",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the test environment
			viper.Set("extra-args", tc.extraArgs)

			// Call the function under test
			err := ValidateArgs()
			if tc.wantErr {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, viper.GetStringSlice("extra-args"), tc.expectedArgs)
			}
		})
	}
}

func TestTrimVersion(t *testing.T) {

	testCases := []struct {
		name            string
		version         string
		expectedVersion string
	}{
		{
			name:            "stable version",
			version:         "v1.28.6",
			expectedVersion: "1.28.6",
		},
		{
			name:            "pre released version",
			version:         "v1.28.6+0fb426",
			expectedVersion: "1.28.6",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trimmedVersion, _ := trimVersion(tc.version)
			assert.Equal(t, tc.expectedVersion, trimmedVersion)
		})
	}
}
