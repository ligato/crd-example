// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"reflect"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ligato/crd-example/pkg/crdexample"
)

// CRD Constants
const (
	Group                   string = "crdexample.io"
	GroupVersion            string = "v1"
	CRDSpecPath             string = "github.com/ligato/crd-example/pkg/apis" + Group
	CRDExampleSingular      string = "crdexample"
	CRDExamplePlural        string = CRDExampleSingular + "s"
	CRDExampleEmbedSingular string = "crdexampleembed"
	CRDExampleEmbedPlural   string = CRDExampleEmbedSingular + "s"
)

var (
	CRDExampleTypeName      = reflect.TypeOf(CrdExample{}).Name()
	FullCRDExampleName      = CRDSpecPath + "/" + GroupVersion + "." + CRDExampleTypeName
	CRDExampleEmbedTypeName = reflect.TypeOf(CrdExampleEmbed{}).Name()
	FullCRDExampleEmbedName = CRDSpecPath + "/" + GroupVersion + "." + CRDExampleEmbedTypeName
)

// CrdExample CRD
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type CrdExample struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            crdexample.CrdExample `json:"spec"`
	Status          CrdExampleStatus      `json:"status,omitempty"`
}

// CrdExampleStatus is the status schema for this CRD
type CrdExampleStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// CrdExampleList is the list schema for this CRD
// -genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type CrdExampleList struct {
	meta.TypeMeta `json:",inline"`
	// +optional
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []CrdExample `json:"items"`
}

// CrdExampleEmbed CRD
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type CrdExampleEmbed struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            crdexample.CrdExample_CrdExampleEmbed `json:"spec"`
	Status          CrdExampleEmbedStatus                 `json:"status,omitempty"`
}

// CrdExampleEmbedStatus is the status schema for this CRD
type CrdExampleEmbedStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// CrdExampleEmbedList is the list schema for this CRD
// -genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type CrdExampleEmbedList struct {
	meta.TypeMeta `json:",inline"`
	// +optional
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []CrdExampleEmbed `json:"items"`
}
