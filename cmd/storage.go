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
	"strconv"

	storage "github.com/MottainaiCI/mottainai-server/pkg/storage"

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
			Usage: "list storages",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				storage_name := c.Args().First()

				fmt.Println("Available storages: ", storage_name)

				var n []storage.Storage
				fetcher.GetJSONOptions("/api/storage/list", map[string]string{}, &n)

				for _, i := range n {
					fmt.Println(strconv.Itoa(i.ID) + " Name:" + i.Name + " Path:" + i.Path)
				}

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
