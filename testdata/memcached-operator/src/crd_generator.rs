/*
Copyright 2025.

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

mod api;

use k8s_openapi::apiextensions_apiserver::pkg::apis::apiextensions::v1::CustomResourceDefinition;
use kube::CustomResourceExt;
use std::fs;
use std::fs::File;

fn main() {
    fs::create_dir_all("target/kubernetes").expect("Error creating directory 'target/kubernetes'");
    write_crd_to_yaml(&api::memcached_types::Memcached::crd());
    // +kubebuilder:scaffold:writers
}

fn write_crd_to_yaml(crd: &CustomResourceDefinition) {
    let file_path = format!(
        "target/kubernetes/{name}-{version}.yaml",
        name = crd.metadata.name.clone().unwrap(),
        version = crd.spec.versions.first().unwrap().name
    );
    let file = File::create(file_path).expect("Error creating YAML file");
    serde_yaml::to_writer(file, crd).expect("Error writing to YAML file");
}
