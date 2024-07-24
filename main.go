package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Recommendation struct {
	HTMLColor string `json:"html_color"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/recommend", func(w http.ResponseWriter, r *http.Request) {
		number := 10 // Default number of recommendations
		if n, ok := r.URL.Query()["Number"]; ok {
			fmt.Sscanf(n[0], "%d", &number)
		}

		var recommendations []Recommendation
		for i := 0; i < number; i++ {
			recommendations = append(recommendations, Recommendation{
				HTMLColor: fmt.Sprintf("#%06x", rand.Intn(0xFFFFFF)),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recommendations)
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

		fmt.Printf("User %s liked color %s\n", userId, htmlColor)
		w.WriteHeader(http.StatusOK)
	})

	rand.Seed(time.Now().UnixNano())
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
