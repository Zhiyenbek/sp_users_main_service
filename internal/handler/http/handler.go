package handler

import (
	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/service"
	"github.com/Zhiyenbek/users-auth-service/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handler struct {
	service *service.Service
	cfg     *config.Configs
	logger  *zap.SugaredLogger
}

type Handler interface {
	InitRoutes() *gin.Engine
}

func New(services *service.Service, logger *zap.SugaredLogger, cfg *config.Configs) Handler {
	return &handler{
		service: services,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/account", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.GetMe)
	router.GET("/candidates", h.GetCandidates)
	router.GET("/candidate/:candidate_public_id", h.GetCandidateByPublicID)
	router.PUT("/candidate", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.UpdateCandidate)
	router.PUT("/candidate/:candidate_public_id", h.UpdateCandidateByPublicID)
	router.DELETE("/candidate/:candidate_public_id", h.DeleteCandidateByPublicID)
	router.DELETE("/candidate", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.DeleteCandidate)
	router.POST("/candidate/skills", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.CreateSkillsForCandidate)
	router.DELETE("/candidate/skills", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.DeleteSkillsFromCandidate)
	router.GET("/candidate/:candidate_public_id/interviews", h.GetCandidateInterviewsByID)
	router.GET("/candidate/interviews", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.GetCandidateInterviews)
	router.GET("/recruiter/:recruiter_public_id", h.GetRecruiter)
	router.GET("/recruiter/:recruiter_public_id/interviews", h.GetRecruiterInterviewsByID)
	router.GET("/recruiter/interviews", middleware.VerifyToken(h.cfg.Token.TokenSecret, h.logger), h.GetRecruiterInterviews)
	router.POST("/company", h.CreateCompany)
	router.GET("/companies", h.GetCompanies)
	router.GET("/company/:public_id", h.GetCompany)
	router.PUT("/company", h.UpdateCompany)
	return router
}

func sendResponse(status int, data interface{}, err error) gin.H {
	var errResponse gin.H
	if err != nil {
		errResponse = gin.H{
			"message": err.Error(),
		}
	} else {
		errResponse = nil
	}

	return gin.H{
		"data":   data,
		"status": status,
		"error":  errResponse,
	}
}
