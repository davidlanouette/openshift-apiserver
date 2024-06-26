package buildclone

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"

	buildapi "github.com/openshift/openshift-apiserver/pkg/build/apis/build"
	"github.com/openshift/openshift-apiserver/pkg/build/apiserver/buildgenerator"
)

// NewStorage creates a new storage object for build generation
func NewStorage(generator *buildgenerator.BuildGenerator) *CloneREST {
	return &CloneREST{generator: generator}
}

// CloneREST is a RESTStorage implementation for a BuildGenerator which supports only
// the Get operation (as the generator has no underlying storage object).
type CloneREST struct {
	generator *buildgenerator.BuildGenerator
}

var _ rest.Creater = &CloneREST{}
var _ rest.Storage = &CloneREST{}

// New creates a new build clone request
func (s *CloneREST) New() runtime.Object {
	return &buildapi.BuildRequest{}
}

func (s *CloneREST) Destroy() {}

// Create instantiates a new build from an existing build
func (s *CloneREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	objectMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	rest.FillObjectMetaSystemFields(objectMeta)

	if err := rest.BeforeCreate(Strategy, ctx, obj); err != nil {
		return nil, err
	}
	if err := createValidation(ctx, obj); err != nil {
		return nil, err
	}

	return s.generator.CloneInternal(ctx, obj.(*buildapi.BuildRequest))
}
