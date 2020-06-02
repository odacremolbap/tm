// Copyright 2020 TriggerMesh, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/triggermesh/tm/pkg/client"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/client/pkg/util"
)

// CfgFile is tm's config file is the path for the Kubernetes config
var CfgFile string

// Config contains the variables for the Kn config
type Config struct {
	DefaultConfigDir string
}

// TmParams for creating commands. Useful for inserting mocks for testing.
type TmParams struct {
	// Output            io.Writer
	KubeCfgPath  string
	ClientConfig clientcmd.ClientConfig
	// NewServingClient  func(namespace string) (clientservingv1.KnServingClient, error)
	// NewSourcesClient  func(namespace string) (v1alpha2.KnSourcesClient, error)
	// NewEventingClient func(namespace string) (clienteventingv1alpha1.KnEventingClient, error)
	// NewDynamicClient  func(namespace string) (clientdynamic.KnDynamicClient, error)

	ConfigSet *ConfigSet

	// // General global options
	// LogHTTP bool

	// // Set this if you want to nail down the namespace
	fixedCurrentNamespace string
}

// Initialize clients
func (p *TmParams) Initialize() {

	confPath := ConfigPath(p.KubeCfgPath)
	configset, err := NewClient(confPath, tmCmd.OutOrStdout())
	err != nil
	if err != nil {
		log.Fatalln(err)
	}

	p.ConfigSet
	clientset.Printer.Format = client.Output
	if debug {
		clientset.Log.SetDebugLevel()
	}
	clientset.Registry.Host = registryHost
	clientset.Registry.Secret = registrySecret
	clientset.Registry.SkipTLS = registrySkipTLS

	// if params.NewServingClient == nil {
	// 	params.NewServingClient = params.newServingClient
	// }

	// if params.NewEventingClient == nil {
	// 	params.NewEventingClient = params.newEventingClient
	// }
}

// func (params *TmParams) newServingClient(namespace string) (clientservingv1.KnServingClient, error) {
// 	restConfig, err := params.RestConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	client, _ := servingv1client.NewForConfig(restConfig)
// 	return clientservingv1.NewKnServingClient(client, namespace), nil
// }

// func (params *TmParams) newEventingClient(namespace string) (clienteventingv1alpha1.KnEventingClient, error) {
// 	restConfig, err := params.RestConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	client, _ := eventing.NewForConfig(restConfig)
// 	return clienteventingv1alpha1.NewKnEventingClient(client, namespace), nil
// }

// RestConfig returns REST config, which can be to use to create specific clientset
func (params *TmParams) RestConfig() (*rest.Config, error) {
	var err error

	if params.ClientConfig == nil {
		params.ClientConfig, err = params.GetClientConfig()
		if err != nil {
			return nil, err
		}
	}

	config, err := params.ClientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	if params.LogHTTP {
		// TODO: When we update to the newer version of client-go, replace with
		// config.Wrap() for future compat.
		config.WrapTransport = util.NewLoggingTransport
	}

	return config, nil
}

// GetClientConfig gets ClientConfig from KubeCfgPath
func (params *TmParams) GetClientConfig() (clientcmd.ClientConfig, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.Precedence = append([]string{newDefaultConfigPath("config.json")}, loadingRules.Precedence...)

	// loadingRules.Precedence = []string{newDefaultConfigPath("config.json")}
	fmt.Printf("loading rules %+v\n", loadingRules)
	if len(params.KubeCfgPath) == 0 {
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}), nil
	}

	_, err := os.Stat(params.KubeCfgPath)
	if err == nil {
		loadingRules.ExplicitPath = params.KubeCfgPath
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}), nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	paths := filepath.SplitList(params.KubeCfgPath)
	if len(paths) > 1 {
		return nil, fmt.Errorf("Can not find config file. '%s' looks like a path. "+
			"Please use the env var KUBECONFIG if you want to check for multiple configuration files", params.KubeCfgPath)
	}
	return nil, fmt.Errorf("Config file '%s' can not be found", params.KubeCfgPath)
}

// Returns a pointer to bool, hard to do better in Golang
func newBoolP(b bool) *bool {
	aBool := b
	return &aBool
}

// Returns default config path based on target OS
func newDefaultConfigPath(subDir string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "tm", subDir)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".tm", subDir)
}
