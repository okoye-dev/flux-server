package services

import (
	"log"

	chatbot "github.com/green-api/whatsapp-chatbot-golang"
	"github.com/okoye-dev/flux-server/internal/bot"
)

// WhatsAppBot represents the WhatsApp bot service
type WhatsAppBot struct {
	bot       *chatbot.Bot
	aiService *bot.AIService
	mainScene *bot.MainBotScene
}

// NewWhatsAppBot creates a new WhatsApp bot instance
func NewWhatsAppBot(instanceID, token string) *WhatsAppBot {
	chatbotInstance := chatbot.NewBot(instanceID, token)
	
	// Initialize AI service
	aiService := bot.NewAIService()
	
	// Initialize main scene with all sub-scenes
	mainScene := bot.NewMainBotScene(aiService)
	
	// Set the main scene as the start scene
	chatbotInstance.SetStartScene(*mainScene)
	
	return &WhatsAppBot{
		bot:       chatbotInstance,
		aiService: aiService,
		mainScene: mainScene,
	}
}

// Start starts the WhatsApp bot using Green API polling
func (w *WhatsAppBot) Start() {
	log.Println("Starting WhatsApp bot with Green API polling...")
	
	// Handle errors from the bot
	go func() {
		for err := range w.bot.ErrorChannel {
			if err != nil {
				log.Printf("WhatsApp bot error: %v", err)
			}
		}
	}()
	
	// Start receiving notifications using Green API polling
	w.bot.StartReceivingNotifications()
	log.Println("WhatsApp bot started successfully and polling for messages...")
}

