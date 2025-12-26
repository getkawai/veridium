package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kawai-network/veridium/pkg/store"
)

// AuthMiddleware validates API keys and injects user context
func AuthMiddleware(kvStore *store.KVStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Missing Authorization header",
					"type":    "invalid_request_error",
					"code":    "missing_api_key",
				},
			})
			c.Abort()
			return
		}

		// Parse "Bearer vk-..." format
		var apiKey string
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Also support direct API key without "Bearer" prefix
			apiKey = authHeader
		}

		// Validate API key and get wallet address
		ctx := context.Background()
		walletAddress, err := kvStore.ValidateAPIKey(ctx, apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Invalid API key",
					"type":    "invalid_request_error",
					"code":    "invalid_api_key",
				},
			})
			c.Abort()
			return
		}

		// Inject user context for downstream handlers
		c.Set("user_address", walletAddress)
		c.Set("api_key", apiKey)

		c.Next()
	}
}

// GetUserAddress extracts the wallet address from gin context
func GetUserAddress(c *gin.Context) (string, bool) {
	address, exists := c.Get("user_address")
	if !exists {
		return "", false
	}

	addressStr, ok := address.(string)
	return addressStr, ok
}

// GetAPIKey extracts the API key from gin context
func GetAPIKey(c *gin.Context) (string, bool) {
	apiKey, exists := c.Get("api_key")
	if !exists {
		return "", false
	}

	apiKeyStr, ok := apiKey.(string)
	return apiKeyStr, ok
}
