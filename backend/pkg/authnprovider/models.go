package authnprovider

// AuthnMetadata contains metadata for authentication.
type AuthnMetadata struct {
	AppMetadata map[string]interface{} `json:"appMetadata,omitempty"`
}

// AuthnResult represents the result of an authentication attempt.
type AuthnResult struct {
	EntityID       string `json:"entityId"`
	EntityCategory string `json:"entityCategory"`
	EntityType     string `json:"entityType"`
	OUID           string `json:"ouId"`

	UserID   string `json:"userId"`
	UserType string `json:"userType"`

	Token                     string              `json:"token"`
	IsAttributeValuesIncluded bool                `json:"isAttributeValuesIncluded"`
	AttributesResponse        *AttributesResponse `json:"attributesResponse,omitempty"`

	ExternalSub     string                 `json:"externalSub,omitempty"`
	ExternalClaims  map[string]interface{} `json:"externalClaims,omitempty"`
	IsExistingUser  bool                   `json:"isExistingUser"`
	IsAmbiguousUser bool                   `json:"isAmbiguousUser"`
}

// GetAttributesMetadata contains metadata for fetching attributes.
type GetAttributesMetadata struct {
	AppMetadata map[string]interface{} `json:"appMetadata,omitempty"`
	Locale      string                 `json:"locale"`
}

// GetAttributesResult represents the result of fetching attributes.
type GetAttributesResult struct {
	EntityID       string `json:"entityId"`
	EntityCategory string `json:"entityCategory"`
	EntityType     string `json:"entityType"`
	OUID           string `json:"ouId"`

	UserID   string `json:"userId"`
	UserType string `json:"userType"`

	AttributesResponse *AttributesResponse `json:"attributeResponse,omitempty"`
}

// AssuranceMetadataResponse contains assurance metadata for an attribute.
type AssuranceMetadataResponse struct {
	IsVerified       bool   `json:"isVerified"`
	VerificationID   string `json:"verificationId,omitempty"`
}

// VerificationResponse contains verification details for an attribute.
type VerificationResponse struct {
	TrustFramework      string `json:"trustFramework,omitempty"`
	Time                string `json:"time,omitempty"`
	VerificationProcess string `json:"verificationProcess,omitempty"`
}

// RequestedAttributes contains the requested attributes and verifications.
type RequestedAttributes struct {
	Attributes    map[string]*AttributeMetadataRequest `json:"attributes,omitempty"`
	Verifications map[string]*VerificationRequest      `json:"verifications,omitempty"`
}

// AttributeMetadataRequest contains metadata request details for an attribute.
type AttributeMetadataRequest struct {
	GenericMetadataRequest   *GenericMetadataRequest   `json:"genericMetadataRequest,omitempty"`
	AssuranceMetadataRequest *AssuranceMetadataRequest `json:"assuranceMetadataRequest,omitempty"`
}

// GenericMetadataRequest contains generic metadata request details.
type GenericMetadataRequest struct {
	Essential bool     `json:"essential,omitempty"`
	Value     string   `json:"value,omitempty"`
	Values    []string `json:"values,omitempty"`
}

// GenericTimeMetadataRequest extends GenericMetadataRequest with time-related metadata.
type GenericTimeMetadataRequest struct {
	GenericMetadataRequest
	MaxAge *int `json:"maxAge,omitempty"`
}

// AssuranceMetadataRequest contains assurance metadata request details.
type AssuranceMetadataRequest struct {
	ShouldVerify   bool   `json:"shouldVerify,omitempty"`
	VerificationID string `json:"verificationId,omitempty"`
}

// VerificationRequest contains verification request details.
type VerificationRequest struct {
	TrustFramework      *GenericMetadataRequest     `json:"trustFramework,omitempty"`
	VerificationProcess *GenericMetadataRequest     `json:"verificationProcess,omitempty"`
	Time                *GenericTimeMetadataRequest `json:"time,omitempty"`
}

// AttributesResponse contains the response with attributes and verifications.
type AttributesResponse struct {
	Attributes    map[string]*AttributeResponse    `json:"attributes,omitempty"`
	Verifications map[string]*VerificationResponse `json:"verifications,omitempty"`
}

// AttributeResponse contains the response for an attribute with its value and assurance metadata.
type AttributeResponse struct {
	Value                     interface{}                `json:"value,omitempty"`
	AssuranceMetadataResponse *AssuranceMetadataResponse `json:"assuranceMetadataResponse,omitempty"`
}
