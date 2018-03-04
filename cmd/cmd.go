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
	"time"

	"github.com/MottainaiCI/mottainai-server/pkg/client"
	"github.com/urfave/cli"
)

func NewClient(host string) *client.Fetcher {
	fetcher := &client.Fetcher{}
	fetcher.BaseURL = host
	return fetcher
}

func StringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func BoolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func IntFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func DurationFlag(name string, value time.Duration, usage string) cli.DurationFlag {
	return cli.DurationFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
