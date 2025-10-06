package services

import (
	"fmt"
	"log"
	"strings"

	chatbot "github.com/green-api/whatsapp-chatbot-golang"
)

// WhatsAppBot represents the WhatsApp bot service
type WhatsAppBot struct {
	bot *chatbot.Bot
}

// NewWhatsAppBot creates a new WhatsApp bot instance
func NewWhatsAppBot(instanceID, token string) *WhatsAppBot {
	bot := chatbot.NewBot(instanceID, token)
	bot.SetStartScene(GreetingScene{})
	
	return &WhatsAppBot{
		bot: bot,
	}
}

// Start starts the WhatsApp bot
func (w *WhatsAppBot) Start() {
	log.Println("Starting WhatsApp bot...")
	
	// Handle errors from the bot
	go func() {
		for err := range w.bot.ErrorChannel {
			if err != nil {
				log.Printf("\nWhatsApp bot error: %v \n", err)
			}
		}
	}()
	
	w.bot.StartReceivingNotifications()
}

// GreetingScene handles the main greeting logic
type GreetingScene struct{}

func (s GreetingScene) Start(bot *chatbot.Bot) {
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
		
		// Check if the remaining message is a greeting
		if s.isGreeting(actualMessage) {
			// Get sender's phone number
			sender, err := notification.Sender()
			if err != nil || sender == "" {
				sender = "there"
			}

			// Send greeting response
			response := fmt.Sprintf("hey, %s", sender)
			notification.AnswerWithText(response)
			log.Printf("Sent greeting to %s: %s", sender, response)
		} else {
			// Send a default response for non-greeting messages
			notification.AnswerWithText("Hi! Say 'Flux hi' to get a personalized greeting!")
		}
	})
}

// isGreeting checks if the message is a greeting
func (s GreetingScene) isGreeting(text string) bool {
	greetings := []string{
		"hi", "hello", "hey", "hiya", "howdy", "greetings",
		"good morning", "good afternoon", "good evening",
		"what's up", "whats up", "sup", "yo",
	}

	for _, greeting := range greetings {
		if strings.Contains(text, greeting) {
			return true
		}
	}
	return false
}

// isGroupChat checks if the message is from a group chat
func (s GreetingScene) isGroupChat(notification *chatbot.Notification) bool {
	// Check the webhook body for group chat indicators
	body := notification.Body
	
	// Check senderData for group chat ID
	if senderData, ok := body["senderData"].(map[string]interface{}); ok {
		if chatId, exists := senderData["chatId"]; exists {
			if chatIdStr, ok := chatId.(string); ok {
				// Group chat IDs end with @g.us
				return strings.HasSuffix(chatIdStr, "@g.us")
			}
		}
	}
	
	return false
}
