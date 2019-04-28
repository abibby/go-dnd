// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"github.com/zwzn/go-dnd/character"
	"github.com/zwzn/go-dnd/event"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "build the character sheet",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		ch, err := character.NewFile(viper.GetString("character-file"))
		check(err)

		err = event.UpdateCharacterFile(ch, viper.GetString("log-file"))
		check(err)

		outFile, err := cmd.Flags().GetString("out-file")
		check(err)

		out := io.Writer(os.Stdout)
		if outFile != "" {
			f, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			check(err)
			out = f
		}
		err = ch.Render(out)
		check(err)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("out-file", "o", "", "where to write the html (default is stdout)")
}
