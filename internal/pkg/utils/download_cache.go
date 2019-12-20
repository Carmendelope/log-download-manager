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
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type DownloadLogState int

const (
	// ExpirationTime time the file is ready to download
	ExpirationTime = 10 * time.Minute
	ExpiredMsg     = "Expired"
	// ReviewTime time to check the status of the operations
	ReviewTime = 2 * time.Minute
	// AliveTime time the operation is stored
	AliveTime =  ExpirationTime +  2 * time.Minute
)

const (
	Queue DownloadLogState = iota + 1
	Generating
	Ready
	Error
	Downloaded
)

var DownloadLogStateFromGRPC = map[grpc_log_download_manager_go.DownloadLogState]DownloadLogState{
	grpc_log_download_manager_go.DownloadLogState_QUEUED:     Queue,
	grpc_log_download_manager_go.DownloadLogState_GENERATING: Generating,
	grpc_log_download_manager_go.DownloadLogState_READY:      Ready,
	grpc_log_download_manager_go.DownloadLogState_DOWNLOADED: Downloaded,
	grpc_log_download_manager_go.DownloadLogState_ERROR:      Error,
}

var DownloadLogStateToGRPC = map[DownloadLogState]grpc_log_download_manager_go.DownloadLogState{
	Queue:      grpc_log_download_manager_go.DownloadLogState_QUEUED,
	Generating: grpc_log_download_manager_go.DownloadLogState_GENERATING,
	Ready:      grpc_log_download_manager_go.DownloadLogState_READY,
	Downloaded: grpc_log_download_manager_go.DownloadLogState_DOWNLOADED,
	Error:      grpc_log_download_manager_go.DownloadLogState_ERROR,
}

func (d DownloadLogState) ToString() string {
	switch d {
	case Queue:
		{
			return "QUEUED"
		}
	case Generating:
		{
			return "GENERATING"
		}
	case Ready:
		{
			return "READY"
		}
	case Error:
		{
			return "ERROR"
		}
	case Downloaded:
		{
			return "DOWNLOADED"
		}
	}
	return ""
}

type DownloadOperation struct {
	OrganizationId string
	RequestId      string
	State          DownloadLogState
	Started        int64 // Creation time in ns
	From           int64
	To             int64
	Expiration     int64
	Info           string
	Url            string
	Directory      string
}

func (d *DownloadOperation) ToGRPC() *grpc_log_download_manager_go.DownloadLogResponse {
	return &grpc_log_download_manager_go.DownloadLogResponse{
		OrganizationId: d.OrganizationId,
		RequestId:      d.RequestId,
		From:           d.From,
		To:             d.To,
		State:          DownloadLogStateToGRPC[d.State],
		Url:            d.Url,
		Expiration:     d.Expiration,
		Info:           d.Info,
	}
}

type DownloadCache struct {
	sync.Mutex
	cache map[string]*DownloadOperation
	url   string
}

func NewDownloadCache(url string, publicHost string) *DownloadCache {
	res := &DownloadCache{
		cache: make(map[string]*DownloadOperation, 0),
		url:   fmt.Sprintf("https://web.%s%s", publicHost, url),
	}
	go res.ReviewOperations()
	return res
}

func (d *DownloadCache) CheckOperations() {

	d.Lock()
	defer d.Unlock()

	for i, ope := range d.cache {
		log.Debug().Str("index", i).Interface("operation", ope).Msg("operation")
		switch ope.State {
		// case Queue, Generating: nothing to do
		case Ready:
			if ope.Expiration  < time.Now().UnixNano(){
				err := RemoveFile(GetZipFilePath(ope.Directory, ope.RequestId))
				if err != nil {
					log.Warn().Str("requestId", ope.RequestId).Msg("error deleting zip file")
				}
				delete(d.cache, i)
				log.Debug().Msg("deleted")
			}
		case Error, Downloaded:
			if time.Unix(0, ope.Started).Add(AliveTime).After(time.Now()) {
				err := RemoveFile(GetZipFilePath(ope.Directory, ope.RequestId))
				if err != nil {
					log.Warn().Str("requestId", ope.RequestId).Msg("error deleting zip file")
				}
				delete(d.cache, i)
				log.Debug().Msg("deleted")
			}
		}
	}
}

func (d *DownloadCache) ReviewOperations() {
	log.Debug().Msg("ReviewOperations")
	ticker := time.NewTicker(ReviewTime)
	for {
		select {
		case <-ticker.C:
			d.CheckOperations()
		}
	}
}

func (d *DownloadCache) Add(organizationId string, requestId string, from int64, to int64, directory string) (*DownloadOperation, derrors.Error) {
	d.Lock()
	defer d.Unlock()

	_, exists := d.cache[requestId]
	if exists {
		return nil, derrors.NewAlreadyExistsError("operation").WithParams(requestId)
	}

	op := &DownloadOperation{
		OrganizationId: organizationId,
		RequestId:      requestId,
		Started:        time.Now().UnixNano(),
		State:          Queue,
		From:           from,
		To:             to,
		Directory:      directory,
	}
	d.cache[requestId] = op

	return op, nil

}

func (d *DownloadCache) Get(requestId string) (*DownloadOperation, derrors.Error) {

	d.Lock()
	defer d.Unlock()

	operation, exists := d.cache[requestId]
	if !exists {
		return nil, derrors.NewNotFoundError("download operation").WithParams(requestId)
	}

	if operation.State == Ready && operation.Expiration < time.Now().UnixNano() {
		operation.Info = ExpiredMsg
	}

	return operation, nil
}

func (d *DownloadCache) Update(requestId string, state DownloadLogState, info string) derrors.Error {
	d.Lock()
	defer d.Unlock()

	log.Debug().Str("requestId", requestId).Interface("state", state).Str("info", info).Msg("updating operation state")

	operation, exists := d.cache[requestId]
	if !exists {
		return derrors.NewNotFoundError("operation").WithParams(requestId)
	}
	operation.State = state
	operation.Info = info

	if state == Ready {
		operation.Expiration = time.Now().Add(ExpirationTime).UnixNano()
		operation.Url = fmt.Sprintf("%s%s.zip", d.url, requestId)
	}

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

func (d *DownloadCache) List(organizationID string) ([]*DownloadOperation, derrors.Error) {

	d.Lock()
	defer d.Unlock()

	list := make([]*DownloadOperation, 0)

	for _, ope := range d.cache {
		if ope.OrganizationId == organizationID {
			if ope.State == Ready && ope.Expiration < time.Now().UnixNano() {
				ope.Info = ExpiredMsg
			}
			list = append(list, ope)
		}
	}

	return list, nil

}

func (d *DownloadCache) Clean() {
	d.cache = make(map[string]*DownloadOperation, 0)
}
