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

package webhook

import (
	"fmt"
	"os"
	"strconv"

	tools "github.com/MottainaiCI/mottainai-cli/common"
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	webhook "github.com/MottainaiCI/mottainai-server/pkg/webhook"
	schema "github.com/MottainaiCI/mottainai-server/routes/schema"
	"github.com/MottainaiCI/mottainai-server/routes/schema/v1"
	tablewriter "github.com/olekukonko/tablewriter"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newWebHookListCommand(config *setting.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List webhooks",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var tlist []webhook.WebHook
			var task_table [][]string
			var quiet bool
			var v *viper.Viper = config.Viper

			fetcher := client.NewTokenClient(v.GetString("master"), v.GetString("apikey"), config)

			req := schema.Request{
				Route:  v1.Schema.GetWebHookRoute("show_all"),
				Target: &tlist,
			}
			err = fetcher.Handle(req)
			tools.CheckError(err)

			quiet, err = cmd.Flags().GetBool("quiet")
			tools.CheckError(err)
			all, err := cmd.Flags().GetBool("all")
			tools.CheckError(err)

			if quiet {
				for _, i := range tlist {
					fmt.Println(i.ID)
				}
				return
			}

			for _, i := range tlist {
				if !all {
					task_table = append(task_table, []string{i.ID, i.Key, i.URL, i.Type, i.OwnerId, strconv.FormatBool(i.HasPipeline()), strconv.FormatBool(i.HasTask()), i.Filter, i.Auth})
				} else {
					t, _ := i.ReadTask()
					p, _ := i.ReadPipeline()
					tstr := fmt.Sprintf("%#v", t)
					pstr := fmt.Sprintf("%#v", p)
					task_table = append(task_table, []string{i.ID, i.Key, i.URL, i.Type, i.OwnerId, pstr, tstr, i.Filter, i.Auth})
				}
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetHeader([]string{"ID", "Key", "URL", "Type", "Owner", "Pipeline", "Task", "Filter", "Auth"})

			for _, v := range task_table {
				table.Append(v)
			}
			table.Render()

		},
	}

	var flags = cmd.Flags()
	flags.BoolP("quiet", "q", false, "Quiet Output")
	flags.BoolP("all", "a", false, "Full list")

	return cmd
}
