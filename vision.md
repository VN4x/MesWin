# MesWin — Vision

## What We're Building

MesWin is a commercial window measurement platform for on-site photo-based sizing of EU rectangular windows (Rehau Synego and common variants). Field staff capture inner/outer photos with a printed calibration target; the system extracts dimensions to ±5 mm accuracy (10 mm tolerance band), produces annotated sketches, and exports results to ERP/CRM workflows.

The product is designed for **speed, accuracy, and sellability** — starting as an internal tool for our factory, then generalizing for self-hosted and SaaS customers (e.g. Fensterbau Frontale).

## Core Principles

1. **Single-binary simplicity** — Go + PocketBase as one deployable backend; minimal ops for customers and us.
2. **Field-first UX** — Guided camera, offline queue, fast review with manual overrides when CV is uncertain.
3. **Commercial from day one** — Multi-tenant data model, export hooks, licensing path, white-label readiness.
4. **Hybrid accuracy** — Strong CV pipeline plus human override; never block the measurer on model failure.
5. **EU window coverage** — Dataset and type selectors cover common rectangular configurations used in European fenestration.

## Target Users

| User | Need |
|------|------|
| Field measurer | Fast capture, clear overlays, works offline, confidence visible |
| Office reviewer | Correct overrides, PDF/JSON export, order history |
| Factory admin | Multi-tenant orgs, users, branding, ERP integration |
| Customer (self-host) | Docker bundle, docs, license key, minimal maintenance |

## Accuracy & Quality Bar

- **Calibration target**: 100×100 mm+ triangle, 10 mm red lines, durable/weatherproof print.
- **Measurement target**: ±5 mm on test set vs manual tape (10 mm hard tolerance).
- **Output**: Dimensions, confidence scores, annotated SVG/PNG sketches, exportable PDF sheets.

## Window Types in Scope

**Single:** Fixed, Tilt & Turn, Side-Hung Casement, Tilt-Only.

**Double:** Fixed + Operable, Two Tilt & Turn, French-Style Double Casement.

**Triple:** Fixed + Center Operable + Fixed, Operable + Fixed + Operable, Triple Fixed/Operable.

**Multi:** 4+ units, transom combinations, with/without mullions.

## Technology Direction

| Layer | Choice | Rationale |
|-------|--------|-----------|
| Backend | Go + PocketBase | Single binary, auth, realtime, files, SQLite → PostgreSQL path |
| Inference | ONNX Runtime (YOLO) in Go | No separate CV microservice at scale; optional Python prototype early |
| Frontend | SvelteKit + Vite + Capacitor | Web + iOS/Android + PWA from one codebase |
| UI | Tailwind / Flowbite-Svelte | Fast forms, guided camera screens |
| Camera | @capacitor/camera + overlays | Native capture with calibration guides |
| Deploy | PocketBase serves static build + App Stores | One VPS bundle for self-host customers |

**Note:** Early Phase 1 may use a Python/FastAPI + YOLO prototype on Cloud Run for rapid iteration; production path consolidates to ONNX in the Go binary.

## Commercial Goals

- Export to existing ERP/CRM (Markdown, PDF, JSON, CSV).
- Multi-tenancy: separate data per customer company.
- White-label branding options.
- Licensing & billing (Stripe or license keys).
- Self-host Docker-compose bundle + optional hosted SaaS.
- GDPR-conscious security and data handling.

## Success Metrics

| Metric | Target |
|--------|--------|
| Internal time savings | 50%+ reduction on measuring jobs vs tape-only |
| CV accuracy | ±5 mm on held-out real-window test set |
| Commercial traction | 5–10 paying customers within 6–12 months post-launch |
| Field validation | Blind test vs manual tape on real Synego windows before MVP sign-off |

## Long-Term Horizon (Phase 5+)

- Pure Go / on-device ONNX inference (no cloud CV dependency).
- AR preview, bulk orders, analytics dashboard.
- Model retraining from user-corrected data (with consent).
- Fensterbau Frontale demo booth and case-study marketing.

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| CV accuracy insufficient | Large labeled dataset (500–1000+ photos per major type), user overrides, hybrid workflow |
| Commercial readiness gap | Ship for our factory first; generalize from real usage |
| Maintenance burden | Go single-binary, PocketBase admin UI, documented self-host bundle |
| Frontend fragmentation | SvelteKit + Capacitor covers web and native from one tree |

## North Star

A measurer opens the app, selects window type, captures inner/outer photos with the calibration target visible, reviews dimensions and sketch in under two minutes per opening, and exports a PDF the office can send to production — with accuracy the factory trusts and deployment simple enough to sell at trade shows.
