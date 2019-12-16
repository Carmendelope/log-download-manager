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
	"github.com/nalej/grpc-application-manager-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"time"
)

const testDir = "./testDir/"
var _ = ginkgo.Describe("Zip Files", func() {


	ginkgo.BeforeEach(func() {
		// create tmp directory
		err := os.MkdirAll(testDir, os.ModePerm)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.AfterEach(func() {
		// create tmp directory
		err := os.RemoveAll(testDir)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.Context("Adding new operation", func() {

		ginkgo.It("should be create a file", func() {

			path := fmt.Sprintf("%stestResponse.file",testDir)
			err := InitializeFile(path, false)
			gomega.Expect(err).To(gomega.Succeed())
			responses := []*grpc_application_manager_go.LogEntryResponse {
				{Msg:"entry 1", Timestamp:time.Now().UnixNano()},
				{Msg:"entry 2", Timestamp:time.Now().UnixNano()},
				{Msg:"entry 3", Timestamp:time.Now().UnixNano()},
			}
			// source, target
			AppendResponses(responses, path)
		})
		ginkgo.It("should be create zip file", func() {

			path := fmt.Sprintf("%stestResponse.file",testDir)
			err := InitializeFile(path, false)
			gomega.Expect(err).To(gomega.Succeed())

			responses := []*grpc_application_manager_go.LogEntryResponse {
				{Msg:"entry 1", Timestamp:time.Now().UnixNano()},
				{Msg:"entry 2", Timestamp:time.Now().UnixNano()},
				{Msg:"entry 3", Timestamp:time.Now().UnixNano()},
			}
			// source, target
			err = AppendResponses(responses, path)
			gomega.Expect(err).To(gomega.Succeed())

			err = ZipFiles(fmt.Sprintf("%stest.zip",testDir), []string{path})
			gomega.Expect(err).To(gomega.Succeed())
		})
	})

})
