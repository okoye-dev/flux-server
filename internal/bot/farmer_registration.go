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
			s.HandleName(notification, text)
		case STATE_REGISTER_CROP:
			s.HandleCrop(notification, text)
		case STATE_REGISTER_MORE_CROPS:
			s.HandleMoreCrops(notification, text)
		case STATE_REGISTER_LOCATION:
			s.HandleLocation(notification, text)
		case STATE_REGISTER_LANGUAGE:
			s.HandleLanguage(notification, text)
		}
	})
}

// startRegistration initiates the registration process
func (s *FarmerRegistrationScene) startRegistration(notification *chatbot.Notification) {
	log.Printf("DEBUG: Starting farmer registration")
	notification.AnswerWithText("🌱 Great! Let's register you as a farmer.\n\nWhat's your full name?")
	notification.UpdateStateData(map[string]interface{}{
		"registration_state": STATE_REGISTER_NAME,
	})
	log.Printf("DEBUG: Set registration state to: %s", STATE_REGISTER_NAME)
}

// handleName processes the name input
func (s *FarmerRegistrationScene) HandleName(notification *chatbot.Notification, name string) {
	log.Printf("DEBUG: HandleName called with: '%s'", name)
	if strings.TrimSpace(name) == "" {
		log.Printf("DEBUG: Empty name provided")
		notification.AnswerWithText("Please enter your full name.")
		return
	}
	
	log.Printf("DEBUG: Setting name to: '%s' and state to: %s", name, STATE_REGISTER_CROP)
	notification.UpdateStateData(map[string]interface{}{
		"name": name,
		"registration_state": STATE_REGISTER_CROP,
	})
	notification.AnswerWithText(fmt.Sprintf("Nice to meet you, %s! 👋\n\nWhat type of crop do you grow? (e.g., maize, rice, wheat, vegetables)", name))
}

// handleCrop processes the crop input
func (s *FarmerRegistrationScene) HandleCrop(notification *chatbot.Notification, crop string) {
	log.Printf("DEBUG: HandleCrop called with: '%s'", crop)
	if strings.TrimSpace(crop) == "" {
		log.Printf("DEBUG: Empty crop provided")
		notification.AnswerWithText("Please tell me what crop you grow.")
		return
	}
	
	// Store the first crop
	crops := []string{crop}
	
	log.Printf("DEBUG: Setting first crop to: '%s' and state to: %s", crop, STATE_REGISTER_MORE_CROPS)
	notification.UpdateStateData(map[string]interface{}{
		"crops": crops,
		"registration_state": STATE_REGISTER_MORE_CROPS,
	})
	notification.AnswerWithText(fmt.Sprintf(MSG_MORE_CROPS_QUESTION, crop))
}

// handleMoreCrops processes additional crop inputs
func (s *FarmerRegistrationScene) HandleMoreCrops(notification *chatbot.Notification, response string) {
	log.Printf("DEBUG: HandleMoreCrops called with: '%s'", response)
	
	stateData := notification.GetStateData()
	crops, ok := stateData["crops"].([]string)
	if !ok {
		log.Printf("DEBUG: No crops found in state, resetting")
		notification.AnswerWithText("Something went wrong. Please start registration again with 'register'.")
		notification.UpdateStateData(map[string]interface{}{"registration_state": STATE_NONE})
		return
	}
	
	response = strings.ToLower(strings.TrimSpace(response))
	
	// Check if user wants to add more crops
	if response == "yes" {
		notification.AnswerWithText(MSG_ADD_MORE_CROPS)
		return
	}
	
	// Check if user is done adding crops
	if response == "no" || response == "done" {
		// Move to location registration
		cropsList := strings.Join(crops, ", ")
		log.Printf("DEBUG: Final crops list: %s, moving to location", cropsList)
		notification.UpdateStateData(map[string]interface{}{
			"registration_state": STATE_REGISTER_LOCATION,
		})
		notification.AnswerWithText(fmt.Sprintf(MSG_CROPS_COMPLETE, cropsList))
		return
	}
	
	// User provided another crop name
	if strings.TrimSpace(response) == "" {
		notification.AnswerWithText("Please tell me the crop name or type 'done' to finish.")
		return
	}
	
	// Add the new crop to the list
	crops = append(crops, response)
	log.Printf("DEBUG: Added crop '%s', total crops: %v", response, crops)
	
	// Update state with new crop list
	notification.UpdateStateData(map[string]interface{}{
		"crops": crops,
	})
	
	// Ask if they want to add more
	cropsList := strings.Join(crops, ", ")
	notification.AnswerWithText(fmt.Sprintf("Great! You grow: %s\n\nDo you grow any other crops? Type 'yes' to add more or 'done' to continue.", cropsList))
}

// handleLocation processes the location input
func (s *FarmerRegistrationScene) HandleLocation(notification *chatbot.Notification, location string) {
	log.Printf("DEBUG: HandleLocation called with: '%s'", location)
	if strings.TrimSpace(location) == "" {
		log.Printf("DEBUG: Empty location provided")
		notification.AnswerWithText("Please tell me your farm location.")
		return
	}
	
	log.Printf("DEBUG: Setting location to: '%s' and state to: %s", location, STATE_REGISTER_LANGUAGE)
	notification.UpdateStateData(map[string]interface{}{
		"location": location,
		"registration_state": STATE_REGISTER_LANGUAGE,
	})
	notification.AnswerWithText(fmt.Sprintf("Perfect! Your farm is in %s. 📍\n\nWhat language do you prefer for advice? (e.g., English, Swahili, French)", location))
}

// handleLanguage processes the language input and completes registration
func (s *FarmerRegistrationScene) HandleLanguage(notification *chatbot.Notification, language string) {
	log.Printf("DEBUG: HandleLanguage called with: '%s'", language)
	if strings.TrimSpace(language) == "" {
		log.Printf("DEBUG: Empty language provided")
		notification.AnswerWithText("Please tell me your preferred language.")
		return
	}
	
	// Get all registration data
	stateData := notification.GetStateData()
	name := stateData["name"].(string)
	crops, ok := stateData["crops"].([]string)
	if !ok {
		// Fallback to single crop if crops array not found
		if crop, exists := stateData["crop"].(string); exists {
			crops = []string{crop}
		} else {
			crops = []string{"Unknown"}
		}
	}
	location := stateData["location"].(string)
	
	// Save farmer profile
	farmerProfile := map[string]interface{}{
		"name": name,
		"crops": crops,
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
	log.Printf("DEBUG: Registration completed for %s, resetting state to NONE", name)
	cropsList := strings.Join(crops, ", ")
	completionMessage := fmt.Sprintf(`✅ Registration Complete!

👤 Name: %s
🌾 Crops: %s
📍 Location: %s
🗣️ Language: %s

You're all set! Now you can:
• Get farming advice with "advice"
• Send feedback with "feedback"
• Check your profile with "status"

Welcome to Farm Assistant! 🌱`, name, cropsList, location, language)
	
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