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

package version

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmdVersion returns a cobra command for fetching versions
func NewCmdVersion(ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the plugin version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			versionInfo := Info{
				GitVersion: gitVersion,
				GitCommit:  gitCommit,
				BuildDate:  buildDate,
				Platform:   platform,
			}

			marshalled, err := json.MarshalIndent(&versionInfo, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(ioStreams.Out, string(marshalled))
			return nil
		},
	}

	cmd.ResetFlags()
	return cmd
}
