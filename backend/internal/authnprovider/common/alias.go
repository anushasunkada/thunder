package common

import pkg "github.com/thunder-id/thunderid/pkg/authnprovider"

type AuthnMetadata = pkg.AuthnMetadata
type AuthnResult = pkg.AuthnResult
type GetAttributesMetadata = pkg.GetAttributesMetadata
type GetAttributesResult = pkg.GetAttributesResult
type RequestedAttributes = pkg.RequestedAttributes
type AttributesResponse = pkg.AttributesResponse
type AttributeResponse = pkg.AttributeResponse
type AssuranceMetadataRequest = pkg.AssuranceMetadataRequest
type AssuranceMetadataResponse = pkg.AssuranceMetadataResponse
type GenericMetadataRequest = pkg.GenericMetadataRequest
type GenericTimeMetadataRequest = pkg.GenericTimeMetadataRequest
type VerificationRequest = pkg.VerificationRequest
type VerificationResponse = pkg.VerificationResponse
type AttributeMetadataRequest = pkg.AttributeMetadataRequest

const (
	ErrorCodeSystemError          = pkg.ErrorCodeSystemError
	ErrorCodeAuthenticationFailed = pkg.ErrorCodeAuthenticationFailed
	ErrorCodeUserNotFound         = pkg.ErrorCodeUserNotFound
	ErrorCodeInvalidToken         = pkg.ErrorCodeInvalidToken
	ErrorCodeNotImplemented       = pkg.ErrorCodeNotImplemented
	ErrorCodeInvalidRequest       = pkg.ErrorCodeInvalidRequest
)
