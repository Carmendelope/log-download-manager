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
 */

package commands

import (
	"github.com/nalej/log-download-manager/internal/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = server.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the server API",
	Long:  `Launch the server API`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		config.Debug = debugLevel
		server := server.NewService(config)
		server.Run()
	},
}

func init() {
	runCmd.Flags().IntVar(&config.Port, "port", 8940, "Port to launch the log-download-manager gRPC")
	runCmd.Flags().IntVar(&config.HttpPort, "httpPort", 8941, "Port to launch the log-download-manager Http")
	runCmd.PersistentFlags().StringVar(&config.ApplicationsManagerAddress, "applicationsManagerAddress", "localhost:8910",
		"Applications Manager address (host:port)")
	runCmd.PersistentFlags().StringVar(&config.DownloadPath, "downloadPath", "/download",
		"download directory path")
	runCmd.PersistentFlags().StringVar(&config.AuthHeader, "authHeader", "", "Authorization Header")
	runCmd.PersistentFlags().StringVar(&config.AuthSecret, "authSecret", "", "Authorization secret")
	runCmd.PersistentFlags().StringVar(&config.ManagementPublicHost, "managementPublicHost", "", "Management publish host")

	rootCmd.AddCommand(runCmd)
}
