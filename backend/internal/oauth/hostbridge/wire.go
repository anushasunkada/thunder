/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import "github.com/thunder-id/thunderid/pkg/oauth/host"

// DCRPartnerFromApplication selects an efficient DCR partner implementation for the configured host contract.
func DCRPartnerFromApplication(h host.DCRApplication) DCRPartner {
	if h == nil {
		return nil
	}
	if th, ok := h.(*ThunderDCRApplication); ok {
		return &thunderDCRPartner{inner: th.ApplicationService}
	}
	return &hostDCRPartner{host: h}
}
