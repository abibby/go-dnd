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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zwzn/go-dnd/event"
)

// lootCmd represents the loot command
var lootCmd = &cobra.Command{
	Use:   "loot",
	Short: "Add money in gp",
	Run: func(cmd *cobra.Command, args []string) {
		amount, err := strconv.Atoi(args[0])
		check(err)

		check(writeEvent(event.Event{
			Type: "money",
			MoneyEvent: event.MoneyEvent{
				Amount: event.Money(amount),
			},
		}))
	},
}

func init() {
	rootCmd.AddCommand(lootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
