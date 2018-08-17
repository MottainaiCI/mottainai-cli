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

package settingcmd

import (
	"fmt"
	"os"

	tools "github.com/MottainaiCI/mottainai-cli/common"
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	tablewriter "github.com/olekukonko/tablewriter"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newSettingListCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List settings",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var tlist []setting.Setting
			var setting_table [][]string
			var quiet bool
			var fetcher *client.Fetcher
			var v *viper.Viper = setting.Configuration.Viper

			fetcher = client.NewTokenClient(v.GetString("master"), v.GetString("apikey"))
			fetcher.GetJSONOptions("/api/settings", map[string]string{}, &tlist)

			quiet, err = cmd.Flags().GetBool("quiet")
			tools.CheckError(err)

			if quiet {
				for _, i := range tlist {
					fmt.Println(i.Key, i.Value)
				}
				return
			}

			for _, i := range tlist {
				setting_table = append(setting_table, []string{i.ID, i.Key, i.Value})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetHeader([]string{"ID", "Key", "Value"})

			for _, v := range setting_table {
				table.Append(v)
			}
			table.Render()

		},
	}

	var flags = cmd.Flags()
	flags.BoolP("quiet", "q", false, "Quiet Output")

	return cmd
}
