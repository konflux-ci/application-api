/*
Copyright 2021-2022 Red Hat, Inc.

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

package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSnapshotEnvironmentBindingValidatingWebhook(t *testing.T) {

	originalBinding := SnapshotEnvironmentBinding{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{"test-key-a": "test-value-a"},
		},
		Spec: SnapshotEnvironmentBindingSpec{
			Application: "test-app-a",
			Environment: "test-env-a",
		},
	}

	tests := []struct {
		testName      string                     // Name of test
		testData      SnapshotEnvironmentBinding // Test data to be passed to webhook function
		expectedError string                     // Expected error message from webhook function
	}{
		{
			testName: "No error when Spec is same.",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"test-key-a": "test-value-a"},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-a",
					Environment: "test-env-a",
				},
			},
			expectedError: "",
		},

		{
			testName: "Error occurs when Spec.Application name is changed.",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"test-key-a": "test-value-a"},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-a-changed",
					Environment: "test-env-a",
				},
			},
			expectedError: "application cannot be updated to test-app-a-changed",
		},

		{
			testName: "Error occurs when Spec.Environment name is changed.",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"test-key-a": "test-value-a"},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-a",
					Environment: "test-env-a-changed",
				},
			},
			expectedError: "environment cannot be updated to test-env-a-changed",
		},

		{
			testName: "Error occurs when existing label is changed.",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"test-key-a": "test-value-b"},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-a",
					Environment: "test-env-a-changed",
				},
			},
			expectedError: "labels cannot be updated to map[test-key-a:test-value-b]",
		},

		{
			testName: "Error occurs when new label is added.",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"test-key-a": "test-value-a",
						"test-key-b": "test-value-b",
					},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-a",
					Environment: "test-env-a-changed",
				},
			},
			expectedError: "labels cannot be updated to map[test-key-a:test-value-a test-key-b:test-value-b]",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualError := test.testData.ValidateUpdate(&originalBinding)

			if test.expectedError == "" {
				assert.Nil(t, actualError)
			} else {
				assert.Contains(t, actualError.Error(), test.expectedError)
			}
		})
	}
}
