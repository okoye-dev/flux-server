package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AIService handles all AI-related operations
type AIService struct{}

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents a content item in Gemini request
type Content struct {
	Parts []Part `json:"parts"`
}

// Part represents a part of content
type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response structure from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

// Candidate represents a candidate response from Gemini
type Candidate struct {
	Content Content `json:"content"`
}

// NewAIService creates a new AI service instance
func NewAIService() *AIService {
	return &AIService{}
}

// FarmerProfile represents a farmer's profile data
type FarmerProfile struct {
	Name     string   `json:"name"`
	Crops    []string `json:"crops"`
	Location string   `json:"location"`
	Language string   `json:"language"`
	Phone    string   `json:"phone"`
}

// WeatherData represents weather information
type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Rainfall    float64 `json:"rainfall"`
	Condition   string  `json:"condition"`
	Date        string  `json:"date"`
}

// MarketData represents market price information
type MarketData struct {
	CropType   string  `json:"crop_type"`
	Price      float64 `json:"price"`
	Currency   string  `json:"currency"`
	Unit       string  `json:"unit"`
	Location   string  `json:"location"`
	Date       string  `json:"date"`
	Trend      string  `json:"trend"` // "up", "down", "stable"
}

// AIAdviceRequest represents the input for AI advice generation
type AIAdviceRequest struct {
	FarmerProfile FarmerProfile `json:"farmer_profile"`
	WeatherData   WeatherData   `json:"weather_data"`
	MarketData    MarketData    `json:"market_data"`
	Season        string        `json:"season"`
}

// AIAdviceResponse represents the AI-generated advice
type AIAdviceResponse struct {
	PlantingAdvice    string `json:"planting_advice"`
	IrrigationAdvice  string `json:"irrigation_advice"`
	HarvestAdvice     string `json:"harvest_advice"`
	MarketAdvice      string `json:"market_advice"`
	GeneralAdvice     string `json:"general_advice"`
	Confidence        int    `json:"confidence"` // 1-100
	GeneratedAt       string `json:"generated_at"`
}

// CallGeminiAI generates personalized farming advice using Gemini API
func (ai *AIService) CallGeminiAI(request AIAdviceRequest) (*AIAdviceResponse, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable not set")
	}

	// Get crops list for advice
	cropsList := strings.Join(request.FarmerProfile.Crops, ", ")
	
	// Create the prompt for Gemini
	prompt := fmt.Sprintf(`You are an expert agricultural advisor. Provide farming advice for a farmer with the following details:

Farmer Profile:
- Name: %s
- Crops: %s
- Location: %s
- Language: %s

Weather Data:
- Temperature: %.1f°C
- Humidity: %.1f%%
- Condition: %s

Market Data:
- Price: %.2f %s per %s
- Trend: %s

Please provide concise, actionable advice in the following format:
1. Planting advice (1-2 sentences)
2. Irrigation advice (1-2 sentences) 
3. Harvest advice (1-2 sentences)
4. Market advice (1-2 sentences)
5. General advice (1-2 sentences)

Keep responses practical and specific to the farmer's location and crops. Use simple language.`, 
		request.FarmerProfile.Name,
		cropsList,
		request.FarmerProfile.Location,
		request.FarmerProfile.Language,
		request.WeatherData.Temperature,
		request.WeatherData.Humidity,
		request.WeatherData.Condition,
		request.MarketData.Price,
		request.MarketData.Currency,
		request.MarketData.Unit,
		request.MarketData.Trend)

	// Create Gemini request
	geminiReq := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API call to Gemini
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode Gemini response: %w", err)
	}

	// Extract AI response text
	var aiResponse string
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		aiResponse = geminiResp.Candidates[0].Content.Parts[0].Text
	} else {
		aiResponse = "Unable to generate advice at this time. Please try again later."
	}

	// Parse the AI response into structured advice
	advice := ai.parseAIResponse(aiResponse, cropsList)
	
	return advice, nil
}

// parseAIResponse parses the AI response into structured advice
func (ai *AIService) parseAIResponse(aiResponse, cropsList string) *AIAdviceResponse {
	// Split response into lines and extract advice
	lines := strings.Split(aiResponse, "\n")
	
	response := &AIAdviceResponse{
		PlantingAdvice:   "Based on current conditions, follow optimal planting schedules for your crops.",
		IrrigationAdvice: "Monitor soil moisture and adjust irrigation based on weather conditions.",
		HarvestAdvice:    fmt.Sprintf("Your %s crops should be ready for harvest based on growth conditions.", cropsList),
		MarketAdvice:     "Monitor market trends and prices for optimal selling timing.",
		GeneralAdvice:    "Continue monitoring your crops regularly and maintain proper farming practices.",
		Confidence:       85,
		GeneratedAt:      time.Now().Format(time.RFC3339),
	}

	// Try to extract specific advice from AI response
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "plant") && i < len(lines)-1 {
			response.PlantingAdvice = strings.TrimSpace(lines[i+1])
		} else if strings.Contains(strings.ToLower(line), "irrigat") && i < len(lines)-1 {
			response.IrrigationAdvice = strings.TrimSpace(lines[i+1])
		} else if strings.Contains(strings.ToLower(line), "harvest") && i < len(lines)-1 {
			response.HarvestAdvice = strings.TrimSpace(lines[i+1])
		} else if strings.Contains(strings.ToLower(line), "market") && i < len(lines)-1 {
			response.MarketAdvice = strings.TrimSpace(lines[i+1])
		} else if strings.Contains(strings.ToLower(line), "general") && i < len(lines)-1 {
			response.GeneralAdvice = strings.TrimSpace(lines[i+1])
		}
	}

	// If we got a good response, use the first part as general advice
	if len(aiResponse) > 50 {
		response.GeneralAdvice = aiResponse
	}

	return response
}

// GetWeatherData fetches weather information for a location
func (ai *AIService) GetWeatherData(location string) (*WeatherData, error) {
	
	// Simulate API call delay
	time.Sleep(1 * time.Second)
	
	// Dummy weather data - replace with actual OpenWeather API call
	weatherData := &WeatherData{
		Temperature: 28.5,
		Humidity:    65.0,
		Rainfall:    12.3,
		Condition:   "Partly Cloudy",
		Date:        time.Now().Format("2006-01-02"),
	}
	
	return weatherData, nil
}

// GetMarketData fetches market price information for a crop
func (ai *AIService) GetMarketData(cropType, location string) (*MarketData, error) {
	
	// Simulate API call delay
	time.Sleep(1 * time.Second)
	
	// Dummy market data - replace with actual market API call
	marketData := &MarketData{
		CropType: cropType,
		Price:    2.45,
		Currency: "$",
		Unit:     "kg",
		Location: location,
		Date:     time.Now().Format("2006-01-02"),
		Trend:    "up",
	}
	
	return marketData, nil
}

// ProcessFeedback processes farmer feedback and returns insights
func (ai *AIService) ProcessFeedback(farmerProfile FarmerProfile, feedback string) (string, error) {
	
	// Simulate processing delay
	time.Sleep(1 * time.Second)
	
	// Dummy feedback processing - replace with actual AI analysis
	cropsList := strings.Join(farmerProfile.Crops, ", ")
	response := fmt.Sprintf("Thank you for your feedback: '%s'. This information will help improve future recommendations for your %s farming in %s.", 
		feedback, cropsList, farmerProfile.Location)
	
	log.Printf("✅ Feedback processed successfully")
	return response, nil
}

// Helper function to get market trend advice
func getMarketTrendAdvice(trend string) string {
	switch trend {
	case "up":
		return "Consider holding your crop for a few more days as prices are rising."
	case "down":
		return "Consider selling soon as prices are declining."
	case "stable":
		return "Prices are stable, you can sell when convenient."
	default:
		return "Monitor market trends closely before making selling decisions."
	}
}
