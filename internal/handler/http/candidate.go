package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhiyenbek/sp_users_main_service/internal/models"
	"github.com/gin-gonic/gin"
)

type GetCandidatesResult struct {
	Candidates []*models.Candidate `json:"candidates"`
	Count      int                 `json:"count"`
}

func (h *handler) GetCandidates(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}
	skills := c.QueryArray("skills")
	searchArgs := &models.SearchArgs{
		PageNum:  pageNum,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Skills:   skills, // Assign skills as an array
	}
	res, count, err := h.service.GetCandidatesBySearch(searchArgs)
	if err != nil {
		var errMsg error
		var code int
		switch {
		case errors.Is(err, models.ErrUsernameExists):
			errMsg = models.ErrUsernameExists
			code = http.StatusBadRequest
		default:
			errMsg = models.ErrInternalServer
			code = http.StatusInternalServerError
		}
		c.JSON(code, sendResponse(-1, nil, errMsg))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, GetCandidatesResult{
		Candidates: res,
		Count:      count,
	}, nil))
}

func (h *handler) GetCandidateByPublicID(c *gin.Context) {
	publicID := c.Param("candidate_public_id")
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
}
