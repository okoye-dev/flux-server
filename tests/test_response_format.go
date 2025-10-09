package main

import (
	"fmt"
	"strings"
)

// AIAdviceResponse represents the structured AI advice response
type AIAdviceResponse struct {
	PlantingAdvice   string `json:"planting_advice"`
	IrrigationAdvice string `json:"irrigation_advice"`
	HarvestAdvice    string `json:"harvest_advice"`
	MarketAdvice     string `json:"market_advice"`
	GeneralAdvice    string `json:"general_advice"`
	Confidence       int    `json:"confidence"`
	GeneratedAt      string `json:"generated_at"`
}

// parseAIResponse parses the AI response into structured advice
func parseAIResponse(aiResponse, cropsList string) *AIAdviceResponse {
	// Clean up the response
	aiResponse = strings.TrimSpace(aiResponse)
	
	// If the response is well-formatted, use it directly
	if strings.Contains(aiResponse, "Planting advice:") || strings.Contains(aiResponse, "1.") {
		return &AIAdviceResponse{
			PlantingAdvice:   extractAdviceSection(aiResponse, "planting", "1."),
			IrrigationAdvice: extractAdviceSection(aiResponse, "irrigation", "2."),
			HarvestAdvice:    extractAdviceSection(aiResponse, "harvest", "3."),
			MarketAdvice:     extractAdviceSection(aiResponse, "market", "4."),
			GeneralAdvice:    extractAdviceSection(aiResponse, "general", "5."),
			Confidence:       90,
			GeneratedAt:      "2025-01-27T10:30:00Z",
		}
	}
	
	// Fallback: return the full response as general advice
	return &AIAdviceResponse{
		PlantingAdvice:   "Based on current conditions, follow optimal planting schedules for your crops.",
		IrrigationAdvice: "Monitor soil moisture and adjust irrigation based on weather conditions.",
		HarvestAdvice:    fmt.Sprintf("Your %s crops should be ready for harvest based on growth conditions.", cropsList),
		MarketAdvice:     "Monitor market trends and prices for optimal selling timing.",
		GeneralAdvice:    aiResponse,
		Confidence:       85,
		GeneratedAt:      "2025-01-27T10:30:00Z",
	}
}

// extractAdviceSection extracts a specific section from the AI response
func extractAdviceSection(response, keyword, number string) string {
	lines := strings.Split(response, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lowerLine := strings.ToLower(line)
		
		// Look for lines that contain both the number and the keyword
		if strings.Contains(lowerLine, number) && strings.Contains(lowerLine, keyword) {
			// Extract the advice part after the colon
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				advice := strings.TrimSpace(parts[1])
				// Remove any markdown formatting
				advice = strings.TrimPrefix(advice, "**")
				advice = strings.TrimSuffix(advice, "**")
				advice = strings.TrimSpace(advice)
				return advice
			}
		}
	}
	
	// Fallback
	return fmt.Sprintf("Follow best practices for %s based on current conditions.", keyword)
}

func main() {
	// Test with the actual AI response we got
	aiResponse := `Here is your farming advice, John Doe:

1.  **Planting advice:** Ensure proper seed depth and spacing for both maize and beans to maximize light exposure and nutrient uptake in sunny conditions. Consider planting drought-resistant varieties if rainfall is unreliable.
2.  **Irrigation advice:** With sunny conditions and warm temperatures, monitor soil moisture daily and irrigate adequately, especially for maize during its tasseling and silking stages, and beans during pod formation. Avoid waterlogging.
3.  **Harvest advice:** Harvest maize when kernels are firm and dented, and beans when pods are dry and rattling, to ensure optimal quality and storage life. Timely harvest reduces pest damage and yield loss.
4.  **Market advice:** Given the stable price, focus on delivering high-quality produce to maintain demand and consider aggregating your yield with other farmers to improve bargaining power.
5.  **General advice:** Regularly scout for pests and diseases, and practice good field hygiene to prevent outbreaks. Crop rotation with other suitable crops in subsequent seasons will help improve soil fertility.`

	fmt.Println("🧪 Testing AI Response Parsing...")
	fmt.Println("==================================================")
	
	// Parse the response
	parsed := parseAIResponse(aiResponse, "maize, beans")
	
	// Display the structured response
	fmt.Printf("🌱 Planting Advice: %s\n\n", parsed.PlantingAdvice)
	fmt.Printf("💧 Irrigation Advice: %s\n\n", parsed.IrrigationAdvice)
	fmt.Printf("🌾 Harvest Advice: %s\n\n", parsed.HarvestAdvice)
	fmt.Printf("💰 Market Advice: %s\n\n", parsed.MarketAdvice)
	fmt.Printf("📋 General Advice: %s\n\n", parsed.GeneralAdvice)
	fmt.Printf("🎯 Confidence: %d%%\n", parsed.Confidence)
	
	fmt.Println("==================================================")
	fmt.Println("✅ Response parsing test complete!")
}
