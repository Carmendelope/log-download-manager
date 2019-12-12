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

package log_manager

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required clients for roles operations.
type Manager struct {
	appManagerClient grpc_application_manager_go.UnifiedLoggingClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(appManagerClient grpc_application_manager_go.UnifiedLoggingClient) Manager {
	return Manager{appManagerClient: appManagerClient}
}

func (m *Manager) download(request *grpc_log_download_manager_go.DownloadLogRequest, requestId string) {
	log.Debug().Str("requestId", requestId).Msg("downloading logs...")
}

// DownloadLog asks for a logs download operation. These logs are going to be stored in a zip file
func (m *Manager) DownloadLog(request *grpc_log_download_manager_go.DownloadLogRequest) (*grpc_log_download_manager_go.DownloadLogResponse, derrors.Error) {

	requestId := uuid.New().String()

	return nil, nil
}

// Check asks for a download operation state
func (m *Manager) Check(request *grpc_log_download_manager_go.DownloadRequestId) (*grpc_log_download_manager_go.DownloadLogResponse, derrors.Error) {
	return nil, nil
}
