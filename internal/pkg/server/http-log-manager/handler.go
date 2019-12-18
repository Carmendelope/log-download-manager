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
	"net/http"
)

// Handler structure for the http requests.
type Handler struct {
	Manager Manager

}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{Manager:manager}
}

func (h *Handler) DownloadFile(prefix string, handler http.Handler) http.Handler {
	return h.Manager.DownloadFile(prefix, handler)
}

func (h *Handler) DownloadFile2() http.Handler {
	return h.Manager.DownloadFile2()
}
