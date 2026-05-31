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

package runtime

import (
	"context"
	"errors"
	"time"

	"github.com/thunder-id/thunderid/internal/enginebridge"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/runtime"
)

type runtimeStoreWrapper struct {
	inner enginebridge.RuntimeStore
}

func (w *runtimeStoreWrapper) StoreFlowContext(
	ctx context.Context, executionID string, data []byte, expiry time.Time,
) error {
	return w.inner.StoreFlowContext(ctx, executionID, data, expiry)
}

func (w *runtimeStoreWrapper) GetFlowContext(ctx context.Context, executionID string) ([]byte, error) {
	data, err := w.inner.GetFlowContext(ctx, executionID)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return nil, runtime.ErrNotFound
	}
	return data, err
}

func (w *runtimeStoreWrapper) UpdateFlowContext(ctx context.Context, executionID string, data []byte) error {
	return w.inner.UpdateFlowContext(ctx, executionID, data)
}

func (w *runtimeStoreWrapper) DeleteFlowContext(ctx context.Context, executionID string) error {
	return w.inner.DeleteFlowContext(ctx, executionID)
}

func (w *runtimeStoreWrapper) StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error {
	return w.inner.StoreAuthCode(ctx, code, data, expiry)
}

func (w *runtimeStoreWrapper) GetAuthCode(ctx context.Context, code string) ([]byte, error) {
	data, err := w.inner.GetAuthCode(ctx, code)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return nil, runtime.ErrNotFound
	}
	return data, err
}

func (w *runtimeStoreWrapper) DeleteAuthCode(ctx context.Context, code string) error {
	err := w.inner.DeleteAuthCode(ctx, code)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return runtime.ErrNotFound
	}
	return err
}

func (w *runtimeStoreWrapper) StoreAuthRequest(
	ctx context.Context, requestID string, data []byte, expiry time.Time,
) error {
	return w.inner.StoreAuthRequest(ctx, requestID, data, expiry)
}

func (w *runtimeStoreWrapper) GetAuthRequest(ctx context.Context, requestID string) ([]byte, error) {
	data, err := w.inner.GetAuthRequest(ctx, requestID)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return nil, runtime.ErrNotFound
	}
	return data, err
}

func (w *runtimeStoreWrapper) DeleteAuthRequest(ctx context.Context, requestID string) error {
	return w.inner.DeleteAuthRequest(ctx, requestID)
}

func (w *runtimeStoreWrapper) StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error {
	return w.inner.StorePAR(ctx, requestURI, data, expiry)
}

func (w *runtimeStoreWrapper) GetPAR(ctx context.Context, requestURI string) ([]byte, error) {
	data, err := w.inner.GetPAR(ctx, requestURI)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return nil, runtime.ErrNotFound
	}
	return data, err
}

func (w *runtimeStoreWrapper) DeletePAR(ctx context.Context, requestURI string) error {
	err := w.inner.DeletePAR(ctx, requestURI)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return runtime.ErrNotFound
	}
	return err
}

func (w *runtimeStoreWrapper) StoreJTI(ctx context.Context, jti string, expiry time.Time) error {
	return w.inner.StoreJTI(ctx, jti, expiry)
}

func (w *runtimeStoreWrapper) ExistsJTI(ctx context.Context, jti string) (bool, error) {
	return w.inner.ExistsJTI(ctx, jti)
}

func (w *runtimeStoreWrapper) StoreAttributeCache(ctx context.Context, id string, data []byte, expiry time.Time) error {
	return w.inner.StoreAttributeCache(ctx, id, data, expiry)
}

func (w *runtimeStoreWrapper) GetAttributeCache(ctx context.Context, id string) ([]byte, error) {
	data, err := w.inner.GetAttributeCache(ctx, id)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return nil, runtime.ErrNotFound
	}
	return data, err
}

func (w *runtimeStoreWrapper) ExtendAttributeCacheExpiry(ctx context.Context, id string, expiry time.Time) error {
	err := w.inner.ExtendAttributeCacheExpiry(ctx, id, expiry)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return runtime.ErrNotFound
	}
	return err
}

func (w *runtimeStoreWrapper) DeleteAttributeCache(ctx context.Context, id string) error {
	err := w.inner.DeleteAttributeCache(ctx, id)
	if errors.Is(err, enginebridge.ErrNotFound) {
		return runtime.ErrNotFound
	}
	return err
}
