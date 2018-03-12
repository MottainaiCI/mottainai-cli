/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>

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

package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	storage "github.com/MottainaiCI/mottainai-server/pkg/storage"
	"github.com/olekukonko/tablewriter"

	"github.com/urfave/cli"
)

var Storage = cli.Command{
	Name:        "storage",
	Usage:       "create, delete, tag, show, list, download, upload, remove",
	Description: `Create, delete, tag, show and list storagess, also download/upload file`,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new storage",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				storage := c.Args().First()
				if len(storage) == 0 {
					log.Fatalln("You need to define a storage name")
				}
				res, err := fetcher.GetOptions("/api/storage/"+storage+"/create", map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "delete a namespace",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				storage := c.Args().First()
				if len(storage) == 0 {
					log.Fatalln("You need to define a storage id")
				}
				res, err := fetcher.GetOptions("/api/storage/"+storage+"/delete", map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "show",
			Usage: "show artefacts belonging to a storage",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				storage := c.Args().First()
				if len(storage) == 0 {
					log.Fatalln("You need to define a storage id")
				}
				fmt.Println("Storage: ", storage)
				var tlist []string

				fetcher.GetJSONOptions("/api/storage/"+storage+"/list", map[string]string{}, &tlist)

				for _, i := range tlist {
					fmt.Println("- " + i)
				}
				return nil
			},
		},

		{
			Name:  "list",
			Usage: "list available storages",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)

				log.Println("Available storages: ")

				var n []storage.Storage
				fetcher.GetJSONOptions("/api/storage/list", map[string]string{}, &n)

				var storage_table [][]string

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

				return nil
			},
		},
		{
			Name:  "download",
			Usage: "<storageid> <target>",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				storage := c.Args().First()
				target := c.Args().Get(1)
				if len(storage) == 0 || len(target) == 0 {
					log.Fatalln("You need to define a storage id and a target")
				}

				fetcher.DownloadArtefactsFromStorage(storage, target)

				return nil
			},
		},
		{
			Name:  "upload",
			Usage: "<storageid> <file> <storagepath>",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				storage := c.Args().First()
				file := c.Args().Get(1)
				path := c.Args().Get(2)

				fetcher.UploadStorageFile(storage, file, path)

				return nil
			},
		},
	},
}
