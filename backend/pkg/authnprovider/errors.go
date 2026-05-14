package authnprovider

// Error codes returned by AuthnProvider implementations (aligned with Thunder authnprovider common).
const (
	ErrorCodeSystemError          = "AUP-0001"
	ErrorCodeAuthenticationFailed = "AUP-0002"
	ErrorCodeUserNotFound         = "AUP-0003"
	ErrorCodeInvalidToken         = "AUP-0004"
	ErrorCodeNotImplemented       = "AUP-0005"
	ErrorCodeInvalidRequest       = "AUP-0006"
)

// ErrorType classifies a ServiceError as client or server side.
type ErrorType string

const (
	ClientError ErrorType = "client_error"
	ServerError ErrorType = "server_error"
)

// ServiceError is the public error shape for AuthnProvider callbacks.
// Implementations return nil on success.
type ServiceError struct {
	Type        ErrorType `json:"type"`
	Code        string    `json:"code"`
	Message     string    `json:"message"`
	Description string    `json:"description,omitempty"`
}

// NewClientError builds a client-side ServiceError.
func NewClientError(code, message, description string) *ServiceError {
	return &ServiceError{
		Type:        ClientError,
		Code:        code,
		Message:     message,
		Description: description,
	}
}

// NewServerError builds a server-side ServiceError.
func NewServerError(code, message, description string) *ServiceError {
	return &ServiceError{
		Type:        ServerError,
		Code:        code,
		Message:     message,
		Description: description,
	}
}
