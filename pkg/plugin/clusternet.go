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
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/kubectl/pkg/cmd/annotate"
	"k8s.io/kubectl/pkg/cmd/apiresources"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/create"
	"k8s.io/kubectl/pkg/cmd/delete"
	"k8s.io/kubectl/pkg/cmd/edit"
	cmdexec "k8s.io/kubectl/pkg/cmd/exec"
	"k8s.io/kubectl/pkg/cmd/get"
	"k8s.io/kubectl/pkg/cmd/label"
	"k8s.io/kubectl/pkg/cmd/logs"
	"k8s.io/kubectl/pkg/cmd/portforward"
	"k8s.io/kubectl/pkg/cmd/scale"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	proxiesapi "github.com/clusternet/clusternet/pkg/apis/proxies/v1alpha1"
	"github.com/clusternet/kubectl-clusternet/pkg/version"
)

const (
	headerPrefix = "Clusternet-"
)

var (
	kubectlClusternet = "kubectl clusternet"
)

// ClusternetOptions provides information required to make requests to Clusternet
type ClusternetOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	clusterID       string
	childKubeConfig string

	// for child cluster impersonation
	as      string
	asUID   string
	asGroup []string
}

// NewClusternetOptions provides an instance of ClusternetOptions with default values
func NewClusternetOptions(streams genericclioptions.IOStreams) *ClusternetOptions {
	o := &ClusternetOptions{
		configFlags: genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag(),
		IOStreams:   streams,
	}
	o.configFlags.WrapConfigFn = o.WrapConfigFn
	return o
}

// Complete fills in fields required to have valid data
func (o *ClusternetOptions) Complete() error {
	// TODO

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ClusternetOptions) Validate() error {
	if len(o.clusterID) != 0 && len(o.childKubeConfig) == 0 {
		return fmt.Errorf("please specify a valid kubeconfig for child cluster through '--child-kubeconfig'")
	}

	if len(o.clusterID) == 0 && len(o.childKubeConfig) != 0 {
		return fmt.Errorf("please specify a valid cluster UUID with '--cluster-id'")
	}

	return nil
}

func (o *ClusternetOptions) Run() error {
	// TODO

	return nil
}

func (o *ClusternetOptions) WrapConfigFn(config *rest.Config) *rest.Config {
	if len(o.childKubeConfig) > 0 {
		var childConfig *rest.Config
		clientConfig, err := clientcmd.LoadFromFile(o.childKubeConfig)
		if err == nil {
			childConfig, err = clientcmd.NewDefaultClientConfig(*clientConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
		}
		if err != nil {
			panic(fmt.Sprintf("error while loading kubeconfig from file %s: %v", o.childKubeConfig, err))
		}

		if config.Impersonate.Extra == nil {
			config.Impersonate.Extra = make(map[string][]string)
		}
		config.Impersonate.UserName = "clusternet"

		if len(o.as) > 0 {
			config.Impersonate.Extra[fmt.Sprintf("%s%s", headerPrefix, transport.ImpersonateUserHeader)] = []string{o.as}
		}
		if len(o.asUID) > 0 {
			config.Impersonate.Extra[fmt.Sprintf("%s%s", headerPrefix, transport.ImpersonateUIDHeader)] = []string{o.asUID}
		}
		if len(o.asGroup) > 0 {
			config.Impersonate.Extra[fmt.Sprintf("%s%s", headerPrefix, transport.ImpersonateGroupHeader)] = o.asGroup
		}

		if len(childConfig.BearerToken) > 0 {
			config.Impersonate.Extra[fmt.Sprintf("%sToken", headerPrefix)] = []string{childConfig.BearerToken}
		}
		if len(childConfig.CertData) > 0 {
			config.Impersonate.Extra[fmt.Sprintf("%sCertificate", headerPrefix)] = []string{base64.StdEncoding.EncodeToString(childConfig.CertData)}
		}
		if len(childConfig.KeyData) > 0 {
			config.Impersonate.Extra[fmt.Sprintf("%sPrivateKey", headerPrefix)] = []string{base64.StdEncoding.EncodeToString(childConfig.KeyData)}
		}

		config.Host = strings.Join([]string{
			strings.TrimRight(config.Host, "/"),
			fmt.Sprintf("apis/%s/sockets/%s/proxy/direct", proxiesapi.SchemeGroupVersion.String(), o.clusterID),
		}, "/")
	}

	return config
}

func (o *ClusternetOptions) AddFlags(pfs *flag.FlagSet) {
	o.configFlags.AddFlags(pfs)

	pfs.StringVar(&o.clusterID, "cluster-id", o.clusterID,
		"[Clusternet] The child/member cluster UUID. Only works with '--child-kubeconfig'.")
	pfs.StringVar(&o.childKubeConfig, "child-kubeconfig", o.childKubeConfig,
		"[Clusternet] Path to the kubeconfig file for a child/member cluster. The apiserver url could be an inner address."+
			" Only works with '--cluster-id'.")
	pfs.StringVar(&o.as, "clusternet-as", o.as,
		"[Clusternet] Username to impersonate for the operation to child cluster. User could be a regular user "+
			"or a service account in a namespace. Only works with '--child-kubeconfig'.")
	pfs.StringVar(&o.asUID, "clusternet-as-uid", o.asUID,
		"[Clusternet] UID to impersonate for the operation to child cluster, this flag can be repeated to "+
			"specify multiple groups. Only works with '--child-kubeconfig'.")
	pfs.StringArrayVar(&o.asGroup, "clusternet-as-group", o.asGroup,
		"[Clusternet] Group to impersonate for the operation to child cluster, this flag can be repeated to "+
			"specify multiple groups. Only works with '--child-kubeconfig'.")
}

// NewCmdClusternet provides a cobra command wrapping ClusternetOptions
func NewCmdClusternet(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewClusternetOptions(streams)

	cmd := &cobra.Command{
		Use:          "clusternet",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	o.AddFlags(cmd.PersistentFlags())

	f := cmdutil.NewFactory(NewClusternetGetter(o.configFlags, &o.clusterID))

	// add subcommands
	cmd.AddCommand(apiresources.NewCmdAPIResources(f, streams))
	cmd.AddCommand(create.NewCmdCreate(f, streams))
	cmd.AddCommand(get.NewCmdGet(kubectlClusternet, f, streams))
	cmd.AddCommand(apply.NewCmdApply(kubectlClusternet, f, streams))
	cmd.AddCommand(delete.NewCmdDelete(f, streams))
	cmd.AddCommand(scale.NewCmdScale(f, streams))
	cmd.AddCommand(edit.NewCmdEdit(f, streams))
	cmd.AddCommand(label.NewCmdLabel(f, streams))
	cmd.AddCommand(annotate.NewCmdAnnotate(kubectlClusternet, f, streams))
	cmd.AddCommand(cmdexec.NewCmdExec(f, streams))
	cmd.AddCommand(logs.NewCmdLogs(f, o.IOStreams))
	cmd.AddCommand(portforward.NewCmdPortForward(f, o.IOStreams))

	// add subcommand version
	cmd.AddCommand(version.NewCmdVersion(streams))

	// replace "kubectl" to "kubectl cluster" in example
	for _, command := range cmd.Commands() {
		command.Example = strings.Replace(command.Example, "kubectl", kubectlClusternet, -1)
	}
	return cmd
}
