# Localised Client Metadata

## Problem Statement

Client name and other human-readable metadata are only displayed in English. Client metadata should change based on the user-selected `ui_locale`.

---

## Goals

- Support language-tagged client metadata following the OIDC Dynamic Client Registration 1.0 convention. Refer https://openid.net/specs/openid-connect-registration-1_0.html#LanguagesAndScripts
- Support to add and update of the language-tagged client metadata in application endpoint.
- Support to add and update of the language-tagged client metadata in thunder DCR endpoint.
- Return the registered localised values for `client_name`, `tos_uri`, `policy_uri`, and `logo_uri` at /flow/meta endpoint along with falling back, base (untagged) value.
- Relying party requested `ui_locale` should be returned in the /flow/meta endpoint. UI MUST reflect the preferred `ui_locale` and change the client metadata displayed on the UI w.r.t the current selected `ui_locale`. Fallback to base (untagged) value if the `ui_locale` is missing or unknown.
- During the application registration or DCR, validate the BCP 47 language tags
- Changes MUST be backward compatible, not breaking already registred applications.

## Non-Goals

- Does not perform automatic translation of client metadata — only registered variants are served.
- Does not extend localisation support to metadata fields beyond `client_name`, `tos_uri`, `policy_uri`, and `logo_uri`.
- Does not affect any other locale-based behaviour in the authentication or token flows (e.g., sms, email templates).

---

## Acceptance Criteria

| # | Area | Scenario | Expected Behaviour |
|---|------|----------|--------------------|
| AC-01 | Registration (App API) | Create application with `client_name#fr` and base `client_name` | Both values stored; response echoes all tagged variants |
| AC-02 | Registration (App API) | Update application to add a new `client_name#de` variant | New variant persisted; existing variants unaffected |
| AC-03 | Registration (DCR) | `POST /register` with `client_name#fr` and base `client_name` | Both values stored; DCR response echoes all registered variants |
| AC-04 | Registration (DCR) | `PUT /register/{client_id}` to add or overwrite a tagged variant | Variant is updated; other variants remain unchanged |
| AC-05 | BCP 47 Validation | Register with a well-formed tag (`client_name#en-US`) | Accepted; tag normalised to lowercase on storage |
| AC-06 | BCP 47 Validation | Register with `client_name#` (empty tag) | Rejected with `400 invalid_client_metadata` |
| AC-07 | BCP 47 Validation | Register with `client_name#en US` or `client_name#en!` (illegal characters) | Rejected with `400 invalid_client_metadata` |
| AC-08 | BCP 47 Validation | Register with `client_name#en#US` (multiple `#`) | Rejected with `400 invalid_client_metadata` |
| AC-09 | BCP 47 Validation | Register with a tag exceeding maximum allowed length | Rejected with `400 invalid_client_metadata` |
| AC-10 | BCP 47 Validation | Register `client_name#fr` and `client_name#FR` in the same request | Last occurrence wins after tag normalisation to lowercase |
| AC-11 | BCP 47 Validation | Tags `fr`, `FR`, `Fr` all resolve to the same stored variant | All treated as equal after normalisation |
| AC-12 | Unsupported fields | Register `redirect_uris#fr` (tagged variant on non-localisable field) | silently ignored |
| AC-13 | URI Validation | Register `logo_uri#fr` with a non-HTTPS or malformed URL | Rejected with `400 invalid_client_metadata`, same rules as base `logo_uri` |
| AC-14 | Storage cap | Register more than the maximum allowed locale variants per field | Rejected with `400 invalid_client_metadata` indicating the limit |
| AC-15 | `/flow/meta` — exact match | `GET /flow/meta?type=APP&id={appId}&flowId={flowId}`; flow context has `ui_locale=fr`; `client_name#fr` registered | `client_name` in response resolves to the `fr` variant |
| AC-16 | `/flow/meta` — subtag fallback | Flow context has `ui_locale=fr-CA`; only `client_name#fr` registered | `client_name` resolves to the `fr` variant (language-only fallback applies) |
| AC-17 | `/flow/meta` — base fallback | Flow context has `ui_locale=de`; no `client_name#de` registered | `client_name` falls back to base (untagged) value |
| AC-18 | `/flow/meta` — no `flowId` provided | `GET /flow/meta?type=APP&id={appId}` with no `flowId` | Base (untagged) values returned for all localisable fields; no error |
| AC-19 | `/flow/meta` — expired or unknown `flowId` | `flowId` provided but flow does not exist or has expired | Degrades gracefully to base values; no `4xx`/`5xx` |
| AC-20 | `/flow/meta` — invalid `ui_locale` in flow context | Flow context contains `ui_locale=!!invalid` | Server degrades gracefully; base values returned; no `4xx`/`5xx` |
| AC-21 | `/flow/meta` — space-separated `ui_locale` | Flow context has `ui_locale=de fr` (OIDC spec list) | First matching variant wins |
| AC-22 | `/flow/meta` — `ui_locale` in response | `GET /flow/meta` with a valid `flowId` whose context has a `ui_locale` | Response includes the `ui_locale` value from the flow context |
| AC-23 | Backwards compatibility — no tagged variants | Existing client with no tagged variants makes an OAuth/OIDC request | Base values returned unchanged; no error or behaviour change |
| AC-24 | Backwards compatibility — no `ui_locale` sent | Consumer that never sends `ui_locale` | Receives base values exactly as before |
| AC-25 | Backwards compatibility — registration | Existing clients require no re-registration or migration | All pre-existing clients continue to work without modification |
| AC-26 | Authorisation | Unauthenticated or unauthorised party attempts to add/overwrite tagged variants | Request rejected with `401`/`403` |
| AC-27 | XSS / Injection | Tag values and URI variants stored and later reflected in responses | Values sanitised; no stored XSS or header injection possible |
| AC-28 | JSON serialisation | Tagged keys (containing `#`) round-trip through `encoding/json` without mangling | Field names preserved exactly in both request parsing and response serialisation |
| AC-29 | DCR read response | `GET /register/{client_id}` | Returns all registered tagged variants; no locale resolution applied |

---

## Edge Cases & Constraints

### Input & Parsing

- **`#` with no tag** — `client_name#` has an empty language tag; must be rejected with `400 invalid_client_metadata`.
- **Multiple `#` in field name** — `client_name#en#US` is ambiguous; the `#` separator is used only once — the tag itself may contain `-` for BCP 47 subtags (e.g., `client_name#en-US`). Reject with `400`.
- **Case sensitivity of tags** — BCP 47 tags are case-insensitive (`fr`, `FR`, `Fr` must resolve to the same variant); normalise to lowercase on write.
- **Whitespace or special characters in tag** — `client_name#en US` or `client_name#en!` must be rejected with `400 invalid_client_metadata`.
- **Extremely long tag** — tags must be bounded in length (e.g., max 35 characters per RFC 5646 practical limits) and rejected cleanly if exceeded.
- **Duplicate variants in the same request** — a request containing both `client_name#fr` and `client_name#FR` resolves to a conflict after normalisation; last-write-wins (the last occurrence in the parsed map is stored).

### Storage & Data Integrity

- **Base field absent but tagged variants present** — the base field (`client_name`) should always be provided and is used as the fallback when no registered variant matches the requested `ui_locale`. If the base is absent and no variant matches, the field is returned empty.
- **Deleting a base value while tagged variants exist** — an update that nulls `client_name` but leaves `client_name#fr` creates an inconsistent state; explicitly permitted but `/flow/meta` will only return the matching tagged variant or empty.
- **Maximum number of locale variants per field** — cap at 20 variants per field per client to prevent unbounded storage growth in the JSON column. Reject further additions with `400 invalid_client_metadata`.
- **Supported fields only** — a tagged variant on a non-localisable field (e.g., `redirect_uris#fr`) is silently ignored; the field is stripped before storage.

### Locale Resolution

- **Region subtag fallback** — if `ui_locale=fr-CA` is requested but only `fr` is registered, `fr` is returned as the resolved value. The resolution chain is: exact match → language-only prefix match → base value.
- **`ui_locale` with multiple values** — space-separated list per OIDC spec; first matching variant wins.
- **`ui_locale` containing an invalid BCP 47 tag at runtime** — server degrades gracefully to the base value; no `4xx`/`5xx` returned.

### Security & Misuse

- **Locale tag as an injection vector** — tag values are stored and reflected in responses; sanitise using the existing `sysutils.SanitizeString()` pattern before storage.
- **URI fields with tagged variants must still pass URI validation** — `logo_uri#fr`, `tos_uri#fr`, `policy_uri#fr` must be validated with the same HTTPS and well-formedness checks as their base fields.
- **Privilege escalation via update** — only an authorised party (the client itself or an admin) may add or overwrite language-tagged variants; enforced by the existing JWT middleware in `system/security/`.

### Backwards Compatibility

- **Existing clients with no tagged variants** — continue to work without modification; no migration or re-registration required. The JSON column addition is additive.
- **Consumers that do not send `ui_locale`** — receive base values exactly as before; no behaviour change.
- **Serialisation of tagged keys in JSON** — Go's `encoding/json` is expected to round-trip keys containing `#` without escaping, consistent with the existing `app_json` column pattern; must be confirmed with a unit test before relying on it.

---

## Technical Notes

### Changes Required

#### 1. Client Model — `backend/internal/application/model/application.go`

Add localised variant maps to `ApplicationDTO` and `ApplicationProcessedDTO`:

```go
// New fields in ApplicationDTO
LocalisedClientName   map[string]string `json:"client_name_localised,omitempty"`   // keyed by normalised BCP 47 tag
LocalisedLogoURL      map[string]string `json:"logo_uri_localised,omitempty"`
LocalisedTosURI       map[string]string `json:"tos_uri_localised,omitempty"`
LocalisedPolicyURI    map[string]string `json:"policy_uri_localised,omitempty"`
```

When deserialising an inbound request, a pre-processing step must scan all JSON keys for the `#` separator, extract the tag, validate it as BCP 47, normalise it to lowercase, and populate the corresponding map. This step must run before standard `json.Unmarshal` so that unknown tagged keys do not surface as errors.

#### 2. DCR Model — `backend/internal/oauth/oauth2/dcr/model.go`

`DCRRegistrationRequest` and `DCRRegistrationResponse` use fixed struct fields. Tagged keys (`client_name#fr`) cannot be captured by fixed fields.

**Approach: custom `UnmarshalJSON` / `MarshalJSON`** — preferred over a catch-all `AdditionalFields` map, which would leak into the response and require a separate filter pass.

**Request (`UnmarshalJSON`)** — decode fixed fields via an alias type to avoid recursion, then scan raw keys for the `#` separator to populate the localised maps:

```go
func (r *DCRRegistrationRequest) UnmarshalJSON(data []byte) error {
    // Step 1: decode all fixed fields (avoids infinite recursion via alias)
    type Alias DCRRegistrationRequest
    if err := json.Unmarshal(data, (*Alias)(r)); err != nil {
        return err
    }
    // Step 2: scan for tagged keys
    var raw map[string]json.RawMessage
    if err := json.Unmarshal(data, &raw); err != nil {
        return err
    }
    for key, val := range raw {
        field, tag, ok := strings.Cut(key, "#")
        if !ok {
            continue
        }
        if !IsValidBCP47Tag(tag) {
            return fmt.Errorf("invalid BCP 47 tag in field %q", key)
        }
        tag = NormaliseBCP47Tag(tag)
        var s string
        if err := json.Unmarshal(val, &s); err != nil {
            continue
        }
        switch field {
        case "client_name":
            if r.LocalisedClientName == nil {
                r.LocalisedClientName = make(map[string]string)
            }
            r.LocalisedClientName[tag] = s
        case "logo_uri":
            // ... same pattern
        case "tos_uri":
            // ... same pattern
        case "policy_uri":
            // ... same pattern
        }
    }
    return nil
}
```

**Response (`MarshalJSON`)** — marshal the struct normally via alias, then inject tagged keys as top-level properties:

```go
func (r DCRRegistrationResponse) MarshalJSON() ([]byte, error) {
    type Alias DCRRegistrationResponse
    base, err := json.Marshal(Alias(r))
    if err != nil {
        return nil, err
    }
    // Unmarshal into a map, inject tagged variants, re-marshal
    var m map[string]interface{}
    if err := json.Unmarshal(base, &m); err != nil {
        return nil, err
    }
    for tag, val := range r.LocalisedClientName {
        m["client_name#"+tag] = val
    }
    // ... same for logo_uri, tos_uri, policy_uri
    return json.Marshal(m)
}
```

The four localised maps (`LocalisedClientName`, `LocalisedLogoURL`, `LocalisedTosURI`, `LocalisedPolicyURI`) are added to both `DCRRegistrationRequest` and `DCRRegistrationResponse` with `json:"-"` tags so standard marshalling ignores them — only the custom methods handle them.

#### 3. BCP 47 Validation Utility — `backend/internal/system/utils/` (new file `locale_util.go`)

A lightweight validation function:

```go
// IsValidBCP47Tag returns true if the tag conforms to BCP 47 (RFC 5646) syntax.
// Accepts tags like "fr", "en-US", "zh-Hant-TW". Rejects empty, whitespace, or special chars.
func IsValidBCP47Tag(tag string) bool

// NormaliseBCP47Tag returns the lowercase-normalised form of a valid tag.
func NormaliseBCP47Tag(tag string) string
```

No external library needed — a regex covering primary language + optional subtags is sufficient. Place alongside the existing `MatchRedirectURIPattern` in `http_util.go` or in a new `locale_util.go`.

#### 4. Application Service — `backend/internal/application/service.go`

- **`validateApplication()`** — add a call to validate all tagged keys in `LocalisedClientName`, `LocalisedLogoURL`, `LocalisedTosURI`, `LocalisedPolicyURI`: BCP 47 tag format, tag length cap, variant count cap (≤ 20 per field), URI validity for tagged URI fields.
- **`CreateApplication()` / `UpdateApplication()`** — normalise all tag keys to lowercase before passing to the store. On update, merge incoming tagged variants with existing ones (do not overwrite unmentioned variants).
- Strip tagged variants from non-localisable fields before storage (silently ignore per AC-12).

#### 5. Storage — `backend/internal/application/store.go`

The localised variant maps are included inside `app_json` (the existing JSON column). No schema migration is required. The `ApplicationProcessedDTO` must serialise the four `*_localised` maps into `app_json` and deserialise them on read.

Verify the JSON column size limit is sufficient for 4 fields × 20 variants × reasonable string length; if using SQLite, `TEXT` is unbounded; if Postgres, `jsonb` is also unbounded — no action needed.

#### 6. Authorize Endpoint — `backend/internal/oauth/oauth2/authz/`

`ui_locale` arrives as a query parameter on `GET /authorize` alongside `client_id`, `scope`, etc. It is not currently extracted.

**`service.go` — `HandleInitialAuthorizationRequest()`**:
- Extract `ui_locale` from the `OAuthParameters` (add `UILocale string` field to the `OAuthParameters` model in `oauth/oauth2/model/parameter.go`).
- Add a new runtime data key in `flow/common/constants.go`:
  ```go
  RuntimeKeyUILocale = "ui_locale"
  ```
- Pass `ui_locale` into the `RuntimeData` map when building `FlowInitContext`, alongside the existing `RuntimeKeyRequiredLocales` entry:
  ```go
  runtimeData[flowcm.RuntimeKeyUILocale] = uiLocale
  ```
- Store it in `authRequestContext` (in `auth_req_store.go`) alongside other OAuth parameters for consistency.

#### 7. `/flow/meta` Endpoint — `backend/internal/flow/flowmeta/`

`ui_locale` is **not** a direct query param here. It must be read from the flow's `RuntimeData` using the `flow_id` that the client already holds from the authorize response.

**`handler.go`**:
- Accept an optional `flowId` query parameter; sanitise with `sysutils.SanitizeString()`.
- Pass it to the service.

**`service.go` — `GetFlowMetadata()`**:
- Accept `flowID *string` as a new parameter.
- If `flowID` is provided, load the flow context from the flow execution store (`flowexec` store) and extract `RuntimeData[RuntimeKeyUILocale]`.
- Use the resolved `ui_locale` to call `ResolveLocalisedValue` for each of `client_name`, `logo_uri`, `tos_uri`, `policy_uri`.
- If `flowID` is absent or the flow context has no `ui_locale`, fall back to base (untagged) values — no error.
- Include the resolved `ui_locale` string in `FlowMetadataResponse`.

**`model.go`** — add `UILocale string` to `FlowMetadataResponse` and `ApplicationMetadata`.

**Locale resolution helper** (in `system/utils/locale_util.go`):
```go
// ResolveLocalisedValue returns the best matching value for the requested locale.
// Resolution order: exact tag match → language prefix match → base value.
func ResolveLocalisedValue(variants map[string]string, base string, uiLocale string) string
```

For space-separated `ui_locale` lists (OIDC spec), iterate in order and return the first match.

**Data flow summary:**
```
GET /authorize?...&ui_locale=fr-CA
  → extracted into OAuthParameters.UILocale
  → stored in FlowContext.RuntimeData["ui_locale"] = "fr-CA"
  → flowId returned to client in redirect query params

GET /flow/meta?type=APP&id={appId}&flowId={flowId}
  → loads FlowContext by flowId
  → reads RuntimeData["ui_locale"] = "fr-CA"
  → resolves client_name, logo_uri, tos_uri, policy_uri against registered variants
  → returns resolved values + ui_locale in response
```

#### 8. DCR Service — `backend/internal/oauth/oauth2/dcr/service.go`

- **`RegisterClient()`** — the custom `UnmarshalJSON` on `DCRRegistrationRequest` (Section 2) handles extraction and validation of tagged fields; the service receives a fully populated `LocalisedClientName` etc. on the request struct and maps them to `ApplicationDTO` before calling the application service.
- **`GET /register/{client_id}`** — return all stored tagged variants in the response; do not perform locale resolution (raw registered data, per AC-28).

### Services & Layers Touched

| Layer | File | Change |
|-------|------|--------|
| Model | `oauth/oauth2/model/parameter.go` | Add `UILocale string` field to `OAuthParameters` |
| Constants | `flow/common/constants.go` | Add `RuntimeKeyUILocale = "ui_locale"` |
| Service | `oauth/oauth2/authz/service.go` | Extract `ui_locale` from authorize request; write to flow `RuntimeData` and `authRequestContext` |
| Model | `application/model/application.go` | Add four `map[string]string` localised variant fields to `ApplicationDTO` |
| Model | `oauth/oauth2/dcr/model.go` | Custom unmarshal/marshal to handle `#`-keyed fields in DCR request/response |
| Validation | `system/utils/locale_util.go` (new) | `IsValidBCP47Tag`, `NormaliseBCP47Tag`, `ResolveLocalisedValue` |
| Service | `application/service.go` | Validate, normalise, merge localised variants on create/update |
| Store | `application/store.go` | Serialise/deserialise localised maps inside `app_json` |
| Service | `oauth/oauth2/dcr/service.go` | Parse tagged keys from DCR request; populate localised maps |
| Handler | `flow/flowmeta/handler.go` | Accept optional `flowId` query param |
| Service | `flow/flowmeta/service.go` | Load flow context by `flowId`; read `ui_locale` from `RuntimeData`; resolve and return localised variants |
| Model | `flow/flowmeta/model.go` | Add `UILocale` to `FlowMetadataResponse` and `ApplicationMetadata` |

### Known Gaps & Decisions Needed

- **JSON key with `#`** — verify early with a unit test that `encoding/json` round-trips `{"client_name#fr": "..."}` without mangling; the existing `app_json` JSON column pattern suggests this will work.
- **DECLARATIVE / COMPOSITE store modes** — the declarative file-based store reads YAML config files; ensure localised variant maps can be expressed in those YAML files and are parsed correctly.
- **OQ-7 unresolved** — whether the Asgardeo JS SDK forwards `ui_locale` to the backend authorize request must be confirmed before frontend integration is considered complete.

### Dependencies

- **No external team dependencies** — DCR and client management are owned within this repo.
- **OIDC DCR 1.0 spec** — OpenID Connect Dynamic Client Registration 1.0 §2 is the normative reference for the `#`-suffix convention.
- **BCP 47 (RFC 5646)** — governs valid language tag syntax; a lightweight regex is sufficient, no heavy i18n dependency needed.

---

## Open Questions

| # | Question | Owner | Status |
|---|----------|-------|--------|
| OQ-1 | Is the base field (`client_name`) required when language-tagged variants are present, or can a client register with only tagged variants? | Product / Backend | Resolved: We should still capture base field. should be used as fallback value when language tagged value is not found. |
| OQ-2 | Should BCP 47 region subtag fallback be supported — i.e., does `ui_locale=fr-CA` resolve to a `fr` variant if `fr-CA` is not registered? | Product | Resolved: We should fallback to language tag without region  |
| OQ-3 | What is the maximum number of locale variants allowed per field per client? | Backend | Resolved: cap at 20 variants per field per client (Edge Cases — Storage & Data Integrity) |
| OQ-4 | If `ui_locale` is a space-separated list (as permitted by OIDC spec), is first-match-wins the correct resolution strategy? | Backend | Resolved: first-match-wins |
| OQ-5 | What is the storage strategy for dynamic tagged variants — embed in the existing `app_json` column, or a separate `client_metadata` key-value table? | Backend | Resolved: embed in `app_json`; additive change, no schema migration required |
| OQ-6 | Should tagged variants on non-localisable fields (e.g., `redirect_uris#fr`) be rejected with an error or silently ignored? | Backend | Resolved: silently ignored; stripped before storage (AC-12) |
| OQ-7 | Does the Asgardeo JS SDK reliably forward `ui_locale` from the login gate to the backend authorization request? If not, does the frontend need to set it explicitly? | Frontend / SDK | Open |
| OQ-8 | Are `tos_uri` and `policy_uri` tagged variants subject to the same URI allow-list or HTTPS enforcement as the base fields? | Security | Resolved: yes, same HTTPS and well-formedness validation applies to all tagged URI fields (AC-13) |
| OQ-9 | Should the DCR read response (`GET /register/{client_id}`) expose all registered tagged variants, or only the resolved value for the current request's locale? | Product | Resolved: It should expose all registered tagged variants. |
| OQ-10 | Is there a migration plan for existing clients — do they need a backfill, or do they simply have no tagged variants and fall back to base by default? | Backend | Resolved: no backfill needed; existing clients have no tagged variants and fall back to base by default (AC-23) |

> **Signal to proceed:** All rows above should be marked **Resolved** before development starts. Any open question at dev kickoff is a spec gap.

---

## References
- https://openid.net/specs/openid-connect-registration-1_0.html#LanguagesAndScripts