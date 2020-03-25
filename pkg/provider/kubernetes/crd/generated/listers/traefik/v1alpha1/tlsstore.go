/*
The MIT License (MIT)

Copyright (c) 2016-2020 Containous SAS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/cgws/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TLSStoreLister helps list TLSStores.
type TLSStoreLister interface {
	// List lists all TLSStores in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.TLSStore, err error)
	// TLSStores returns an object that can list and get TLSStores.
	TLSStores(namespace string) TLSStoreNamespaceLister
	TLSStoreListerExpansion
}

// tLSStoreLister implements the TLSStoreLister interface.
type tLSStoreLister struct {
	indexer cache.Indexer
}

// NewTLSStoreLister returns a new TLSStoreLister.
func NewTLSStoreLister(indexer cache.Indexer) TLSStoreLister {
	return &tLSStoreLister{indexer: indexer}
}

// List lists all TLSStores in the indexer.
func (s *tLSStoreLister) List(selector labels.Selector) (ret []*v1alpha1.TLSStore, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TLSStore))
	})
	return ret, err
}

// TLSStores returns an object that can list and get TLSStores.
func (s *tLSStoreLister) TLSStores(namespace string) TLSStoreNamespaceLister {
	return tLSStoreNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TLSStoreNamespaceLister helps list and get TLSStores.
type TLSStoreNamespaceLister interface {
	// List lists all TLSStores in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.TLSStore, err error)
	// Get retrieves the TLSStore from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.TLSStore, error)
	TLSStoreNamespaceListerExpansion
}

// tLSStoreNamespaceLister implements the TLSStoreNamespaceLister
// interface.
type tLSStoreNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TLSStores in the indexer for a given namespace.
func (s tLSStoreNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.TLSStore, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TLSStore))
	})
	return ret, err
}

// Get retrieves the TLSStore from the indexer for a given namespace and name.
func (s tLSStoreNamespaceLister) Get(name string) (*v1alpha1.TLSStore, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("tlsstore"), name)
	}
	return obj.(*v1alpha1.TLSStore), nil
}
