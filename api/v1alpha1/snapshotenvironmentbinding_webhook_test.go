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
			expectedError: "application field cannot be updated to test-app-a-changed",
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
			expectedError: "environment field cannot be updated to test-env-a-changed",
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

func TestSnapshotEnvironmentBindingValidateCreate(t *testing.T) {
	tests := []struct {
		testName      string                       // Name of test
		testData      SnapshotEnvironmentBinding   // Test data to be passed to ValidateCreate function
		existingSEBs  []SnapshotEnvironmentBinding // Existing SnapshotEnvironmentBindings for the namespace
		expectedError string                       // Expected error message from ValidateCreate function
	}{
		{
			testName: "No error when no existing SnapshotEnvironmentBindings with the same combination",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"test-key-b": "test-value-b"},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-b",
					Environment: "test-env-b",
				},
			},
			existingSEBs:  []SnapshotEnvironmentBinding{},
			expectedError: "",
		},

		{
			testName: "Error occurs when an existing SnapshotEnvironmentBinding has the same combination",
			testData: SnapshotEnvironmentBinding{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"test-key-c": "test-value-c"},
				},
				Spec: SnapshotEnvironmentBindingSpec{
					Application: "test-app-c",
					Environment: "test-env-c",
				},
			},
			existingSEBs: []SnapshotEnvironmentBinding{
				{
					ObjectMeta: v1.ObjectMeta{
						Labels: map[string]string{"test-key-d": "test-value-d"},
					},
					Spec: SnapshotEnvironmentBindingSpec{
						Application: "test-app-c",
						Environment: "test-env-c",
					},
				},
			},
			expectedError: "duplicate combination of Application (test-app-c) and Environment (test-env-c)",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualError := test.testData.ValidateCreate()

			if test.expectedError == "" {
				assert.Nil(t, actualError)
			} else {
				assert.Contains(t, actualError.Error(), test.expectedError)
			}
		})
	}
}
