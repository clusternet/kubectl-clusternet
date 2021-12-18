/*
Copyright 2021 The Clusternet Authors.

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

package plugin

import (
	"net/http"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/clusternet/clusternet/pkg/wrappers/clientgo"
)

// ClusternetGetter wraps rest.Config for Clusternet.
type ClusternetGetter struct {
	Delegate  genericclioptions.RESTClientGetter
	ClusterID *string
}

var _ genericclioptions.RESTClientGetter = &ClusternetGetter{}

// ToRESTConfig implements RESTClientGetter.
// Returns a REST client configuration based on a provided path
// to a .kubeconfig file, loading rules, and config flag overrides.
// Expects the AddFlags method to have been called.
func (f *ClusternetGetter) ToRESTConfig() (*rest.Config, error) {
	clientConfig, err := f.Delegate.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	err = setKubernetesDefaults(clientConfig)
	if err != nil {
		return nil, err
	}

	// apply Clusternet wrapper
	clientConfig.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		if f.ClusterID != nil && len(*f.ClusterID) > 0 {
			return rt
		}
		return clientgo.NewClusternetTransport(clientConfig.Host, rt)
	})

	return clientConfig, nil
}

func (f *ClusternetGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return f.Delegate.ToRawKubeConfigLoader()
}

func (f *ClusternetGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return f.Delegate.ToDiscoveryClient()
}

// ToRESTMapper returns a mapper.
func (f *ClusternetGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return f.Delegate.ToRESTMapper()
}

func NewClusternetGetter(delegate genericclioptions.RESTClientGetter, clusterID *string) *ClusternetGetter {
	return &ClusternetGetter{
		Delegate:  delegate,
		ClusterID: clusterID,
	}
}

// setKubernetesDefaults sets default values on the provided client config for accessing the
// Kubernetes API or returns an error if any of the defaults are impossible or invalid.
func setKubernetesDefaults(config *rest.Config) error {
	// TODO remove this hack.  This is allowing the GetOptions to be serialized.
	config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}

	if config.APIPath == "" {
		config.APIPath = "/api"
	}
	if config.NegotiatedSerializer == nil {
		// This codec factory ensures the resources are not converted. Therefore, resources
		// will not be round-tripped through internal versions. Defaulting does not happen
		// on the client.
		config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	}
	return rest.SetKubernetesDefaults(config)
}
