package handlers

import (
	"net/http"
	"strconv"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type updateUser struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type responseUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

func newResponseUser(u models.User) *responseUser {
	return &responseUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Role:  u.Role,
	}
}

func (h *Handler) GetGyId(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("user.GetMe: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// get by id
	usr, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("user.GetMe: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// transform models.User -> user
	resp := *newResponseUser(*usr)

	// write response
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateById(c *gin.Context) {
	// get id из URL параметра
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("user.GetMe: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// get user
	usr := updateUser{}
	if err := c.ShouldBindJSON(&usr); err != nil {
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.User -> user
	err = h.userService.Update(c.Request.Context(), models.User{
		ID:    id,
		Name:  usr.Name,
		Phone: usr.Phone,
	})
	if err != nil {
		h.logger.Error("user_handler.UpdateMe: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.String(http.StatusOK, "User updated")
}
