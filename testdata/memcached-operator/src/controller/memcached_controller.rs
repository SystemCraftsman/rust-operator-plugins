use crate::api::memcached_types::Memcached;
use crate::controller::{ContextData, Error, Reconciler};
use async_trait::async_trait;
use kube::{Api, Resource, ResourceExt};
use kube::runtime::controller::Action;
use std::sync::Arc;
use std::time::Duration;
use k8s_openapi::api::apps::v1::Deployment;
use k8s_openapi::api::core::v1::{Container, PodSpec, PodTemplateSpec};
use k8s_openapi::apimachinery::pkg::apis::meta::v1::{ObjectMeta, OwnerReference, LabelSelector};
use kube::api::PostParams;
use log::info;

pub struct MemcachedReconciler;

#[async_trait]
impl Reconciler<Memcached> for MemcachedReconciler {
    async fn reconcile(memcached: Arc<Memcached>, ctx: Arc<ContextData>) -> Result<Action, Error> {
        let client = &ctx.as_ref().client;
        let namespace = memcached.namespace().unwrap();
        let name = memcached.name_any();
        let deployment_name = format!("{}-deployment", name);
        let deployments: Api<Deployment> = Api::namespaced(client.clone(), &namespace);

        // Try to get the existing Deployment
        match deployments.get(&deployment_name).await {
            Ok(existing_deployment) => {
                let current_replicas = existing_deployment.spec.as_ref().and_then(|s| s.replicas).unwrap_or(1);
                let desired_replicas = memcached.spec.size;

                if current_replicas != desired_replicas {
                    let mut updated = existing_deployment.clone();
                    if let Some(spec) = &mut updated.spec {
                        spec.replicas = Some(desired_replicas);
                    }
                    deployments.replace(&deployment_name, &PostParams::default(), &updated).await?;
                    info!("Updated Deployment '{}' to {} replicas", deployment_name, desired_replicas);
                }
            }
            Err(kube::Error::Api(e)) if e.code == 404 => {
                // Doesn't exist, so create it
                let new_deploy = create_deployment(&memcached)?;
                deployments.create(&PostParams::default(), &new_deploy).await?;
                info!("Created Deployment '{}'", deployment_name);
            }
            Err(e) => return Err(e.into()),
        }

        Ok(Action::requeue(Duration::from_secs(60)))
    }

    fn error_policy(obj: Arc<Memcached>, err: &Error, _ctx: Arc<ContextData>) -> Action {
        eprintln!("Reconciliation error:\n{:?}.\n{:?}", err, obj);
        Action::requeue(Duration::from_secs(5))
    }
}

fn create_deployment(memcached: &Memcached) -> Result<Deployment, Error> {
    let name = memcached.name_any();
    let labels = [("app", name.as_str())]
        .iter()
        .cloned()
        .map(|(k, v)| (k.to_string(), v.to_string()))
        .collect();

    Ok(Deployment {
        metadata: ObjectMeta {
            name: Some(format!("{}-deployment", name)),
            namespace: memcached.namespace(),
            owner_references: Some(vec![OwnerReference {
                api_version: Memcached::api_version(&()).to_string(),
                kind: "Memcached".to_string(),
                name: name.clone(),
                uid: memcached.meta().uid.clone().unwrap_or_default(),
                controller: Some(true),
                block_owner_deletion: Some(true),
            }]),
            ..Default::default()
        },
        spec: Some(k8s_openapi::api::apps::v1::DeploymentSpec {
            replicas: Some(memcached.spec.size),
            selector: LabelSelector {
                match_labels: Some([("app".to_string(), name.clone())].into_iter().collect()),
                ..Default::default()
            },
            template: PodTemplateSpec {
                metadata: Some(ObjectMeta {
                    labels: Some(labels),
                    ..Default::default()
                }),
                spec: Some(PodSpec {
                    containers: vec![Container {
                        name: "memcached".into(),
                        image: Some("memcached:1.4.36-alpine".into()),
                        ports: Some(vec![k8s_openapi::api::core::v1::ContainerPort {
                            container_port: memcached.spec.container_port,
                            ..Default::default()
                        }]),
                        ..Default::default()
                    }],
                    ..Default::default()
                }),
            },
            ..Default::default()
        }),
        ..Default::default()
    })
}
