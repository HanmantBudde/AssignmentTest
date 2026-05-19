package todo

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *Store
}

func NewHandler(s *Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/todos/:id", h.get)
}

func (h *Handler) get(c *gin.Context) {
	t, err := h.store.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondStoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, t)
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
