package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	chatbot "github.com/green-api/whatsapp-chatbot-golang"
)

// AdviceDeliveryScene handles AI advice delivery flow
type AdviceDeliveryScene struct {
	aiService *AIService
}

// NewAdviceDeliveryScene creates a new advice delivery scene
func NewAdviceDeliveryScene(aiService *AIService) *AdviceDeliveryScene {
	return &AdviceDeliveryScene{
		aiService: aiService,
	}
}

// Start begins the advice delivery scene
func (s AdviceDeliveryScene) Start(bot *chatbot.Bot) {
	bot.IncomingMessageHandler(func(notification *chatbot.Notification) {
		// Get the message text
		text, err := notification.Text()
		if err != nil {
			log.Printf("Error getting message text: %v", err)
			return
		}

		// Check if this is from a group chat and ignore it
		if s.isGroupChat(notification) {
			log.Printf("Ignoring group chat message: %s", text)
			return
		}

		// Convert to lowercase for case-insensitive matching
		actualMessage := strings.ToLower(strings.TrimSpace(text))

		// Handle advice command
		if strings.Contains(actualMessage, CMD_ADVICE) {
			s.handleAdviceRequest(notification)
			return
		}
	})
}

// handleAdviceRequest processes advice requests
func (s *AdviceDeliveryScene) handleAdviceRequest(notification *chatbot.Notification) {
	log.Printf("Processing advice request")
	
	// Get farmer profile from state
	stateData := notification.GetStateData()
	farmerProfileData, ok := stateData["farmer_profile"].(map[string]interface{})
	if !ok {
		// If no profile found, ask to register first
		notification.AnswerWithText("❌ Please register first using 'register' to get personalized advice.")
		return
	}
	
	// Convert to FarmerProfile struct
	farmerProfile := FarmerProfile{
		Name:     getStringFromMap(farmerProfileData, "name"),
		Location: getStringFromMap(farmerProfileData, "location"),
		Language: getStringFromMap(farmerProfileData, "language"),
		Phone:    getStringFromMap(farmerProfileData, "phone"),
	}
	
	// Handle both single crop and multiple crops
	if crops, ok := farmerProfileData["crops"].([]string); ok {
		farmerProfile.Crops = crops
	} else if crop, ok := farmerProfileData["crop"].(string); ok {
		farmerProfile.Crops = []string{crop}
	} else {
		farmerProfile.Crops = []string{"Unknown"}
	}
	
	// Set state to waiting for advice
	notification.SetStateData(map[string]interface{}{
		"state": STATE_WAITING_ADVICE,
		"farmer_profile": farmerProfile,
	})
	
	// Send initial processing message
	notification.AnswerWithText(MSG_ADVICE_REQUEST)
	
	// Generate AI advice with loading messages
	s.generateAndSendAdviceWithLoading(notification, farmerProfile)
}

// generateAndSendAdviceWithLoading generates AI advice with loading messages
func (s *AdviceDeliveryScene) generateAndSendAdviceWithLoading(notification *chatbot.Notification, profile FarmerProfile) {
	log.Printf("🤖 Generating AI advice for farmer: %s", profile.Name)
	
	// Send loading message 1
	notification.AnswerWithText(MSG_AI_LOADING_1)
	time.Sleep(3 * time.Second)
	
	// Send loading message 2
	notification.AnswerWithText(MSG_AI_LOADING_2)
	time.Sleep(3 * time.Second)
	
	// Send loading message 3
	notification.AnswerWithText(MSG_AI_LOADING_3)
	time.Sleep(2 * time.Second)
	
	// Now generate the actual advice
	s.generateAndSendAdvice(notification, profile)
}

// generateAndSendAdvice generates AI advice and sends it to the farmer
func (s *AdviceDeliveryScene) generateAndSendAdvice(notification *chatbot.Notification, profile FarmerProfile) {
	log.Printf("🤖 Generating AI advice for farmer: %s", profile.Name)
	
	// Fetch weather data
	weatherData, err := s.aiService.GetWeatherData(profile.Location)
	if err != nil {
		log.Printf("Error fetching weather data: %v", err)
		weatherData = &WeatherData{
			Temperature: 25.0,
			Humidity:    60.0,
			Rainfall:    10.0,
			Condition:   "Sunny",
			Date:        "2025-10-06",
		}
	}
	
	// Fetch market data (use first crop for market data)
	var primaryCrop string
	if len(profile.Crops) > 0 {
		primaryCrop = profile.Crops[0]
	} else {
		primaryCrop = "maize" // fallback
	}
	
	marketData, err := s.aiService.GetMarketData(primaryCrop, profile.Location)
	if err != nil {
		log.Printf("Error fetching market data: %v", err)
		marketData = &MarketData{
			CropType: primaryCrop,
			Price:    2.50,
			Currency: "$",
			Unit:     "kg",
			Location: profile.Location,
			Date:     "2025-10-06",
			Trend:    "stable",
		}
	}
	
	// Prepare AI request
	aiRequest := AIAdviceRequest{
		FarmerProfile: profile,
		WeatherData:   *weatherData,
		MarketData:    *marketData,
		Season:        s.getCurrentSeason(),
	}
	
	// Call Gemini AI
	aiResponse, err := s.aiService.CallGeminiAI(aiRequest)
	if err != nil {
		log.Printf("Error calling Gemini AI: %v", err)
		notification.AnswerWithText("❌ Sorry, I couldn't generate advice right now. Please try again later.")
		return
	}
	
	// Format and send the advice
	adviceMessage := s.formatAdviceMessage(aiResponse, weatherData, marketData)
	notification.AnswerWithText(adviceMessage)
	
	// Update state back to idle
	notification.SetStateData(map[string]interface{}{
		"state": STATE_IDLE,
		"farmer_profile": profile,
	})
	
	log.Printf("✅ Advice delivered successfully to farmer: %s", profile.Name)
}

// formatAdviceMessage formats the AI response into a readable message
func (s *AdviceDeliveryScene) formatAdviceMessage(aiResponse *AIAdviceResponse, weather *WeatherData, market *MarketData) string {
	message := fmt.Sprintf(`🌱 *Your Personalized Farming Advice*

🌤️ *Weather Conditions:*
• Temperature: %.1f°C
• Humidity: %.1f%%
• Condition: %s
• Rainfall: %.1fmm

💰 *Market Information:*
• %s Price: %s%.2f per %s
• Trend: %s

🤖 *AI Recommendations:*

🌱 *Planting:* %s

💧 *Irrigation:* %s

🌾 *Harvest:* %s

📈 *Market Strategy:* %s

💡 *General Advice:* %s

*Confidence: %d%% | Generated: %s*`,
		weather.Temperature,
		weather.Humidity,
		weather.Condition,
		weather.Rainfall,
		market.CropType,
		market.Currency,
		market.Price,
		market.Unit,
		strings.Title(market.Trend),
		aiResponse.PlantingAdvice,
		aiResponse.IrrigationAdvice,
		aiResponse.HarvestAdvice,
		aiResponse.MarketAdvice,
		aiResponse.GeneralAdvice,
		aiResponse.Confidence,
		aiResponse.GeneratedAt[:10], // Just the date part
	)
	
	return message
}

// getCurrentSeason returns the current season (dummy implementation)
func (s *AdviceDeliveryScene) getCurrentSeason() string {
	// In a real implementation, this would calculate based on location and date
	return "Rainy Season"
}

// isGroupChat checks if the message is from a group chat
func (s *AdviceDeliveryScene) isGroupChat(notification *chatbot.Notification) bool {
	body := notification.Body
	
	if senderData, ok := body["senderData"].(map[string]interface{}); ok {
		if chatId, exists := senderData["chatId"]; exists {
			if chatIdStr, ok := chatId.(string); ok {
				return strings.HasSuffix(chatIdStr, "@g.us")
			}
		}
	}
	
	return false
}

// getStringFromMap safely extracts a string value from a map
func getStringFromMap(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}
