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

package storage

import (
	"log"
	"os"
	"strconv"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	storage "github.com/MottainaiCI/mottainai-server/pkg/storage"
	tablewriter "github.com/olekukonko/tablewriter"
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

func newStorageListCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List storages",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var n []storage.Storage
			var storage_table [][]string

			// TODO: replace this with NewClient(host) method
			var fetcher *client.Fetcher

			fetcher = &client.Fetcher{BaseURL: v.GetString("master")}

			fetcher.GetJSONOptions("/api/storage/list", map[string]string{}, &n)

			log.Println("Available storages: ")

			for _, i := range n {
				storage_table = append(storage_table, []string{strconv.Itoa(i.ID), i.Name, i.Path})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Path"})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			for _, v := range storage_table {
				table.Append(v)
			}
			table.Render()

		},
	}

	return cmd
}
