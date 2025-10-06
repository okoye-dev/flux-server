package bot

import (
	"fmt"
	"log"
	"regexp"
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


// Start begins the farmer registration scene (for polling mode - not used in webhook mode)
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
			log.Printf("Ignoring group chat message: %s", text)
			return
		}

		// Convert to lowercase for case-insensitive matching
		lowerText := strings.ToLower(strings.TrimSpace(text))

		// Only respond to messages that start with "Flux"
		if !strings.HasPrefix(lowerText, "flux") {
			return
		}

		// Remove "flux" prefix and get the actual message
		actualMessage := strings.TrimSpace(strings.TrimPrefix(lowerText, "flux"))

		// Handle registration command
		if strings.Contains(actualMessage, CMD_REGISTER) {
			s.handleRegistrationStart(notification)
			return
		}

		// Handle registration data input
		if s.isRegistrationData(text) {
			s.handleRegistrationData(notification, text)
			return
		}
	})
}

// handleRegistrationStart initiates the registration process
func (s *FarmerRegistrationScene) handleRegistrationStart(notification *chatbot.Notification) {
	log.Printf("Starting farmer registration process")
	
	// Set state to registering
	notification.SetStateData(map[string]interface{}{
		"state": STATE_REGISTERING,
		"step":  "waiting_for_data",
	})
	
	notification.AnswerWithText(MSG_REGISTER_START)
}

// handleRegistrationData processes the registration data
func (s *FarmerRegistrationScene) handleRegistrationData(notification *chatbot.Notification, text string) {
	log.Printf("Processing registration data: %s", text)
	
	// Parse the registration data
	profile, err := s.parseRegistrationData(text)
	if err != nil {
		notification.AnswerWithText("❌ Invalid format. Please provide your information in this format:\n\nName: [Your Name]\nCrop: [Crop Type]\nLocation: [Your Location]\nLanguage: [English/Local]")
		return
	}
	
	// Get sender phone number
	sender, err := notification.Sender()
	if err != nil {
		sender = "unknown"
	}
	profile.Phone = sender
	
	// Store the farmer profile (in a real app, this would go to database)
	err = s.storeFarmerProfile(profile)
	if err != nil {
		log.Printf("Error storing farmer profile: %v", err)
		notification.AnswerWithText("❌ Error saving your profile. Please try again.")
		return
	}
	
	// Update state to idle
	notification.SetStateData(map[string]interface{}{
		"state": STATE_IDLE,
		"farmer_profile": profile,
	})
	
	log.Printf("✅ Farmer registration completed for: %s", profile.Name)
	notification.AnswerWithText(MSG_REGISTRATION_COMPLETE)
}

// parseRegistrationData extracts farmer information from text
func (s *FarmerRegistrationScene) parseRegistrationData(text string) (*FarmerProfile, error) {
	profile := &FarmerProfile{}
	
	// Use regex to extract information
	patterns := map[string]*regexp.Regexp{
		FIELD_NAME:     regexp.MustCompile(`(?i)name:\s*(.+)`),
		FIELD_CROP:     regexp.MustCompile(`(?i)crop:\s*(.+)`),
		FIELD_LOCATION: regexp.MustCompile(`(?i)location:\s*(.+)`),
		FIELD_LANGUAGE: regexp.MustCompile(`(?i)language:\s*(.+)`),
	}
	
	fields := map[string]*string{
		FIELD_NAME:     &profile.Name,
		FIELD_CROP:     &profile.Crop,
		FIELD_LOCATION: &profile.Location,
		FIELD_LANGUAGE: &profile.Language,
	}
	
	// Extract each field
	for field, pattern := range patterns {
		matches := pattern.FindStringSubmatch(text)
		if len(matches) < 2 {
			return nil, fmt.Errorf("missing field: %s", field)
		}
		*fields[field] = strings.TrimSpace(matches[1])
	}
	
	// Validate required fields
	if profile.Name == "" || profile.Crop == "" || profile.Location == "" || profile.Language == "" {
		return nil, fmt.Errorf("missing required fields")
	}
	
	return profile, nil
}

// storeFarmerProfile stores the farmer profile (dummy implementation)
func (s *FarmerRegistrationScene) storeFarmerProfile(profile *FarmerProfile) error {
	log.Printf("💾 Storing farmer profile: %+v", profile)
	
	// In a real implementation, this would save to database
	// For now, just log the profile
	log.Printf("✅ Farmer profile stored successfully")
	
	return nil
}

// isRegistrationData checks if the message contains registration data
func (s *FarmerRegistrationScene) isRegistrationData(text string) bool {
	// Check if text contains all required fields
	requiredFields := []string{"name:", "crop:", "location:", "language:"}
	lowerText := strings.ToLower(text)
	
	for _, field := range requiredFields {
		if !strings.Contains(lowerText, field) {
			return false
		}
	}
	
	return true
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
