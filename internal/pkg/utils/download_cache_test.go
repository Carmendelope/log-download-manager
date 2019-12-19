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
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Utils", func() {

	var downloadCache  *DownloadCache
	var organizationID string

	ginkgo.BeforeSuite(func() {
		downloadCache = NewDownloadCache("/test/", "nalej.tech")
	})

	ginkgo.BeforeEach(func() {
		downloadCache.Clean()
		organizationID = uuid.New().String()
	})

	ginkgo.Context("Adding new operation", func() {
		ginkgo.It("should be able to add a new one", func() {
			requestID := uuid.New().String()
			_, err := downloadCache.Add(organizationID, requestID, 0, 0)
			gomega.Expect(err).To(gomega.Succeed())

			ope, err := downloadCache.Get(requestID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ope.RequestId).Should(gomega.Equal(requestID))


		})
		ginkgo.It("should not be able to add one operation twice", func() {
			requestID := uuid.New().String()
			_, err := downloadCache.Add(organizationID, requestID, 0, 0)
			gomega.Expect(err).To(gomega.Succeed())

			_, err = downloadCache.Add(organizationID, requestID, 0, 0)
			gomega.Expect(err).NotTo(gomega.Succeed())



		})
	})
	ginkgo.Context("Removing an operation", func() {
		ginkgo.It("should be able to remove an operation", func() {
			requestID := uuid.New().String()
			_, err := downloadCache.Add(organizationID, requestID, 0, 0)
			gomega.Expect(err).To(gomega.Succeed())

			err = downloadCache.Remove(requestID)
			gomega.Expect(err).To(gomega.Succeed())


		})
		ginkgo.It("should not be able to remove one operation if it does not exists", func() {
			requestID := uuid.New().String()
			err := downloadCache.Remove(requestID)
			gomega.Expect(err).NotTo(gomega.Succeed())



		})
	})
	ginkgo.Context("Updating an operation", func() {
		ginkgo.It("should be able to update an operation", func() {
			requestID := uuid.New().String()
			_, err := downloadCache.Add(organizationID, requestID, 0, 0)
			gomega.Expect(err).To(gomega.Succeed())

			err = downloadCache.Update(requestID, Generating, "")
			gomega.Expect(err).To(gomega.Succeed())

			ope, err := downloadCache.Get(requestID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ope.RequestId).Should(gomega.Equal(requestID))
			gomega.Expect(ope.State).Should(gomega.Equal(Generating))


		})
		ginkgo.It("should not be able to update one operation if it does not exists", func() {
			requestID := uuid.New().String()
			err := downloadCache.Update(requestID, Generating, "")
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
	ginkgo.Context("Listing operations", func() {
		ginkgo.It("should be able to list operations", func() {
			num := 5
			for i:= 0; i < num; i++ {
				_, err := downloadCache.Add(organizationID, uuid.New().String(), 0, 0)
				gomega.Expect(err).To(gomega.Succeed())
			}

			list, err := downloadCache.List(organizationID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(num))

		})
	})
	ginkgo.Context("Listing operations", func() {
		ginkgo.It("should be able to list an empty list of operations", func() {
			num := 5
			for i:= 0; i < num; i++ {
				_, err := downloadCache.Add(organizationID, uuid.New().String(), 0, 0)
				gomega.Expect(err).To(gomega.Succeed())
			}

			list, err := downloadCache.List(uuid.New().String())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))

		})
	})
})