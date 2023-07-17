//
// Copyright 2023 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnvironmentCreateValidatingWebhook(t *testing.T) {

	badIngressDomain := "AbADIngr3ssDomaiN.CoM"

	orgEnv := Environment{
		ObjectMeta: v1.ObjectMeta{
			Name: "kubernetes-environment",
		},
		Spec: EnvironmentSpec{
			UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
				ClusterType: ConfigurationClusterType_Kubernetes,
				KubernetesClusterCredentials: KubernetesClusterCredentials{
					IngressDomain: "domain",
				},
			},
		},
	}

	tests := []struct {
		name     string
		newEnv   Environment
		err      string
		warnings []string
	}{
		{
			name: "environment ingress domain is empty when its Kubernetes",
			err:  MissingIngressDomain,
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType:                  ConfigurationClusterType_Kubernetes,
						KubernetesClusterCredentials: KubernetesClusterCredentials{},
					},
				},
			},
		},
		{
			name: "environment ingress domain not DNS 1123 compliant",
			err:  fmt.Sprintf(InvalidDNS1123Subdomain, badIngressDomain),
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_Kubernetes,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: badIngressDomain,
						},
					},
				},
			},
		},
		{
			name:   "environment ingress domain is good",
			newEnv: orgEnv,
		},
		{
			name: "environment unstable config is empty",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: nil,
				},
			},
		},
		{
			name: "environment's ingress domain is empty when its OpenShift",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							TargetNamespace: "mynamespace",
						},
					},
				},
			},
		},
		{
			name: "environment's ingress domain is provided when its OpenShift",
			err:  fmt.Sprintf(InvalidDNS1123Subdomain, badIngressDomain),
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: badIngressDomain,
						},
					},
				},
			},
		}, {
			name: "environment name mush have DNS-1123 format  (test 1)",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment-1",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: "domain",
						},
					},
				},
			},
		}, {
			name: "environment name mush have DNS-1123 format (test 2)",
			err:  "invalid environment name: Kubernetes-environment, an environment resource name must start with a lower case alphabetical character, be under 63 characters, and can only consist of lower case alphanumeric characters or ‘-’",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "Kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: "domain",
						},
					},
				},
			},
		}, {
			name: "environment name mush have DNS-1123 format  (test 3)",
			err:  "invalid environment name: kubernetesEnvironment, an environment resource name must start with a lower case alphabetical character, be under 63 characters, and can only consist of lower case alphanumeric characters or ‘-’",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetesEnvironment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: "domain",
						},
					},
				},
			},
		}, {
			name: "environment name mush have DNS-1123 format  (test 4)",
			err:  "invalid environment name: abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde, an environment resource name must start with a lower case alphabetical character, be under 63 characters, and can only consist of lower case alphanumeric characters or ‘-’",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: strings.Repeat("abcde", 13),
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: badIngressDomain,
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			warnings, err := test.newEnv.ValidateCreate()

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}

			if len(test.warnings) > 0 {
				assert.Equal(t, test.warnings, warnings)
			}
		})
	}
}

func TestEnvironmentUpdateValidatingWebhook(t *testing.T) {

	badIngressDomain := "AbADIngr3ssDomaiN.CoM"

	orgEnv := Environment{
		ObjectMeta: v1.ObjectMeta{
			Name: "kubernetes-environment",
		},
		Spec: EnvironmentSpec{
			UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
				ClusterType: ConfigurationClusterType_Kubernetes,
				KubernetesClusterCredentials: KubernetesClusterCredentials{
					IngressDomain: "domain",
				},
			},
		},
	}

	tests := []struct {
		name     string
		newEnv   Environment
		err      string
		warnings []string
	}{
		{
			name: "environment ingress domain is empty when its Kubernetes",
			err:  MissingIngressDomain,
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType:                  ConfigurationClusterType_Kubernetes,
						KubernetesClusterCredentials: KubernetesClusterCredentials{},
					},
				},
			},
		},
		{
			name: "environment ingress domain not DNS 1123 compliant",
			err:  fmt.Sprintf(InvalidDNS1123Subdomain, badIngressDomain),
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_Kubernetes,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: badIngressDomain,
						},
					},
				},
			},
		},
		{
			name:   "environment ingress domain is good",
			newEnv: orgEnv,
		},
		{
			name: "environment unstable config is empty",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: nil,
				},
			},
		},
		{
			name: "environment's ingress domain is empty when its OpenShift",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType:                  ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{},
					},
				},
			},
		},
		{
			name: "environment's ingress domain is provided when its OpenShift",
			newEnv: Environment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kubernetes-environment",
				},
				Spec: EnvironmentSpec{
					UnstableConfigurationFields: &UnstableEnvironmentConfiguration{
						ClusterType: ConfigurationClusterType_OpenShift,
						KubernetesClusterCredentials: KubernetesClusterCredentials{
							IngressDomain: "domain",
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			warnings, err := test.newEnv.ValidateUpdate(&orgEnv)

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}

			if len(test.warnings) > 0 {
				assert.Equal(t, test.warnings, warnings)
			}
		})
	}
}
