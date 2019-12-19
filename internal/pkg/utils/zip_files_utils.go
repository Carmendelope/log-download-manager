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
	"archive/zip"
	"fmt"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

const logWithMetadata = "[%s][DESCRIPTOR-%s][INSTANCE-%s][SERVICE-GROUP-%s][SERVICE-%s]:%s\n"

////////////////////
// InitializeFile create the file and write the header
func InitializeFile(target string, includeMetadata bool) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Debug().Msg("file created successfully")
	return nil
}

func AppendResponses(responses []*grpc_application_manager_go.LogEntryResponse, target string, includeMetadata bool) error {
	log.Debug().Bool("includeMedatada", includeMetadata).Msg("AppendResponse")
	f, err := os.OpenFile(target, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	var writeErr error
	for _, response := range responses {
		if includeMetadata{
			_, writeErr = f.WriteString(fmt.Sprintf(logWithMetadata, time.Unix(0, response.Timestamp), response.AppDescriptorName, response.AppInstanceName,
				response.ServiceGroupName, response.ServiceName, response.Msg))
		}else {
			_, writeErr = f.WriteString(fmt.Sprintf("%s\n", response.Msg))
		}
		if writeErr != nil {
			f.Close()
			return writeErr
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Debug().Msg("file appended successfully")
	return nil
}

func RemoveFile (path string) error {
	return os.Remove(path)
}
// ZipFiles compresses one or many files into a single zip archive file.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	//header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func GetFilePath(filesDirectory string, requestId string) string {
	return fmt.Sprintf("%s%s.file", filesDirectory, requestId)
}
func GetZipFilePath(filesDirectory string, requestId string) string {
	return fmt.Sprintf("%s%s.zip", filesDirectory, requestId)
}