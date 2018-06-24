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

package node

import (
	"os"
	"strconv"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	nodes "github.com/MottainaiCI/mottainai-server/pkg/nodes"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	tablewriter "github.com/olekukonko/tablewriter"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newNodeListCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List nodes",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var n []nodes.Node
			var node_table [][]string
			var fetcher *client.Fetcher
			var v *viper.Viper = setting.Configuration.Viper

			fetcher = client.NewClient(v.GetString("master"))
			fetcher.GetJSONOptions("/api/nodes", map[string]string{}, &n)

			for _, i := range n {
				node_table = append(node_table, []string{strconv.Itoa(i.ID), i.Hostname, i.User, i.Pass, i.Key, i.NodeID})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetHeader([]string{"ID", "Hostname", "User", "Pass", "Key", "UUID"})

			for _, v := range node_table {
				table.Append(v)
			}
			table.Render()
		},
	}

	return cmd
}
