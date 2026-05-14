// Package oauthdeps holds the dependency contracts used to wire OAuth.
// It lives under pkg/oauth/deps as a separate package so internal/oauth can
// import these types without an import cycle with pkg/oauth.
//
// Application (DCRApplication) and Inbound (InboundOAuth) are pkg/oauth/host
// contracts so embedders outside this module can supply implementations without
// importing internal packages. Thunder defaults are constructed via
// internal/oauth/hostbridge (for example NewThunderApplication and NewThunderInbound).
//
// Other dependency fields remain aliases to Thunder internal interfaces until
// they are narrowed similarly.
//
// Transactioner must be supplied on Dependencies: internal/oauth.Initialize no
// longer fetches it from the global DB provider (embedders can pass any
// implementation). Other subpackages may still read config/DB globals until
// further refactors.
package oauthdeps
