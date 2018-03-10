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
	"strconv"
	"time"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"

	"github.com/urfave/cli"
)

var Task = cli.Command{
	Name:        "task",
	Usage:       "create, remove, stop, start, show, list, attach, artefacts, download, execute",
	Description: `Task interface`,
	//Action:      runTask,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new task",
			Flags: []cli.Flag{
				BoolFlag("json", "Decode parameters as JSON from --file"),
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

					fmt.Printf("File contents: %s", content)

					if err := json.Unmarshal(content, &dat); err != nil {
						panic(err)
					}

					res, err := fetcher.Form("/api/tasks", dat)
					if err != nil {
						panic(err)
					}
					fmt.Println(fetcher.BaseURL + "/tasks/display/" + string(res))
					return nil
				}

				for _, n := range c.FlagNames() {
					dat[n] = c.String(n)
				}
				fmt.Println(dat)

				res, err := fetcher.Form("/api/tasks", dat)
				if err != nil {
					panic(err)
				}
				fmt.Println(fetcher.BaseURL + "/tasks/display/" + string(res))

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

				res, err := fetcher.GetOptions("/api/tasks/start"+task, map[string]string{})
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

				res, err := fetcher.GetOptions("/api/tasks/stop/"+task, map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(string(res))

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
			Name:  "attach",
			Usage: "attach to a task output",
			Action: func(c *cli.Context) {
				host := c.GlobalString("master")
				fetcher := NewClient(host)
				task := c.Args().First()

				var pos = 0

				for {
					time.Sleep(time.Second + 2)
					var t citasks.Task
					fetcher.GetJSONOptions("/api/tasks/"+task, map[string]string{}, &t)
					if t.Status != "running" {
						if t.Status == "done" && pos == 0 {
							buff, err := fetcher.GetOptions("/artefact/"+task+"/build.log", map[string]string{})
							if err != nil {
								panic(err)
							}
							printBuff(buff)
						} else {
							fmt.Println("Can't attach to any log")
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
			Action: func(c *cli.Context) error {
				host := c.GlobalString("master")
				fetcher := &client.Fetcher{}
				fetcher.BaseURL = host
				var tlist []citasks.Task
				fetcher.GetJSONOptions("/api/tasks", map[string]string{}, &tlist)

				for _, i := range tlist {
					fmt.Println(strconv.Itoa(i.ID) + " " + i.Status)
				}
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
