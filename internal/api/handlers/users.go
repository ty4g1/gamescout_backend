package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ty4g1/gamescout_backend/internal/repository"
	"github.com/ty4g1/gamescout_backend/internal/utils"
)

type UserHandler struct {
	ur *repository.UserRepository
	gr *repository.GameRepository
}

func NewUserHandler(ur *repository.UserRepository, gr *repository.GameRepository) *UserHandler {
	return &UserHandler{
		ur: ur,
		gr: gr,
	}
}

func (uh *UserHandler) AddUser(c *gin.Context) {
	// Parse id
	var req struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := uh.ur.AddUser(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uh *UserHandler) UpdatePreference(c *gin.Context) {
	// Parse id
	var req struct {
		ID       string `json:"id" binding:"required"`
		Likes    []int  `json:"likes" binding:"required"`
		Dislikes []int  `json:"dislikes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uh.ur.UpdateUserSwipes(c.Request.Context(), req.ID, append(req.Likes, req.Dislikes...))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update user swipes: %v", err)})
		return
	}

	additions, err := uh.gr.GetFeatureVecByAppIDs(c.Request.Context(), req.Likes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feature vectors"})
		return
	}

	subtractions, err := uh.gr.GetFeatureVecByAppIDs(c.Request.Context(), req.Dislikes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feature vectors"})
		return
	}

	preferences, err := uh.ur.GetUserPreference(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user preferences %v", err)})
		return
	}

	netAddition, err := utils.SumRows(additions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sum likes: %v", err)})
		return
	}
	netSubtraction, err := utils.SumRows(subtractions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sum dislikes: %v", err)})
		return
	}

	preferences, err = utils.AddVectors(preferences, netAddition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add likes: %v", err)})
		return
	}

	preferences, err = utils.SubtractVectors(preferences, netSubtraction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to subtract dislikes: %v", err)})
		return
	}

	preferences = utils.NormalizeVector(preferences)

	err = uh.ur.UpdateUserPreference(c.Request.Context(), req.ID, preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update user preferences: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"preferences": preferences})
}
