package src

import (
	"fmt"
	localmachinery "github.com/SystemCraftsman/rust-operator-plugins/pkg/machinery"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/constants"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

const (
	defaultControllerPath = "src/controller.rs"
)

var _ machinery.Template = &Controller{}

type Controller struct {
	machinery.TemplateMixin
}

// SetTemplateDefaults implements file.Template
func (f *Controller) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultControllerPath)
	}

	f.TemplateBody = fmt.Sprintf(controllerTemplate,
		localmachinery.NewMarkerFor(f.Path, constants.ModuleMarker),
	)

	return nil
}

var _ machinery.Inserter = &ControllerUpdater{}

type ControllerUpdater struct { //nolint:maligned
	machinery.ResourceMixin

	// Flags to indicate which parts need to be included when updating the file
	WireResource, WireController, WireWebhook bool

	// Deprecated - The flag should be removed from go/v5
	// IsLegacyPath indicates if webhooks should be scaffolded under the API.
	// Webhooks are now decoupled from APIs based on controller-runtime updates and community feedback.
	// This flag ensures backward compatibility by allowing scaffolding in the legacy/deprecated path.
	IsLegacyPath bool
}

// GetPath implements file.Builder
func (*ControllerUpdater) GetPath() string {
	return defaultMainPath
}

// GetIfExistsAction implements file.Builder
func (*ControllerUpdater) GetIfExistsAction() machinery.IfExistsAction {
	return machinery.OverwriteFile
}

// GetMarkers implements file.Inserter
func (f *ControllerUpdater) GetMarkers() []machinery.Marker {
	return []machinery.Marker{
		machinery.NewMarkerFor(defaultControllerPath, constants.ModuleMarker),
	}
}

const (
	controllerModuleImportCodeFragment = `pub mod %s_controller;
`
)

// GetCodeFragments implements file.Inserter
func (f *ControllerUpdater) GetCodeFragments() machinery.CodeFragmentsMap {
	fragments := make(machinery.CodeFragmentsMap, 3)

	// If resource is not being provided we are creating the file, not updating it
	if f.Resource == nil {
		return fragments
	}

	// Generate module import code fragments
	imports := make([]string, 0)
	if f.WireController {
		imports = append(imports, fmt.Sprintf(controllerModuleImportCodeFragment, f.Resource.Kind))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(imports) != 0 {
		fragments[machinery.NewMarkerFor(defaultApiPath, constants.ModuleMarker)] = imports
	}

	return fragments
}

// nolint:lll
var controllerTemplate = `%s

use async_trait::async_trait;
use futures::stream::StreamExt;
use k8s_openapi::NamespaceResourceScope;
use kube::runtime::controller::Action;
use kube::runtime::Controller;
use kube::{Api, Client, Resource};
use serde::de::DeserializeOwned;
use std::fmt::Debug;
use std::hash::Hash;
use std::marker;
use std::sync::Arc;

#[async_trait]
pub trait Reconciler<K: Resource<Scope = NamespaceResourceScope>> {
    async fn reconcile(obj: Arc<K>, ctx: Arc<ContextData>) -> Result<Action, Error>;
    fn error_policy(obj: Arc<K>, err: &Error, _ctx: Arc<ContextData>) -> Action;
}
pub struct ControllerRunner<K: Resource<Scope = NamespaceResourceScope>> {
    _resource_marker: marker::PhantomData<K>,
}

impl<
        K: Resource<Scope = NamespaceResourceScope>
            + Clone
            + DeserializeOwned
            + Debug
            + Send
            + Sync
            + 'static,
    > ControllerRunner<K>
{
    pub async fn run<T: Reconciler<K>>()
    where
        <K as Resource>::DynamicType: Default,
        <K as Resource>::DynamicType: std::cmp::Eq,
        <K as Resource>::DynamicType: Hash,
        <K as Resource>::DynamicType: Clone,
        <K as kube::Resource>::DynamicType: Debug,
        <K as kube::Resource>::DynamicType: Unpin,
    {
        let client: Client = Client::try_default()
            .await
            .expect("Expected a valid KUBECONFIG environment variable.");
        let context: Arc<ContextData> = Arc::new(ContextData::new(client.clone()));
        let crd_api: Api<K> = Api::all(client);

        Controller::new(crd_api, Default::default())
            .run(<T>::reconcile, <T>::error_policy, context)
            .for_each(|reconciliation_result| async move {
                match reconciliation_result {
                    Ok(resource) => {
                        println!("Reconciliation successful. Resource: {:?}", resource);
                    }
                    Err(reconciliation_err) => {
                        eprintln!("Reconciliation error: {:?}", reconciliation_err)
                    }
                }
            })
            .await;
    }
}

pub struct ContextData {
    client: Client,
}

impl ContextData {
    pub fn new(client: Client) -> Self {
        ContextData { client }
    }
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("Kubernetes reported error: {source}")]
    KubeError {
        #[from]
        source: kube::Error,
    },
    #[error("Invalid Echo CRD: {0}")]
    UserInputError(String),
}
`
