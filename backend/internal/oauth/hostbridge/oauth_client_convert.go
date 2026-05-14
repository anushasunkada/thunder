/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import (
	"github.com/thunder-id/thunderid/internal/cert"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	oauth2const "github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
	"github.com/thunder-id/thunderid/pkg/oauth/oauthclient"
)

func oauthClientFromInbound(o *inboundmodel.OAuthClient) *oauthclient.Client {
	if o == nil {
		return nil
	}
	grants := make([]string, 0, len(o.GrantTypes))
	for _, g := range o.GrantTypes {
		grants = append(grants, string(g))
	}
	responses := make([]string, 0, len(o.ResponseTypes))
	for _, r := range o.ResponseTypes {
		responses = append(responses, string(r))
	}
	c := &oauthclient.Client{
		ID:                                 o.ID,
		OUID:                               o.OUID,
		ClientID:                           o.ClientID,
		RedirectURIs:                       append([]string(nil), o.RedirectURIs...),
		GrantTypes:                         grants,
		ResponseTypes:                      responses,
		TokenEndpointAuthMethod:            string(o.TokenEndpointAuthMethod),
		PKCERequired:                       o.PKCERequired,
		PublicClient:                       o.PublicClient,
		RequirePushedAuthorizationRequests: o.RequirePushedAuthorizationRequests,
		Scopes:                             append([]string(nil), o.Scopes...),
		ScopeClaims:                        copyScopeClaims(o.ScopeClaims),
		AcrValues:                          append([]string(nil), o.AcrValues...),
	}
	if o.Token != nil {
		c.Token = tokenClientFromInbound(o.Token)
	}
	if o.UserInfo != nil {
		c.UserInfo = userInfoClientFromInbound(o.UserInfo)
	}
	if o.Certificate != nil {
		c.Certificate = &oauthclient.Certificate{Type: string(o.Certificate.Type), Value: o.Certificate.Value}
	}
	return c
}

func oauthClientToInbound(c *oauthclient.Client) (*inboundmodel.OAuthClient, error) {
	if c == nil {
		return nil, nil
	}
	o := &inboundmodel.OAuthClient{
		ID:                                 c.ID,
		OUID:                               c.OUID,
		ClientID:                           c.ClientID,
		RedirectURIs:                       append([]string(nil), c.RedirectURIs...),
		GrantTypes:                         parseGrantTypes(c.GrantTypes),
		ResponseTypes:                      parseResponseTypes(c.ResponseTypes),
		TokenEndpointAuthMethod:            oauth2const.TokenEndpointAuthMethod(c.TokenEndpointAuthMethod),
		PKCERequired:                       c.PKCERequired,
		PublicClient:                       c.PublicClient,
		RequirePushedAuthorizationRequests: c.RequirePushedAuthorizationRequests,
		Scopes:                             append([]string(nil), c.Scopes...),
		ScopeClaims:                        copyScopeClaims(c.ScopeClaims),
		AcrValues:                          append([]string(nil), c.AcrValues...),
	}
	if c.Token != nil {
		o.Token = tokenClientToInbound(c.Token)
	}
	if c.UserInfo != nil {
		o.UserInfo = userInfoClientToInbound(c.UserInfo)
	}
	if c.Certificate != nil {
		o.Certificate = &inboundmodel.Certificate{
			Type:  cert.CertificateType(c.Certificate.Type),
			Value: c.Certificate.Value,
		}
	}
	return o, nil
}

func tokenClientFromInbound(t *inboundmodel.OAuthTokenConfig) *oauthclient.OAuthTokenConfig {
	if t == nil {
		return nil
	}
	out := &oauthclient.OAuthTokenConfig{}
	if t.AccessToken != nil {
		out.AccessToken = &oauthclient.AccessTokenConfig{
			ValidityPeriod: t.AccessToken.ValidityPeriod,
			UserAttributes: append([]string(nil), t.AccessToken.UserAttributes...),
		}
	}
	if t.IDToken != nil {
		out.IDToken = &oauthclient.IDTokenConfig{
			ValidityPeriod: t.IDToken.ValidityPeriod,
			UserAttributes: append([]string(nil), t.IDToken.UserAttributes...),
			ResponseType:   string(t.IDToken.ResponseType),
			EncryptionAlg:  t.IDToken.EncryptionAlg,
			EncryptionEnc:  t.IDToken.EncryptionEnc,
		}
	}
	return out
}

func tokenClientToInbound(t *oauthclient.OAuthTokenConfig) *inboundmodel.OAuthTokenConfig {
	if t == nil {
		return nil
	}
	out := &inboundmodel.OAuthTokenConfig{}
	if t.AccessToken != nil {
		out.AccessToken = &inboundmodel.AccessTokenConfig{
			ValidityPeriod: t.AccessToken.ValidityPeriod,
			UserAttributes: append([]string(nil), t.AccessToken.UserAttributes...),
		}
	}
	if t.IDToken != nil {
		out.IDToken = &inboundmodel.IDTokenConfig{
			ValidityPeriod: t.IDToken.ValidityPeriod,
			UserAttributes: append([]string(nil), t.IDToken.UserAttributes...),
			ResponseType:   inboundmodel.IDTokenResponseType(t.IDToken.ResponseType),
			EncryptionAlg:  t.IDToken.EncryptionAlg,
			EncryptionEnc:  t.IDToken.EncryptionEnc,
		}
	}
	return out
}

func userInfoClientFromInbound(u *inboundmodel.UserInfoConfig) *oauthclient.UserInfoConfig {
	if u == nil {
		return nil
	}
	return &oauthclient.UserInfoConfig{
		ResponseType:   string(u.ResponseType),
		UserAttributes: append([]string(nil), u.UserAttributes...),
		SigningAlg:     u.SigningAlg,
		EncryptionAlg:  u.EncryptionAlg,
		EncryptionEnc:  u.EncryptionEnc,
	}
}

func userInfoClientToInbound(u *oauthclient.UserInfoConfig) *inboundmodel.UserInfoConfig {
	if u == nil {
		return nil
	}
	return &inboundmodel.UserInfoConfig{
		ResponseType:   inboundmodel.UserInfoResponseType(u.ResponseType),
		UserAttributes: append([]string(nil), u.UserAttributes...),
		SigningAlg:     u.SigningAlg,
		EncryptionAlg:  u.EncryptionAlg,
		EncryptionEnc:  u.EncryptionEnc,
	}
}
