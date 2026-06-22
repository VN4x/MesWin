# ADR 001: CV runtime — ONNX in Go (not Python/FastAPI)

**Status:** Accepted  
**Date:** 2026-06-22

## Context

mesWin needs computer vision for window measurement. Two approaches were considered:

1. **Python/FastAPI** microservice (Ultralytics YOLO, OpenCV) on Cloud Run/Docker
2. **ONNX export** + **ONNX Runtime in Go**, embedded in the PocketBase binary

## Decision

**Use ONNX in Go as the only production inference path.**

Python remains in `cv/` strictly for training and exporting `.onnx` files.

## Rationale

| Factor | ONNX in Go | Python/FastAPI |
|--------|------------|----------------|
| Deployment | Single binary / one Docker image | Second service, networking, versioning |
| Commercial self-host | Simple bundle for customers | Extra container, ops burden |
| Latency | In-process, no HTTP hop | Network round-trip per photo |
| Offline / field | Server-side inference without cloud CV | Depends on service availability |
| Alignment with stack | Matches PocketBase + Go architecture | Splits stack across languages at runtime |

Training ergonomics favor Python, but that work is **offline and infrequent**. Runtime inference is **hot path** and must match the sellable single-binary product.

## Consequences

- Add `yalue/onnxruntime_go` (or purego build tag) to `backend/internal/cv/`
- Implement homography + dimension math in Go (`gocv` or stdlib image as needed)
- `cv/` contains only training scripts + export; no FastAPI
- Phase 1 deliverable: ONNX model + Go inference replacing current stub engine

## Alternatives rejected

- **Python/FastAPI microservice:** rejected for operational complexity and weak single-binary story
- **Dual path (prototype Python, later Go):** rejected — team commits to one runtime from the start
