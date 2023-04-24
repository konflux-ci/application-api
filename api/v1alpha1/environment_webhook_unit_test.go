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
		name   string
		newEnv Environment
		err    string
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.newEnv.ValidateCreate()

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
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
		name   string
		newEnv Environment
		err    string
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.newEnv.ValidateUpdate(&orgEnv)

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}
		})
	}
}
