package bot

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// AIService handles all AI-related operations
type AIService struct{}

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
	
	// Simulate API call delay
	time.Sleep(2 * time.Second)
	
	// Get crops list for advice
	cropsList := strings.Join(request.FarmerProfile.Crops, ", ")
	
	// Dummy response - replace with actual Gemini API call
	response := &AIAdviceResponse{
		PlantingAdvice:   fmt.Sprintf("Based on current weather conditions, it's optimal to plant %s in the next 3-5 days. Soil moisture levels are ideal.", cropsList),
		IrrigationAdvice: "With current humidity at 65%, reduce irrigation frequency to every 3 days. Monitor soil moisture closely.",
		HarvestAdvice:    fmt.Sprintf("Your %s crops should be ready for harvest in approximately 45-60 days based on current growth conditions.", cropsList),
		MarketAdvice:     fmt.Sprintf("Current market prices for your crops (%s) are favorable. %s", cropsList, getMarketTrendAdvice(request.MarketData.Trend)),
		GeneralAdvice:    "This would be some info from our AI based on what you asked. Monitor your crops regularly and maintain proper spacing for optimal yield.",
		Confidence:       85,
		GeneratedAt:      time.Now().Format(time.RFC3339),
	}
	
	return response, nil
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
