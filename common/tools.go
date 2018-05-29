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
package common

import (
	"fmt"
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"strings"
)

func PrintBuff(buff []byte) {
	data := string(buff)
	data = strings.TrimSpace(data)
	if len(data) > 0 {
		fmt.Println(data)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: pass settings in input.
func BuildCmdArgs(cmd *cobra.Command, msg string) string {
	var ans string = "mottainai-cli "

	if cmd == nil {
		panic("Invalid command")
	}

	if cmd.Flag("master").Changed {
		ans += "--master " + v.GetString("master") + " "
	}
	if v.GetString("profile") != "" {
		ans += "--profile " + v.GetString("profile") + " "
	}

	ans += msg

	return ans
}
