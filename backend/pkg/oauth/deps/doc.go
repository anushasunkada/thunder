// Package oauthdeps holds the dependency contracts used to wire OAuth.
// It lives under pkg/oauth/deps as a separate package so internal/oauth can
// import these types without an import cycle with pkg/oauth.
//
// Types are currently aliases to Thunder internal service interfaces so the
// server and pkg/oauth share one dependency surface. Narrowing these to
// pkg-local interfaces (with public models) is a follow-up for embedders
// outside this module.
package oauthdeps
