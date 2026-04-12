package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m0xyu/learning-go-shop/internal/utils"
)

func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderResponse, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.SuccessResponse(c, "Order created successfully", orderResponse)
}

func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to retrieve orders", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", orders, *meta)
}

func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID", err)
		return
	}

	orderResponse, err := s.orderService.GetOrder(userID, uint(orderID))
	if err != nil {
		utils.BadRequestResponse(c, "Failed to retrieve order", err)
		return
	}

	utils.SuccessResponse(c, "Order retrieved successfully", orderResponse)
}
