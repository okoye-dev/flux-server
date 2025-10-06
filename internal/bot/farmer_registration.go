package bot

import (
	"fmt"
	"log"
	"strings"

	chatbot "github.com/green-api/whatsapp-chatbot-golang"
)

// FarmerRegistrationScene handles farmer registration flow
type FarmerRegistrationScene struct {
	aiService *AIService
}

// NewFarmerRegistrationScene creates a new farmer registration scene
func NewFarmerRegistrationScene(aiService *AIService) *FarmerRegistrationScene {
	return &FarmerRegistrationScene{
		aiService: aiService,
	}
}

// Start begins the farmer registration scene
func (s FarmerRegistrationScene) Start(bot *chatbot.Bot) {
	bot.IncomingMessageHandler(func(notification *chatbot.Notification) {
		// Get the message text
		text, err := notification.Text()
		if err != nil {
			log.Printf("Error getting message text: %v", err)
			return
		}

		// Check if this is from a group chat and ignore it
		if s.isGroupChat(notification) {
			return
		}

		// Convert to lowercase for case-insensitive matching
		lowerText := strings.ToLower(strings.TrimSpace(text))

		// Handle registration command
		if strings.Contains(lowerText, "register") {
			s.startRegistration(notification)
			return
		}

		// Handle ongoing registration based on state
		stateData := notification.GetStateData()
		currentState := stateData["registration_state"]

		switch currentState {
		case STATE_REGISTER_NAME:
			s.handleName(notification, text)
		case STATE_REGISTER_CROP:
			s.handleCrop(notification, text)
		case STATE_REGISTER_LOCATION:
			s.handleLocation(notification, text)
		case STATE_REGISTER_LANGUAGE:
			s.handleLanguage(notification, text)
		}
	})
}

// startRegistration initiates the registration process
func (s *FarmerRegistrationScene) startRegistration(notification *chatbot.Notification) {
	notification.AnswerWithText("🌱 Great! Let's register you as a farmer.\n\nWhat's your full name?")
	notification.UpdateStateData(map[string]interface{}{
		"registration_state": STATE_REGISTER_NAME,
	})
}

// handleName processes the name input
func (s *FarmerRegistrationScene) handleName(notification *chatbot.Notification, name string) {
	if strings.TrimSpace(name) == "" {
		notification.AnswerWithText("Please enter your full name.")
		return
	}
	
	notification.UpdateStateData(map[string]interface{}{
		"name": name,
		"registration_state": STATE_REGISTER_CROP,
	})
	notification.AnswerWithText(fmt.Sprintf("Nice to meet you, %s! 👋\n\nWhat type of crop do you grow? (e.g., maize, rice, wheat, vegetables)", name))
}

// handleCrop processes the crop input
func (s *FarmerRegistrationScene) handleCrop(notification *chatbot.Notification, crop string) {
	if strings.TrimSpace(crop) == "" {
		notification.AnswerWithText("Please tell me what crop you grow.")
		return
	}
	
	notification.UpdateStateData(map[string]interface{}{
		"crop": crop,
		"registration_state": STATE_REGISTER_LOCATION,
	})
	notification.AnswerWithText(fmt.Sprintf("Got it! You grow %s. 🌾\n\nWhere is your farm located? (e.g., city, region, state)", crop))
}

// handleLocation processes the location input
func (s *FarmerRegistrationScene) handleLocation(notification *chatbot.Notification, location string) {
	if strings.TrimSpace(location) == "" {
		notification.AnswerWithText("Please tell me your farm location.")
		return
	}
	
	notification.UpdateStateData(map[string]interface{}{
		"location": location,
		"registration_state": STATE_REGISTER_LANGUAGE,
	})
	notification.AnswerWithText(fmt.Sprintf("Perfect! Your farm is in %s. 📍\n\nWhat language do you prefer for advice? (e.g., English, Swahili, French)", location))
}

// handleLanguage processes the language input and completes registration
func (s *FarmerRegistrationScene) handleLanguage(notification *chatbot.Notification, language string) {
	if strings.TrimSpace(language) == "" {
		notification.AnswerWithText("Please tell me your preferred language.")
		return
	}
	
	// Get all registration data
	stateData := notification.GetStateData()
	name := stateData["name"].(string)
	crop := stateData["crop"].(string)
	location := stateData["location"].(string)
	
	// Save farmer profile
	farmerProfile := map[string]interface{}{
		"name": name,
		"crop": crop,
		"location": location,
		"language": language,
	}
	
	// Update state with complete profile
	notification.UpdateStateData(map[string]interface{}{
		"farmer_profile": farmerProfile,
		"registration_state": STATE_NONE,
		"registered": true,
	})
	
	// Send completion message
	completionMessage := fmt.Sprintf(`✅ Registration Complete!

👤 Name: %s
🌾 Crop: %s
📍 Location: %s
🗣️ Language: %s

You're all set! Now you can:
• Get farming advice with "advice"
• Send feedback with "feedback"
• Check your profile with "status"

Welcome to Farm Assistant! 🌱`, name, crop, location, language)
	
	notification.AnswerWithText(completionMessage)
}

// isGroupChat checks if the message is from a group chat
func (s *FarmerRegistrationScene) isGroupChat(notification *chatbot.Notification) bool {
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