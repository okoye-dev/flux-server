# WhatsApp Bot Demo Flow

## Complete Demo Flow Implementation

### 1. Environment Setup

Make sure your `.env` file has:

```env
API_KEY=AIzaSyBESrT9GiRbt27-J42qJ_U48zSfrIruffw
# ... other environment variables
```

### 2. Demo Flow Steps

#### Step 1: Farmer Registration

```
User: "register"
Bot: "🌱 Great! Let's register you as a farmer.

What's your full name?"

User: "John Doe"
Bot: "Nice to meet you, John Doe! 👋

What type of crop do you grow? (e.g., maize, rice, wheat, vegetables)"

User: "maize"
Bot: "🌾 Great! You grow maize.

Do you grow any other crops?
• Type "yes" to add more crops
• Type "no" to continue with location"

User: "yes"
Bot: "🌱 What other crop do you grow?
(e.g., maize, rice, wheat, vegetables, beans, etc.)

Type "done" when you're finished adding crops."

User: "beans"
Bot: "Great! You grow: maize, beans

Do you grow any other crops? Type "yes" to add more or "done" to continue."

User: "done"
Bot: "✅ Perfect! You grow: maize, beans

Now, where is your farm located? (e.g., city, region, state)"

User: "Kano, Nigeria"
Bot: "Perfect! Your farm is in Kano, Nigeria. 📍

What language do you prefer for advice? (e.g., English, Swahili, French)"

User: "English"
Bot: "✅ Registration Complete!

👤 Name: John Doe
🌾 Crops: maize, beans
📍 Location: Kano, Nigeria
🗣️ Language: English

You're all set! Now you can:
• Get farming advice with "advice"
• Send feedback with "feedback"
• Check your profile with "status"

Welcome to Farm Assistant! 🌱"
```

#### Step 2: Get AI-Tailored Advice

```
User: "advice"
Bot: "🤖 Getting your personalized farming advice...

This may take a moment while I analyze:
• Your farm profile
• Current weather conditions
• Market prices
• Best practices"

[AI processes with Gemini API using farmer's data]

Bot: "🌱 **AI-Powered Farming Advice**

🌱 **Planting Advice:**
Based on current weather conditions in Kano, it's optimal to plant maize and beans in the next 3-5 days. Soil moisture levels are ideal for both crops.

💧 **Irrigation Advice:**
With current humidity at 65%, reduce irrigation frequency to every 3 days. Monitor soil moisture closely for both maize and beans.

🌾 **Harvest Advice:**
Your maize and beans crops should be ready for harvest in approximately 45-60 days based on current growth conditions in Kano.

💰 **Market Advice:**
Current market prices for your crops (maize, beans) are favorable. Monitor local Kano markets for optimal selling timing.

📋 **General Advice:**
Continue monitoring your crops regularly and maintain proper farming practices. Consider crop rotation between maize and beans for better soil health.

*Confidence: 85% | Generated: 2025-01-27T10:30:00Z*"
```

#### Step 3: Get Market Insights

```
User: "market"
Bot: "💰 **Market Insights**

🌾 **Rice** in Kano markets: ₦900 per bag
🌽 **Maize** in Lagos markets: ₦650 per bag
🍅 **Tomatoes** in Abuja markets: ₦1,200 per basket
🥜 **Groundnuts** in Kaduna markets: ₦800 per bag

*Prices updated 2 hours ago*

For detailed market analysis, type "go" to access our web app."
```

#### Step 4: Access Web App

```
User: "go"
Bot: "🌐 **Access Our Web App**

Visit: https://flux-farm-app.vercel.app

Our web app provides:
• 📊 Detailed market analysis
• 🌤️ Advanced weather forecasts
• 📚 Learning & Advisory Center
• 🤖 AI-curated learning modules
• 👨‍🌾 Expert consultations
• 💡 Daily farming tips

*Bookmark this link for easy access!*"
```

#### Step 5: Check Status

```
User: "status"
Bot: "👤 **Your Farmer Profile**

📝 **Name:** John Doe
🌱 **Crops:** maize, beans
📍 **Location:** Kano, Nigeria
🗣️ **Language:** English
📱 **Phone:** Not specified

You can:
• Get advice with "advice"
• Send feedback with "feedback"
• Update your profile anytime"
```

### 3. Key Features Implemented

✅ **Multiple Crop Registration**: Farmers can register multiple crops
✅ **Real Gemini AI Integration**: Uses actual API key for AI advice
✅ **Market Insights**: Provides current market prices
✅ **Web App Access**: Links to the web application
✅ **State Management**: Stores farmer information temporarily
✅ **Weather Integration**: Considers weather data in advice
✅ **Location-Based**: Tailored advice for specific locations

### 4. AI Integration Details

The bot now:

- Uses the real Gemini API with your API key
- Sends farmer profile, weather, and market data to Gemini
- Receives structured advice back from the AI
- Formats the response for WhatsApp display
- Handles API errors gracefully

### 5. Commands Available

- `register` - Register as a farmer (multiple crops supported)
- `advice` - Get AI-tailored farming advice
- `market` - Get market prices and insights
- `go` - Access the web app
- `status` - Check your farmer profile
- `feedback` - Send feedback
- `help` - Show all commands

The demo flow is now complete and ready for testing! 🚀
