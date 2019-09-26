// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ApplicationStatusLister helps list ApplicationStatuses.
type ApplicationStatusLister interface {
	// List lists all ApplicationStatuses in the indexer.
	List(selector labels.Selector) (ret []*v1.ApplicationStatus, err error)
	// ApplicationStatuses returns an object that can list and get ApplicationStatuses.
	ApplicationStatuses(namespace string) ApplicationStatusNamespaceLister
	ApplicationStatusListerExpansion
}

// applicationStatusLister implements the ApplicationStatusLister interface.
type applicationStatusLister struct {
	indexer cache.Indexer
}

// NewApplicationStatusLister returns a new ApplicationStatusLister.
func NewApplicationStatusLister(indexer cache.Indexer) ApplicationStatusLister {
	return &applicationStatusLister{indexer: indexer}
}

// List lists all ApplicationStatuses in the indexer.
func (s *applicationStatusLister) List(selector labels.Selector) (ret []*v1.ApplicationStatus, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ApplicationStatus))
	})
	return ret, err
}

// ApplicationStatuses returns an object that can list and get ApplicationStatuses.
func (s *applicationStatusLister) ApplicationStatuses(namespace string) ApplicationStatusNamespaceLister {
	return applicationStatusNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ApplicationStatusNamespaceLister helps list and get ApplicationStatuses.
type ApplicationStatusNamespaceLister interface {
	// List lists all ApplicationStatuses in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ApplicationStatus, err error)
	// Get retrieves the ApplicationStatus from the indexer for a given namespace and name.
	Get(name string) (*v1.ApplicationStatus, error)
	ApplicationStatusNamespaceListerExpansion
}

// applicationStatusNamespaceLister implements the ApplicationStatusNamespaceLister
// interface.
type applicationStatusNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ApplicationStatuses in the indexer for a given namespace.
func (s applicationStatusNamespaceLister) List(selector labels.Selector) (ret []*v1.ApplicationStatus, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ApplicationStatus))
	})
	return ret, err
}

// Get retrieves the ApplicationStatus from the indexer for a given namespace and name.
func (s applicationStatusNamespaceLister) Get(name string) (*v1.ApplicationStatus, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("applicationstatus"), name)
	}
	return obj.(*v1.ApplicationStatus), nil
}