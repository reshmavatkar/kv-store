package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ----------- Handler with Router Methods -----------

type Handler struct {
	store StoreClient
}

func NewHandler(store StoreClient) *Handler {
	return &Handler{store: store}
}

// ----------- Named Request Structs -----------

type PutRequestBody struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type GetRequestParam struct {
	Key string `uri:"key" binding:"required"`
}

type DeleteRequestParam struct {
	Key string `uri:"key" binding:"required"`
}

func (h *Handler) PutValue(c *gin.Context) {
	var req PutRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.store.Put(c.Request.Context(), req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store value"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetValue(c *gin.Context) {
	var req GetRequestParam
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing key"})
		return
	}
	val, err := h.store.Get(c.Request.Context(), req.Key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": req.Key, "value": val})
}

func (h *Handler) DeleteValue(c *gin.Context) {
	var req DeleteRequestParam
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing key"})
		return
	}
	if err := h.store.Delete(c.Request.Context(), req.Key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
