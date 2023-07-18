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
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var snapshotenvironmentbindinglog = logf.Log.WithName("snapshotenvironmentbinding-resource")

func (r *snapshotEnvironmentBindingWebhookHandler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(&SnapshotEnvironmentBinding{}).
		WithValidator(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-appstudio-redhat-com-v1alpha1-snapshotenvironmentbinding,mutating=true,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=snapshotenvironmentbindings,verbs=create;update,versions=v1alpha1,name=msnapshotenvironmentbinding.kb.io,admissionReviewVersions=v1

var _ webhook.CustomDefaulter = &snapshotEnvironmentBindingWebhookHandler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (h *snapshotEnvironmentBindingWebhookHandler) Default(ctx context.Context, obj runtime.Object) error {
	binding, ok := obj.(*SnapshotEnvironmentBinding)
	if !ok {
		return fmt.Errorf("runtime object is not of type SnapshotEnvironmentBinding")
	}

	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", binding.Name).WithValues("namespace", binding.Namespace)
	snapshotenvironmentbindinglog.Info("default", "name", binding.Name)

	return nil
}

// change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:path=/validate-appstudio-redhat-com-v1alpha1-snapshotenvironmentbinding,mutating=false,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=snapshotenvironmentbindings,verbs=create;update,versions=v1alpha1,name=vsnapshotenvironmentbinding.kb.io,admissionReviewVersions=v1
type snapshotEnvironmentBindingWebhookHandler struct {
	client.Client
}

var _ webhook.CustomValidator = &snapshotEnvironmentBindingWebhookHandler{}

func (h *snapshotEnvironmentBindingWebhookHandler) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	binding, ok := obj.(*SnapshotEnvironmentBinding)
	if !ok {
		return nil, fmt.Errorf("runtime object is not of type SnapshotEnvironmentBinding")
	}

	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", binding.Name).WithValues("namespace", binding.Namespace)
	snapshotenvironmentbindinglog.Info("validating create")

	// Retrieve the list of existing SnapshotEnvironmentBindings from the namespace
	existingSEBs := SnapshotEnvironmentBindingList{}

	if err := h.Client.List(context.Background(), &existingSEBs, &client.ListOptions{Namespace: binding.Namespace}); err != nil {
		return nil, fmt.Errorf("failed to list existing SnapshotEnvironmentBindings: %v", err)
	}

	// Check if any existing SEB has the same Application/Environment combination
	for _, existingSEB := range existingSEBs.Items {
		if existingSEB.Spec.Application == binding.Spec.Application && existingSEB.Spec.Environment == binding.Spec.Environment {
			return nil, fmt.Errorf("duplicate combination of Application (%s) and Environment (%s)", binding.Spec.Application, binding.Spec.Environment)
		}
	}

	return nil, nil
}

func (h *snapshotEnvironmentBindingWebhookHandler) ValidateUpdate(ctx context.Context, obj, oldObj runtime.Object) (admission.Warnings, error) {
	newBinding, ok := obj.(*SnapshotEnvironmentBinding)
	if !ok {
		return nil, fmt.Errorf("runtime object is not of type SnapshotEnvironmentBinding")
	}

	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", newBinding.Name).WithValues("namespace", newBinding.Namespace)
	snapshotenvironmentbindinglog.Info("validating update")

	switch old := oldObj.(type) {
	case *SnapshotEnvironmentBinding:
		if !reflect.DeepEqual(newBinding.Spec.Application, old.Spec.Application) {
			return nil, fmt.Errorf("application field cannot be updated to %+v", newBinding.Spec.Application)
		}

		if !reflect.DeepEqual(newBinding.Spec.Environment, old.Spec.Environment) {
			return nil, fmt.Errorf("environment field cannot be updated to %+v", newBinding.Spec.Environment)
		}

		// Retrieve the list of existing SnapshotEnvironmentBindings from the namespace
		existingSEBs := SnapshotEnvironmentBindingList{}

		if err := h.Client.List(context.Background(), &existingSEBs, &client.ListOptions{Namespace: newBinding.Namespace}); err != nil {
			return nil, fmt.Errorf("failed to list existing SnapshotEnvironmentBindings: %v", err)
		}

		// Check if any existing SEB has the same Application/Environment combination
		for _, existingSEB := range existingSEBs.Items {
			if old.Spec.Application == newBinding.Spec.Application && old.Spec.Environment == newBinding.Spec.Environment && existingSEB.Spec.Application == newBinding.Spec.Application && existingSEB.Spec.Environment == newBinding.Spec.Environment {
				return nil, fmt.Errorf("duplicate combination of Application (%s) and Environment (%s)", newBinding.Spec.Application, newBinding.Spec.Environment)
			}
		}

	default:
		return fmt.Errorf("runtime object is not of type SnapshotEnvironmentBinding")
	}

	return nil, nil

}

func (h *snapshotEnvironmentBindingWebhookHandler) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	binding, ok := obj.(*SnapshotEnvironmentBinding)
	if !ok {
		return nil, fmt.Errorf("runtime object is not of type SnapshotEnvironmentBinding")
	}

	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", binding.Name).WithValues("namespace", binding.Namespace)
	snapshotenvironmentbindinglog.Info("validating delete")
	return nil, nil
}
