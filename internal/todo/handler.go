package todo

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *Store
}

func NewHandler(s *Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/todos", h.create)
	r.GET("/todos", h.list)
	r.GET("/todos/:id", h.get)
	r.PUT("/todos/:id", h.update)
	r.DELETE("/todos/:id", h.delete)
}

type createRequest struct {
	Text    string    `json:"text" binding:"required"`
	DueDate time.Time `json:"due_date" binding:"required"`
}

func (h *Handler) create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		respondError(c, http.StatusBadRequest, "text must not be empty")
		return
	}
	t, err := h.store.Create(c.Request.Context(), req.Text, req.DueDate)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) get(c *gin.Context) {
	t, err := h.store.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondStoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) list(c *gin.Context) {
	includeCompleted := false
	if v := c.Query("include_completed"); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes":
			includeCompleted = true
		}
	}
	out, err := h.store.List(c.Request.Context(), includeCompleted)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, out)
}

type updateRequest struct {
	Text      *string    `json:"text,omitempty"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Completed *bool      `json:"completed,omitempty"`
}

func (h *Handler) update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.Text == nil && req.DueDate == nil && req.Completed == nil {
		respondError(c, http.StatusBadRequest, "no fields to update")
		return
	}
	if req.Text != nil && strings.TrimSpace(*req.Text) == "" {
		respondError(c, http.StatusBadRequest, "text must not be empty")
		return
	}
	t, err := h.store.Update(c.Request.Context(), c.Param("id"), req.Text, req.DueDate, req.Completed)
	if err != nil {
		respondStoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) delete(c *gin.Context) {
	if err := h.store.Delete(c.Request.Context(), c.Param("id")); err != nil {
		respondStoreError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func respondStoreError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	respondError(c, http.StatusInternalServerError, err.Error())
}
