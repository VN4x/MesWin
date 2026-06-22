package cv

import (
	"context"
	"errors"
	"time"
)

// Result holds dimension extraction output from the CV pipeline.
type Result struct {
	WidthMM     float64            `json:"width_mm"`
	HeightMM    float64            `json:"height_mm"`
	DiagonalMM  float64            `json:"diagonal_mm,omitempty"`
	Confidence  float64            `json:"confidence"`
	Warnings    []string           `json:"warnings,omitempty"`
	SketchSVG   string             `json:"sketch_svg,omitempty"`
	Dimensions  map[string]float64 `json:"dimensions,omitempty"`
	ProcessedAt time.Time          `json:"processed_at"`
}

// Engine processes window photos and returns measurements.
type Engine interface {
	Process(ctx context.Context, innerPhotoPath, outerPhotoPath, windowType string) (*Result, error)
}

// StubEngine returns placeholder measurements until ONNX model is wired.
type StubEngine struct{}

func NewStubEngine() *StubEngine {
	return &StubEngine{}
}

func (s *StubEngine) Process(ctx context.Context, innerPhotoPath, outerPhotoPath, windowType string) (*Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if innerPhotoPath == "" && outerPhotoPath == "" {
		return nil, errors.New("at least one photo path is required")
	}

	return &Result{
		WidthMM:    1200,
		HeightMM:   1400,
		Confidence: 0.42,
		Warnings: []string{
			"stub engine active — replace with ONNX inference",
			"manual review required",
		},
		SketchSVG: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 240"><rect x="10" y="10" width="180" height="220" fill="none" stroke="#e11" stroke-width="2"/><text x="100" y="130" text-anchor="middle" font-size="12">stub sketch</text></svg>`,
		Dimensions: map[string]float64{
			"width_mm":  1200,
			"height_mm": 1400,
		},
		ProcessedAt: time.Now().UTC(),
	}, nil
}
