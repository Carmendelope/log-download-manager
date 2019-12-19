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

package interceptor

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/derrors"
	"net/http"
)


const (
	UserID = "user-id"
	OrganizationID = "organization-id"
)

var permission = []string{"ORG", "APPS"}

// PersonalClaim is the claim that include system information.
type PersonalClaim struct {
	UserID         string   `json:"userID,omitempty"`
	Primitives     []string `json:"access,omitempty"`
	RoleName       string   `json:"role,omitempty"`
	OrganizationID string   `json:"organizationID,omitempty"`
}

// Claim joins the personal claim and the standard JWT claim.
type Claim struct {
	jwt.StandardClaims
	PersonalClaim
}

type Interceptor struct {
	secret     string
	authHeader string
}

func NewInterceptor(secret string, authHeader string) *Interceptor {
	return &Interceptor{secret: secret, authHeader:authHeader}
}

func (i *Interceptor) Validate(r *http.Request) derrors.Error {

	vErr := i.checkJWT(r)
	if vErr != nil {
		return vErr
	}

	return nil
}

func (i *Interceptor) checkJWT(r *http.Request) derrors.Error {

	// extract Token
	token := r.Header.Get(i.authHeader)
	if token == "" {
		return derrors.NewUnauthenticatedError("token is not supplied")
	}

	// check Token
	// validateToken function validates the token
	tk, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(i.secret), nil
	})

	if err != nil {
		return derrors.NewUnauthenticatedError("token is not valid", err)
	}

	userAuth := tk.Claims.(*Claim)

	r.Header.Add(UserID, userAuth.UserID)
	r.Header.Add(OrganizationID, userAuth.OrganizationID)

	// check the role permissions
	found := false
	for i:=0; i<len(permission) && !found; i++ {
		for j:=0; j< len(userAuth.Primitives) && !found; j++ {
			if permission[i] == userAuth.Primitives[j] {
				found = true
			}
		}
	}

	if ! found {
		return derrors.NewUnauthenticatedError("unauthorized method")
	}

	return nil
}

