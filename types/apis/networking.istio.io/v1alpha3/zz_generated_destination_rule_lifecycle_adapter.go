package v1alpha3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type DestinationRuleLifecycle interface {
	Create(obj *DestinationRule) (runtime.Object, error)
	Remove(obj *DestinationRule) (runtime.Object, error)
	Updated(obj *DestinationRule) (runtime.Object, error)
}

type destinationRuleLifecycleAdapter struct {
	lifecycle DestinationRuleLifecycle
}

func (w *destinationRuleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*DestinationRule))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *destinationRuleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*DestinationRule))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *destinationRuleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*DestinationRule))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDestinationRuleLifecycleAdapter(name string, clusterScoped bool, client DestinationRuleInterface, l DestinationRuleLifecycle) DestinationRuleHandlerFunc {
	adapter := &destinationRuleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *DestinationRule) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
