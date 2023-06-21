package handler

import (
	"compare/internal/models"
	"compare/internal/service"
	"compare/pkg/logging"
	"github.com/gin-gonic/gin"
)

var logger = logging.GetLogger()

type Handler struct {
	Engine  *gin.Engine
	Service *service.Service
}

func NewHandler(engine *gin.Engine, service *service.Service) *Handler {
	return &Handler{
		Engine:  engine,
		Service: service,
	}
}

// ===================================================================>

func (h *Handler) InitRoutes() {
	api := h.Engine.Group("/api/v1")
	api.POST("/humo_payment", h.HumoPayment)
	api.POST("/partner_payment", h.PartnerPayment)
	//api.GET("/compare", h.Compare)
}

// ===================================================================>

func (h *Handler) HumoPayment(c *gin.Context) {
	var hp []models.HumoPayment

	err := c.ShouldBindJSON(&hp)
	if err != nil {
		logger.Error(err)
	}

	err = h.Service.HumoPayment(hp)
	if err != nil {
		logger.Error(err)
		return
	}

	c.JSON(201, gin.H{
		"msg": "humo payment created",
	})
}

func (h *Handler) PartnerPayment(c *gin.Context) {
	var pp []models.PartnerPayment

	err := c.ShouldBindJSON(&pp)
	if err != nil {
		logger.Error(err)
	}

	err = h.Service.PartnerPayment(pp)
	if err != nil {
		logger.Error(err)
		return
	}

	err = h.Service.SaveToFile(pp)
	if err != nil {
		logger.Error(err)
		return
	}

	c.JSON(201, gin.H{
		"message": "partner payment created",
		"code":    28,
	})
}

//func (h *Handler) Compare(c *gin.Context) {
//	out, err := h.Service.Compare()
//	if err != nil {
//		logger.Error(err)
//		c.JSON(500, gin.H{
//			"msg": "internal server error",
//		})
//		return
//	}
//
//	//err = h.Service.SaveOutput()
//	//if err != nil {
//	//	logger.Error(err)
//	//	c.JSON(500, gin.H{
//	//		"msg": "internal server error",
//	//	})
//	//	return
//	//}
//
//	c.JSON(200, gin.H{
//		"msg": out,
//	})
//}
