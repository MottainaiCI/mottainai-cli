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
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

func newTaskArtefactsCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "artefacts <taskid> [OPTIONS]",
		Short: "Show artefacts of a task",
		Args:  cobra.RangeArgs(1, 1),
		Run: func(cmd *cobra.Command, args []string) {
			var fetcher *client.Fetcher
			var tlist []string

			id := args[0]
			if len(id) == 0 {
				log.Fatalln("You need to define a task id")
			}

			fmt.Println("Artefacts for:", id)
			fetcher = client.NewClient(v.GetString("master"))
			fetcher.GetJSONOptions("/api/tasks/"+id+"/artefacts", map[string]string{}, &tlist)

			for _, i := range tlist {
				fmt.Println("- " + i)
			}
		},
	}

	return cmd
}
