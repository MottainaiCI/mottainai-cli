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

package plan

import (
	"fmt"
	"os"
	"sort"

	tools "github.com/MottainaiCI/mottainai-cli/common"
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	tablewriter "github.com/olekukonko/tablewriter"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newPlanListCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List plans",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var tlist []citasks.Plan
			var task_table [][]string
			var quiet bool
			var fetcher *client.Fetcher
			var v *viper.Viper = setting.Configuration.Viper

			fetcher = client.NewTokenClient(v.GetString("master"), v.GetString("apikey"))
			fetcher.GetJSONOptions("/api/tasks/planned", map[string]string{}, &tlist)

			sort.Slice(tlist[:], func(i, j int) bool {
				return tlist[i].CreatedTime > tlist[j].CreatedTime
			})

			quiet, err = cmd.Flags().GetBool("quiet")
			tools.CheckError(err)

			if quiet {
				for _, i := range tlist {
					fmt.Println(i.ID)
				}
				return
			}

			for _, i := range tlist {
				task_table = append(task_table, []string{i.ID, i.Planned, i.Directory})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetHeader([]string{"ID", "Planned", "Dir"})

			for _, v := range task_table {
				table.Append(v)
			}
			table.Render()

		},
	}

	var flags = cmd.Flags()
	flags.BoolP("quiet", "q", false, "Quiet Output")

	return cmd
}
