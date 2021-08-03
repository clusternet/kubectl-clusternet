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
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// ClusternetOptions provides information required to make requests to Clusternet
type ClusternetOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

// NewClusternetOptions provides an instance of ClusternetOptions with default values
func NewClusternetOptions(streams genericclioptions.IOStreams) *ClusternetOptions {
	return &ClusternetOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdClusternet provides a cobra command wrapping ClusternetOptions
func NewCmdClusternet(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewClusternetOptions(streams)

	cmd := &cobra.Command{
		Use:          "clusternet [sub-command] [flags]",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
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

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete fills in fields required to have valid data
func (o *ClusternetOptions) Complete(cmd *cobra.Command, args []string) error {
	// TODO

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ClusternetOptions) Validate() error {
	// TODO

	return nil
}

func (o *ClusternetOptions) Run() error {
	// TODO

	return nil
}
