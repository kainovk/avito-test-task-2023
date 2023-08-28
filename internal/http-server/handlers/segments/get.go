package segments

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/models/segment"
)

type GetResponse struct {
	Segments []string `json:"segments"`
}

type SegmentGetter interface {
	GetSegments() ([]*segment.Segment, error)
}

func NewSegmentGetter(log *slog.Logger, segmentGetter SegmentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segments.get.NewSegmentGetter"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		segments, err := segmentGetter.GetSegments()
		if err != nil {
			log.Error("failed to get user segments", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get user segments"))
			return
		}

		log.Info("user segments retrieved")

		segmentSlugs := make([]string, len(segments))
		for i, seg := range segments {
			segmentSlugs[i] = seg.Slug
		}

		response := GetResponse{
			Segments: segmentSlugs,
		}

		render.JSON(w, r, response)
	}
}
