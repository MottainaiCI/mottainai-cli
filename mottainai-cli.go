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
package main

import (
	"fmt"

	cli "github.com/MottainaiCI/mottainai-cli/cmd"
	common "github.com/MottainaiCI/mottainai-cli/common"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
)

func main() {

	v := setting.Configuration.Viper
	v.SetDefault("master", "http://localhost:8080")
	v.SetDefault("profile", "")
	// Temporary for config print
	v.SetDefault("config", "")

	// Initialize Default Viper Configuration
	setting.GenDefault(v)

	// Define env variable
	v.SetEnvPrefix(common.MCLI_ENV_PREFIX)
	v.AutomaticEnv()

	// Set Config paths list
	v.AddConfigPath(common.MCLI_LOCAL_PATH)
	v.AddConfigPath(fmt.Sprintf("$HOME/%s", common.MCLI_HOME_PATH))

	// Set config file name (without extension)
	v.SetConfigName(common.MCLI_CONFIG_NAME)

	v.SetTypeByDefaultValue(true)

	cli.Execute()
}
