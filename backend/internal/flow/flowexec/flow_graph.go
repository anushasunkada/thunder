/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
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

package flowexec

import (
	"errors"

	"github.com/thunder-id/thunderid/pkg/flowgraphbridge"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

func coreGraphFromProvider(fg thunderidengine.FlowGraph) (core.GraphInterface, *serviceerror.ServiceError) {
	if fg == nil {
		return nil, &serviceerror.InternalServerError
	}
	graph, ok := flowgraphbridge.CoreGraphFromFlowGraph(fg)
	if !ok {
		return nil, &serviceerror.InternalServerError
	}
	return graph, nil
}

func serviceErrorFromFlowGraph(err error) *serviceerror.ServiceError {
	if err == nil {
		return nil
	}
	if svcErr, ok := flowgraphbridge.ServiceErrorFromErr(err); ok {
		return svcErr
	}
	if errors.Is(err, thunderidengine.ErrFlowNotFound) {
		return &serviceerror.InternalServerError
	}
	return &serviceerror.InternalServerError
}
