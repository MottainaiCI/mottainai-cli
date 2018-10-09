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
	"strconv"
	"time"

	tools "github.com/MottainaiCI/mottainai-cli/common"
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newTaskAttachCommand(config *setting.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "attach <taskid> [OPTIONS]",
		Short: "attach to a task output",
		Args:  cobra.RangeArgs(1, 1),
		Run: func(cmd *cobra.Command, args []string) {
			var fetcher *client.Fetcher
			var v *viper.Viper = config.Viper

			fetcher = client.NewTokenClient(v.GetString("master"), v.GetString("apikey"), config)

			id := args[0]
			if len(id) == 0 {
				log.Fatalln("You need to define a task id")
			}
			var pos = 0

			for {
				time.Sleep(time.Second + 2)
				var t citasks.Task
				fetcher.GetJSONOptions("/api/tasks/"+id, map[string]string{}, &t)
				if t.Status != "running" {
					if t.Status == "done" && pos == 0 {
						buff, err := fetcher.GetOptions("/artefact/"+id+"/build_"+t.ID+".log", map[string]string{})
						if err != nil {
							panic(err)
						}
						tools.PrintBuff(buff)
					} else {
						fmt.Println("Build status: " + t.Status + " Can't attach to any live stream.")
					}
					return
				}

				buff, err := fetcher.GetOptions("/api/tasks/stream_output/"+id+"/"+strconv.Itoa(pos), map[string]string{})
				if err != nil {
					panic(err)
				}
				pos += len(buff)
				tools.PrintBuff(buff)
			}

		},
	}

	return cmd
}
