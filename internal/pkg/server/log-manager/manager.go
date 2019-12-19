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
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/log-download-manager/internal/pkg/entities"
	"github.com/nalej/log-download-manager/internal/pkg/utils"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required clients for roles operations.
type Manager struct {
	appManagerClient  grpc_application_manager_go.UnifiedLoggingClient
	opeCache          *utils.DownloadCache
	DownloadDirectory string
}

// NewManager creates a Manager using a set of clients.
func NewManager(appManagerClient grpc_application_manager_go.UnifiedLoggingClient, opeCache *utils.DownloadCache, downloadDirectory string) Manager {
	return Manager{
		appManagerClient:  appManagerClient,
		opeCache:          opeCache,
		DownloadDirectory: downloadDirectory,
	}
}

func (m *Manager) getFilePath(requestId string) string {
	return fmt.Sprintf("%s/%s.file", m.DownloadDirectory, requestId)
}
func (m *Manager) getZipFilePath(requestId string) string {
	return fmt.Sprintf("%s/%s.zip", m.DownloadDirectory, requestId)
}

// download generates the zip file with the log entries
func (m *Manager) download(request *grpc_log_download_manager_go.DownloadLogRequest, requestId string) {
	log.Debug().Str("requestId", requestId).Msg("downloading logs...")

	// 1.- update the status of the operation
	updateErr := m.opeCache.Update(requestId, utils.Generating, "")
	if updateErr != nil {
		log.Error().Err(updateErr).Msg("error updating the operation state")
	}

	// 2.- create the search request
	searchRequest := entities.NewSearchRequest(request)

	for {
		// check it the connection already exists
		ctx, cancel := utils.GetContext()
		// 3.- Search
		response, err := m.appManagerClient.Search(ctx, searchRequest)
		if err != nil {
			m.opeCache.Update(requestId, utils.Error, err.Error())
			cancel()
			break
		} else {
			cancel()
			log.Debug().Int("responses", len(response.Entries)).Msg("entries retrieved")
			if len(response.Entries) > 0 {
				// 4.- Copy the log entries in a file ordered
				err = utils.AppendResponses(entities.Sort(response.Entries, request.Order.Order), utils.GetFilePath(m.DownloadDirectory, requestId), request.IncludeMetadata)
				if err != nil {
					updateErr := m.opeCache.Update(requestId, utils.Error, err.Error())
					if updateErr != nil {
						log.Error().Err(updateErr).Msg("error updating the operation state")
					}
					break
				}
			} else {
				// 5.- If there is no more entries -> create zip file
				zipErr := utils.ZipFiles(utils.GetZipFilePath(m.DownloadDirectory, requestId), []string{utils.GetFilePath(m.DownloadDirectory, requestId)})

				if zipErr != nil {
					updateErr := m.opeCache.Update(requestId, utils.Error, zipErr.Error())
					if updateErr != nil {
						log.Error().Err(updateErr).Msg("error updating the operation state")
					}
				} else {
					updateErr := m.opeCache.Update(requestId, utils.Ready, "file generated")
					if updateErr != nil {
						log.Error().Err(updateErr).Msg("error updating the operation state")
					}
					utils.RemoveFile(utils.GetFilePath(m.DownloadDirectory, requestId))
				}
				break
			}
			if request.Order.Order == grpc_common_go.Order_ASC {
				searchRequest.From = response.To + 1000000
			} else {
				searchRequest.To = response.From - 1000000
			}
		}

	}
}

// DownloadLog asks for a logs download operation. These logs are going to be stored in a zip file
func (m *Manager) DownloadLog(request *grpc_log_download_manager_go.DownloadLogRequest) (*grpc_log_download_manager_go.DownloadLogResponse, derrors.Error) {

	log.Debug().Interface("request", request).Msg("DownloadLog request")
	requestId := uuid.New().String()
	op, err := m.opeCache.Add(request.OrganizationId, requestId, request.From, request.To)
	if err != nil {
		return nil, err
	}

	// Create the file
	utils.InitializeFile(utils.GetFilePath(m.DownloadDirectory, requestId), request.IncludeMetadata)

	go m.download(request, requestId)

	return &grpc_log_download_manager_go.DownloadLogResponse{
		OrganizationId: request.OrganizationId,
		RequestId:      op.RequestId,
		From:           op.From,
		To:             op.To,
		State:          utils.DownloadLogStateToGRPC[op.State],
	}, nil
}

// Check asks for a download operation state
func (m *Manager) Check(request *grpc_log_download_manager_go.DownloadRequestId) (*grpc_log_download_manager_go.DownloadLogResponse, derrors.Error) {
	operation, err := m.opeCache.Get(request.RequestId)
	if err != nil {
		return nil, conversions.ToDerror(err)
	}

	return entities.NewDownloadLogResponse(request, operation), nil
}

// List retrieves a list of LogResponses
func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_log_download_manager_go.DownloadLogResponseList, derrors.Error) {
	list, err := m.opeCache.List(organizationID.OrganizationId)
	if err != nil {
		return nil, conversions.ToDerror(err)
	}
	logResponseList := make([]*grpc_log_download_manager_go.DownloadLogResponse, len(list))
	for i, ope := range list {
		logResponseList[i] = ope.ToGRPC()
	}
	return &grpc_log_download_manager_go.DownloadLogResponseList{
		Responses: logResponseList,
	}, nil
}
