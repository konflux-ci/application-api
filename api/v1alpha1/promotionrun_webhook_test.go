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
)

func TestPromotionRunValidatingWebhook(t *testing.T) {

	originalPromotionRun := PromotionRun{
		Spec: PromotionRunSpec{
			Snapshot:    "test-snapshot-a",
			Application: "test-app-a",
			ManualPromotion: ManualPromotionConfiguration{
				TargetEnvironment: "test-env-a",
			},
		},
	}

	tests := []struct {
		testName      string       // Name of test
		testData      PromotionRun // Test data to be passed to webhook function
		expectedError string       // Expected error message from webhook function
		warnings      []string
	}{
		{
			testName: "No error when Spec is same.",
			testData: PromotionRun{
				Spec: PromotionRunSpec{
					Snapshot:    "test-snapshot-a",
					Application: "test-app-a",
					ManualPromotion: ManualPromotionConfiguration{
						TargetEnvironment: "test-env-a",
					},
				},
			},
			expectedError: "",
		},

		{
			testName: "Error occurs when Spec.Snapshot is changed.",
			testData: PromotionRun{
				Spec: PromotionRunSpec{
					Snapshot:    "test-snapshot-a-changed",
					Application: "test-app-a",
					ManualPromotion: ManualPromotionConfiguration{
						TargetEnvironment: "test-env-a",
					},
				},
			},
			expectedError: "spec cannot be updated to {Snapshot:test-snapshot-a-changed Application:test-app-a ManualPromotion:{TargetEnvironment:test-env-a} AutomatedPromotion:{InitialEnvironment:}}",
		},

		{
			testName: "Error occurs when Spec.Application is changed.",
			testData: PromotionRun{
				Spec: PromotionRunSpec{
					Snapshot:    "test-snapshot-a",
					Application: "test-app-a-changed",
					ManualPromotion: ManualPromotionConfiguration{
						TargetEnvironment: "test-env-a",
					},
				},
			},
			expectedError: "spec cannot be updated to {Snapshot:test-snapshot-a Application:test-app-a-changed ManualPromotion:{TargetEnvironment:test-env-a} AutomatedPromotion:{InitialEnvironment:}}",
		},

		{
			testName: "Error occurs when Spec.Application is changed.",
			testData: PromotionRun{
				Spec: PromotionRunSpec{
					Snapshot:    "test-snapshot-a",
					Application: "test-app-a-changed",
					ManualPromotion: ManualPromotionConfiguration{
						TargetEnvironment: "test-env-a-changed",
					},
				},
			},
			expectedError: "spec cannot be updated to {Snapshot:test-snapshot-a Application:test-app-a-changed ManualPromotion:{TargetEnvironment:test-env-a-changed} AutomatedPromotion:{InitialEnvironment:}}",
		},

		{
			testName: "Error occurs when Spec.AutomatedPromotion is added.",
			testData: PromotionRun{
				Spec: PromotionRunSpec{
					Snapshot:    "test-snapshot-a",
					Application: "test-app-a-changed",
					AutomatedPromotion: AutomatedPromotionConfiguration{
						InitialEnvironment: "test-env-a",
					},
				},
			},
			expectedError: "spec cannot be updated to {Snapshot:test-snapshot-a Application:test-app-a-changed ManualPromotion:{TargetEnvironment:} AutomatedPromotion:{InitialEnvironment:test-env-a}}",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			warnings, actualError := test.testData.ValidateUpdate(&originalPromotionRun)

			if test.expectedError == "" {
				assert.Nil(t, actualError)
			} else {
				assert.Contains(t, actualError.Error(), test.expectedError)
			}

			if len(test.warnings) > 0 {
				assert.Equal(t, test.warnings, warnings)
			}
		})
	}
}
