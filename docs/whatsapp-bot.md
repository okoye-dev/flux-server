# WhatsApp Bot Integration

This application includes a WhatsApp bot that responds to greetings with personalized messages.

## Features

- Only responds to messages that start with "Flux"
- Responds to various greeting messages (hi, hello, hey, etc.) after "Flux"
- Sends personalized responses like "hey, [phone_number]"
- Handles case-insensitive greetings
- Provides helpful responses for non-greeting messages
- Ignores all other messages to prevent spam

## Setup

1. **Get Green API Credentials**

   - Sign up at [Green API](https://green-api.com/)
   - Create a new instance to get your `INSTANCE_ID` and `TOKEN`

2. **Configure Environment Variables**

   ```bash
   # Copy the example environment file
   cp env.example .env

   # Edit .env and add your WhatsApp credentials
   WHATSAPP_ENABLED=true
   WHATSAPP_INSTANCE_ID=your-instance-id
   WHATSAPP_TOKEN=your-token
   ```

3. **Run the Application**
   ```bash
   go run cmd/main.go
   ```

## How It Works

The bot listens for incoming WhatsApp messages and:

1. **Greeting Detection**: Recognizes various greeting patterns:

   - "hi", "hello", "hey", "hiya", "howdy", "greetings"
   - "good morning", "good afternoon", "good evening"
   - "what's up", "whats up", "sup", "yo"

2. **Personalized Response**: When a greeting is detected, responds with:

   ```
   hey, [sender_phone_number]
   ```

3. **Default Response**: For non-greeting messages, responds with:
   ```
   Hi! Say 'hi' to get a personalized greeting!
   ```

## Configuration Options

- `API_URL`: THE APIURL from green-api
- `WHATSAPP_ENABLED`: Set to `true` to enable the bot, `false` to disable
- `WHATSAPP_INSTANCE_ID`: Your Green API instance ID
- `WHATSAPP_TOKEN`: Your Green API token

## Example Usage

1. Start the application with WhatsApp bot enabled
2. Send a WhatsApp message to your bot instance
3. Try these messages:
   - "Flux hi" → "hey, +1234567890"
   - "Flux hello there" → "hey, +1234567890"
   - "Flux how are you?" → "Hi! Say 'Flux hi' to get a personalized greeting!"
   - "Regular message" → (ignored, no response)

## Troubleshooting

- Ensure your Green API instance is active and properly configured
- Check that your phone number is authorized in the Green API dashboard
- Verify that the webhook URL is correctly set in your Green API instance settings
- Check application logs for any error messages
