package bot

import (
	"fmt"
	"log"
	"strings"

	chatbot "github.com/green-api/whatsapp-chatbot-golang"
)

// MainBotScene handles the main bot flow and command routing
type MainBotScene struct {
	aiService              *AIService
	registrationScene      *FarmerRegistrationScene
	adviceScene           *AdviceDeliveryScene
	feedbackScene         *FeedbackCollectionScene
}

// NewMainBotScene creates a new main bot scene
func NewMainBotScene(aiService *AIService) *MainBotScene {
	return &MainBotScene{
		aiService:         aiService,
		registrationScene: NewFarmerRegistrationScene(aiService),
		adviceScene:      NewAdviceDeliveryScene(aiService),
		feedbackScene:    NewFeedbackCollectionScene(aiService),
	}
}


// Start begins the main bot scene (for polling mode - not used in webhook mode)
func (s MainBotScene) Start(bot *chatbot.Bot) {
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
		actualMessage := strings.ToLower(strings.TrimSpace(text))

		// Route to appropriate handler based on command
		s.routeCommand(notification, actualMessage)
	})
}

// routeCommand routes commands to appropriate handlers
func (s *MainBotScene) routeCommand(notification *chatbot.Notification, message string) {
	
	// Handle different commands
	switch {
	case strings.Contains(message, CMD_START) || message == "":
		s.handleStart(notification)
	case strings.Contains(message, CMD_HI) || strings.Contains(message, CMD_HEY):
		s.handleGreeting(notification)
	case strings.Contains(message, CMD_REGISTER):
		s.registrationScene.startRegistration(notification)
	case strings.Contains(message, CMD_ADVICE):
		s.adviceScene.handleAdviceRequest(notification)
	case strings.Contains(message, CMD_FEEDBACK):
		s.feedbackScene.handleFeedbackRequest(notification, message)
	case strings.Contains(message, CMD_HELP):
		s.handleHelp(notification)
	case strings.Contains(message, CMD_STATUS):
		s.handleStatus(notification)
	default:
		s.handleInvalidCommand(notification)
	}
}

// handleStart handles the start command
func (s *MainBotScene) handleStart(notification *chatbot.Notification) {
	// Get sender info
	sender, err := notification.Sender()
	if err != nil {
		sender = "there"
		log.Printf("Error getting sender: %v", err)
	}
	
	welcomeMessage := fmt.Sprintf("Hey, %s! %s", sender, MSG_WELCOME)
	notification.AnswerWithText(welcomeMessage)
}

// handleGreeting handles hi/hey commands
func (s *MainBotScene) handleGreeting(notification *chatbot.Notification) {
	// Get sender info
	sender, err := notification.Sender()
	if err != nil {
		sender = "there"
		log.Printf("Error getting sender: %v", err)
	}
	
	greeting := fmt.Sprintf("Hey %s! 👋\n\n%s", sender, MSG_WELCOME)
	notification.AnswerWithText(greeting)
}

// handleHelp handles the help command
func (s *MainBotScene) handleHelp(notification *chatbot.Notification) {
	notification.AnswerWithText(MSG_HELP)
}

// handleStatus handles the status command
func (s *MainBotScene) handleStatus(notification *chatbot.Notification) {
	
	// Get farmer profile from state
	stateData := notification.GetStateData()
	farmerProfile, ok := stateData["farmer_profile"].(FarmerProfile)
	if !ok {
		notification.AnswerWithText("❌ You're not registered yet. Use 'register' to get started!")
		return
	}
	
	statusMessage := fmt.Sprintf(`👤 **Your Farmer Profile**

📝 **Name:** %s
🌱 **Crop:** %s
📍 **Location:** %s
🗣️ **Language:** %s
📱 **Phone:** %s

You can:
• Get advice with "advice"
• Send feedback with "feedback"
• Update your profile anytime`,
		farmerProfile.Name,
		farmerProfile.Crop,
		farmerProfile.Location,
		farmerProfile.Language,
		farmerProfile.Phone,
	)
	
	notification.AnswerWithText(statusMessage)
}

// handleInvalidCommand handles invalid commands
func (s *MainBotScene) handleInvalidCommand(notification *chatbot.Notification) {
	notification.AnswerWithText(MSG_INVALID_COMMAND)
}

// isGroupChat checks if the message is from a group chat
func (s *MainBotScene) isGroupChat(notification *chatbot.Notification) bool {
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
