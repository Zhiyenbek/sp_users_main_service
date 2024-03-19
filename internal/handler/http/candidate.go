package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GetCandidatesResult struct {
	Candidates []*models.Candidate `json:"candidates"`
	Count      int                 `json:"count"`
}
type skillsReq struct {
	Skills []string `json:"skills"`
}
type InterviewResponse struct {
	Interview []*models.InterviewResults `json:"interviews"`
	Count     int                        `json:"count"`
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

	searchArgs := &models.SearchArgs{
		PageNum:  pageNum,
		PageSize: pageSize,
		Search:   c.Query("search"),
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

func (h *handler) UpdateCandidate(c *gin.Context) {
	req := &models.Candidate{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Errorf("failed to parse request body when updating candidate. %s\n", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	publicID := c.GetString("public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	err := h.service.UpdateCandidateByID(publicID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	res, err := h.service.GetCandidateByPublicID(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, res, nil))
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

func (h *handler) CreateSkillsForCandidate(c *gin.Context) {
	req := &skillsReq{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("failed to parse request body when signing up candidate. %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}
	if len(req.Skills) == 0 {
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	publicID := c.GetString("public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	err := h.service.AddSkillsToCandidate(publicID, req.Skills)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
}

func (h *handler) DeleteSkillsFromCandidate(c *gin.Context) {
	req := &skillsReq{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("failed to parse request body when signing up candidate. %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}
	if len(req.Skills) == 0 {
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	publicID := c.GetString("public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	err := h.service.AddSkillsToCandidate(publicID, req.Skills)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
}

func (h *handler) UpdateCandidateByPublicID(c *gin.Context) {
	req := &models.Candidate{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Errorf("failed to parse request body when updating candidate. %s\n", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	publicID := c.Param("candidate_public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrUserNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	err := h.service.UpdateCandidateByID(publicID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	res, err := h.service.GetCandidateByPublicID(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, res, nil))
}

func (h *handler) DeleteCandidate(c *gin.Context) {
	publicID := c.GetString("public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusUnauthorized, sendResponse(-1, nil, models.ErrPermissionDenied))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	err := h.service.DeleteCandidateByID(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}

func (h *handler) DeleteCandidateByPublicID(c *gin.Context) {
	publicID := c.Param("candidate_public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrUserNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	err := h.service.DeleteCandidateByID(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}

func (h *handler) GetCandidateInterviewsByID(c *gin.Context) {
	publicID := c.Param("candidate_public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
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

	res, count, err := h.service.CandidatesService.GetInterviewsByPublicID(publicID, searchArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, InterviewResponse{
		Interview: res,
		Count:     count,
	}, nil))
}

func (h *handler) GetCandidateInterviews(c *gin.Context) {
	publicID := c.GetString("public_id")
	if err := h.service.CandidatesService.Exists(publicID); err != nil {
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

	res, count, err := h.service.CandidatesService.GetInterviewsByPublicID(publicID, searchArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, InterviewResponse{
		Interview: res,
		Count:     count,
	}, nil))
}
