// Package flowdeps holds the dependency contracts used to wire flow execution.
// It lives under pkg/flow/deps as a separate package so internal/flow/flowexec
// can import these types without an import cycle with pkg/flow.
//
// Types are currently aliases to Thunder internal service interfaces so the
// server and pkg/flow share one dependency surface. Narrowing these to
// pkg-local interfaces (with public models) is a follow-up for embedders
// outside this module.
//
// internal/flow/core does not consume these aliases in its Initialize signature:
// flowdeps imports internal/flow/executor, which imports core, so core cannot
// import flowdeps without a cycle.
package flowdeps
