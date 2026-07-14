package handlers

import (
	"fmt"
	"net/http"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("user_handler.Register: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	if err := h.userService.Register(c.Request.Context(), req); err != nil {
		h.logger.Error("user_handler.Register: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.String(http.StatusOK, "a verification code has been sent to your email, please check your email")
}

// принимаем код, после проверяем и создаем запрос для сохранения в БД
func (h *Handler) Verify(c *gin.Context) {
	var req models.Verification
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.logger.Error("user_handler.Verify: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// sending code
	id, err := h.userService.Verification(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("user_handler.Verify: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("verification was successful, user id %d", id))
}

func (h *Handler) Login(c *gin.Context) {
	// get email & password
	var req models.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.logger.Error("user_handler.Login: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, errs.ErrBadRequestBody)
		return
	}

	// login
	tokens, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("user_handler.Login: ", zap.String("error", err.Error()))
		errsToHttp(c.Writer, err)
		return
	}

	// answer
	c.JSON(http.StatusOK, tokens)
}
