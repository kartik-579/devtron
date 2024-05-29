/*
 * Copyright (c) 2024. Devtron Inc.
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

package main

import (
	"github.com/devtron-labs/authenticator/client"
	"github.com/devtron-labs/authenticator/middleware"
	"github.com/devtron-labs/devtron/api/apiToken"
	"github.com/google/wire"
)

// AuthWireSet:	 set of components used to initialise authentication with dex
var AuthWireSet = wire.NewSet(
	client.GetRuntimeConfig,
	client.NewK8sClient,
	client.BuildDexConfig,
	client.GetSettings,
	apiToken.ApiTokenSecretWireSet,
	middleware.NewSessionManager,
	middleware.NewUserLogin,
)
