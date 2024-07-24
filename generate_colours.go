package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"strings"
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

// getColorProperties converts a hex color to its RGB components
func getColorProperties(color string) (int, int, int) {
	colorHex := strings.TrimPrefix(color, "#")
	r, _ := strconv.ParseInt(colorHex[0:2], 16, 32)
	g, _ := strconv.ParseInt(colorHex[2:4], 16, 32)
	b, _ := strconv.ParseInt(colorHex[4:6], 16, 32)
	return int(r), int(g), int(b)
}

// getColorLabels generates multiple labels for a color based on its properties
func getColorLabels(color string) []string {
	r, g, b := getColorProperties(color)
	labels := []string{}

	// Basic hue-based categories
	if r > g && r > b {
		labels = append(labels, "warm", "red")
	} else if b > r && b > g {
		labels = append(labels, "cool", "blue")
	} else if g > r && g > b {
		labels = append(labels, "green")
	} else {
		labels = append(labels, "neutral")
	}

	// Brightness-based categories
	brightness := (r + g + b) / 3
	if brightness > 200 {
		labels = append(labels, "bright")
	} else if brightness < 100 {
		labels = append(labels, "dark")
	} else {
		labels = append(labels, "medium")
	}

	// Saturation-based categories
	max := max(r, g, b)
	min := min(r, g, b)
	if max == min {
		labels = append(labels, "gray")
	} else {
		saturation := float64(max-min) / float64(max)
		if saturation > 0.5 {
			labels = append(labels, "vivid")
		} else {
			labels = append(labels, "pastel")
		}
	}

	return labels
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
			Labels:     getColorLabels(color),
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


// max returns the maximum of three integers
func max(a, b, c int) int {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}