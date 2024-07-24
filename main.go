package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	client "github.com/zhenghaoz/gorse/client"
)

type Recommendation struct {
	HTMLColor string `json:"html_color"`
}

func main() {
	gorseMasterHost := os.Getenv("GORSE_SERVER_HOST")
	gorseMasterPort := os.Getenv("GORSE_SERVER_PORT")
	gorseApiKey := os.Getenv("GORSE_API_KEY")
	gorseUrl := fmt.Sprintf("http://%s:%s", gorseMasterHost, gorseMasterPort)
	fmt.Println("Proxy for gorse server at ", gorseUrl)
	gorse := client.NewGorseClient(gorseUrl, gorseApiKey)

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/recommend", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		number := 10 // Default number of recommendations
		if n := r.URL.Query().Get("Number"); n != "" {
			number, _ = strconv.Atoi(n)
		}

		recommendations, err := gorse.GetRecommend(context.Background(), userId, "", number)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var colorRecommendations []Recommendation
		for _, rec := range recommendations {
			colorRecommendations = append(colorRecommendations, Recommendation{
				HTMLColor: rec,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(colorRecommendations)
	})

	http.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userId := data["userId"]
		htmlColor := data["html_color"]

		_, err = gorse.InsertFeedback(context.Background(), []client.Feedback{
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
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
