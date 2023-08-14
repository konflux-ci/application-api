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
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var snapshotenvironmentbindinglog = logf.Log.WithName("snapshotenvironmentbinding-resource")

var snapshotEnvironmentBindingClientFromManager client.Client

func (r *SnapshotEnvironmentBinding) SetupWebhookWithManager(mgr ctrl.Manager) error {
	snapshotEnvironmentBindingClientFromManager = mgr.GetClient()

	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-appstudio-redhat-com-v1alpha1-snapshotenvironmentbinding,mutating=true,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=snapshotenvironmentbindings,verbs=create;update,versions=v1alpha1,name=msnapshotenvironmentbinding.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &SnapshotEnvironmentBinding{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SnapshotEnvironmentBinding) Default() {
	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", r.Name).WithValues("namespace", r.Namespace)
	snapshotenvironmentbindinglog.Info("default", "name", r.Name)
}

// change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-appstudio-redhat-com-v1alpha1-snapshotenvironmentbinding,mutating=false,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=snapshotenvironmentbindings,verbs=create;update,versions=v1alpha1,name=vsnapshotenvironmentbinding.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &SnapshotEnvironmentBinding{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *SnapshotEnvironmentBinding) ValidateCreate() error {
	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", r.Name).WithValues("namespace", r.Namespace)
	snapshotenvironmentbindinglog.Info("validating create")

	if err := validateSEB(r); err != nil {
		return fmt.Errorf("invalid SnapshotEnvironmentBinding: %v", err)
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *SnapshotEnvironmentBinding) ValidateUpdate(old runtime.Object) error {
	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", r.Name).WithValues("namespace", r.Namespace)
	snapshotenvironmentbindinglog.Info("validating update")

	switch old := old.(type) {
	case *SnapshotEnvironmentBinding:
		if !reflect.DeepEqual(r.Spec.Application, old.Spec.Application) {
			return fmt.Errorf("application field cannot be updated to %+v", r.Spec.Application)
		}

		if !reflect.DeepEqual(r.Spec.Environment, old.Spec.Environment) {
			return fmt.Errorf("environment field cannot be updated to %+v", r.Spec.Environment)
		}
		if err := validateSEB(r); err != nil {
			return fmt.Errorf("invalid SnapshotEnvironmentBinding: %v", err)
		}

	default:
		return fmt.Errorf("runtime object is not of type SnapshotEnvironmentBinding")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *SnapshotEnvironmentBinding) ValidateDelete() error {
	snapshotenvironmentbindinglog := snapshotenvironmentbindinglog.WithValues("controllerKind", "SnapshotEnvironmentBinding").WithValues("name", r.Name).WithValues("namespace", r.Namespace)
	snapshotenvironmentbindinglog.Info("validating delete")

	return nil
}

func validateSEB(newBinding *SnapshotEnvironmentBinding) error {

	if snapshotEnvironmentBindingClientFromManager == nil {
		return fmt.Errorf("webhook not initialized")
	}

	// Retrieve the list of existing SnapshotEnvironmentBindings from the namespace
	existingSEBs := SnapshotEnvironmentBindingList{}

	if err := snapshotEnvironmentBindingClientFromManager.List(context.Background(), &existingSEBs, &client.ListOptions{Namespace: newBinding.Namespace}); err != nil {
		return fmt.Errorf("failed to list existing SnapshotEnvironmentBindings: %v", err)
	}

	// Check if any existing SEB has the same Application/Environment combination
	for _, existingSEB := range existingSEBs.Items {
		if existingSEB.Spec.Application == newBinding.Spec.Application && existingSEB.Spec.Environment == newBinding.Spec.Environment {
			return fmt.Errorf("duplicate combination of Application (%s) and Environment (%s). Duplicated by: %s", newBinding.Spec.Application, newBinding.Spec.Environment, existingSEB.Name)
		}
	}

	return nil
}
