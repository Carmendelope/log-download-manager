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
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	DefaultTimeout  = time.Minute
	UserID = "user-id"
)

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

func GetUserFromContext(ctx context.Context) string {
	userID := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		user, found := md[UserID]
		if found {
			userID = user[0]
		}
	}
	log.Debug().Str("userID", userID).Msg("user identifier")

	return userID
}