# mesWin — Vision & Requirements

## 1. Project Overview

**mesWin** is a mobile/web app that allows window/door sales technicians to measure existing windows and doors **by taking photos** instead of using measuring tape and paper.

### Core Value

- Replace hours of manual measuring + data entry with a **point-and-shoot** photo workflow.
- Achieve **±5 mm accuracy** (10 mm tolerance) — sufficient for Rehau Synego uPVC systems with triple glazing.
- Generate ready-to-use **dimension sheets + sketches** for sales offers and manufacturing.
- Start as an **internal tool** for our Estonia factory; evolve into a **sellable product** (self-hosted / SaaS) for other manufacturers (target: Fensterbau Frontale Nuremberg).

### Key Innovation

Printed **calibration target** + **ONNX YOLO** computer vision running in the **Go backend** (no separate CV microservice in production).

---

## 2. Business Context

- Small uPVC windows & doors factory (Rehau Synego system dominant).
- Frequent need to measure **inner + outer dimensions** (wall depth / reveals).
- Current pain: manual tape measuring + handwritten notes + ERP entry = very time-consuming and error-prone.
- Goal: save significant time per customer visit (multiple windows per order).

---

## 3. Core User Flow (Shopping-Cart Style)

1. Create new **Measuring Order** — customer data: name, address, contact, notes.
2. Add **Window Items** (like a shopping cart):
   - Select window type from predefined list.
   - Take photos (guided camera with calibration overlays).
   - Review AI results + manual overrides.
3. Save order → generate export (**PDF / Markdown / JSON**) with photos, dimensions, sketches.
4. In office: review, edit, import to existing ERP/CRM (start with CSV/JSON, later API).

**Inner + outer support:** separate or combined photos per window for wall construction analysis.

---

## 4. Technical Stack (Final)

| Layer | Technology |
|-------|------------|
| Frontend | SvelteKit + Vite + Capacitor (web + native iOS/Android + PWA) |
| Backend | Go + PocketBase (embedded, single binary) |
| Computer Vision | Ultralytics YOLO → ONNX → ONNX Runtime in Go (`yalue/onnxruntime_go` or purego) |
| Database | PocketBase (SQLite initially, PostgreSQL option later) |
| Storage | PocketBase file fields for photos |
| Deployment | Single binary / small Docker bundle for self-hosting and commercialization |
| UI | Tailwind + Flowbite-Svelte |

---

## 5. Calibration Target

- Physical **90° triangle** (100×100 mm or 150×150 mm recommended).
- Clear **10 mm red markings**.
- Durable, weatherproof print for field use.

---

## 6. Accuracy & Output

| Requirement | Target |
|-------------|--------|
| Measurement accuracy | ±5 mm on test set |
| Hard tolerance | 10 mm |
| Confidence | Per-dimension scores; manual override when low |
| Sketches | Annotated SVG/PNG with dimension labels |
| Exports | PDF sheet, Markdown, JSON/CSV for ERP |

---

## 7. Window Types in Scope

**Single:** Fixed, Tilt & Turn, Side-Hung Casement, Tilt-Only.

**Double:** Fixed + Operable, Two Tilt & Turn, French-Style Double Casement.

**Triple:** Fixed + Center Operable + Fixed, Operable + Fixed + Operable, Triple Fixed/Operable.

**Multi:** 4+ units, transom combinations, with/without mullions.

Dataset priority: **Tilt & Turn** variants first.

---

## 8. Target Users

| User | Need |
|------|------|
| Field sales technician | Fast photo capture, guided overlays, offline, confidence visible |
| Office reviewer | Overrides, PDF/JSON export, order history |
| Factory admin | Users, org settings, ERP export, branding (later) |
| External customer (future) | Self-host bundle, license key, minimal ops |

---

## 9. Commercial Goals

- Multi-tenant data model from day one (`organization_id` on all business records).
- Export hooks for ERP/CRM (CSV/JSON first).
- White-label branding (logo, colors) for resale.
- Licensing & billing (Stripe or license keys).
- GDPR-conscious photo retention and training-data consent.
- Demo-ready for **Fensterbau Frontale**.

---

## 10. Success Metrics

| Metric | Target |
|--------|--------|
| Time savings (internal) | 50%+ reduction vs tape-only per visit |
| CV accuracy | ±5 mm on held-out real Synego windows |
| Commercial | 5–10 paying customers within 6–12 months post-launch |
| MVP sign-off | Blind tape comparison passes on majority of test openings |

---

## 11. Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| CV accuracy | Large labeled dataset, calibration target, user overrides, hybrid workflow |
| Commercial readiness | Ship for Estonia factory first; generalize from real usage |
| Maintenance | Go single-binary + PocketBase admin UI |
| Field connectivity | Offline photo queue + sync on reconnect |

---

## 12. North Star

A technician arrives at a customer site, creates a measuring order, photographs each window (inner and outer) with the calibration target visible, reviews dimensions and sketches in under two minutes per opening, and exports a PDF the office can turn into a Synego offer — with accuracy the factory trusts.
