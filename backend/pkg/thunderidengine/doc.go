// Package thunderidengine defines the public embeddable contract for OAuth 2 / OIDC and flow
// execution: host provider interfaces, shared models, OAuth runtime Options, and flow
// ExecutorInterface.
//
// Runtime wiring (Engine, HTTP routes) is not part of this package yet; hosts implement the
// types declared here and will connect them to the engine in a follow-up change.
package thunderidengine
