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
	"log"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

func newNamespaceUploadCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "upload <namespace> <file> <storage-path> [OPTIONS]",
		Short: "Upload file to a namespace",
		Args:  cobra.RangeArgs(3, 3),
		Run: func(cmd *cobra.Command, args []string) {
			var fetcher *client.Fetcher

			storage := args[0]
			file := args[1]
			path := args[2]
			if len(storage) == 0 || len(file) == 0 || len(path) == 0 {
				log.Fatalln("You need to define a storage id, a file and a target storage path.")
			}

			fetcher = client.NewClient(v.GetString("master"))
			fetcher.UploadNamespaceFile(storage, file, path)
		},
	}

	return cmd
}