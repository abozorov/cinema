package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	AgeLimit    int    `json:"age_limit"`
}

type updateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	AgeLimit    int    `json:"age_limit"`
}

type movieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	AgeLimit    int    `json:"age_limit"`
}

func newMovieResponse(m models.Movie) *movieResponse {
	return &movieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		AgeLimit:    m.AgeLimit,
	}
}

// Create - создание фильма
func (h *Handler) CreateMovie(c *gin.Context) {
	// get movie from body
	req := createMovieRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("movie.Create: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// creating & transform request -> models.Movie
	movie := models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		AgeLimit:    req.AgeLimit,
	}

	id, err := h.movieService.Create(c.Request.Context(), movie)
	if err != nil {
		h.logger.Error("movie.Create: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// transform models.Movie -> response
	c.JSON(http.StatusCreated, fmt.Sprintf("movie id: %d", id))
}

// GetById - получение фильма по id
func (h *Handler) GetMovieById(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("movie.GetById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// get by id
	movie, err := h.movieService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("movie.GetById: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// transform models.Movie -> response
	resp := *newMovieResponse(*movie)

	c.JSON(http.StatusOK, resp)
}

// Update - обновление фильма по id
func (h *Handler) UpdateMovie(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("movie.Update: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// get movie from body
	req := updateMovieRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("movie.Update: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// creating & transform request -> models.Movie
	err = h.movieService.Update(c.Request.Context(), models.Movie{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		AgeLimit:    req.AgeLimit,
	})
	if err != nil {
		h.logger.Error("movie.Update: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.String(http.StatusOK, "Movie updated")
}

// List - получение списка фильмов
func (h *Handler) ListMovies(c *gin.Context) {
	movies, err := h.movieService.List(c.Request.Context())
	if err != nil {
		h.logger.Error("movie.List: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// transform []models.Movie -> []response
	resp := make([]movieResponse, 0, len(movies))
	for _, m := range movies {
		resp = append(resp, *newMovieResponse(m))
	}

	c.JSON(http.StatusOK, resp)
}
