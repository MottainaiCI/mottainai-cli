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
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newTaskMonitorCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "monitor id1 id2 id3 ...",
		Short: "Monitor tasks and propagate exit status",
		Run: func(cmd *cobra.Command, args []string) {
			var v *viper.Viper = setting.Configuration.Viper

			fetcher := client.NewTokenClient(v.GetString("master"), v.GetString("apikey"))

			var tasks = make(map[string]bool)
			for _, id := range args {
				tasks[id] = false
			}
			MonitorTasks(fetcher, tasks)
		},
	}

	return cmd
}
