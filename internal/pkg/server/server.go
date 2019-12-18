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

package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/log-download-manager/internal/pkg/server/http-log-manager"
	"github.com/nalej/log-download-manager/internal/pkg/server/log-manager"
	"github.com/nalej/log-download-manager/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

// Service structure with the configuration and the gRPC server.
type Service struct {
	Configuration Config
	OpeCache      *utils.DownloadCache
	Router        *mux.Router
}

// NewService creates a new log download manager service.
func NewService(conf Config) *Service {
	return &Service{
		Configuration: conf,
		Router:        mux.NewRouter(),
		OpeCache:      utils.NewDownloadCache(),
	}
}

// Clients structure with the gRPC clients for remote services.
type Clients struct {
	AppManagerClient grpc_application_manager_go.UnifiedLoggingClient
}

func (s *Service) GetClients() (*Clients, derrors.Error) {

	appManagerConn, err := grpc.Dial(s.Configuration.ApplicationsManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with application managerr")
	}

	appManagerClient := grpc_application_manager_go.NewUnifiedLoggingClient(appManagerConn)

	return &Clients{appManagerClient}, nil
}
func (s *Service) LaunchGRPC() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// Clients
	clients, cErr := s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("Cannot create clients")
	}

	// Create handlers
	appManager := log_manager.NewManager(clients.AppManagerClient, s.OpeCache)
	appHandler := log_manager.NewHandler(appManager)

	grpcServer := grpc.NewServer()
	grpc_log_download_manager_go.RegisterLogDownloadManagerServer(grpcServer, appHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

func (s *Service) LaunchHTTP() error {
	log.Debug().Msg("Server started on port 8941")

	httpLogManager := http_log_manager.NewManager(s.OpeCache)
	httpLogHandler := http_log_manager.NewHandler(httpLogManager)

	s.Router.PathPrefix("/logs/download/").Handler(httpLogHandler.DownloadFile("/logs/download/", http.FileServer(http.Dir(s.Configuration.DownloadPath))))

	log.Info().Int("port", s.Configuration.HttpPort).Msg("Launching Http server")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Configuration.HttpPort), s.Router); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

func (s *Service) Run() error {
	// Configuration
	cErr := s.Configuration.Validate()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("invalid configuration")
	}
	s.Configuration.Print()

	go s.LaunchGRPC()
	return s.LaunchHTTP()
}
