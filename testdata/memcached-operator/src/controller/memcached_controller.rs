use crate::api::memcached_types::Memcached;
use crate::controller::{ContextData, Error, Reconciler};
use async_trait::async_trait;
use kube::ResourceExt;
use kube::runtime::controller::Action;
use std::sync::Arc;
use std::time::Duration;

pub struct MemcachedReconciler;

#[async_trait]
impl Reconciler<Memcached> for MemcachedReconciler {
    async fn reconcile(obj: Arc<Memcached>, _ctx: Arc<ContextData>) -> Result<Action, Error> {
        // TODO(user): your logic here
        println!("reconcile request: {}", obj.name_any());
        Ok(Action::requeue(Duration::from_secs(3600)))
    }

    fn error_policy(obj: Arc<Memcached>, err: &Error, _ctx: Arc<ContextData>) -> Action {
        eprintln!("Reconciliation error:\n{:?}.\n{:?}", err, obj);
        Action::requeue(Duration::from_secs(5))
    }
}
