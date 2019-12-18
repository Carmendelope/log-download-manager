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

package server

import (
	"github.com/nalej/derrors"
	"github.com/nalej/log-download-manager/version"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Debug level is active.
	Debug bool
	// Port where the gRPC API service will listen requests.
	Port int
	// HttpPort where the HTTP service will listen request
	HttpPort int
	// ApplicationsManagerAddress with the host:port to connect to the Applications manager.
	ApplicationsManagerAddress string
	// DownloadPath with the path where the logs are going to be stored
	DownloadPath string
}

func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 || conf.HttpPort <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}

	if conf.ApplicationsManagerAddress == "" {
		return derrors.NewInvalidArgumentError("applicationsManagerAddress must be set")
	}

	if conf.DownloadPath == "" {
		return derrors.NewInvalidArgumentError("DownloadDir must be set")
	}

	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Int("HttpPort", conf.HttpPort).Msg("Http port")
	log.Info().Str("URL", conf.ApplicationsManagerAddress).Msg("Applications Manager")
	log.Info().Str("DownloadPath", conf.DownloadPath).Msg("download Path")

}
