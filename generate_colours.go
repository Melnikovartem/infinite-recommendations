package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
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
func hexToRGB(color string) (int, int, int) {
	colorHex := strings.TrimPrefix(color, "#")
	r, _ := strconv.ParseInt(colorHex[0:2], 16, 32)
	g, _ := strconv.ParseInt(colorHex[2:4], 16, 32)
	b, _ := strconv.ParseInt(colorHex[4:6], 16, 32)
	return int(r), int(g), int(b)
}

// getHueLabel classifies the color into basic hue categories
func getHueLabel(r, g, b int) string {
	if r > g && r > b {
		return "red"
	} else if g > r && g > b {
		return "green"
	} else if b > r && b > g {
		return "blue"
	}
	return "neutral"
}

// getBrightnessLabel classifies the color based on brightness
func getBrightnessLabel(r, g, b int) string {
	brightness := (r*299 + g*587 + b*114) / 1000
	if brightness > 200 {
		return "bright"
	}
	return "dark"
}

// getSaturationLabel classifies the color based on saturation
func getSaturationLabel(r, g, b int) string {
	max := max(r, g, b)
	min := min(r, g, b)
	if max == min {
		return "unsaturated"
	}
	saturation := float64(max-min) / float64(max)
	if saturation > 0.5 {
		return "vibrant"
	}
	return "dull"
}

// getColorLabels generates multiple labels for a color based on its properties
func getColorLabels(color string) []string {
	r, g, b := hexToRGB(color)
	return []string{
		getHueLabel(r, g, b),
		getBrightnessLabel(r, g, b),
		getSaturationLabel(r, g, b),
	}
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

func sendItemToGorse(item Item) error {
	url := "http://localhost:8088/api/item" // Replace with your Gorse API endpoint

	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send item to Gorse, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Generate 100 color items
	items := generateItems(1000)

	// Send items to Gorse
	for _, item := range items {
		err := sendItemToGorse(item)
		if err != nil {
			fmt.Printf("Error sending item %s to Gorse: %v\n", item.ItemId, err)
		} else {
			fmt.Printf("Successfully sent item %s to Gorse\n", item.ItemId)
		}
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
