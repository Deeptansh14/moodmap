package routes

import (
	"log"
	"moodmap/utils"
	"net/http"
	"strings"

	"moodmap/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func PostMood(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract username from claims
	claims := token.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Parse body (without username)
	var body struct {
		Emoji   string  `json:"emoji"`
		Message string  `json:"message"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}

	if err := c.BindJSON(&body); err != nil {
		log.Println("binding error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Delete any existing mood from this user
	_, err = utils.DB.Exec("DELETE FROM moods WHERE username = $1", username)
	if err != nil {
		log.Println("DB delete error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete previous mood"})
		return
	}

	// Insert the new mood
	_, err = utils.DB.Exec(`
		INSERT INTO moods (username, emoji, message, location)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326))
	`, username, body.Emoji, body.Message, body.Lng, body.Lat)

	if err != nil {
		log.Println("DB insert error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save mood"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mood posted!"})
}
func GetAllMoods(c *gin.Context) {
	rows, err := utils.DB.Query(`
		SELECT username, emoji, message,
		       ST_Y(location::geometry) AS lat,
		       ST_X(location::geometry) AS lng,
		       timestamp
		FROM moods
		ORDER BY timestamp DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()

	var moods []models.Mood

	for rows.Next() {
		var m models.Mood
		if err := rows.Scan(&m.Username, &m.Emoji, &m.Message, &m.Lat, &m.Lng, &m.Timestamp); err != nil {
			continue
		}
		moods = append(moods, m)
	}

	// Apply DBSCAN clustering
	clustered := utils.DBSCAN(moods)

	// Convert to response format
	var response []gin.H
	for _, mood := range clustered {
		response = append(response, gin.H{
			"username":   mood.Username,
			"emoji":      mood.Emoji,
			"message":    mood.Message,
			"lat":        mood.Lat,
			"lng":        mood.Lng,
			"timestamp":  mood.Timestamp,
			"cluster_id": mood.ClusterID,
		})
	}

	c.JSON(http.StatusOK, response)
}
