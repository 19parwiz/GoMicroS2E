package handler

import (
	protos "github.com/19parwiz/api-gateway/pkg/protos/gen/golang"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"strconv"
)

type Items struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type Order struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Items  []Items `json:"items"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	log.Println("user_id: ", userID)

	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		log.Println(err)
		return
	}
	items := make([]*protos.OrderItem, len(order.Items))
	for i, item := range order.Items {

		items[i] = &protos.OrderItem{
			ProductId: uint64(item.ProductID),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  uint64(item.Quantity),
		}
	}

	req := &protos.CreateOrderRequest{
		UserId: userID.(uint64),
		Items:  items,
	}

	log.Printf("Create Order: %v", req)

	resp, err := h.Clients.Order.CreateOrder(c.Request.Context(), req)
	if err != nil {
		code, msg := mapGRPCErrorToHTTP(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	jsonBytes, err := protojson.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal response"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}

func (h *Handler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	req := &protos.ListOrdersRequest{
		UserId: userID.(uint64),
		Page:   page,
		Limit:  limit,
	}

	resp, err := h.Clients.Order.ListOrders(c.Request.Context(), req)
	if err != nil {
		code, msg := mapGRPCErrorToHTTP(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	jsonBytes, err := protojson.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal response"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}
