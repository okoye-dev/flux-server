package bot

// Bot Commands
const (
	CMD_START     = "start"
	CMD_REGISTER  = "register"
	CMD_ADVICE    = "advice"
	CMD_FEEDBACK  = "feedback"
	CMD_HELP      = "help"
	CMD_STATUS    = "status"
)

// Bot Messages
const (
	MSG_WELCOME = `🌱 Welcome to Flux Farm Assistant!

I'm here to help you with:
• 📝 Farmer registration
• 🌤️ Weather-based advice
• 💰 Market price insights
• 🤖 AI-powered recommendations

Type "Flux help" to see all commands.`

	MSG_HELP = `📋 Available Commands:

• "Flux register" - Register as a farmer
• "Flux advice" - Get farming advice
• "Flux feedback" - Send feedback
• "Flux status" - Check your profile
• "Flux help" - Show this help

Just type "Flux" followed by any command!`

	MSG_REGISTER_START = `📝 Let's register you as a farmer!

Please provide your:
1. Full name
2. Crop type (e.g., maize, rice, wheat)
3. Location (city/region)
4. Language preference (English/Local)

Type your information in this format:
Name: [Your Name]
Crop: [Crop Type]
Location: [Your Location]
Language: [English/Local]`

	MSG_ADVICE_REQUEST = `🤖 Getting your personalized farming advice...

This may take a moment while I analyze:
• Your farm profile
• Current weather conditions
• Market prices
• Best practices`

	MSG_FEEDBACK_REQUEST = `📝 Share your feedback!

You can tell me about:
• "Planted" - I've planted my crops
• "Harvested" - I've harvested
• "Pest problem" - I have pest issues
• "Weather issue" - Weather problems
• "Market update" - Market information
• Or any other updates

Just type your feedback after "Flux feedback"`

	MSG_STATUS_CHECK = `👤 Checking your farmer profile...`

	MSG_INVALID_COMMAND = `❌ I didn't understand that command.

Type "Flux help" to see available commands.`

	MSG_REGISTRATION_COMPLETE = `✅ Registration complete!

Your farmer profile has been saved. You can now:
• Get personalized advice with "Flux advice"
• Update your status with "Flux feedback"
• Check your profile with "Flux status"`

	MSG_AI_PROCESSING = `🤖 Processing your request with AI...`
)

// Bot States
const (
	STATE_IDLE           = "idle"
	STATE_REGISTERING    = "registering"
	STATE_WAITING_ADVICE = "waiting_advice"
	STATE_COLLECTING_FEEDBACK = "collecting_feedback"
)

// AI Service Types
const (
	AI_TYPE_GEMINI = "gemini"
	AI_TYPE_WEATHER = "weather"
	AI_TYPE_MARKET = "market"
)

// Farmer Profile Fields
const (
	FIELD_NAME     = "name"
	FIELD_CROP     = "crop"
	FIELD_LOCATION = "location"
	FIELD_LANGUAGE = "language"
)
