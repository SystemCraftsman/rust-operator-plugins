/*
Copyright 2025 System Craftsman LLC

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

package controller

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

var _ machinery.Template = &Controllers{}

// Controllers scaffolds the file that defines the controller for a CRD or a builtin resource
// nolint:maligned
type Controllers struct {
	machinery.TemplateMixin
	machinery.ResourceMixin
	machinery.BoilerplateMixin

	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *Controllers) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("src", "controller", "%[kind]_controller.rs")
	}

	f.Path = f.Resource.Replacer().Replace(f.Path)
	log.Println(f.Path)

	f.TemplateBody = controllerTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		f.IfExistsAction = machinery.Error
	}

	return nil
}

//nolint:lll
const controllerTemplate = `{{ .Boilerplate }}

use crate::api::{{ lower .Resource.Kind }}_types::{{ .Resource.Kind }};
use crate::controller::{ContextData, Error, Reconciler};
use async_trait::async_trait;
use kube::runtime::controller::Action;
use kube::ResourceExt;
use std::sync::Arc;
use std::time::Duration;


pub struct {{ .Resource.Kind }}Reconciler;

#[async_trait]
impl Reconciler<{{ .Resource.Kind }}> for {{ .Resource.Kind }}Reconciler {
    async fn reconcile(obj: Arc<{{ .Resource.Kind }}>, _ctx: Arc<ContextData>) -> Result<Action, Error> {
        // TODO(user): your logic here
		println!("reconcile request: {}", obj.name_any());
        Ok(Action::requeue(Duration::from_secs(60)))
    }

    fn error_policy(obj: Arc<{{ .Resource.Kind }}>, err: &Error, _ctx: Arc<ContextData>) -> Action {
		eprintln!("Reconciliation error:\n{:?}.\n{:?}", err, obj);
        Action::requeue(Duration::from_secs(5))
    }
}
`
