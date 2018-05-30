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

package namespace

import (
	"os"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	tablewriter "github.com/olekukonko/tablewriter"
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

func newNamespaceListCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List namespaces",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var tlist []string
			var ns_table [][]string
			var fetcher *client.Fetcher

			fetcher = client.NewClient(v.GetString("master"))

			fetcher.GetJSONOptions("/api/namespace/list", map[string]string{}, &tlist)

			for _, i := range tlist {
				ns_table = append(ns_table, []string{i})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name"})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			for _, v := range ns_table {
				table.Append(v)
			}
			table.Render()
		},
	}

	return cmd
}
