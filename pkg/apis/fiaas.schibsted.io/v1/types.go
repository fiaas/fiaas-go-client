/*
Copyright 2017-2021 The FIAAS Authors

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

package v1

import (
	"bytes"
	"encoding/gob"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application is a top-level type. A client is created for it.
type Application struct {
	metav1.TypeMeta   `json:",inline"` // apiVersion, kind
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApplicationSpec `json:"spec"`
}

// ApplicationSpec contains data used to create a CRD.
// See https://github.com/fiaas/fiaas-deploy-daemon/blob/master/docs/crd/examples/fiaas-deploy-daemon.yaml
// Note: using an anonymous interface{} type for Config results in badly generated code
type ApplicationSpec struct {
	Application           string                         `json:"application"`
	Image                 string                         `json:"image"`
	Config                Config                         `json:"config"`
	AdditionalLabels      *AdditionalLabelsOrAnnotations `json:"additional_labels,omitempty"`
	AdditionalAnnotations *AdditionalLabelsOrAnnotations `json:"additional_annotations,omitempty"`
}

// Config stores fiaas.yml
// Reference: https://github.com/kubernetes/code-generator/issues/50
type Config map[string]interface{}

type AdditionalLabelsOrAnnotations struct {
	Global                  map[string]string `json:"global,omitempty"`
	Deployment              map[string]string `json:"deployment,omitempty"`
	HorizontalPodAutoscaler map[string]string `json:"horizontal_pod_autoscaler,omitempty"`
	Ingress                 map[string]string `json:"ingress,omitempty"`
	Service                 map[string]string `json:"service,omitempty"`
	ServiceAccount          map[string]string `json:"service_account,omitempty"`
	Pod                     map[string]string `json:"pod,omitempty"`
	Status                  map[string]string `json:"status,omitempty"`
}

// DeepCopyInto is necessary to be able to use a map with an anonymous interface as type for Config
// kudos to https://gist.github.com/soroushjp/0ec92102641ddfc3ad5515ca76405f4d
func (in *Config) DeepCopyInto(out *Config) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	enc.Encode(in)

	dec.Decode(&out)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList is a top-level list type. The client methods for lists are automatically created.
// You are not supposed to create a separated client for this one.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

//
// application-status.fiaas.schibsted.io:
//

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resourceName=application-statuses

// ApplicationStatus is a top-level type. A client is created for it.
type ApplicationStatus struct {
	metav1.TypeMeta   `json:",inline"` // apiVersion, kind
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Logs              []string `json:"logs"`
	Result            string   `json:"result"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationStatusList is a top-level list type. The client methods for lists are automatically created.
// You are not supposed to create a separated client for this one.
type ApplicationStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationStatus `json:"items"`
}
