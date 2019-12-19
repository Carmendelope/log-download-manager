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

package http_log_manager

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/log-download-manager/internal/pkg/server/interceptor"
	"github.com/nalej/log-download-manager/internal/pkg/utils"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Manager structure with the required clients for http log download operations.
type Manager struct {
	opeCache          *utils.DownloadCache
	interceptor       *interceptor.Interceptor
	pathPrefix        string
	DownloadDirectory string
}

func NewManager(opeCache *utils.DownloadCache, secret string, authHeader string, pathPrefix string, dir string) Manager {
	return Manager{
		opeCache:          opeCache,
		interceptor:       interceptor.NewInterceptor(secret, authHeader),
		pathPrefix:        pathPrefix,
		DownloadDirectory: dir,
	}
}

func (m *Manager) SplitPath(path string) (string, string, derrors.Error) {
	p := strings.TrimPrefix(path, m.pathPrefix)
	split := strings.Split(p, ".")
	if len(split) != 2 {
		return "", "", derrors.NewInvalidArgumentError("invalid path").WithParams(path)
	}
	return p, split[0], nil
}

// ValidToDownload check if the file can be downloaded
func (m *Manager) ValidToDownload(ope *utils.DownloadOperation) derrors.Error {

	switch ope.State {
	case utils.Error, utils.Queue, utils.Generating:
		errMsg := fmt.Sprintf("download operation is not ready. State (%s)", ope.State.ToString())
		if ope.Info != "" {
			errMsg = fmt.Sprintf("%s - %s", errMsg, ope.Info)
		}
		return derrors.NewPermissionDeniedError(errMsg)
	case utils.Ready:
		// check if it is expired
		if ope.Expiration < time.Now().UnixNano() {
			errMsg := "download operation is not ready. State (EXPIRED)"
			return derrors.NewPermissionDeniedError(errMsg)
		}
	}

	return nil
}

func (m *Manager) DownloadFile() http.Handler {
	h := http.FileServer(http.Dir(m.DownloadDirectory))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vErr := m.interceptor.Validate(r)
		if vErr != nil {
			http.Error(w, vErr.Error(), http.StatusUnauthorized)
			return
		}

		file, requestId, err := m.SplitPath(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ope, getErr := m.opeCache.Get(requestId)
		if getErr != nil {
			http.Error(w, getErr.Error(), http.StatusInternalServerError)
			return
		}
		vOpeErr := m.ValidToDownload(ope)
		if vOpeErr != nil {
			http.Error(w, getErr.Error(), http.StatusUnauthorized)
			return
		}

		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = file
		h.ServeHTTP(w, r2)


		err = m.opeCache.Update(requestId, utils.Downloaded, "")
		/*
		err = m.opeCache.Remove(requestId)
		if err != nil {
			log.Warn().Str("requestId", requestId).Msg("unable to delete operation")
		}

		removeErr := utils.RemoveFile(utils.GetZipFilePath(m.DownloadDirectory, requestId))

		if removeErr != nil {
			log.Warn().Str("file", file).Msg("unable to remove file")
		}
		*/

	})
}
