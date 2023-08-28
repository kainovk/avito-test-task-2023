package segments

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/storage"
)

type DeleteResponse struct {
	resp.Response
}

type SegmentDeleter interface {
	DeleteSegmentBySlug(segmentName string) error
}

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

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := segmentDeleter.DeleteSegmentBySlug(slug)
		if errors.Is(err, storage.ErrSegmentNotExists) {
			log.Info("segment does not exist", slog.String("slug", slug))

			render.JSON(w, r, resp.Error("segment does not exist"))
			return
		}
		if err != nil {
			log.Error("failed to delete segment", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete segment"))
			return
		}

		log.Info("segment deleted", slog.String("slug", slug))

		render.JSON(w, r, DeleteResponse{
			Response: resp.OK(),
		})
	}
}
