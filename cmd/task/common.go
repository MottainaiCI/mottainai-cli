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

package task

import (
	"regexp"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	nodes "github.com/MottainaiCI/mottainai-server/pkg/nodes"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"

	"fmt"
	"os"

	"github.com/mudler/anagent"
)

func GenerateTasks(c *client.Fetcher, dat map[string]interface{}, hostreg string) map[string]bool {
	reg, err := regexp.Compile(hostreg)
	if err != nil {
		panic(err)
	}
	var created = make(map[string]bool)

	var n []nodes.Node
	var q []string
	c.GetJSONOptions("/api/nodes", map[string]string{}, &n)

	for _, i := range n {
		// Make a Regex to say we only want
		if reg.MatchString(i.Hostname + i.NodeID) {
			q = append(q, i.Hostname+i.NodeID)
			fmt.Println("Node matched regex: ", i.Hostname+i.NodeID)
		}

	}
	for _, queue := range q {
		dat["queue"] = queue
		res, err := c.GenericForm("/api/tasks", dat)
		if err != nil {
			panic(err)
		}
		tid := string(res)
		fmt.Println("Task "+tid+" has been created for", queue)
		created[tid] = false
	}
	return created
}
func MonitorTasks(f *client.Fetcher, created map[string]bool) {
	agent := anagent.New()
	var done int
	var res = 0
	agent.Map(f)
	for k, _ := range created {
		fmt.Println("Tracking ", k)
	}
	agent.TimerSeconds(int64(1), true, func(c *client.Fetcher) {

		if done >= len(created) {
			agent.Stop()
		}

		for k, v := range created {
			var t citasks.Task
			c.GetJSONOptions("/api/tasks/"+k, map[string]string{}, &t)
			if t.ID == "" && !v {
				// There is no task anymore associated with it!
				done++
				res = 1 // Error :( something went wrong!
				fmt.Println("Error: No task associated with id ", k)
			}
			if t.IsDone() && !v {
				done++
				created[k] = true
				fmt.Println("Task ", k, "Done")

				if !t.IsSuccess() {
					res = 1
					fmt.Println("Task ", k, "Fail")
				} else {
					fmt.Println("Task ", k, "Success")
				}
			}
		}

	})

	agent.Start()
	os.Exit(res)
}
