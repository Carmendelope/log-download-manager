/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// This is an example of an executable command.

package commands

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

var cmdSaySomething = &cobra.Command{
	Use:   "something [string to print]",
	Short: "Print anything to the screen",
	Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msgf("We say: %s", strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(cmdSaySomething)
}
