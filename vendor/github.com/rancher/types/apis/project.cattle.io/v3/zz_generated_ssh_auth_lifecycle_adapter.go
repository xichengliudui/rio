package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type SSHAuthLifecycle interface {
	Create(obj *SSHAuth) (runtime.Object, error)
	Remove(obj *SSHAuth) (runtime.Object, error)
	Updated(obj *SSHAuth) (runtime.Object, error)
}

type sshAuthLifecycleAdapter struct {
	lifecycle SSHAuthLifecycle
}

func (w *sshAuthLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*SSHAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *sshAuthLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*SSHAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *sshAuthLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*SSHAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewSSHAuthLifecycleAdapter(name string, clusterScoped bool, client SSHAuthInterface, l SSHAuthLifecycle) SSHAuthHandlerFunc {
	adapter := &sshAuthLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *SSHAuth) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
