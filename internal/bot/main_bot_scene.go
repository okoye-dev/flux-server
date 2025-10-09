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

		// Check for ongoing registration first
		stateData := notification.GetStateData()
		registrationState := stateData["registration_state"]
		
		log.Printf("DEBUG: Message '%s' - Current registration state: %v", text, registrationState)
		
		// If user is in the middle of registration, let the registration scene handle it
		if registrationState != nil && registrationState != STATE_NONE {
			log.Printf("DEBUG: User is in registration state %v, handling ongoing registration", registrationState)
			// Handle ongoing registration directly
			s.handleOngoingRegistration(notification, text)
			return
		}
		
		log.Printf("DEBUG: No ongoing registration, processing as normal command")

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
	farmerProfileData, ok := stateData["farmer_profile"].(map[string]interface{})
	if !ok {
		notification.AnswerWithText("❌ You're not registered yet. Use 'register' to get started!")
		return
	}
	
	// Handle both single crop and multiple crops
	var cropsDisplay string
	if crops, ok := farmerProfileData["crops"].([]string); ok {
		cropsDisplay = strings.Join(crops, ", ")
	} else if crop, ok := farmerProfileData["crop"].(string); ok {
		cropsDisplay = crop
	} else {
		cropsDisplay = "Not specified"
	}

	// Get other profile data with fallbacks
	name, _ := farmerProfileData["name"].(string)
	location, _ := farmerProfileData["location"].(string)
	language, _ := farmerProfileData["language"].(string)
	phone, _ := farmerProfileData["phone"].(string)

	statusMessage := fmt.Sprintf(`👤 **Your Farmer Profile**

📝 **Name:** %s
🌱 **Crops:** %s
📍 **Location:** %s
🗣️ **Language:** %s
📱 **Phone:** %s

You can:
• Get advice with "advice"
• Send feedback with "feedback"
• Update your profile anytime`,
		name,
		cropsDisplay,
		location,
		language,
		phone,
	)
	
	notification.AnswerWithText(statusMessage)
}

// handleInvalidCommand handles invalid commands
func (s *MainBotScene) handleInvalidCommand(notification *chatbot.Notification) {
	notification.AnswerWithText(MSG_INVALID_COMMAND)
}

// handleOngoingRegistration handles messages during registration
func (s *MainBotScene) handleOngoingRegistration(notification *chatbot.Notification, text string) {
	stateData := notification.GetStateData()
	currentState := stateData["registration_state"]

	log.Printf("DEBUG: Handling ongoing registration - State: %v, Message: '%s'", currentState, text)

	switch currentState {
	case STATE_REGISTER_NAME:
		log.Printf("DEBUG: Processing name input: '%s'", text)
		s.registrationScene.HandleName(notification, text)
	case STATE_REGISTER_CROP:
		log.Printf("DEBUG: Processing crop input: '%s'", text)
		s.registrationScene.HandleCrop(notification, text)
	case STATE_REGISTER_MORE_CROPS:
		log.Printf("DEBUG: Processing more crops input: '%s'", text)
		s.registrationScene.HandleMoreCrops(notification, text)
	case STATE_REGISTER_LOCATION:
		log.Printf("DEBUG: Processing location input: '%s'", text)
		s.registrationScene.HandleLocation(notification, text)
	case STATE_REGISTER_LANGUAGE:
		log.Printf("DEBUG: Processing language input: '%s'", text)
		s.registrationScene.HandleLanguage(notification, text)
	default:
		log.Printf("DEBUG: Unknown registration state %v, resetting", currentState)
		// Unknown state, reset to main menu
		notification.UpdateStateData(map[string]interface{}{"registration_state": STATE_NONE})
		notification.AnswerWithText("Registration reset. Type 'register' to start again.")
	}
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
