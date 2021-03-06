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
	"context"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/log-download-manager/internal/pkg/entities"
	"github.com/nalej/log-download-manager/internal/pkg/utils"
)

// Handler structure for the user requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// DownloadLog asks for a logs download operation. These logs are going to be stored in a zip file
func (h *Handler) DownloadLog(ctx context.Context, request *grpc_log_download_manager_go.DownloadLogRequest) (*grpc_log_download_manager_go.DownloadLogResponse, error) {

	vErr := entities.ValidDownloadLogRequest(request)
	if vErr != nil {
		return nil, conversions.ToDerror(vErr)
	}

	response, err := h.Manager.DownloadLog(request, utils.GetUserFromContext(ctx))
	if err != nil {
		return nil, conversions.ToDerror(err)
	}
	return response, nil
}

// Check asks for a download operation state
func (h *Handler) Check(ctx context.Context, request *grpc_log_download_manager_go.DownloadRequestId) (*grpc_log_download_manager_go.DownloadLogResponse, error) {

	vErr := entities.ValidDownloadRequestId(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	response, err := h.Manager.Check(request, utils.GetUserFromContext(ctx))
	if err != nil {
		return nil, conversions.ToDerror(err)
	}
	return response, nil
}


func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_log_download_manager_go.DownloadLogResponseList, error) {

	vErr := entities.ValidOrganizationId(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	response, err := h.Manager.List(organizationID, utils.GetUserFromContext(ctx))
	if err != nil {
		return nil, conversions.ToDerror(err)
	}
	return response, nil
}
