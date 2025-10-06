# Deployment Guide

## 🚀 Railway (Recommended)

1. **Deploy**: Go to [Railway.app](https://railway.app) → New Project → Deploy from GitHub
2. **Set Environment Variables** in Railway dashboard:
   ```bash
   SUPABASE_URL=https://code.supabase.co
   SUPABASE_ANON_KEY=xxx.xxx.xxx
   SUPABASE_SERVICE_ROLE_KEY=xxx.xxx.yyy
   JWT_SECRET=xxx+xxx/xxx/xxx+Q9sp5EPg==
   API_URL=https://0000.api.greenapi.com
   WHATSAPP_ENABLED=true
   WHATSAPP_INSTANCE_ID=0000
   WHATSAPP_TOKEN=abcd
   PORT=8080
   ENVIRONMENT=production
   ```
3. **Get your webhook URL**: `https://your-app.railway.app/webhook/whatsapp`
4. **Update Green API**:
   ```bash
   curl -X POST "https://000.api.greenapi.com/waInstance0000/setSettings/abcd" \
     -H "Content-Type: application/json" \
     -d '{"webhookUrl": "https://your-app.railway.app/webhook/whatsapp", "incomingWebhook": "yes", "outgoingMessageWebhook": "yes", "outgoingAPIMessageWebhook": "yes"}'
   ```

## 🎯 Alternative: Render

1. Go to [Render.com](https://render.com) → New Web Service
2. Connect GitHub repo
3. Set same environment variables
4. Webhook URL: `https://your-app.onrender.com/webhook/whatsapp`

## ✅ Test

Send "Flux hi" to your WhatsApp → Should get "hey, [phone_number]"
