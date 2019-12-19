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
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/log-download-manager/internal/pkg/utils"
	"sort"
	"time"
)

func NewSearchRequest(request *grpc_log_download_manager_go.DownloadLogRequest) *grpc_application_manager_go.SearchRequest {
	to := request.To
	if to == 0 {
		to = time.Now().UnixNano()
	}
	nfirst := true
	if request.Order.Order == grpc_common_go.Order_DESC {
		nfirst = false
	}
	return &grpc_application_manager_go.SearchRequest{
		OrganizationId:         request.OrganizationId,
		AppDescriptorId:        request.AppDescriptorId,
		AppInstanceId:          request.AppInstanceId,
		ServiceGroupId:         request.ServiceGroupId,
		ServiceGroupInstanceId: request.ServiceGroupInstanceId,
		ServiceId:              request.ServiceId,
		ServiceInstanceId:      request.ServiceInstanceId,
		MsgQueryFilter:         request.MsgQueryFilter,
		From:                   request.From,
		To:                     to,
		IncludeMetadata:        request.IncludeMetadata,
		NFirst:                 nfirst,
	}
}

func NewDownloadLogResponse(request *grpc_log_download_manager_go.DownloadRequestId, opeInfo *utils.DownloadOperation) *grpc_log_download_manager_go.DownloadLogResponse {
	return &grpc_log_download_manager_go.DownloadLogResponse{
		OrganizationId: request.OrganizationId,
		RequestId:      request.RequestId,
		From:           opeInfo.From,
		To:             opeInfo.To,
		State:          utils.DownloadLogStateToGRPC[opeInfo.State],
		Expiration:     opeInfo.Expiration,
		Info:           opeInfo.Info,
		Url:            opeInfo.Url,
	}

}

func Sort(elements []*grpc_application_manager_go.LogEntryResponse, order grpc_common_go.Order) []*grpc_application_manager_go.LogEntryResponse {
	if len(elements) == 0 {
		return elements
	}
	if order == grpc_common_go.Order_ASC {
		sort.SliceStable(elements, func(i, j int) bool {

			return elements[i].Timestamp < elements[j].Timestamp

		})
	} else {
		sort.SliceStable(elements, func(i, j int) bool {

			return elements[i].Timestamp > elements[j].Timestamp

		})
	}
	return elements
}

