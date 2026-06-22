# MesWin — Skills

Domain skills for AI assistants and developers working in this repository. Load the relevant skill before making changes in that area.

---

## Skill: `meswin-architecture`

**Use when:** Planning features, choosing where code lives, reviewing cross-cutting changes.

### Stack

- **Backend:** Go 1.22+, PocketBase (embedded or sidecar), custom routes via Go hooks or separate HTTP handlers.
- **DB:** PocketBase/SQLite default; design for PostgreSQL migration (org-scoped queries).
- **CV:** ONNX Runtime in Go (production); Python/YOLO only for training and optional Phase 1 prototype.
- **Frontend:** SvelteKit 2, Vite, Tailwind, Flowbite-Svelte, Capacitor 6+.
- **Auth:** PocketBase auth; JWT/session via JS SDK on client.

### Layout (conventional)

```
/
├── backend/          # Go module: PocketBase hooks, ONNX, exports
├── pb/               # PocketBase migrations, schema
├── frontend/         # SvelteKit app
├── cv/               # Training scripts, ONNX export (not runtime Python in prod)
├── docs/             # vision, taskflow, API
└── deploy/           # Docker, compose, install scripts
```

### Rules

- Business logic that touches measurements or exports lives in **Go**, not PocketBase JS hooks.
- Never store raw measurement truth only in client state; persist server-side with override audit fields.
- All collections include `organization_id` for multi-tenant queries.

---

## Skill: `meswin-pocketbase`

**Use when:** Defining collections, relations, rules, or admin-facing data.

### Core collections

| Collection | Key fields |
|------------|------------|
| `organizations` | name, slug, branding_json, license_key |
| `users` | role (field, reviewer, admin), organization_id |
| `customers` | name, address, contact, organization_id |
| `measuring_orders` | customer_id, status, cart-style line refs, organization_id |
| `window_items` | order_id, window_type, inner_photos[], outer_photos[], cv_result_json, overrides_json, confidence |
| `files` | storage ref, type (photo, sketch, export), parent relation |

### Access rules

- Users read/write only records where `organization_id` matches their org.
- Field role: create orders and window items; cannot delete org.
- Reviewer: override dimensions; trigger export.
- Admin: manage users and branding.

### Patterns

- Use PocketBase file fields for photos; store CV JSON separately for queryability.
- Keep `overrides_json` with `{ field, old, new, user_id, at }` audit trail.

---

## Skill: `meswin-cv`

**Use when:** Training models, ONNX export, Go inference, sketch generation.

### Detection classes

`calibration_target`, `frame`, `sash`, `corner`, `mullion`, `openable`, `fixed`.

### Pipeline steps

1. Detect calibration target → compute scale (mm/px).
2. Detect frame/sashes/corners → fit quadrilateral / homography.
3. Inner + outer photo pairing per window item.
4. Extract width, height, diagonals, segment breakdown.
5. Render SVG/PNG sketch with dimension labels.
6. Emit JSON: `{ dimensions, confidence, sketch_url, warnings[] }`.

### Accuracy

- Target ±5 mm; flag items below confidence threshold for manual review.
- Always return confidence per dimension; never silently round away uncertainty.

### Training

- Roboflow for labeling; Ultralytics YOLOv8/v11 for training.
- Export to ONNX opset compatible with onnxruntime-go.
- Version models (`models/yolo-windows-v{N}.onnx`) and record in config.

---

## Skill: `meswin-go-backend`

**Use when:** Writing handlers, export logic, ONNX integration, tests.

### Conventions

- Standard library + minimal deps; `chi` or PocketBase router patterns as established in repo.
- Context timeouts on all CV calls (even in-process ONNX).
- Idempotent upload handlers (hash or client idempotency key).

### Key endpoints (illustrative)

```
POST   /api/orders                    # create cart-style order
POST   /api/orders/:id/windows        # add window item
POST   /api/windows/:id/photos        # upload inner/outer
POST   /api/windows/:id/process       # trigger CV
PATCH  /api/windows/:id/overrides     # manual corrections
GET    /api/orders/:id/export.pdf     # PDF with sketches
GET    /api/orders/:id/export.json    # ERP payload
```

### Export

- PDF: photos + sketch + dimension table + customer header.
- JSON/CSV: stable field names documented in `docs/export-schema.md` when created.

---

## Skill: `meswin-sveltekit-frontend`

**Use when:** UI flows, camera, offline sync, Capacitor builds.

### Conventions

- Svelte 5 runes or project-established reactivity pattern — match existing files.
- PocketBase JS SDK in `lib/pb.ts` singleton; typed wrappers for collections.
- Tailwind + Flowbite-Svelte for forms; no one-off CSS unless necessary.

### Critical screens

1. **Org select / login**
2. **Order list + new order**
3. **Window type picker** (EU types from vision.md)
4. **Guided camera** — inner/outer mode, overlay component, capture queue
5. **Processing / review** — confidence badges, override inputs
6. **Export** — format picker, share sheet on native

### Offline

- IndexedDB or Capacitor Storage for pending photos and draft orders.
- Background sync on `online` event; show sync status in header.
- Never discard failed uploads silently.

### Capacitor

- `@capacitor/camera` for capture; test permissions on iOS and Android.
- Build: `npm run build && npx cap sync`.

---

## Skill: `meswin-commercial`

**Use when:** Licensing, multi-tenant branding, GDPR, customer deploy bundles.

### Multi-tenancy

- Subdomain or org slug in app config; all API calls scoped by org.
- Branding: logo URL, primary color in `organizations.branding_json`.

### Licensing

- Support license key validation endpoint or Stripe subscription webhook.
- Graceful degradation when license expired (read-only exports).

### GDPR

- Photo retention policy configurable per org.
- Explicit consent before using corrections for model retraining.
- Document data processor role in customer bundle.

### Deploy bundle

- `docker-compose.yml`: PocketBase + volume mounts.
- Env template: `PB_ENCRYPTION_KEY`, `LICENSE_KEY`, `MODEL_PATH`.
- One-page quickstart in `deploy/README.md`.

---

## Skill: `meswin-testing`

**Use when:** Writing tests or validating accuracy.

### Backend

- Table-driven Go tests for dimension math and export builders.
- Golden-file tests for JSON export shape.

### CV

- Held-out image set with known tape measurements; CI checks max error ≤ 10 mm, target 5 mm avg.

### Frontend

- Playwright for critical flows (login → add window → mock upload).
- Manual device checklist for camera overlays on iOS/Android.

---

## Quick Reference: Window Type Enum

Use consistent slugs in DB and UI:

```
single_fixed
single_tilt_turn
single_side_hung
single_tilt_only
double_fixed_operable
double_two_tilt_turn
double_french_casement
triple_fixed_center_operable_fixed
triple_operable_fixed_operable
triple_mixed
multi_4plus
multi_transom
multi_mullion
```
