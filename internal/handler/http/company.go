package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/gin-gonic/gin"
)

type GetCompaniesResult struct {
	Companies []*models.Company `json:"companies"`
	Count     int               `json:"count"`
}

func (h *handler) CreateCompany(c *gin.Context) {
	company := &models.Company{}
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	if err := h.service.CompanyService.CreateCompany(company); err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
}

func (h *handler) UpdateCompany(c *gin.Context) {
	publicID := c.Param("public_id")

	company := &models.Company{}
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	company.PublicID = publicID

	if err := h.service.CompanyService.UpdateCompany(company); err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}

func (h *handler) GetCompany(c *gin.Context) {
	publicID := c.Param("public_id")
	if err := h.service.CompanyService.Exists(publicID); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrCompanyNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	company, err := h.service.CompanyService.GetCompany(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, company, nil))
}

func (h *handler) GetCompanies(c *gin.Context) {
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

	companies, count, err := h.service.CompanyService.GetCompanies(searchArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, GetCompaniesResult{
		Companies: companies,
		Count:     count,
	}, nil))
}
