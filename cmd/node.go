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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	nodes "github.com/MottainaiCI/mottainai-server/pkg/nodes"
	"github.com/olekukonko/tablewriter"

	"github.com/urfave/cli"
)

var Node = cli.Command{
	Name:        "node",
	Usage:       "create, remove, list",
	Description: `Create, remove and list nodes`,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new node",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)

				res, err := fetcher.GetOptions("/api/nodes/add", map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove a node",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				node := c.Args().First()
				if len(node) == 0 {
					log.Fatalln("You need to define a node id")
				}
				res, err := fetcher.GetOptions("/api/nodes/delete/"+node, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},

		{
			Name:  "show",
			Usage: "show a node",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := &client.Fetcher{}
				node := c.Args().First()
				fetcher.BaseURL = host
				if len(node) == 0 {
					log.Fatalln("You need to define a node id")
				}
				var n nodes.Node
				fetcher.GetJSONOptions("/api/nodes/show/"+node, map[string]string{}, &n)

				//fmt.Println(t)

				b, err := json.MarshalIndent(n, "", "  ")
				if err != nil {
					log.Fatalln("error:", err)
				}
				fmt.Println(string(b))
				//for _, i := range tlist {
				//	fmt.Println(strconv.Itoa(i.ID) + " " + i.Status)
				//}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list nodes",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := &client.Fetcher{}
				fetcher.BaseURL = host
				var n []nodes.Node
				fetcher.GetJSONOptions("/api/nodes", map[string]string{}, &n)

				var node_table [][]string

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

				return nil
			},
		},
	},
}
