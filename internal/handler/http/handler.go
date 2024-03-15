package handler

import (
	"github.com/Zhiyenbek/sp_users_main_service/config"
	"github.com/Zhiyenbek/sp_users_main_service/internal/service"
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
	router.GET("/account", h.GetMe)
	router.GET("/candidates", h.GetCandidates)
	router.GET("/candidate/:candidate_public_id", h.GetCandidateByPublicID)
	router.PUT("/candidate", middleware.VerifyToken(h.cfg.Token.TokenSecret), h.UpdateCandidate)
	router.PUT("/candidate/:candidate_public_id", h.UpdateCandidateByPublicID)
	router.DELETE("/candidate/:candidate_public_id", h.DeleteCandidateByPublicID)
	router.DELETE("/candidate", middleware.VerifyToken(h.cfg.Token.TokenSecret), h.DeleteCandidate)
	router.POST("/candidate/skills", middleware.VerifyToken(h.cfg.Token.TokenSecret), h.CreateSkillsForCandidate)
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
