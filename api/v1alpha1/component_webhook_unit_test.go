//
// Copyright 2022-2023 Red Hat, Inc.
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

func TestComponentCreateValidatingWebhook(t *testing.T) {

	tests := []struct {
		name    string
		newComp Component
		err     string
	}{
		{
			name: "component metadata.name is invalid",
			err:  "invalid component name",
			newComp: Component{
				ObjectMeta: v1.ObjectMeta{
					Name: "1-test-component",
				},
				Spec: ComponentSpec{
					ComponentName: "component1",
					Application:   "application1",
				},
			},
		},
		{
			name: "component cannot be created due to bad URL",
			err:  "invalid URI for request" + InvalidSchemeGitSourceURL,
			newComp: Component{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-component",
				},
				Spec: ComponentSpec{
					ComponentName: "component1",
					Application:   "application1",
					Source: ComponentSource{
						ComponentSourceUnion: ComponentSourceUnion{
							GitSource: &GitSource{
								URL: "badurl",
							},
						},
					},
				},
			},
		},
		{
			name: "component needs to have one source specified",
			err:  MissingGitOrImageSource,
			newComp: Component{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-component",
				},
				Spec: ComponentSpec{
					ComponentName: "component1",
					Application:   "application1",
					Source: ComponentSource{
						ComponentSourceUnion: ComponentSourceUnion{
							GitSource: &GitSource{},
						},
					},
				},
			},
		},
		{
			name: "valid component with invalid git vendor src",
			err:  fmt.Errorf(InvalidGithubVendorURL, "http://url", SupportedGitRepo).Error(),
			newComp: Component{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-component",
				},
				Spec: ComponentSpec{
					ComponentName: "component1",
					Application:   "application1",
					Source: ComponentSource{
						ComponentSourceUnion: ComponentSourceUnion{
							GitSource: &GitSource{
								URL: "http://url",
							},
						},
					},
				},
			},
		},
		{
			name: "valid component with invalid git scheme src",
			err:  "invalid URI for request" + InvalidSchemeGitSourceURL,
			newComp: Component{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-component",
				},
				Spec: ComponentSpec{
					ComponentName: "component1",
					Application:   "application1",
					Source: ComponentSource{
						ComponentSourceUnion: ComponentSourceUnion{
							GitSource: &GitSource{
								URL: "git@github.com:devfile-samples/devfile-sample-java-springboot-basic.git",
							},
						},
					},
				},
			},
		},
		{
			name: "valid component with container image",
			newComp: Component{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-component",
				},
				Spec: ComponentSpec{
					ComponentName:  "component1",
					Application:    "application1",
					ContainerImage: "image",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.newComp.ValidateCreate()

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}
		})
	}
}

func TestComponentUpdateValidatingWebhook(t *testing.T) {

	originalComponent := Component{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-component",
		},
		Spec: ComponentSpec{
			ComponentName: "component",
			Application:   "application",
			Source: ComponentSource{
				ComponentSourceUnion: ComponentSourceUnion{
					GitSource: &GitSource{
						URL:     "http://link",
						Context: "context",
					},
				},
			},
		},
	}

	tests := []struct {
		name       string
		updateComp Component
		err        string
	}{
		{
			name: "component name cannot be changed",
			err:  fmt.Errorf(ComponentNameUpdateError, "component1").Error(),
			updateComp: Component{
				Spec: ComponentSpec{
					ComponentName: "component1",
				},
			},
		},
		{
			name: "application name cannot be changed",
			err:  fmt.Errorf(ApplicationNameUpdateError, "application1").Error(),
			updateComp: Component{
				Spec: ComponentSpec{
					ComponentName: "component",
					Application:   "application1",
				},
			},
		},
		{
			name: "git src cannot be changed",
			err: fmt.Errorf(GitSourceUpdateError, GitSource{
				URL:     "http://link1",
				Context: "context",
			}).Error(),
			updateComp: Component{
				Spec: ComponentSpec{
					ComponentName: "component",
					Application:   "application",
					Source: ComponentSource{
						ComponentSourceUnion: ComponentSourceUnion{
							GitSource: &GitSource{
								URL:     "http://link1",
								Context: "context",
							},
						},
					},
				},
			},
		},
		{
			name: "container image can be changed",
			updateComp: Component{
				Spec: ComponentSpec{
					ComponentName:  "component",
					Application:    "application",
					ContainerImage: "image1",
				},
			},
		},
		{
			name: "not component",
			err:  InvalidComponentError,
			updateComp: Component{
				Spec: ComponentSpec{
					ComponentName: "component1",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.err == "" {
				originalComponent = Component{
					Spec: ComponentSpec{
						ComponentName:  "component",
						Application:    "application",
						ContainerImage: "image",
						Source: ComponentSource{
							ComponentSourceUnion: ComponentSourceUnion{
								GitSource: &GitSource{
									Context: "context",
								},
							},
						},
					},
				}
			}
			var err error
			if test.name == "not component" {
				originalApplication := Application{
					Spec: ApplicationSpec{
						DisplayName: "My App",
					},
				}
				_, err = test.updateComp.ValidateUpdate(&originalApplication)
			} else {
				_, err = test.updateComp.ValidateUpdate(&originalComponent)
			}

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}
		})
	}
}

func TestComponentDeleteValidatingWebhook(t *testing.T) {

	tests := []struct {
		name    string
		newComp Component
		err     string
	}{
		{
			name:    "ValidateDelete should return nil, it's unimplemented",
			err:     "",
			newComp: Component{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.newComp.ValidateDelete()

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}
		})
	}
}
