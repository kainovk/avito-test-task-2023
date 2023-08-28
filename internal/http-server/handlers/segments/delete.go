package segments

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/storage"
)

// DeleteResponse is the response structure for segment deletion.
// @Summary Delete {object}
// @Description Response structure for segment deletion.
// @Accept json
// @Produce json
type DeleteResponse struct {
	response.Response
}

type SegmentDeleter interface {
	DeleteSegmentBySlug(segmentName string) error
}

// NewSegmentDeleter handles the HTTP request for deleting a segment by slug.
//
// @Summary Delete a segment
// @Description Delete a segment by its slug.
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Segment slug to delete"
// @Success 200 {object} DeleteResponse
// @Router /segments/{slug} [delete]
func NewSegmentDeleter(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segments.delete.NewSegmentDeleter"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			log.Info("slug param is empty")

			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		err := segmentDeleter.DeleteSegmentBySlug(slug)
		if errors.Is(err, storage.ErrSegmentNotExists) {
			log.Info("segment does not exist", slog.String("slug", slug))

			render.JSON(w, r, response.Error("segment does not exist"))
			return
		}
		if err != nil {
			log.Error("failed to delete segment", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete segment"))
			return
		}

		log.Info("segment deleted", slog.String("slug", slug))

		render.JSON(w, r, DeleteResponse{
			Response: response.OK(),
		})
	}
}
