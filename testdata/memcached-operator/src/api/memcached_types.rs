use k8s_openapi::serde::{Deserialize, Serialize};
use kube::CustomResource;
use schemars::JsonSchema;

#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, JsonSchema)]
#[kube(
    kind = "Memcached",
    group = "cache.example.com",
    version = "v1alpha1",
    namespaced,
    status = "MemcachedStatus"
)]
pub struct MemcachedSpec {
    pub size: i32,

    #[serde(rename = "containerPort")]
    pub container_port: i32,
}

#[derive(Deserialize, Serialize, Clone, Debug, JsonSchema)]
pub struct MemcachedStatus {
    pub conditions: Vec<Condition>,
}

#[derive(Clone, Debug, Deserialize, Serialize, JsonSchema, Default)]
pub struct Condition {
    pub type_: String,
    pub status: String,
    pub reason: Option<String>,
    pub message: Option<String>,
}
