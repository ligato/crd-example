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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/ligato/crd-example/pkg/apis/crdexample.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CrdExampleEmbedLister helps list CrdExampleEmbeds.
type CrdExampleEmbedLister interface {
	// List lists all CrdExampleEmbeds in the indexer.
	List(selector labels.Selector) (ret []*v1.CrdExampleEmbed, err error)
	// CrdExampleEmbeds returns an object that can list and get CrdExampleEmbeds.
	CrdExampleEmbeds(namespace string) CrdExampleEmbedNamespaceLister
	CrdExampleEmbedListerExpansion
}

// crdExampleEmbedLister implements the CrdExampleEmbedLister interface.
type crdExampleEmbedLister struct {
	indexer cache.Indexer
}

// NewCrdExampleEmbedLister returns a new CrdExampleEmbedLister.
func NewCrdExampleEmbedLister(indexer cache.Indexer) CrdExampleEmbedLister {
	return &crdExampleEmbedLister{indexer: indexer}
}

// List lists all CrdExampleEmbeds in the indexer.
func (s *crdExampleEmbedLister) List(selector labels.Selector) (ret []*v1.CrdExampleEmbed, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CrdExampleEmbed))
	})
	return ret, err
}

// CrdExampleEmbeds returns an object that can list and get CrdExampleEmbeds.
func (s *crdExampleEmbedLister) CrdExampleEmbeds(namespace string) CrdExampleEmbedNamespaceLister {
	return crdExampleEmbedNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CrdExampleEmbedNamespaceLister helps list and get CrdExampleEmbeds.
type CrdExampleEmbedNamespaceLister interface {
	// List lists all CrdExampleEmbeds in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.CrdExampleEmbed, err error)
	// Get retrieves the CrdExampleEmbed from the indexer for a given namespace and name.
	Get(name string) (*v1.CrdExampleEmbed, error)
	CrdExampleEmbedNamespaceListerExpansion
}

// crdExampleEmbedNamespaceLister implements the CrdExampleEmbedNamespaceLister
// interface.
type crdExampleEmbedNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CrdExampleEmbeds in the indexer for a given namespace.
func (s crdExampleEmbedNamespaceLister) List(selector labels.Selector) (ret []*v1.CrdExampleEmbed, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CrdExampleEmbed))
	})
	return ret, err
}

// Get retrieves the CrdExampleEmbed from the indexer for a given namespace and name.
func (s crdExampleEmbedNamespaceLister) Get(name string) (*v1.CrdExampleEmbed, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("crdexampleembed"), name)
	}
	return obj.(*v1.CrdExampleEmbed), nil
}
