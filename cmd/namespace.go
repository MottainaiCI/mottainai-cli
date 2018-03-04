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

	"github.com/urfave/cli"
)

var Namespace = cli.Command{
	Name:        "namespace",
	Usage:       "create|delete|tag|show|list",
	Description: `Create, delete, tag, show and list namespaces`,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new namespace",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				ns := c.Args().First()

				res, err := fetcher.GetOptions("/api/namespace/"+ns+"/create", map[string]string{})
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
				ns := c.Args().First()

				res, err := fetcher.GetOptions("/api/namespace/"+ns+"/delete", map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

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

				fmt.Println("Namespace: ", ns)
				var tlist []string

				fetcher.GetJSONOptions("/api/namespace/"+ns+"/list", map[string]string{}, &tlist)

				for _, i := range tlist {
					fmt.Println("- " + i)
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

				res, err := fetcher.GetOptions("/api/namespace/"+ns+"/tag/"+from, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list namespaces",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				ns := c.Args().First()

				fmt.Println("Available namespaces: ", ns)
				var tlist []string

				fetcher.GetJSONOptions("/api/namespace/list", map[string]string{}, &tlist)

				for _, i := range tlist {
					fmt.Println("- " + i)
				}
				return nil
			},
		},
	},
}
