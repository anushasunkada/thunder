/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import (
	"fmt"

	"github.com/thunder-id/thunderid/internal/application/model"
	"github.com/thunder-id/thunderid/internal/cert"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	oauth2const "github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
	"github.com/thunder-id/thunderid/pkg/oauth/app"
)

func applicationDTOFromAppCreate(in *app.ApplicationCreate) (*model.ApplicationDTO, error) {
	if in == nil {
		return nil, fmt.Errorf("application create payload is required")
	}
	oauthCfg, err := oauthConfigWithSecretFromRegistration(&in.OAuth)
	if err != nil {
		return nil, err
	}
	dto := &model.ApplicationDTO{
		ID:        in.ID,
		OUID:      in.OUID,
		Name:      in.Name,
		URL:       in.URL,
		LogoURL:   in.LogoURL,
		TosURI:    in.TosURI,
		PolicyURI: in.PolicyURI,
		Contacts:  append([]string(nil), in.Contacts...),
		InboundAuthConfig: []inboundmodel.InboundAuthConfigWithSecret{
			{Type: inboundmodel.OAuthInboundAuthType, OAuthConfig: oauthCfg},
		},
	}
	if in.Certificate != nil {
		dto.Certificate = &inboundmodel.Certificate{
			Type:  cert.CertificateType(in.Certificate.Type),
			Value: in.Certificate.Value,
		}
	}
	return dto, nil
}

func applicationCreatedFromDTO(dto *model.ApplicationDTO) *app.ApplicationCreated {
	if dto == nil {
		return nil
	}
	c := applicationCreateFromDTO(dto)
	if c == nil {
		return nil
	}
	out := &app.ApplicationCreated{
		ID:          c.ID,
		OUID:        c.OUID,
		Name:        c.Name,
		URL:         c.URL,
		LogoURL:     c.LogoURL,
		TosURI:      c.TosURI,
		PolicyURI:   c.PolicyURI,
		Contacts:    append([]string(nil), c.Contacts...),
		Certificate: copyAppCertificate(c.Certificate),
		OAuth:       oauthRegistrationFromOAuthConfig(dto.InboundAuthConfig[0].OAuthConfig),
	}
	return out
}

func applicationCreateFromDTO(dto *model.ApplicationDTO) *app.ApplicationCreate {
	if dto == nil || len(dto.InboundAuthConfig) == 0 || dto.InboundAuthConfig[0].OAuthConfig == nil {
		return nil
	}
	oc := dto.InboundAuthConfig[0].OAuthConfig
	reg := oauthRegistrationFromOAuthConfig(oc)
	out := &app.ApplicationCreate{
		ID:        dto.ID,
		OUID:      dto.OUID,
		Name:      dto.Name,
		URL:       dto.URL,
		LogoURL:   dto.LogoURL,
		TosURI:    dto.TosURI,
		PolicyURI: dto.PolicyURI,
		Contacts:  append([]string(nil), dto.Contacts...),
		OAuth:     reg,
	}
	if dto.Certificate != nil {
		out.Certificate = &app.Certificate{
			Type:  string(dto.Certificate.Type),
			Value: dto.Certificate.Value,
		}
	}
	return out
}

func applicationCreatedToDTO(res *app.ApplicationCreated) (*model.ApplicationDTO, error) {
	if res == nil {
		return nil, fmt.Errorf("application created payload is required")
	}
	in := &app.ApplicationCreate{
		ID:          res.ID,
		OUID:        res.OUID,
		Name:        res.Name,
		URL:         res.URL,
		LogoURL:     res.LogoURL,
		TosURI:      res.TosURI,
		PolicyURI:   res.PolicyURI,
		Contacts:    append([]string(nil), res.Contacts...),
		Certificate: copyAppCertificate(res.Certificate),
		OAuth:       res.OAuth,
	}
	return applicationDTOFromAppCreate(in)
}

func copyAppCertificate(c *app.Certificate) *app.Certificate {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
}

func oauthRegistrationFromOAuthConfig(oc *inboundmodel.OAuthConfigWithSecret) app.OAuthRegistration {
	if oc == nil {
		return app.OAuthRegistration{}
	}
	grants := make([]string, 0, len(oc.GrantTypes))
	for _, g := range oc.GrantTypes {
		grants = append(grants, string(g))
	}
	responses := make([]string, 0, len(oc.ResponseTypes))
	for _, r := range oc.ResponseTypes {
		responses = append(responses, string(r))
	}
	reg := app.OAuthRegistration{
		ClientID:                           oc.ClientID,
		ClientSecret:                       oc.ClientSecret,
		RedirectURIs:                       append([]string(nil), oc.RedirectURIs...),
		GrantTypes:                         grants,
		ResponseTypes:                      responses,
		TokenEndpointAuthMethod:            string(oc.TokenEndpointAuthMethod),
		PKCERequired:                       oc.PKCERequired,
		PublicClient:                       oc.PublicClient,
		RequirePushedAuthorizationRequests: oc.RequirePushedAuthorizationRequests,
		Scopes:                             append([]string(nil), oc.Scopes...),
		ScopeClaims:                        copyScopeClaims(oc.ScopeClaims),
		AcrValues:                          append([]string(nil), oc.AcrValues...),
	}
	if oc.Token != nil {
		reg.Token = tokenConfigFromInbound(oc.Token)
	}
	if oc.UserInfo != nil {
		reg.UserInfo = userInfoConfigFromInbound(oc.UserInfo)
	}
	if oc.Certificate != nil {
		reg.Certificate = &app.Certificate{Type: string(oc.Certificate.Type), Value: oc.Certificate.Value}
	}
	return reg
}

func oauthConfigWithSecretFromRegistration(reg *app.OAuthRegistration) (*inboundmodel.OAuthConfigWithSecret, error) {
	if reg == nil {
		return nil, fmt.Errorf("oauth registration is required")
	}
	oc := &inboundmodel.OAuthConfigWithSecret{
		ClientID:                           reg.ClientID,
		ClientSecret:                       reg.ClientSecret,
		RedirectURIs:                       append([]string(nil), reg.RedirectURIs...),
		GrantTypes:                         parseGrantTypes(reg.GrantTypes),
		ResponseTypes:                      parseResponseTypes(reg.ResponseTypes),
		TokenEndpointAuthMethod:            oauth2const.TokenEndpointAuthMethod(reg.TokenEndpointAuthMethod),
		PKCERequired:                       reg.PKCERequired,
		PublicClient:                       reg.PublicClient,
		RequirePushedAuthorizationRequests: reg.RequirePushedAuthorizationRequests,
		Scopes:                             append([]string(nil), reg.Scopes...),
		ScopeClaims:                        copyScopeClaims(reg.ScopeClaims),
		AcrValues:                          append([]string(nil), reg.AcrValues...),
	}
	if reg.Token != nil {
		oc.Token = tokenConfigToInbound(reg.Token)
	}
	if reg.UserInfo != nil {
		oc.UserInfo = userInfoConfigToInbound(reg.UserInfo)
	}
	if reg.Certificate != nil {
		oc.Certificate = &inboundmodel.Certificate{
			Type:  cert.CertificateType(reg.Certificate.Type),
			Value: reg.Certificate.Value,
		}
	}
	return oc, nil
}

func parseGrantTypes(in []string) []oauth2const.GrantType {
	out := make([]oauth2const.GrantType, 0, len(in))
	for _, g := range in {
		if g == "" {
			continue
		}
		out = append(out, oauth2const.GrantType(g))
	}
	return out
}

func parseResponseTypes(in []string) []oauth2const.ResponseType {
	out := make([]oauth2const.ResponseType, 0, len(in))
	for _, r := range in {
		if r == "" {
			continue
		}
		out = append(out, oauth2const.ResponseType(r))
	}
	return out
}

func copyScopeClaims(m map[string][]string) map[string][]string {
	if m == nil {
		return nil
	}
	out := make(map[string][]string, len(m))
	for k, v := range m {
		out[k] = append([]string(nil), v...)
	}
	return out
}

func tokenConfigFromInbound(t *inboundmodel.OAuthTokenConfig) *app.OAuthTokenConfig {
	if t == nil {
		return nil
	}
	out := &app.OAuthTokenConfig{}
	if t.AccessToken != nil {
		out.AccessToken = &app.AccessTokenConfig{
			ValidityPeriod: t.AccessToken.ValidityPeriod,
			UserAttributes: append([]string(nil), t.AccessToken.UserAttributes...),
		}
	}
	if t.IDToken != nil {
		out.IDToken = &app.IDTokenConfig{
			ValidityPeriod: t.IDToken.ValidityPeriod,
			UserAttributes: append([]string(nil), t.IDToken.UserAttributes...),
			ResponseType:   string(t.IDToken.ResponseType),
			EncryptionAlg:  t.IDToken.EncryptionAlg,
			EncryptionEnc:  t.IDToken.EncryptionEnc,
		}
	}
	return out
}

func tokenConfigToInbound(t *app.OAuthTokenConfig) *inboundmodel.OAuthTokenConfig {
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

func userInfoConfigFromInbound(u *inboundmodel.UserInfoConfig) *app.UserInfoConfig {
	if u == nil {
		return nil
	}
	return &app.UserInfoConfig{
		ResponseType:   string(u.ResponseType),
		UserAttributes: append([]string(nil), u.UserAttributes...),
		SigningAlg:     u.SigningAlg,
		EncryptionAlg:  u.EncryptionAlg,
		EncryptionEnc:  u.EncryptionEnc,
	}
}

func userInfoConfigToInbound(u *app.UserInfoConfig) *inboundmodel.UserInfoConfig {
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
