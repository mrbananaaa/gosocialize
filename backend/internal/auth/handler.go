package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
	"github.com/mrbananaaa/gosocialize/internal/platform/httpx"
	"github.com/mrbananaaa/gosocialize/internal/platform/validator"
	"github.com/mrbananaaa/gosocialize/pkg/logger"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email" example:"mail.me@here.dev"`
	Username string `json:"username" validate:"required" example:"21charmaxyoo"`
	Password string `json:"password" validate:"required" example:"@Itskindahardtodefine"`
	Name     string `json:"name" validate:"required" example:"John Doe Is Ya Name"`
}

type SignUpResponse struct {
	ID        uuid.UUID `json:"id" format:"uuid" example:"b14136a3-284f-4620-82dc-290478530cee"`
	Email     string    `json:"email" format:"email" example:"mail.me@here.dev"`
	Username  string    `json:"username" example:"21charmaxyoo"`
	Name      string    `json:"name" example:"John Doe Is Ya Name"`
	CreatedAt time.Time `json:"created_at" format:"dateTime" example:"2026-06-24T07:29:51.705430194Z"`
	UpdatedAt time.Time `json:"updated_at" format:"dateTime" example:"2026-06-24T07:29:51.705430194Z"`
}

// SignUp godoc
//
// @Summary register new user
// @Description register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body SignUpRequest true "signup request body"
// @Success 201 {object} httpx.Response{data=SignUpResponse}
// @Failure 400 {object} httpx.ErrorResponse
// @Router /auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to parse request body", logger.ErrorField(err))
		httpx.Error(w, apperr.DecodeBodyErr(err))
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Error("failed to valdiate request body", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	u, err := h.authService.Register(r.Context(), RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		logger.Error("Failed to register user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.Created(w, SignUpResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

type SignInRequest struct {
	Username string `json:"username" validate:"required" example:"21charmaxyoo"`
	Password string `json:"password" validate:"required" example:"@Itskindahardtodefine"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Signin godoc
//
// @Summary user login
// @Description login and get the access token
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body SignInRequest true "signin request body"
// @Success 200 {object} httpx.Response{data=SignInResponse}
// @Failure 400 {object} httpx.ErrorResponse
// @Router /auth/signin [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to parse request body", logger.ErrorField(err))
		httpx.Error(w, apperr.DecodeBodyErr(err))
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Error("failed to valdiate request body", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	t, err := h.authService.Login(r.Context(), LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("Failed to logging in user", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	httpx.SetRefreshToken(w, t.RefreshToken)

	httpx.OK(w, SignInResponse{
		AccessToken: t.AccessToken,
	})
}

type RefreshResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Refresh godoc
//
// @Summary refresh token
// @Description get a new access token using refresh token from cookie
// @Tags auth
// @Produce json
// @Success 200 {object} httpx.Response{data=RefreshResponse}
// @Failure 401 {object} httpx.ErrorResponse
// @Router /auth/refresh [get]
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh, err := httpx.GetRefreshToken(r)
	if err != nil {
		logger.Error("couldn't get refresh token from cookies", logger.ErrorField(err))
		httpx.Error(w, err)
		return
	}

	token, err := h.authService.Refresh(r.Context(), RefreshInput{
		RefreshToken: refresh,
	})
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.OK(w, RefreshResponse{
		AccessToken: token.AccessToken,
	})
}

// SignOut godoc
//
// @Summary logout
// @Description logout
// @Tags auth
// @Produce json
// @Success 200 {object} httpx.Response
// @Router /auth/signout [post]
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	httpx.CleanRefreshToken(w)

	httpx.OK(w, nil)
}
