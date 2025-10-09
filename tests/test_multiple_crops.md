# Multiple Crop Registration Test

## Overview

This document demonstrates how the multiple crop registration feature works in the WhatsApp bot.

## Test Flow

### 1. Start Registration

User sends: `register`

Bot responds:

```
🌱 Great! Let's register you as a farmer.

What's your full name?
```

### 2. Enter Name

User sends: `John Doe`

Bot responds:

```
Nice to meet you, John Doe! 👋

What type of crop do you grow? (e.g., maize, rice, wheat, vegetables)
```

### 3. Enter First Crop

User sends: `maize`

Bot responds:

```
🌾 Great! You grow maize.

Do you grow any other crops?
• Type "yes" to add more crops
• Type "no" to continue with location
```

### 4. Add More Crops

User sends: `yes`

Bot responds:

```
🌱 What other crop do you grow?
(e.g., maize, rice, wheat, vegetables, beans, etc.)

Type "done" when you're finished adding crops.
```

### 5. Add Second Crop

User sends: `beans`

Bot responds:

```
Great! You grow: maize, beans

Do you grow any other crops? Type "yes" to add more or "done" to continue.
```

### 6. Add Third Crop

User sends: `tomatoes`

Bot responds:

```
Great! You grow: maize, beans, tomatoes

Do you grow any other crops? Type "yes" to add more or "done" to continue.
```

### 7. Finish Adding Crops

User sends: `done`

Bot responds:

```
✅ Perfect! You grow: maize, beans, tomatoes

Now, where is your farm located? (e.g., city, region, state)
```

### 8. Continue with Location and Language

The registration continues with location and language as before.

### 9. Final Registration Complete

Bot responds:

```
✅ Registration Complete!

👤 Name: John Doe
🌾 Crops: maize, beans, tomatoes
📍 Location: Nairobi, Kenya
🗣️ Language: English

You're all set! Now you can:
• Get farming advice with "advice"
• Send feedback with "feedback"
• Check your profile with "status"

Welcome to Farm Assistant! 🌱
```

## Database Changes

### New Table: farmer_crops

```sql
CREATE TABLE farmer_crops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id BIGINT NOT NULL,
    crop_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_farmer_crops_farmer FOREIGN KEY (farmer_id) REFERENCES farmers(id) ON DELETE CASCADE,
    CONSTRAINT fk_farmer_crops_crop FOREIGN KEY (crop_id) REFERENCES crops(id) ON DELETE CASCADE,
    CONSTRAINT unique_farmer_crop UNIQUE (farmer_id, crop_id)
);
```

### Updated Models

- `FarmerProfile` now has `Crops []string` instead of `Crop string`
- Added `FarmerCrop` junction table model
- Added `FarmerWithCrops` and `FarmerCropWithDetails` models

## Key Features

1. **Multiple Crop Support**: Farmers can register multiple crops during registration
2. **Flexible Flow**: Users can add as many crops as they want
3. **Backward Compatibility**: System handles both single crop (legacy) and multiple crop data
4. **Database Normalization**: Uses proper junction table for many-to-many relationship
5. **Auto Crop Creation**: New crops are automatically created in the database if they don't exist

## Status Command

When users check their status with `status`, they now see:

```
👤 **Your Farmer Profile**

📝 **Name:** John Doe
🌱 **Crops:** maize, beans, tomatoes
📍 **Location:** Nairobi, Kenya
🗣️ **Language:** English
📱 **Phone:** +254700000000

You can:
• Get advice with "advice"
• Send feedback with "feedback"
• Update your profile anytime
```

## AI Advice Integration

The AI service now provides advice for all registered crops:

- Planting advice considers all crops
- Harvest advice covers all crops
- Market advice provides information for all crops
