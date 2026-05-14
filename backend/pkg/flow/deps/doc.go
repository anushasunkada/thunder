// Package flowdeps holds the dependency contracts used to wire flow execution.
// It lives under pkg/flow/deps as a separate package so internal/flow/flowexec
// can import these types without an import cycle with pkg/flow.
//
// Inbound is pkg/flow/host.InboundFlow so embedders outside this module can supply
// entity inbound resolution without importing internal/inboundclient. Thunder
// defaults are constructed via internal/flow/hostbridge.NewThunderInboundFlow.
//
// Other dependency fields remain aliases to Thunder internal service interfaces
// until they are narrowed similarly.
//
// internal/flow/core does not consume these aliases in its Initialize signature:
// flowdeps imports internal/flow/executor, which imports core, so core cannot
// import flowdeps without a cycle.
package flowdeps
