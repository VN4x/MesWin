package api

import (
	"context"
	"net/http"
	"time"

	"github.com/vn4x/meswin/internal/cv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// Handler registers custom mesWin API routes.
type Handler struct {
	app      core.App
	cvEngine cv.Engine
}

func NewHandler(app core.App, cvEngine cv.Engine) *Handler {
	return &Handler{app: app, cvEngine: cvEngine}
}

func (h *Handler) Register(e *core.ServeEvent) {
	e.Router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/meswin/windows/:id/process",
		Handler: func(c echo.Context) error {
			return h.processWindow(c)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(h.app),
			apis.RequireRecordAuth(),
		},
	})

	e.Router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/meswin/orders/:id/export.json",
		Handler: func(c echo.Context) error {
			return h.exportOrderJSON(c)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(h.app),
			apis.RequireRecordAuth(),
		},
	})
}

func (h *Handler) processWindow(c echo.Context) error {
	windowID := c.PathParam("id")
	if windowID == "" {
		return apis.NewBadRequestError("missing window id", nil)
	}

	record, err := h.app.Dao().FindRecordById("window_items", windowID)
	if err != nil {
		return apis.NewNotFoundError("window item not found", err)
	}

	windowType := record.GetString("window_type")
	innerPhotos := record.GetStringSlice("inner_photos")
	outerPhotos := record.GetStringSlice("outer_photos")

	var innerPath, outerPath string
	if len(innerPhotos) > 0 {
		innerPath = innerPhotos[0]
	}
	if len(outerPhotos) > 0 {
		outerPath = outerPhotos[0]
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	result, err := h.cvEngine.Process(ctx, innerPath, outerPath, windowType)
	if err != nil {
		return apis.NewBadRequestError("cv processing failed", err)
	}

	record.Set("cv_result", result)
	record.Set("status", "review")
	record.Set("confidence", result.Confidence)

	if err := h.app.Dao().SaveRecord(record); err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to save cv result", err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"window_id": windowID,
		"result":    result,
	})
}

func (h *Handler) exportOrderJSON(c echo.Context) error {
	orderID := c.PathParam("id")
	if orderID == "" {
		return apis.NewBadRequestError("missing order id", nil)
	}

	order, err := h.app.Dao().FindRecordById("measuring_orders", orderID)
	if err != nil {
		return apis.NewNotFoundError("measuring order not found", err)
	}

	windows, err := h.app.Dao().FindRecordsByFilter(
		"window_items",
		"order = {:order}",
		"-created",
		0,
		0,
		map[string]any{"order": orderID},
	)
	if err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to load window items", err)
	}

	items := make([]map[string]any, 0, len(windows))
	for _, w := range windows {
		items = append(items, map[string]any{
			"id":           w.Id,
			"window_type":  w.GetString("window_type"),
			"label":        w.GetString("label"),
			"status":       w.GetString("status"),
			"cv_result":    w.Get("cv_result"),
			"overrides":    w.Get("overrides"),
			"inner_photos": w.GetStringSlice("inner_photos"),
			"outer_photos": w.GetStringSlice("outer_photos"),
		})
	}

	payload := map[string]any{
		"order": map[string]any{
			"id":       order.Id,
			"status":   order.GetString("status"),
			"notes":    order.GetString("notes"),
			"customer": order.GetString("customer"),
			"created":  order.GetString("created"),
			"updated":  order.GetString("updated"),
		},
		"windows":     items,
		"exported_at": time.Now().UTC(),
	}

	return c.JSON(http.StatusOK, payload)
}

// ensure models import is used for RequestInfo typing in future middleware.
var _ = models.RequestInfo{}
