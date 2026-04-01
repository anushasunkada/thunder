# Feature: LOGIN_ID_INPUT Flow Component

## Problem Statement

Authentication flows often need to support multiple login identifier types — e.g., mobile number, email, national ID — within the same sign-in step. Currently, the Asgardeo UI SDK's flow graph model has no way to express this pattern: the server cannot instruct the SDK to render a type-selector UI alongside a contextual input field in a single flow step.

The existing workaround would require either:
- Multiple server round-trips (one `ACTION TRIGGER` per login type), introducing noticeable latency on every tab switch, or
- Hardcoding login type behavior in the SDK, coupling UI logic to application-level configuration (as in the esignet `loginid-types` design).

Neither is acceptable. The flow graph model should be expressive enough to represent this without special-casing it.

---

## Goals & Non-Goals

### Goals

- Add a `LOGIN_ID_INPUT` component type to the embedded flow graph model.
- Allow the server to declare any number of login ID types (mobile, email, NRC, VID, etc.) and their per-type configuration within a single component.
- Support per-type prefix configuration — a single prefix string or a list of selectable prefixes (e.g., country dial codes).
- Support per-type postfix strings.
- The SDK assembles the final identifier value as `prefix + rawInput + postfix` before submitting to the server. The server receives only the assembled value.
- Support `default: true` on one login ID type to pre-select it on render.
- Use semantic icon names (resolved by the SDK via lucide-react, consistent with the existing `ICON` component type) rather than asset-path strings.
- Support per-type `maxLength` and `regex` validation, consistent with existing form validation in `AuthOptionFactory`. Validation failures use a generic i18n message — no per-type custom message needed.
- Postfix is appended to the submitted value silently; it is not displayed to the user in the input field.
- Support i18n via the existing `{{t(...)}}` template literal convention for labels and placeholders.
- The feature is additive and backward-compatible — no existing component types or flow contracts are modified.

### Non-Goals

- This feature does not introduce any new server-side API endpoints. The assembled identifier is submitted through the existing `EmbeddedSignInFlowRequest.inputs` map.
- This feature does not add country-flag rendering for dial code prefixes in v1 of implementation.

---

## Acceptance Criteria

1. **New component type registered**: `EmbeddedFlowComponentType.LoginIdInput = 'LOGIN_ID_INPUT'` exists in `packages/javascript/src/models/v2/embedded-flow-v2.ts`.

2. **Type definitions complete**: A `LoginIdType` interface and updated `EmbeddedFlowComponent` discriminated union are exported from `@asgardeo/javascript`.

3. **Rendered correctly**: When `AuthOptionFactory` receives a `LOGIN_ID_INPUT` component, it renders:
   - The component-level `label` as a heading above the grid.
   - A grid of selector buttons, one per login ID type.
   - The active `LoginIdType.label` as an input field label below the grid.
   - A single input field whose `inputType`, prefix selector (if applicable), placeholder, maxLength, and validation change to match the selected login ID type.
   - The type marked `default: true` is pre-selected on mount. If no type is marked default, the first type is selected.

4. **Prefix handling**:
   - If `prefixes` is a single string, it is displayed as a static prefix label.
   - If `prefixes` is an array of objects, it is rendered as a dropdown/selector. The user picks one prefix; the selected value is used in assembly.
   - If `prefixes` is absent or empty, no prefix UI is shown.

5. **Value assembly**: On form submission, the value submitted for the component's `ref` key is `selectedPrefix.value + rawInput + postfix`. If prefix or postfix is absent, those parts are omitted (no trailing/leading empty string concatenation edge cases).

6. **Validation**: `regex` and `maxLength` from the active login ID type (with prefix-level overrides applied) are enforced against the raw input. Failures surface via the existing `formErrors` / `fieldErrors` mechanism using a generic i18n fallback message. No `validationMessage` field on `LoginIdType` is required.

7. **Postfix display**: The postfix is never shown in the input field. It is appended silently during value assembly on submit.

8. **Icon resolution**: The `icon` field on each login ID type is resolved to a lucide-react icon by name, consistent with how the existing `ICON` component type resolves icons. An unrecognized icon name renders no icon without throwing.

9. **i18n**: `label` and `placeholder` fields support `{{t(key)}}` template literals, resolved by the existing template resolver in `AuthOptionFactory`.

10. **Single-type degenerate case**: If only one login ID type is provided, the selector tab row is not rendered — only the input is shown.

11. **Tab switch behavior**: Switching login ID type swaps the input instantly (no animation). The raw input value is cleared on type switch; existing validation errors are also cleared.

12. **Accessibility**: The type selector uses `role="tablist"` / `role="tab"` (ARIA tab pattern) and is keyboard-navigable (arrow keys between tabs, Enter/Space to select). The input has a proper `aria-label` derived from the active type's label.

13. **Unit tests**: Coverage for value assembly logic (prefix + input + postfix permutations), prefix-switch re-validation, and the degenerate single-type case.

---

## Technical Notes

### New type: `LoginIdType`

Add to [packages/javascript/src/models/v2/embedded-flow-v2.ts](../../packages/javascript/src/models/v2/embedded-flow-v2.ts):

```typescript
export interface LoginIdPrefix {
  /** Display label for this prefix (e.g. country name or ISO code) */
  label: string;
  /** The actual prefix value to concatenate (e.g. "+91") */
  value: string;
  /** Overrides outer maxLength for this prefix selection */
  maxLength?: number;
  /** Overrides outer regex for this prefix selection */
  regex?: string;
}

export interface LoginIdType {
  /** Unique identifier for this login ID type (e.g. "mobile", "email") */
  id: string;
  /** Semantic icon name resolved by the SDK (lucide-react icon name) */
  icon?: string;
  /** Display label. Supports {{t(key)}} template literals */
  label: string;
  /** Input placeholder. Supports {{t(key)}} template literals */
  placeholder?: string;
  /** Static prefix string, or list of selectable prefix objects */
  prefixes?: string | LoginIdPrefix[];
  /** Static string appended to the raw input value before submission */
  postfix?: string;
  /** Maximum character length of the raw input (before prefix/postfix) */
  maxLength?: number;
  /** Regex pattern the raw input must satisfy */
  regex?: string;
  /** HTML input type hint for the browser. Defaults to "text" if absent. */
  inputType?: 'text' | 'email' | 'tel' | 'number';
  /** Whether this type is pre-selected on render */
  default?: boolean;
}
```

Extend `EmbeddedFlowComponentType`:

```typescript
/** Login ID type selector with contextual input for multi-identifier authentication */
LoginIdInput = 'LOGIN_ID_INPUT',
```

Add `loginIdTypes` as an optional field on `EmbeddedFlowComponent`:

```typescript
/** Present when type === LOGIN_ID_INPUT */
loginIdTypes?: LoginIdType[];
```

### Value assembly in `AuthOptionFactory`

[packages/react/src/components/presentation/auth/AuthOptionFactory.tsx](../../packages/react/src/components/presentation/auth/AuthOptionFactory.tsx) adds one new case in `createAuthComponentFromFlow()` that renders `<LoginIdInput>`. All state lives inside `useLoginIdInput` — `AuthOptionFactory` is only responsible for passing `onInputChange` and the component definition through, as it does for all other component types.

Value assembly happens at submit time (not on every keystroke) inside `useLoginIdInput`:
```
finalValue = (selectedPrefix?.value ?? "") + rawInput + (activeType.postfix ?? "")
```

The assembled `finalValue` is what gets written into `formValues[ref]` before the existing submit path runs. The server derives the login ID type from the assembled value. No additional type metadata field is submitted.

### Icon resolution

Re-use the existing icon resolution map used by `EmbeddedFlowComponentType.Icon`. The `icon` field on `LoginIdType` follows the same lookup — no new resolution mechanism needed.

### Prefix selector component

`PHONE_INPUT` is declared in `EmbeddedFlowComponentType` but has no React renderer yet. Rather than waiting for it, the prefix selector should be implemented as a standalone `PrefixSelector` component inside `LoginIdInput`, so that `PHONE_INPUT`'s future renderer can reuse it without rebuilding the pattern.

```
packages/react/src/components/presentation/auth/LoginIdInput/
  LoginIdInput.tsx          ← top-level component (grid + input)
  PrefixSelector.tsx        ← reusable prefix dropdown, used here and by future PHONE_INPUT
  useLoginIdInput.ts        ← state: activeTypeId, selectedPrefix, rawInput; assembles finalValue on submit
```

`PrefixSelector` should be a custom styled dropdown (not a native `<select>`) to stay consistent with the SDK's component library.

### Packages touched

| Package | Change |
|---|---|
| `@asgardeo/javascript` | New types: `LoginIdType`, `LoginIdPrefix`; new enum value `LoginIdInput` |
| `@asgardeo/react` | New `LoginIdInput` component; new case in `AuthOptionFactory`; value assembly logic |
| `@asgardeo/i18n` | Verify a generic validation error key exists (e.g. `errors.invalidFormat`); add it if absent |

### Example flow graph payload

```json
{
  "type": "LOGIN_ID_INPUT",
  "id": "login-id-field",
  "ref": "username",
  "label": "{{t(loginId.selector.label)}}",
  "loginIdTypes": [
    {
      "id": "mobile",
      "icon": "Smartphone",
      "label": "{{t(loginId.mobile.label)}}",
      "placeholder": "{{t(loginId.mobile.placeholder)}}",
      "prefixes": [
        { "label": "IND", "value": "+91", "maxLength": 10 },
        { "label": "KHM", "value": "+855", "maxLength": 9 }
      ],
      "postfix": "@phone",
      "regex": "^[0-9]+$",
      "default": true
    },
    {
      "id": "email",
      "icon": "Mail",
      "label": "{{t(loginId.email.label)}}",
      "placeholder": "{{t(loginId.email.placeholder)}}",
      "maxLength": 254
    },
    {
      "id": "nrc",
      "icon": "IdCard",
      "label": "{{t(loginId.nrc.label)}}",
      "placeholder": "{{t(loginId.nrc.placeholder)}}",
      "postfix": "@NRC"
    }
  ]
}
```

---

## UX Design

The component is built on the same primitives as `TEXT_INPUT`: `FormControl`, `InputLabel`, and `TextField`. There are two distinct labels: the component-level `label` (rendered above the button grid as a section heading) and each `LoginIdType.label` (rendered below the grid as the input field label, updating when the active type changes).

### Layout

**3 types → 1 row × 3 columns:**
```
┌─────────────────────────────────────────────┐
│ Select a preferred ID to Login              │
│                                              │
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ │ 📱 Mobile   │ │ ✉ Email     │ │ 🪪 NRC      │
│ └─────────────┘ └─────────────┘ └─────────────┘
│  Label                                       │
│ ┌──────────┬────────────────────────────┐   │
│ │ +91 ▾   │  Enter mobile number       │   │
│ └──────────┴────────────────────────────┘   │
│                                              │
│  ⚠ Error message                            │
└─────────────────────────────────────────────┘
```

**4 types → 2 rows × 2 columns:**
```
┌─────────────────────────────────────────────┐
│ Select a preferred ID to Login               │
│                                              │
│ ┌───────────────────┐ ┌───────────────────┐ │
│ │ 📱 Mobile         │ │ ✉ Email           │ │
│ └───────────────────┘ └───────────────────┘ │
│ ┌───────────────────┐ ┌───────────────────┐ │
│ │ 🪪 NRC            │ │ 🪪 VID            │ │
│ └───────────────────┘ └───────────────────┘ │
│   Label                                     │
│ ┌──────────┬────────────────────────────┐   │
│ │ +91 ▾   │  Enter mobile number       │   │
│ └──────────┴────────────────────────────┘   │
│                                              │
│  ⚠ Error message                            │
└─────────────────────────────────────────────┘
```

- The component-level label above the grid uses `InputLabel` with `variant="block"` — same as `TEXT_INPUT`.
- The per-type input label below the grid also uses `InputLabel` with `variant="block"` and updates when the active type changes.
- The error message below the input uses `FormControl`'s helper text slot, identical to `TEXT_INPUT`.
- The entire component is wrapped in `FormControl` so spacing, error state propagation, and BEM class structure stay consistent.

### Tab Row (Login ID Type Selector)

- Rendered as a grid of independent `Button` components with `variant="outline"` — each button is separate with its own border, not joined into a button group.
- Column count is computed at render time from `loginIdTypes.length`:
  - **≤ 3 types**: `columns = count` — all buttons on one row, equal width.
  - **> 3 types**: `columns = 2` — buttons fill a 2-column grid, wrapping into as many rows as needed.
- Layout uses CSS Grid: `display: grid; grid-template-columns: repeat(columns, 1fr); gap: theme.vars.spacing.unit`. Each button stretches to fill its cell (`1fr`), so rows are always visually balanced.
- The active button switches to `variant="solid"` with `color="primary"` to indicate selection.
- Each button shows the lucide-react icon (if provided) as `startIcon`, followed by the label text.
- Uses `role="tablist"` / `role="tab"` with `aria-selected` on each button.
- On narrow viewports the 2-column grid naturally compresses — no breakpoint overrides needed. If labels overflow a cell, text truncates with ellipsis and the full label is available via `title` attribute.

### Prefix Selector (`PrefixSelector`)

- Rendered inline, to the left of the `TextField`, inside the same input container — visually appears as a prefixed segment of the field (matching the `startIcon` padding pattern in `TextField.styles.ts`).
- When `prefixes` is a single string: renders as a static non-interactive label with the same padding/border as the input, separated by a divider.
- When `prefixes` is an array: renders as a custom dropdown button (not a native `<select>`) that opens a listbox above/below via Floating UI — consistent with how `Select` primitive uses Floating UI for positioning.
- The prefix container and the text input share a single outer border (one unified field boundary), not two separate bordered elements side by side.

### Visual States

Mirrors `TEXT_INPUT` states exactly:

| State | Appearance |
|---|---|
| **Default** | 1px solid `theme.vars.colors.border` |
| **Focused** | Primary blue border + 20% opacity blue box-shadow |
| **Error** | 1px solid `theme.vars.colors.error.main` |
| **Focused + Error** | Error red border + 20% opacity red box-shadow |
| **Disabled** | 0.6 opacity, `theme.vars.colors.background.disabled` |

The error state applies to the entire input container (prefix + text field unified boundary), not just the text portion.

### Transitions

Consistent with `TextField`: `border-color` and `box-shadow` transition at `0.2s ease`.

### BEM Class Structure

```
asgardeo-login-id-input                     ← FormControl root
  asgardeo-login-id-input__label            ← component-level InputLabel (above grid)
  asgardeo-login-id-input__grid             ← button grid
    asgardeo-login-id-input__tab            ← individual type button
    asgardeo-login-id-input__tab--active    ← selected type modifier
  asgardeo-login-id-input__input-label      ← per-type InputLabel (below grid, updates on type switch)
  asgardeo-login-id-input__field-row        ← prefix + input container
    asgardeo-login-id-input__prefix         ← PrefixSelector
    asgardeo-login-id-input__input          ← TextField
  asgardeo-login-id-input__helper-text      ← error / helper text
  asgardeo-login-id-input__helper-text--error
```

---

## Edge Cases or Gaps

- **No `default: true` set**: First type in the array is selected. If the array is empty, the component renders nothing and logs a warning (consistent with how `OU_SELECT` handles missing `rootOuId`).

- **Multiple `default: true`**: First one wins; subsequent `default: true` flags are ignored.

- **Single login ID type**: Selector tab row is suppressed. Only the input is rendered. The single type's prefix/postfix still apply.

- **Prefix is empty string `""`**: Treated the same as absent — no prefix UI, no prefix concatenation.

- **User switches login type mid-input**: Raw input is cleared on type switch. This avoids submitting e.g. an email address with a phone postfix appended, and avoids surfacing stale validation errors from the previous type.

- **`maxLength` enforcement with prefix**: `maxLength` applies to the raw input only (before prefix/postfix), not the assembled value. This must be explicit in implementation to avoid off-by-one truncation. When the user switches prefix within the same login type, the input is **not** cleared — instead, validation re-runs immediately against the new prefix's `maxLength` (or the outer `maxLength` if the new prefix does not define one), surfacing an error if the existing input now exceeds the limit.

- **Regex applied to raw input or assembled value**: Regex should validate the raw input (before assembly), since the postfix is a known static string and validating the assembled value would require escaping it into the regex. Document this clearly.

- **Server sends `LOGIN_ID_INPUT` with `loginIdTypes: null`**: Treat as empty array — render nothing, emit a `logger.warn`.

- **Framework packages beyond `@asgardeo/react`**: Vue (`@asgardeo/vue`) and other framework packages each have their own component rendering layer. This document covers `@asgardeo/react` only. Other frameworks will need equivalent implementations tracked separately.
