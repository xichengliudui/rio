package config

import (
	"context"
	"sync"

	"github.com/rancher/norman/pkg/kv"
	"github.com/rancher/types/apis/core/v1"
	"istio.io/api/mesh/v1alpha1"
	metav1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Factory struct {
	sync.Mutex

	configMapNamespace string
	configMapName      string
	configMapKey       string
	template           string
	meshConfig         *v1alpha1.MeshConfig
}

func NewConfigFactory(ctx context.Context, configMap v1.ConfigMapInterface, configMapNamespace, configMapName, configMapKey string) *Factory {
	f := &Factory{
		configMapNamespace: configMapNamespace,
		configMapName:      configMapName,
		configMapKey:       configMapKey,
	}
	configMap.Controller().AddHandler(ctx, "istio-config-cache", f.sync)
	return f
}

func (c *Factory) sync(key string, cm *metav1.ConfigMap) (runtime.Object, error) {
	ns, name := kv.Split(key, "/")
	if ns != c.configMapNamespace && name != c.configMapName {
		return nil, nil
	}

	if cm == nil {
		c.Lock()
		c.template = ""
		c.meshConfig = nil
		c.Unlock()
		return nil, nil
	}

	val, ok := cm.Data[c.configMapKey]
	if !ok {
		return nil, nil
	}

	meshConfig, template, err := DoConfigAndTemplate(val)
	if err != nil {
		return nil, err
	}

	c.Lock()
	c.template = template
	c.meshConfig = meshConfig
	c.Unlock()

	return nil, nil
}

func (c *Factory) TemplateAndConfig() (*v1alpha1.MeshConfig, string) {
	c.Lock()
	defer c.Unlock()

	return c.meshConfig, c.template
}
