# CV — Training & ONNX export only

**Runtime inference is Go + ONNX only.** This folder is for offline model work.

## Workflow

1. Label images in **Roboflow** (500–1000+ per major window type; start with Tilt & Turn).
2. Train with **Ultralytics YOLOv8/v11** (`train.py` to be added).
3. Export to ONNX → `../models/yolo-windows-v1.onnx`.
4. Go backend loads the model via `yalue/onnxruntime_go` in `backend/internal/cv/`.

## Detection classes

`calibration_target`, `frame`, `sash`, `corner`, `mullion`, `openable`, `fixed`

## What does NOT live here

- No FastAPI service
- No Cloud Run deployment
- No production HTTP CV endpoint

Homography, dimension extraction, and sketch rendering are implemented in Go.

## Accuracy target

±5 mm on held-out test set (10 mm hard tolerance).
