package users

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/models/segment"
)

type GetSegmentsResponse struct {
	Segments []string `json:"segments"`
}

type UserSegmentsGetter interface {
	GetUserSegments(userID int64) ([]*segment.Segment, error)
}

func NewUserSegmentsGetter(log *slog.Logger, userSegmentsGetter UserSegmentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.get-segments.NewUserSegmentsGetter"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userIDStr := chi.URLParam(r, "user_id")
		if userIDStr == "" {
			log.Info("user_id param is empty")

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Error("failed to parse user_id")

			render.JSON(w, r, resp.Error("invalid request"))
		}

		segments, err := userSegmentsGetter.GetUserSegments(int64(userID))
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

		response := GetSegmentsResponse{
			Segments: segmentSlugs,
		}

		render.JSON(w, r, response)
	}
}
