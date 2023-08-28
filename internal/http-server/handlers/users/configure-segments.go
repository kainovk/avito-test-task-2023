package users

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	resp "avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ConfigureSegmentsRequest struct {
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToDelete []string `json:"segments_to_delete"`
}

type ConfigureSegmentsResponse struct {
	resp.Response
}

type UserSegmentConfigurer interface {
	ConfigureUserSegments(userID int64, segAdd []string, segDel []string) error
}

func NewUserSegmentConfigurer(log *slog.Logger, userSegmentConfigurer UserSegmentConfigurer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.configure-segments.NewUserSegmentConfigurer"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req ConfigureSegmentsRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

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

		err = userSegmentConfigurer.ConfigureUserSegments(int64(userID), req.SegmentsToAdd, req.SegmentsToDelete)
		if err != nil {
			log.Error("failed to configure user segments", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to configure user segments"))
			return
		}

		log.Info("user segments updated")

		render.JSON(w, r, ConfigureSegmentsResponse{
			Response: resp.OK(),
		})
	}
}
