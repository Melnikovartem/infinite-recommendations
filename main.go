package main

import (
	"context"
	"encoding/json"
	"fmt"
	cors "github.com/rs/cors"
	client "github.com/zhenghaoz/gorse/client"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

type Recommendation struct {
	HTMLColor string `json:"html_color"`
}

type HTTPHandler struct {
	mu    *http.ServeMux
	gorse *client.GorseClient
}

func main() {
	gorseMasterHost := os.Getenv("GORSE_SERVER_HOST")
	gorseMasterPort := os.Getenv("GORSE_SERVER_PORT")
	gorseApiKey := os.Getenv("GORSE_API_KEY")
	gorseUrl := fmt.Sprintf("http://%s:%s", gorseMasterHost, gorseMasterPort)
	fmt.Println("Proxy for gorse server at ", gorseUrl)

	handler := HTTPHandler{
		mu:    http.NewServeMux(),
		gorse: client.NewGorseClient(gorseUrl, gorseApiKey),
	}

	handler.mu.Handle("/", http.FileServer(http.Dir("./public")))

	handler.mu.HandleFunc("/recommend", handler.recommendUser)

	handler.mu.HandleFunc("/like", handler.likeColour)

	err := http.ListenAndServe(":8080", cors.Default().Handler(handler.mu))
	if err != nil {
		fmt.Println("error server not running at http://localhost:8080 ", err.Error())
		return
	}
	fmt.Println("server running at http://localhost:8080")
}

func getRandomIndex() string {
	total := 100000
	index := rand.IntN(total)
	// Start from white (255, 255, 255) and go to black (0, 0, 0)
	value := int(float64(0xFFFFFF) * (float64(index) / float64(total-1)))
	return fmt.Sprintf("#%06x", value)
}

func (handler *HTTPHandler) recommendUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	number := 10 // Default number of recommendations
	if n := r.URL.Query().Get("n"); n != "" {
		number, _ = strconv.Atoi(n)
	}

	recommendations, err := handler.gorse.GetRecommend(context.Background(), userId, "", number*4/5)
	if err != nil {
		fmt.Printf("error getting recommendations %s\n", err.Error())
		recommendations = []string{}
	}

	for range number - len(recommendations) {
		recommendations = append(recommendations, getRandomIndex())
	}

	rand.Shuffle(len(recommendations), func(i, j int) {
		recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
	})

	var colorRecommendations []Recommendation
	for _, rec := range recommendations {
		colorRecommendations = append(colorRecommendations, Recommendation{
			HTMLColor: rec,
		})
	}
	go handler.informRead(userId, recommendations)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(colorRecommendations)
}

func (handler *HTTPHandler) informRead(userId string, itemIds []string) {
	var feedback []client.Feedback
	for _, item := range itemIds {
		feedback = append(feedback, client.Feedback{
			FeedbackType: "read",
			UserId:       userId,
			ItemId:       item,
		})
	}
	handler.gorse.InsertFeedback(context.Background(), feedback)
}

func (handler *HTTPHandler) likeColour(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := data["userId"]
	htmlColor := data["html_color"]

	_, err = handler.gorse.InsertFeedback(context.Background(), []client.Feedback{
		{
			FeedbackType: "like",
			UserId:       userId,
			ItemId:       htmlColor,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User %s liked color %s\n", userId, htmlColor)
	w.WriteHeader(http.StatusOK)
}
