package models

import (
	"time"

	"github.com/google/uuid"
)

// Farmer represents a farmer in the system
type Farmer struct {
	ID          int64      `json:"id" db:"id"`
	AuthUserID  *uuid.UUID `json:"auth_user_id" db:"auth_user_id"` // Links to auth.users.id
	Name        string     `json:"name" db:"name"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	CropType    string     `json:"crop_type" db:"crop_type"`
	LocationID  int64      `json:"location_id" db:"location_id"`
	Language    string     `json:"language" db:"language"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// ExtensionOfficer represents an extension officer in the system
type ExtensionOfficer struct {
	ID                  int64      `json:"id" db:"id"`
	AuthUserID          *uuid.UUID `json:"auth_user_id" db:"auth_user_id"` // Links to auth.users.id
	Name                string     `json:"name" db:"name"`
	PhoneNumber         string     `json:"phone_number" db:"phone_number"`
	AssignedLocationID  *int64     `json:"assigned_location_id" db:"assigned_location_id"`
}

// Role represents a user role in the system
type Role struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// UserProfile represents a user profile linked to auth.users
type UserProfile struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	AuthUserID  *uuid.UUID      `json:"auth_user_id" db:"auth_user_id"` // Links to auth.users.id
	RoleID      *uuid.UUID      `json:"role_id" db:"role_id"`           // FK to roles.id
	DisplayName *string         `json:"display_name" db:"display_name"`
	Phone       *string         `json:"phone" db:"phone"`
	Metadata    map[string]any  `json:"metadata" db:"metadata"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

// Location represents a geographical location
type Location struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Country   *string   `json:"country" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Crop represents a crop type
type Crop struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	ScientificName *string   `json:"scientific_name" db:"scientific_name"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// FarmerCrop represents the junction table between farmers and crops
type FarmerCrop struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FarmerID  int64     `json:"farmer_id" db:"farmer_id"`   // FK to farmers.id
	CropID    uuid.UUID `json:"crop_id" db:"crop_id"`       // FK to crops.id
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// FarmHarvest represents a farm harvest record
type FarmHarvest struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserProfileID *uuid.UUID `json:"user_profile_id" db:"user_profile_id"` // FK to user_profiles.id
	CropID        *uuid.UUID `json:"crop_id" db:"crop_id"`                 // FK to crops.id
	Quantity      *float64   `json:"quantity" db:"quantity"`
	HarvestedAt   *time.Time `json:"harvested_at" db:"harvested_at"`
}

// AuthUser represents the auth.users table fields we care about
type AuthUser struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	Email             *string                `json:"email" db:"email"`
	Phone             *string                `json:"phone" db:"phone"`
	ConfirmedAt       *time.Time             `json:"confirmed_at" db:"confirmed_at"`
	PhoneConfirmedAt  *time.Time             `json:"phone_confirmed_at" db:"phone_confirmed_at"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	LastSignInAt      *time.Time             `json:"last_sign_in_at" db:"last_sign_in_at"`
	RawUserMetaData   map[string]any         `json:"raw_user_meta_data" db:"raw_user_meta_data"`
	UserMetadata      map[string]any         `json:"user_metadata" db:"user_metadata"`
}

// UserProfileWithRole represents a user profile with role information
type UserProfileWithRole struct {
	UserProfile
	Role *Role `json:"role,omitempty"`
}

// UserProfileWithAuth represents a user profile with auth information
type UserProfileWithAuth struct {
	UserProfile
	AuthUser *AuthUser `json:"auth_user,omitempty"`
}

// FarmHarvestWithDetails represents a farm harvest with related data
type FarmHarvestWithDetails struct {
	FarmHarvest
	UserProfile *UserProfile `json:"user_profile,omitempty"`
	Crop        *Crop        `json:"crop,omitempty"`
}

// CreateUserProfileRequest represents the request to create a user profile
type CreateUserProfileRequest struct {
	AuthUserID  uuid.UUID      `json:"auth_user_id" validate:"required"`
	RoleID      *uuid.UUID     `json:"role_id,omitempty"`
	DisplayName *string        `json:"display_name,omitempty"`
	Phone       *string        `json:"phone,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// UpdateUserProfileRequest represents the request to update a user profile
type UpdateUserProfileRequest struct {
	RoleID      *uuid.UUID     `json:"role_id,omitempty"`
	DisplayName *string        `json:"display_name,omitempty"`
	Phone       *string        `json:"phone,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// CreateFarmHarvestRequest represents the request to create a farm harvest
type CreateFarmHarvestRequest struct {
	UserProfileID *uuid.UUID `json:"user_profile_id,omitempty"`
	CropID        *uuid.UUID `json:"crop_id,omitempty"`
	Quantity      *float64   `json:"quantity,omitempty"`
	HarvestedAt   *time.Time `json:"harvested_at,omitempty"`
}

// FarmerWithCrops represents a farmer with their associated crops
type FarmerWithCrops struct {
	Farmer
	Crops []Crop `json:"crops,omitempty"`
}

// FarmerCropWithDetails represents a farmer-crop relationship with crop details
type FarmerCropWithDetails struct {
	FarmerCrop
	Crop *Crop `json:"crop,omitempty"`
}

// CreateFarmerCropRequest represents the request to add a crop to a farmer
type CreateFarmerCropRequest struct {
	FarmerID int64     `json:"farmer_id" validate:"required"`
	CropID   uuid.UUID `json:"crop_id" validate:"required"`
}

// AddCropsToFarmerRequest represents the request to add multiple crops to a farmer
type AddCropsToFarmerRequest struct {
	FarmerID int64       `json:"farmer_id" validate:"required"`
	CropIDs  []uuid.UUID `json:"crop_ids" validate:"required,min=1"`
}