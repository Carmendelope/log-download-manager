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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/grpc-organization-go"
)

const emptyOrganizationId = "organization_id cannot be empty"
const emptyRequestId = "request_id cannot be empty"

func ValidDownloadLogRequest(request *grpc_log_download_manager_go.DownloadLogRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidDownloadRequestId(request *grpc_log_download_manager_go.DownloadRequestId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.RequestId == "" {
		return derrors.NewInvalidArgumentError(emptyRequestId)
	}

	return nil
}

func ValidOrganizationId (organizationID *grpc_organization_go.OrganizationId) derrors.Error{
	if organizationID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}