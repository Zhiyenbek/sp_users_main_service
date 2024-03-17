package handler

import (
	"errors"
	"net/http"

	"github.com/Zhiyenbek/sp_users_main_service/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *handler) GetMe(c *gin.Context) {
	publicID := c.GetString("public_id")
	role := c.GetString("role")
	switch role {
	case "candidate":
		if err := h.service.CandidatesService.Exists(publicID); err != nil {
			if errors.Is(err, models.ErrPermissionDenied) {
				c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
				return
			}
			c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
			return
		}
		res, err := h.service.GetCandidateByPublicID(publicID)
		if err != nil {
			var errMsg error
			var code int
			switch {
			case errors.Is(err, models.ErrUserNotFound):
				errMsg = models.ErrUserNotFound
				code = http.StatusNotFound
			default:
				errMsg = models.ErrInternalServer
				code = http.StatusInternalServerError
			}
			c.JSON(code, sendResponse(-1, nil, errMsg))
			return
		}

		c.JSON(http.StatusOK, sendResponse(0, res, nil))
	case "recruiter":
		if err := h.service.RecruiterService.Exists(publicID); err != nil {
			if errors.Is(err, models.ErrPermissionDenied) {
				c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
				return
			}
			c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
			return
		}
		res, err := h.service.GetRecruiter(publicID)
		if err != nil {
			var errMsg error
			var code int
			switch {
			case errors.Is(err, models.ErrUserNotFound):
				errMsg = models.ErrUserNotFound
				code = http.StatusNotFound
			default:
				errMsg = models.ErrInternalServer
				code = http.StatusInternalServerError
			}
			c.JSON(code, sendResponse(-1, nil, errMsg))
			return
		}

		c.JSON(http.StatusOK, sendResponse(0, res, nil))
	}

}
