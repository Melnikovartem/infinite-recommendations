package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/zhenghaoz/gorse/base"
)

type Item struct {
	ItemId     string
	IsHidden   bool
	Categories []string
	Timestamp  time.Time
	Labels     []string
	Comment    string
}

func generateColor() string {
	return fmt.Sprintf("#%06x", rand.Intn(0xFFFFFF))
}

// getColorLabel categorizes the color based on its hexadecimal value
func getColorLabel(color string) string {
	// Remove the '#' from the color code
	colorHex := strings.TrimPrefix(color, "#")
	// Convert the hex color to RGB values
	r, _ := strconv.ParseInt(colorHex[0:2], 16, 32)
	g, _ := strconv.ParseInt(colorHex[2:4], 16, 32)
	b, _ := strconv.ParseInt(colorHex[4:6], 16, 32)

	// Simple categorization based on RGB values
	if (r > g && r > b) || (r > 150 && g < 100 && b < 100) {
		return "warm"
	} else if (b > r && b > g) || (b > 150 && r < 100 && g < 100) {
		return "cool"
	} else if (g > r && g > b) || (g > 150 && r < 100 && b < 100) {
		return "green"
	}
	return "neutral"
}

func generateItems(n int) []Item {
	rand.Seed(time.Now().UnixNano())
	items := make([]Item, n)

	for i := 0; i < n; i++ {
		color := generateColor()
		items[i] = Item{
			ItemId:     color,
			IsHidden:   false,
			Categories: []string{"color"},
			Timestamp:  time.Now(),
			Labels:     []string{getColorLabel(color)},
			Comment:    "Generated color",
		}
	}

	return items
}

func main() {
	// Generate 100 color items
	items := generateItems(100)

	// Print items
	for _, item := range items {
		fmt.Println(item)
	}
}
