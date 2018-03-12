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

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var Namespace = cli.Command{
	Name:        "namespace",
	Usage:       "create, delete, tag, show, list, download",
	Description: `Create, delete, tag, show, list and download namespaces`,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new namespace",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				ns := c.Args().First()
				if len(ns) == 0 {
					log.Fatalln("You need to define a namespace name ")
				}
				res, err := fetcher.GetOptions("/api/namespace/"+ns+"/create", map[string]string{})
				if err != nil {
					return err
				}
				log.Println(string(res))

				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "delete a namespace",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				ns := c.Args().First()
				if len(ns) == 0 {
					log.Fatalln("You need to define a namespace name ")
				}
				res, err := fetcher.GetOptions("/api/namespace/"+ns+"/delete", map[string]string{})
				if err != nil {
					return err
				}
				log.Println(string(res))

				return nil
			},
		},
		{
			Name:  "show",
			Usage: "show artefacts belonging to namespace",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				ns := c.Args().First()
				if len(ns) == 0 {
					log.Fatalln("You need to define a namespace name ")
				}
				fmt.Println("Namespace: ", ns)
				var tlist []string

				fetcher.GetJSONOptions("/api/namespace/"+ns+"/list", map[string]string{}, &tlist)

				for _, i := range tlist {
					log.Println("- " + i)
				}
				return nil
			},
		},
		{
			Name:  "tag",
			Usage: "tag a namespace",
			Flags: []cli.Flag{
				StringFlag("from, f", "2442563452", "Task id"),
			},
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				from := c.String("from")
				fetcher := NewClient(host)
				ns := c.Args().First()
				if len(ns) == 0 || len(from) == 0 {
					log.Fatalln("You need to define a namespace name and a task id")
				}
				res, err := fetcher.GetOptions("/api/namespace/"+ns+"/tag/"+from, map[string]string{})
				if err != nil {
					return err
				}
				log.Println(string(res))

				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list namespaces",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)

				fmt.Println("Available namespaces: ")
				var tlist []string

				fetcher.GetJSONOptions("/api/namespace/list", map[string]string{}, &tlist)

				var ns_table [][]string

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

				return nil
			},
		},
		{
			Name:  "download",
			Usage: "<namespace> <target>",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				ns := c.Args().First()
				target := c.Args().Get(1)
				if len(ns) == 0 || len(target) == 0 {
					log.Fatalln("You need to define a namespace name and a target")
				}
				fetcher.DownloadArtefactsFromNamespace(ns, target)

				return nil
			},
		},
	},
}
