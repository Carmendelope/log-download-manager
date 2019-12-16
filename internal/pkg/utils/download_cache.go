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

package utils

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type DownloadLogState int

const (
	Queue DownloadLogState = iota + 1
	Generating
	Ready
	Error
)

var DownloadLogStateFromGRPC = map[grpc_log_download_manager_go.DownloadLogState]DownloadLogState{
	grpc_log_download_manager_go.DownloadLogState_QUEUED:     Queue,
	grpc_log_download_manager_go.DownloadLogState_GENERATING: Generating,
	grpc_log_download_manager_go.DownloadLogState_READY:      Ready,
	grpc_log_download_manager_go.DownloadLogState_ERROR:      Error,
}

var DownloadLogStateToGRPC = map[DownloadLogState]grpc_log_download_manager_go.DownloadLogState{
	Queue:      grpc_log_download_manager_go.DownloadLogState_QUEUED,
	Generating: grpc_log_download_manager_go.DownloadLogState_GENERATING,
	Ready:      grpc_log_download_manager_go.DownloadLogState_READY,
	Error:      grpc_log_download_manager_go.DownloadLogState_ERROR,
}

type DownloadOperation struct {
	RequestId string
	State     DownloadLogState
	// Creation time in ns
	Started    int64
	From       int64
	To         int64
	Expiration int64
}

type DownloadCache struct {
	sync.Mutex
	cache map[string]*DownloadOperation
}

func NewDownloadCache() *DownloadCache {
	return &DownloadCache{
		cache: make(map[string]*DownloadOperation, 0),
	}
}

func (d *DownloadCache) Add(requestId string, from int64, to int64) (*DownloadOperation, derrors.Error) {
	d.Lock()
	defer d.Unlock()

	_, exists := d.cache[requestId]
	if exists {
		return nil, derrors.NewAlreadyExistsError("operation").WithParams(requestId)
	}

	op := &DownloadOperation{
		RequestId: requestId,
		Started:   time.Now().UnixNano(),
		State:     Queue,
		From:      from,
		To:        to,
	}
	d.cache[requestId] = op

	return op, nil

}

func (d *DownloadCache) Get(requestId string) (*DownloadOperation, derrors.Error) {

	d.Lock()
	defer d.Unlock()

	operation, exists := d.cache[requestId]
	if !exists {
		return nil, derrors.NewNotFoundError("operation").WithParams(requestId)
	}
	return operation, nil
}

func (d *DownloadCache) Update(requestId string, state DownloadLogState) derrors.Error {
	d.Lock()
	defer d.Unlock()

	log.Debug().Str("requestId", requestId).Interface("state", state).Msg("updating operation state")

	operation, exists := d.cache[requestId]
	if !exists {
		return derrors.NewNotFoundError("operation").WithParams(requestId)
	}
	operation.State = state

	return nil
}

func (d *DownloadCache) Remove(requestId string) derrors.Error {

	d.Lock()
	defer d.Unlock()

	_, exists := d.cache[requestId]
	if !exists {
		return derrors.NewNotFoundError("operation").WithParams(requestId)
	}
	delete(d.cache, requestId)

	return nil
}
