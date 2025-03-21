use k8s_openapi::serde::{Deserialize, Serialize};
use kube::CustomResource;
use schemars::JsonSchema;

#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, JsonSchema)]
#[kube(kind = "Memcached", group = "cache", version = "v1alpha1", namespaced)]
#[kube(status = "MemcachedStatus")]
pub struct MemcachedSpec {
    // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster

    // foo is an example field of Memcached. Edit memcached_types.rs to remove/update
    foo: String,
}

#[derive(Deserialize, Serialize, Clone, Debug, JsonSchema)]
pub struct MemcachedStatus {
    // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
}
