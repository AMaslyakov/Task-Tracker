package API

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "session_id"
	sessionMaxAge     = 1800
)

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/login [post]
func Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login payload"})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user, passwordHash, err := DBLoginByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		log.Printf("failed to fetch login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		log.Printf("failed to generate session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	if err := DBCreateSession(ctx, sessionID, user.ID); err != nil {
		log.Printf("failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}
	if err := DBMarkLoginSuccess(ctx, user.ID); err != nil {
		log.Printf("failed to update login metadata: %v", err)
	}

	c.SetCookie(sessionCookieName, sessionID, sessionMaxAge, "/", "", false, true)
	c.JSON(http.StatusOK, LoginResponse{User: user})
}

// Me godoc
// @Summary Get current user
// @Tags auth
// @Produce json
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/me [get]
func Me(c *gin.Context) {
	ctx := c.Request.Context()

	sessionID, err := c.Cookie(sessionCookieName)
	if err != nil || sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	user, err := DBUserBySession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		log.Printf("failed to fetch current user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch current user"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{User: user})
}

// Logout godoc
// @Summary Logout user
// @Tags auth
// @Success 204
// @Router /api/logout [post]
func Logout(c *gin.Context) {
	sessionID, err := c.Cookie(sessionCookieName)
	if err == nil && sessionID != "" {
		if err := DBDeleteSession(c.Request.Context(), sessionID); err != nil {
			log.Printf("failed to delete session: %v", err)
		}
	}

	c.SetCookie(sessionCookieName, "", -1, "/", "", false, true)
	c.Status(http.StatusNoContent)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(sessionCookieName)
		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		user, err := DBUserBySession(c.Request.Context(), sessionID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
				return
			}
			log.Printf("failed to authenticate request: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate request"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
