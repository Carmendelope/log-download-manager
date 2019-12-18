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
	"github.com/nalej/log-download-manager/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

// Manager structure with the required clients for http log download operations.
type Manager struct {
	opeCache *utils.DownloadCache
}

func NewManager(opeCache *utils.DownloadCache) Manager {
	return Manager{
		opeCache:opeCache,
	}
}

func (s *Manager) DownloadFile(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Debug().Msg("aqui meto el interceptor")

		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}