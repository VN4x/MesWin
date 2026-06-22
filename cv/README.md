# Training data & ONNX export (Phase 1)

Place Ultralytics training scripts and Roboflow export ingest here.

Production inference runs in `backend/internal/cv/` via ONNX Runtime.

## Classes

`calibration_target`, `frame`, `sash`, `corner`, `mullion`, `openable`, `fixed`

## Next steps

1. Export labeled dataset from Roboflow (YOLOv8 format).
2. Train / fine-tune YOLOv8 or v11.
3. Export to ONNX → `models/yolo-windows-v1.onnx`.
4. Wire `backend/internal/cv/onnx_engine.go` (replace stub).
