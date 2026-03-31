# ACR Values Support in Thunder

## Problem statement

Thunder currently ignores the `acr_values` parameter in authorization requests entirely — it is not parsed, not propagated to the flow engine, and not reflected in issued ID tokens. The `acr_values` parameter is a standards-defined mechanism (OpenID Connect Core 1.0 §3.1.2.1) that lets relying parties express the required level of assurance at which an end-user must be authenticated. Without this support, Thunder cannot fulfill requests from RPs that need to enforce specific authentication methods (e.g., OTP-only, password-only, or a preference ordering of methods).

## Goals

- Maintain ACR-AMR mapping in configuration.
- Allow relying parties to register an `allowedAcrs` list on the application via the `/application` endpoint; each entry is validated against the ACR-AMR mapping configuration at registration time.
- Parse the `acr_values` query parameter from the `/oauth2/authorize` request, validate each requested ACR against the application's registered `allowedAcrs`. Ignore unknown or unregistered acr_values, and propagate the validated ACR list into the flow execution context via `RuntimeData`.
- Introduce a new `acr_options` variant on PROMPT nodes in the flow graph JSON to declaratively mark which prompt represents a ACR chooser. Multiple ACR options map to multiple actions in this variant prompt node.
- At `/execute` endpoint, when the next node to be returned is a `acr_options` prompt node variant, filter out any ACR option (prompt/action) whose ACR was not requested, and reorder the remaining options to match the preference order expressed in `acr_values`.
- Skip the `acr_options` variant on PROMPT node and move to next node, if only one valid ACR remains after the validation.
- If none of the requested acr_values are valid, simply fallback to all the registered acr values.
- Auth executor MUST validate if the AMR in the authentication request matches the selected ACR (stored in RuntimeData) from the `acr_options` prompt node. So that authentication request with incorrect / out of order AMR hints malicious use of the /execute call. 
- Include the `acr` claim in the issued ID token to reflect the ACR of the authentication method the user actually statisfied. At the end of the flow execution, statisfied ACR must be added in the auth assertion JWT.
- OIDC Discovery 1.0 specifies the field as `acr_values_supported`. Publish the supported ACR values in the OpenID Connect discovery well-known endpoint.

## Non-Goals

- Support for `acr` inside the `claims` request parameter `id_token` field is out of scope.

## Acceptance Criteria

| # | Area | Criteria |
|---|------|----------|
| AC-1 | ACR-AMR Config | ACR-AMR mapping can be defined in server configuration (e.g., `loa1: [pwd]`, `loa2: [pwd, otp]`). |
| AC-2 | ACR-AMR Config | Server fails to start (or rejects config reload) if the ACR-AMR mapping contains duplicate ACR keys or empty AMR lists. |
| AC-3 | Application Registration | Registering an application with `allowedAcrs` entries that all exist in the ACR-AMR config succeeds. |
| AC-4 | Application Registration | Registering an application with an `allowedAcrs` entry not in the ACR-AMR config returns a `400` validation error. |
| AC-5 | Application Registration | `allowedAcrs` is madatory; Registering an application without `allowedAcrs`  returns a `400` validation error. |
| AC-6 | `/oauth2/authorize` | `acr_values` query parameter is parsed from the authorization request. |
| AC-7 | `/oauth2/authorize` | Only ACR values present in both the request's `acr_values` and the application's `allowedAcrs` are propagated into `RuntimeData`. |
| AC-8 | `/oauth2/authorize` | ACR values in the request that are not in `allowedAcrs` are silently ignored with no error returned to the client. |
| AC-9 | `/oauth2/authorize` | If none of the requested `acr_values` match `allowedAcrs`, no ACR filtering is applied and all registered ACR options are presented in the flow. |
| AC-10 | Flow graph | A PROMPT node with the `acr_options` variant is valid in the flow graph JSON; each action in the node corresponds to one ACR option. |
| AC-11 | `/flow/execute` | When the next node is an `acr_options` prompt node, the response only includes actions whose ACR is in the validated ACR list from `RuntimeData`. |
| AC-12 | `/flow/execute` | The filtered actions are ordered to match the preference order of the original `acr_values` parameter. |
| AC-13 | `/flow/execute` | When exactly one valid ACR option remains after filtering, the `acr_options` node is skipped, that ACR is auto-selected, and the next node is returned instead. |
| AC-14 | Auth executor | The auth executor reads the selected ACR from `RuntimeData` and verifies the method's AMR satisfies it per the ACR-AMR config; a mismatch returns an error. |
| AC-15 | Auth executor | A `/flow/execute` call that bypasses the `acr_options` node and submits credentials for an AMR inconsistent with the selected ACR is rejected. |
| AC-16 | ID token | The issued ID token contains an `acr` claim reflecting the ACR the user satisfied. |
| AC-17 | Auth assertion JWT | At flow completion, the satisfied ACR value is recorded in the auth assertion JWT so the token builder can include it in the ID token. |
| AC-18 | Discovery | `GET /.well-known/openid-configuration` response includes an `acr_values_supported` array. |
| AC-19 | Discovery | `acr_values_supported` contains exactly the ACR values defined in the server's ACR-AMR configuration. 

## Technical Notes

TODO

## Open Questions

1. Adding allowedAcrs to the application model (backend/internal/application/model/oauth_app.go) requires a schema migration. But what should be the default acr values for the existing applications?