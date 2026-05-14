/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package main service registration delegates to internal/serverbootstrap.
package main

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/serverbootstrap"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

func registerServices(mux *http.ServeMux, cacheManager cache.CacheManagerInterface) jwt.JWTServiceInterface {
	return serverbootstrap.RegisterServices(mux, cacheManager, nil)
}

func unregisterServices() {
	serverbootstrap.UnregisterServices()
}
