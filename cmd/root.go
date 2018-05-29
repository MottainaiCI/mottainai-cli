/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>
                         Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	namespace "github.com/MottainaiCI/mottainai-cli/cmd/namespace"
	node "github.com/MottainaiCI/mottainai-cli/cmd/node"
	plan "github.com/MottainaiCI/mottainai-cli/cmd/plan"
	storage "github.com/MottainaiCI/mottainai-cli/cmd/storage"
	task "github.com/MottainaiCI/mottainai-cli/cmd/task"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
)

const (
	cliName = `Mottainai CLI
Copyright (c) 2017-2018 Mottainai

Command line interface for Mottainai clusters`

	cliExamples = `$> mottainai-cli -m http://127.0.0.1:8080 task create --json task.json

$> mottainai-cli -m http://127.0.0.1:8080 namespace list
`
)

var rootCmd = &cobra.Command{
	Short:        cliName,
	Version:      setting.MOTTAINAI_VERSION,
	Example:      cliExamples,
	Args:         cobra.OnlyValidArgs,
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// TODO: Load profile data and override master if present.

		/*
			fmt.Printf("SONO PRE-RUN\n")
			if cmd.Flag("master").Changed {
				fmt.Printf("Master changed.\n")
			}
		*/
	},
}

func init() {

	var pflags = rootCmd.PersistentFlags()

	pflags.StringP("master", "m", "http://localhost:8080", "MottainaiCI webUI URL")
	pflags.StringP("profile", "p", "", "Use specific profile for call API.")

	viper.BindPFlag("master", rootCmd.PersistentFlags().Lookup("master"))
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	rootCmd.AddCommand(
		task.NewTaskCommand(),
		node.NewNodeCommand(),
		namespace.NewNamespaceCommand(),
		plan.NewPlanCommand(),
		storage.NewStorageCommand(),
	)
}

func Execute() {

	// Start command execution
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
