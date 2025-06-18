package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type RegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type MoodPayload struct {
	Emoji   string  `json:"emoji"`
	Message string  `json:"message"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

var emojis = []string{"😊", "😢", "😠", "😂", "😴", "😎", "😱", "😍"}

// India bounding box
func randomInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	client := &http.Client{}

	for i := 1; i <= 100; i++ {
		username := fmt.Sprintf("user%d", i)
		password := "password123"

		// Register
		regBody, _ := json.Marshal(RegisterPayload{Username: username, Password: password})
		resp, err := client.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(regBody))
		if err != nil {
			fmt.Printf("❌ Register failed for %s: %v\n", username, err)
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// Login
		loginBody, _ := json.Marshal(LoginPayload{Username: username, Password: password})
		resp, err = client.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			fmt.Printf("❌ Login failed for %s: %v\n", username, err)
			continue
		}
		var loginResp LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		resp.Body.Close()
		if err != nil || loginResp.Token == "" {
			fmt.Printf("❌ Failed to get token for %s\n", username)
			continue
		}

		// Post Mood
		payload := MoodPayload{
			Emoji:   emojis[rand.Intn(len(emojis))],
			Message: fmt.Sprintf("Feeling mood #%d", i),
			Lat:     randomInRange(8.0, 37.0),
			Lng:     randomInRange(68.0, 97.0),
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "http://localhost:8080/mood", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("❌ Request creation failed for %s: %v\n", username, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Mood post failed for %s: %v\n", username, err)
			continue
		}
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
		fmt.Printf("✅ Mood seeded for %s (%s)\n", username, res.Status)
	}
}
