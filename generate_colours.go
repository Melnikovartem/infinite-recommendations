package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
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

// docker run -d --rm -v $(pwd):/app -w /app golang:1.22-alpine go run generate_colours.go

func main() {
	numberOfColors := 1000
	fmt.Printf("Start sending %d colors", numberOfColors)

	// Send items to Gorse
	for i := 0; i < numberOfColors; i++ {
		color := generateColor(i, numberOfColors)
		item := Item{
			ItemId:     color,
			IsHidden:   false,
			Categories: []string{"color"},
			Timestamp:  time.Now(),
			Labels:     getColorLabels(color),
			Comment:    "Generated color",
		}
		err := sendItemToGorse(item)
		if err != nil {
			fmt.Printf("Error sending item %s to Gorse: %v\n", item.ItemId, err)
		} else if i%(numberOfColors/1000) == 0 || i == (numberOfColors-1) {
			fmt.Printf("%d: Successfully sent item %s to Gorse\n", i, item.ItemId)
		}
	}
}

func generateColor(index, total int) string {
	// Start from white (255, 255, 255) and go to black (0, 0, 0)
	value := int(float64(0xFFFFFF) * (float64(index) / float64(total-1)))
	return fmt.Sprintf("#%06x", value)
}

func sendItemToGorse(item Item) error {
	url := "http://0.0.0.0:8087/api/item" // Replace with your Gorse API endpoint

	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "pIuN9WrgAz25xya0RAzUGnqMwfzY5Fb4")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send item to Gorse, status code: %d", resp.StatusCode)
	}

	return nil
}

// hexToRGB converts a hex color to its RGB components
func hexToRGB(color string) (int, int, int) {
	colorHex := strings.TrimPrefix(color, "#")
	r, _ := strconv.ParseInt(colorHex[0:2], 16, 32)
	g, _ := strconv.ParseInt(colorHex[2:4], 16, 32)
	b, _ := strconv.ParseInt(colorHex[4:6], 16, 32)
	return int(r), int(g), int(b)
}

// getDetailedHueLabel classifies the color into more specific hue categories
func getDetailedHueLabel(r, g, b int) string {
	h, _, _ := rgbToHSV(r, g, b)
	switch {
	case h < 15 || h >= 345:
		return "red"
	case h >= 15 && h < 45:
		return "orange"
	case h >= 45 && h < 75:
		return "yellow"
	case h >= 75 && h < 165:
		return "green"
	case h >= 165 && h < 195:
		return "cyan"
	case h >= 195 && h < 255:
		return "blue"
	case h >= 255 && h < 285:
		return "purple"
	case h >= 285 && h < 345:
		return "magenta"
	}
	return "neutral"
}

// getDetailedBrightnessLabel classifies the color based on more granular brightness levels
func getDetailedBrightnessLabel(r, g, b int) string {
	brightness := (r*299 + g*587 + b*114) / 1000
	switch {
	case brightness < 64:
		return "very_dark"
	case brightness >= 64 && brightness < 128:
		return "dark"
	case brightness >= 128 && brightness < 192:
		return "medium"
	case brightness >= 192 && brightness < 240:
		return "light"
	default:
		return "very_light"
	}
}

// getDetailedSaturationLabel classifies the color based on more granular saturation levels
func getDetailedSaturationLabel(r, g, b int) string {
	_, s, _ := rgbToHSV(r, g, b)
	switch {
	case s < 0.1:
		return "grayscale"
	case s >= 0.1 && s < 0.3:
		return "muted"
	case s >= 0.3 && s < 0.6:
		return "moderate"
	case s >= 0.6 && s < 0.8:
		return "vibrant"
	default:
		return "vivid"
	}
}

// getTemperatureLabel classifies the color based on its temperature
func getTemperatureLabel(r, g, b int) string {
	temperature := (r - b) / 2
	switch {
	case temperature < -30:
		return "very_cool"
	case temperature >= -30 && temperature < -10:
		return "cool"
	case temperature >= -10 && temperature < 10:
		return "neutral"
	case temperature >= 10 && temperature < 30:
		return "warm"
	default:
		return "very_warm"
	}
}

// getPastelLabel determines if a color is pastel
func getPastelLabel(r, g, b int) string {
	_, s, v := rgbToHSV(r, g, b)
	if s < 0.5 && v > 0.7 {
		return "pastel"
	}
	return "not_pastel"
}

// getShadeLabel classifies the color based on its shade
func getShadeLabel(r, g, b int) string {
	_, _, v := rgbToHSV(r, g, b)
	switch {
	case v < 0.2:
		return "shadow"
	case v >= 0.2 && v < 0.4:
		return "deep"
	case v >= 0.4 && v < 0.6:
		return "mid_tone"
	case v >= 0.6 && v < 0.8:
		return "light"
	default:
		return "pale"
	}
}

// getToneLabel classifies the color based on its tone
func getToneLabel(r, g, b int) string {
	_, s, v := rgbToHSV(r, g, b)
	if s < 0.1 {
		if v < 0.5 {
			return "charcoal"
		} else {
			return "silver"
		}
	}
	return "chromatic"
}

// getIntensityLabel classifies the color based on its intensity
func getIntensityLabel(r, g, b int) string {
	intensity := (r + g + b) / 3
	switch {
	case intensity < 64:
		return "subdued"
	case intensity >= 64 && intensity < 128:
		return "moderate"
	case intensity >= 128 && intensity < 192:
		return "bright"
	default:
		return "intense"
	}
}

// getUndertoneLabel classifies the color based on its undertone
func getUndertoneLabel(r, g, b int) string {
	if r > b {
		return "warm_undertone"
	} else if b > r {
		return "cool_undertone"
	}
	return "neutral_undertone"
}

// getComplexityLabel classifies the color based on its complexity
func getComplexityLabel(r, g, b int) string {
	diff := math.Abs(float64(r-g)) + math.Abs(float64(g-b)) + math.Abs(float64(b-r))
	if diff < 30 {
		return "pure"
	} else if diff < 90 {
		return "simple"
	} else {
		return "complex"
	}
}

// getMoodLabel assigns a mood to the color
func getMoodLabel(r, g, b int) string {
	h, s, v := rgbToHSV(r, g, b)
	if v < 0.3 {
		return "somber"
	} else if v > 0.7 && s < 0.3 {
		return "airy"
	} else if s > 0.7 && v > 0.7 {
		return "energetic"
	} else if h >= 0 && h < 60 {
		return "warm"
	} else if h >= 180 && h < 300 {
		return "cool"
	}
	return "balanced"
}

// getColorLabels generates multiple detailed labels for a color based on its properties
func getColorLabels(color string) []string {
	r, g, b := hexToRGB(color)
	return []string{
		getDetailedHueLabel(r, g, b),
		getDetailedBrightnessLabel(r, g, b),
		getDetailedSaturationLabel(r, g, b),
		getTemperatureLabel(r, g, b),
		getPastelLabel(r, g, b),
		getShadeLabel(r, g, b),
		getToneLabel(r, g, b),
		getIntensityLabel(r, g, b),
		getUndertoneLabel(r, g, b),
		getComplexityLabel(r, g, b),
		getMoodLabel(r, g, b),
	}
}

// Helper function to convert RGB to HSV
func rgbToHSV(r, g, b int) (float64, float64, float64) {
	rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	var h, s, v float64
	v = max

	if delta == 0 {
		h, s = 0, 0
	} else {
		s = delta / max
		switch max {
		case rf:
			h = (gf - bf) / delta
			if gf < bf {
				h += 6
			}
		case gf:
			h = ((bf - rf) / delta) + 2
		case bf:
			h = ((rf - gf) / delta) + 4
		}
		h *= 60
	}

	return h, s, v
}
