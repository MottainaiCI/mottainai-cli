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
	"time"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var Task = cli.Command{
	Name:        "task",
	Usage:       "create, clone, remove, stop, start, show, list, attach, artefacts, download, execute",
	Description: `Task interface`,
	//Action:      runTask,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new task",
			Flags: []cli.Flag{
				StringFlag("json", "", "Decode parameters from a JSON file ( e.g. /path/to/file.json )"),
				StringFlag("source, s", "", "Repository url ( e.g. https://github.com/foo/bar.git )"),
				StringFlag("directory, d", "", "Directory inside repository url ( e.g. /test )"),
				StringFlag("script", "", "Entrypoint script"),
				StringFlag("task, t", "docker_execute", "Task type ( default: docker_execute )"),
				StringFlag("storage", "", "Storage ID"),
				StringFlag("image, i", "", "Image used from the task ( e.g. my/docker-image:latest )"),
				StringFlag("namespace, n", "", "Specify a namespace the task will be started on"),
				StringFlag("storage_path, sp", "storage", "Specify the storage path in the task"),
				StringFlag("artefact_path, ap", "artefacts", "Specify the artefact path in the task"),
				StringFlag("tag_namespace, tn", "", "Automatically to the specified namespace on success"),
				StringFlag("prune, pr", "yes", "Perform pruning actions after execution"),
				StringFlag("cache_image, ci", "yes", "Cache image after execution inside the host for later reuse."),
			},
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				dat := make(map[string]interface{})

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

				res, err := fetcher.GenericForm("/api/tasks", dat)
				if err != nil {
					panic(err)
				}
				tid := string(res)
				fmt.Println("-------------------------")
				fmt.Println("Task " + tid + " has been created")
				fmt.Println("-------------------------")
				fmt.Println("Live log:", " mottainai-cli --master "+host+" task attach "+tid)
				fmt.Println("Information:", " mottainai-cli --master "+host+" task show "+tid)
				fmt.Println("URL:", " "+fetcher.BaseURL+"/tasks/display/"+tid)
				fmt.Println("Build Log:", " "+fetcher.BaseURL+"/artefact/"+tid+"/build_"+tid+".log")
				fmt.Println("-------------------------")

				return nil
			},
		},
		{
			Name:  "clone",
			Usage: "clone a task",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				res, err := fetcher.GetOptions("/api/tasks/clone/"+task, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove a task",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				res, err := fetcher.GetOptions("/api/tasks/delete/"+task, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "start",
			Usage: "start a task",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				res, err := fetcher.GetOptions("/api/tasks/start/"+task, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "stop a task",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				_, err := fetcher.GetOptions("/api/tasks/stop/"+task, map[string]string{})
				if err != nil {
					log.Fatalln(err)
					return err
				}
				fmt.Println("Request sent successfully")

				return nil
			},
		},
		{
			Name:  "show",
			Usage: "show a task",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				task := c.Args().First()
				fetcher := NewClient(host)
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				var t citasks.Task
				fetcher.GetJSONOptions("/api/tasks/"+task, map[string]string{}, &t)

				//fmt.Println(t)

				b, err := json.MarshalIndent(t, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}
				fmt.Println(string(b))
				//for _, i := range tlist {
				//	fmt.Println(strconv.Itoa(i.ID) + " " + i.Status)
				//}
				return nil
			},
		},
		{
			Name:  "artefacts",
			Usage: "shows artefacts of a task",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				fmt.Println("Artefacts for:", task)
				var tlist []string
				fetcher.GetJSONOptions("/api/tasks/"+task+"/artefacts", map[string]string{}, &tlist)

				for _, i := range tlist {
					fmt.Println("- " + i)
				}
				return nil
			},
		},
		{
			Name:  "log",
			Usage: "show log of a task",
			Action: func(c *cli.Context) {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				buff, err := fetcher.GetOptions("/api/tasks/stream_output/"+task+"/0", map[string]string{})
				if err != nil {
					panic(err)
				}
				printBuff(buff)
			},
		},
		{
			Name:  "attach",
			Usage: "attach to a task output",
			Action: func(c *cli.Context) {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()
				if len(task) == 0 {
					log.Fatalln("You need to define a task id")
				}
				var pos = 0

				for {
					time.Sleep(time.Second + 2)
					var t citasks.Task
					fetcher.GetJSONOptions("/api/tasks/"+task, map[string]string{}, &t)
					if t.Status != "running" {
						if t.Status == "done" && pos == 0 {
							buff, err := fetcher.GetOptions("/artefact/"+task+"/build_"+strconv.Itoa(t.ID)+".log", map[string]string{})
							if err != nil {
								panic(err)
							}
							printBuff(buff)
						} else {
							fmt.Println("Build status: " + t.Status + " Can't attach to any live stream.")
						}
						return
					}

					buff, err := fetcher.GetOptions("/api/tasks/stream_output/"+task+"/"+strconv.Itoa(pos), map[string]string{})
					if err != nil {
						panic(err)
					}
					pos += len(buff)
					printBuff(buff)
				}

			},
		},
		{
			Name:  "list",
			Usage: "list tasks",
			Flags: []cli.Flag{
				BoolFlag("quiet, q", "Quiet output"),
			},
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := &client.Fetcher{}
				fetcher.BaseURL = host
				var tlist []citasks.Task
				fetcher.GetJSONOptions("/api/tasks", map[string]string{}, &tlist)

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
					t, _ := time.Parse("20060102150405", i.CreatedTime)
					t2, _ := time.Parse("20060102150405", i.EndTime)
					task_table = append(task_table, []string{strconv.Itoa(i.ID), i.Status, i.Result, t.String(), t2.String(), i.Source, i.Directory})
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
				table.SetCenterSeparator("|")
				table.SetHeader([]string{"ID", "Status", "Result", "Created", "End", "Source", "Dir"})

				for _, v := range task_table {
					table.Append(v)
				}
				table.Render()

				return nil
			},
		},
		{
			Name:  "download",
			Usage: "<taskid> <target>",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				id := c.Args().First()
				target := c.Args().Get(1)
				if len(id) == 0 || len(target) == 0 {
					log.Fatalln("You need to define a task id and a target")
				}

				fetcher.DownloadArtefactsFromTask(id, target)

				return nil
			},
		},
		{
			Name:  "execute",
			Usage: "<taskid>",
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				id := c.Args().First()

				var t citasks.Task
				err := fetcher.GetJSONOptions("/api/tasks/"+id, map[string]string{}, &t)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				fn := citasks.DefaultTaskHandler().Handler(t.TaskName)
				setting.GenDefault()
				setting.Configuration.AppURL = host
				fn(id)

				return nil
			},
		},
	},
}
