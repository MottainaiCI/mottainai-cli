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

package task

import (
	"fmt"
	"log"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

func newTaskInspectCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "inspect <taskid> [OPTIONS]",
		Short: "Inspect a task for debugging",
		Args:  cobra.RangeArgs(1, 1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: replace this with NewClient(host) method from mottainai-server
			//var fetcher *client.Fetcher
			//fetcher = &client.Fetcher{BaseURL: v.GetString("master")}

			id := args[0]
			if len(id) == 0 {
				log.Fatalln("You need to define a task id")
			}

			fetcher := client.NewFetcher(id)
			fetcher.BaseURL = v.GetString("master")

			//fetcher.Doc(id)

			th := citasks.DefaultTaskHandler()

			// TODO: Fix this
			// github.com/MottainaiCI/mottainai-cli/cmd/task
			// cmd/task/inspect.go:51:29: cannot use fetcher
			// (type *"github.com/MottainaiCI/mottainai-cli/vendor/github.com/MottainaiCI/mottainai-server/pkg/client".Fetcher)
			// as type *"github.com/MottainaiCI/mottainai-server/pkg/client".Fetcher in argument to th.FetchTask

			task_info := th.FetchTask(fetcher)
			fmt.Println(task_info)
		},
	}

	return cmd
}
