package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *handler) GetRecruiter(c *gin.Context) {
	publicID := c.Param("recruiter_public_id")
	if err := h.service.RecruiterService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrUserNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	res, err := h.service.RecruiterService.GetRecruiter(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, res, nil))
}
func (h *handler) GetRecruiterInterviewsByID(c *gin.Context) {
	publicID := c.Param("recruiter_public_id")
	if err := h.service.RecruiterService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrUserNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}

	searchArgs := &models.SearchArgs{
		PageNum:  pageNum,
		PageSize: pageSize,
		Search:   c.Query("search"),
	}

	res, count, err := h.service.RecruiterService.GetInterviewsByPublicID(publicID, searchArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, InterviewResponse{
		Interview: res,
		Count:     count,
	}, nil))
}

func (h *handler) GetRecrutierInterviews(c *gin.Context) {
	publicID := c.Param("recruiter_public_id")
	if err := h.service.RecruiterService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrUserNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}

	searchArgs := &models.SearchArgs{
		PageNum:  pageNum,
		PageSize: pageSize,
		Search:   c.Query("search"),
	}

	res, count, err := h.service.RecruiterService.GetInterviewsByPublicID(publicID, searchArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, InterviewResponse{
		Interview: res,
		Count:     count,
	}, nil))
}

func (h *handler) GetRecruiterInterviews(c *gin.Context) {
	publicID := c.GetString("public_id")
	if err := h.service.RecruiterService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrUserNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = models.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = models.DefaultPageSize
	}

	searchArgs := &models.SearchArgs{
		PageNum:  pageNum,
		PageSize: pageSize,
		Search:   c.Query("search"),
	}

	res, count, err := h.service.RecruiterService.GetInterviewsByPublicID(publicID, searchArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, InterviewResponse{
		Interview: res,
		Count:     count,
	}, nil))
}
