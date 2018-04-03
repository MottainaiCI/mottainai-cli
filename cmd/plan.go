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
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var Plan = cli.Command{
	Name:        "plan",
	Usage:       "create, remove, list, show",
	Description: `Plan interface`,
	//Action:      runTask,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new plan",
			Flags: []cli.Flag{
				StringFlag("json", "file.json", "Decode parameters from a JSON file"),
				StringFlag("source, s", "https://github.com/foo/bar.git", "Repository url"),
				StringFlag("directory, d", "/test, /example", "Directory inside repository url"),
				StringFlag("script", "/foo/bar", "Entrypoint script"),
				StringFlag("task, t", "docker_execute", "Task name"),
				StringFlag("storage", "my_storage_id", "Task name"),
				StringFlag("image, i", "whatever/foo", "Image used from the task"),
				StringFlag("namespace, n", "test", "Specify a namespace the task will be started on"),
				StringFlag("storage_path, sp", "storage", "Specify the storage path in the task"),
				StringFlag("artefact_path, ap", "artefacts", "Specify the artefact path in the task"),
				StringFlag("tag_namespace, tn", "whatever", "Automatically to a namespace on success"),
				StringFlag("planned, pl", "@every 1m", "Plan task creation with cron syntax"),
			},
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				dat := make(map[string]string)

				if c.IsSet("json") {
					content, err := ioutil.ReadFile(c.String("json"))
					if err != nil {
						panic(err)
					}
					if err := json.Unmarshal(content, &dat); err != nil {
						panic(err)
					}
				} else {
					for _, n := range c.FlagNames() {
						dat[n] = c.String(n)
					}
				}

				res, err := fetcher.Form("/api/tasks/plan", dat)
				if err != nil {
					panic(err)
				}
				tid := string(res)
				fmt.Println("-------------------------")
				fmt.Println("Plan " + tid + " has been created")
				fmt.Println("-------------------------")
				fmt.Println("Information:", " mottainai-cli --master "+host+" plan show "+tid)
				fmt.Println("-------------------------")

				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove a plan",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a plan id")
				}
				res, err := fetcher.GetOptions("/api/tasks/plan/delete/"+task, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "show",
			Usage: "show a plan",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				task := c.Args().First()
				fetcher := NewClient(host)
				if len(task) == 0 {
					log.Fatalln("You need to define a plan id")
				}
				var t citasks.Plan
				fetcher.GetJSONOptions("/api/tasks/plan/"+task, map[string]string{}, &t)
				b, err := json.MarshalIndent(t, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}
				fmt.Println(string(b))
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list plans",
			Flags: []cli.Flag{
				BoolFlag("quiet, q", "Quiet output"),
			},
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := &client.Fetcher{}
				fetcher.BaseURL = host
				var tlist []citasks.Plan
				fetcher.GetJSONOptions("/api/tasks/planned", map[string]string{}, &tlist)

				sort.Slice(tlist[:], func(i, j int) bool {
					return tlist[i].CreatedTime > tlist[j].CreatedTime
				})
				if c.Bool("quiet") {
					for _, i := range tlist {
						fmt.Println(i.ID)
					}
					return nil
				}

				var task_table [][]string

				for _, i := range tlist {
					task_table = append(task_table, []string{strconv.Itoa(i.ID), i.Planned, i.Directory})
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
				table.SetCenterSeparator("|")
				table.SetHeader([]string{"ID", "Planned", "Dir"})

				for _, v := range task_table {
					table.Append(v)
				}
				table.Render()

				return nil
			},
		},
	},
}
