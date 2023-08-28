package users

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/storage"
)

type SaveRequest struct {
	Name string `json:"name" validate:"required"`
}

type SaveResponse struct {
	response.Response
}

type UserSaver interface {
	SaveUser(name string) error
}

// NewUserSaver handles the HTTP request for saving a user.
//
// @Summary Save a user
// @Description Save a new user with the provided name.
// @Tags users
// @Accept json
// @Produce json
// @Param request body SaveRequest true "Request body"
// @Success 200 {object} SaveResponse
// @Router /users [post]
func NewUserSaver(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.save.NewUserSaver"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SaveRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = userSaver.SaveUser(req.Name)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("name", req.Name))

			render.JSON(w, r, response.Error("user already exists"))
			return
		}
		if err != nil {
			log.Error("failed to create user", sl.Err(err))

			render.JSON(w, r, response.Error("failed to create user"))
			return
		}

		log.Info("user created")

		render.JSON(w, r, SaveResponse{
			Response: response.OK(),
		})
	}
}
