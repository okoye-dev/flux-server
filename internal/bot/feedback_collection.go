package bot

import (
	"log"
	"strings"

	chatbot "github.com/green-api/whatsapp-chatbot-golang"
)

// FeedbackCollectionScene handles feedback collection flow
type FeedbackCollectionScene struct {
	aiService *AIService
}

// NewFeedbackCollectionScene creates a new feedback collection scene
func NewFeedbackCollectionScene(aiService *AIService) *FeedbackCollectionScene {
	return &FeedbackCollectionScene{
		aiService: aiService,
	}
}

// Start begins the feedback collection scene
func (s FeedbackCollectionScene) Start(bot *chatbot.Bot) {
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

		// Handle feedback command
		if strings.Contains(actualMessage, CMD_FEEDBACK) {
			s.handleFeedbackRequest(notification, actualMessage)
			return
		}
	})
}

// handleFeedbackRequest processes feedback requests
func (s *FeedbackCollectionScene) handleFeedbackRequest(notification *chatbot.Notification, message string) {
	log.Printf("Processing feedback request: %s", message)
	
	// Get farmer profile from state
	stateData := notification.GetStateData()
	farmerProfile, ok := stateData["farmer_profile"].(FarmerProfile)
	if !ok {
		// If no profile found, ask to register first
		notification.AnswerWithText("❌ Please register first using 'Flux register' to provide feedback.")
		return
	}
	
	// Extract feedback content
	feedbackContent := s.extractFeedbackContent(message)
	if feedbackContent == "" {
		// If no specific feedback, show help
		notification.AnswerWithText(MSG_FEEDBACK_REQUEST)
		return
	}
	
	// Set state to collecting feedback
	notification.SetStateData(map[string]interface{}{
		"state": STATE_COLLECTING_FEEDBACK,
		"farmer_profile": farmerProfile,
	})
	
	// Process the feedback
	s.processAndStoreFeedback(notification, farmerProfile, feedbackContent)
}

// processAndStoreFeedback processes and stores farmer feedback
func (s *FeedbackCollectionScene) processAndStoreFeedback(notification *chatbot.Notification, profile FarmerProfile, feedback string) {
	log.Printf("📝 Processing feedback from %s: %s", profile.Name, feedback)
	
	// Process feedback with AI
	aiResponse, err := s.aiService.ProcessFeedback(profile, feedback)
	if err != nil {
		log.Printf("Error processing feedback: %v", err)
		notification.AnswerWithText("❌ Error processing your feedback. Please try again.")
		return
	}
	
	// Store feedback (in a real app, this would go to database)
	err = s.storeFeedback(profile, feedback, aiResponse)
	if err != nil {
		log.Printf("Error storing feedback: %v", err)
		notification.AnswerWithText("❌ Error saving your feedback. Please try again.")
		return
	}
	
	// Send acknowledgment
	notification.AnswerWithText(aiResponse)
	
	// Update state back to idle
	notification.SetStateData(map[string]interface{}{
		"state": STATE_IDLE,
		"farmer_profile": profile,
	})
	
	log.Printf("✅ Feedback processed and stored for farmer: %s", profile.Name)
}

// extractFeedbackContent extracts feedback content from the message
func (s *FeedbackCollectionScene) extractFeedbackContent(message string) string {
	// Remove the "feedback" command and get the actual feedback
	parts := strings.SplitN(message, " ", 2)
	if len(parts) < 2 {
		return ""
	}
	
	feedback := strings.TrimSpace(parts[1])
	
	// Check for common feedback patterns
	commonFeedback := map[string]string{
		"planted":        "I have planted my crops",
		"harvested":      "I have harvested my crops",
		"pest problem":   "I have pest problems",
		"weather issue":  "I have weather-related issues",
		"market update":  "I have market information to share",
		"good yield":     "I had a good yield this season",
		"poor yield":     "I had a poor yield this season",
		"irrigation":     "I need irrigation advice",
		"fertilizer":     "I need fertilizer advice",
	}
	
	lowerFeedback := strings.ToLower(feedback)
	for key, value := range commonFeedback {
		if strings.Contains(lowerFeedback, key) {
			return value
		}
	}
	
	// Return the original feedback if no pattern matches
	return feedback
}

// storeFeedback stores feedback in the system (dummy implementation)
func (s *FeedbackCollectionScene) storeFeedback(profile FarmerProfile, feedback, aiResponse string) error {
	log.Printf("💾 Storing feedback for farmer %s: %s", profile.Name, feedback)
	
	// In a real implementation, this would save to database
	// For now, just log the feedback
	log.Printf("✅ Feedback stored successfully")
	
	return nil
}

// isGroupChat checks if the message is from a group chat
func (s *FeedbackCollectionScene) isGroupChat(notification *chatbot.Notification) bool {
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
