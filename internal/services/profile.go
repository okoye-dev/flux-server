package services

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/okoye-dev/flux-server/internal/models"
	"github.com/supabase-community/supabase-go"
)

// SignupData contains additional signup information
type SignupData struct {
	PhoneNumber        string
	CropType           string
	LocationID         int64
	Language           string
	AssignedLocationID int64
}

// ProfileService handles user profile operations
type ProfileService struct {
	client *supabase.Client
}

// NewProfileService creates a new profile service
func NewProfileService() (*ProfileService, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseAnonKey == "" {
		return nil, ErrSupabaseConfigMissing
	}

	client, err := supabase.NewClient(supabaseURL, supabaseAnonKey, nil)
	if err != nil {
		return nil, err
	}

	return &ProfileService{client: client}, nil
}

// CreateUserProfile creates a user profile after successful signup
func (s *ProfileService) CreateUserProfile(authUserID, username, roleName string, signupData *SignupData) (*models.UserProfile, error) {
	// Parse auth user ID
	authUUID, err := uuid.Parse(authUserID)
	if err != nil {
		return nil, err
	}

	// Default to farmer role if not specified
	if roleName == "" {
		roleName = "farmer"
	}

	// Get role ID from role name
	roleID, err := s.GetRoleIDByName(roleName)
	if err != nil {
		return nil, err
	}

	// Generate a new UUID for the profile
	profileID := uuid.New()

	// Create user profile
	now := time.Now()
	profile := models.UserProfile{
		ID:          profileID,
		AuthUserID:  &authUUID,
		RoleID:      roleID,
		DisplayName: &username,
		Metadata:    map[string]any{}, // Empty metadata - role and username are in dedicated fields
		CreatedAt:   now,
	}

	// Insert into user_profiles table
	var result []models.UserProfile
	_, err = s.client.From("user_profiles").Insert(profile, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrProfileCreationFailed
	}

	// Create role-specific record
	err = s.createRoleSpecificRecord(authUserID, roleName, username, signupData)
	if err != nil {
		// Log error but don't fail - user profile is created
		// TODO: Add proper logging
		_ = err
	}

	return &result[0], nil
}

// GetUserProfile retrieves a user profile by auth user ID
func (s *ProfileService) GetUserProfile(authUserID string) (*models.UserProfile, error) {
	var result []models.UserProfile
	_, err := s.client.From("user_profiles").Select("*", "", false).Eq("auth_user_id", authUserID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, ErrProfileNotFound
}

// GetRoleIDByName gets role ID by role name
func (s *ProfileService) GetRoleIDByName(roleName string) (*uuid.UUID, error) {
	var result []models.Role
	_, err := s.client.From("roles").Select("id", "", false).Eq("name", roleName).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0].ID, nil
	}

	return nil, ErrRoleNotFound
}

// createRoleSpecificRecord creates farmer or extension officer record based on role
func (s *ProfileService) createRoleSpecificRecord(authUserID, roleName, username string, signupData *SignupData) error {
	authUUID, err := uuid.Parse(authUserID)
	if err != nil {
		return err
	}

	switch roleName {
	case "farmer":
		return s.createFarmerRecord(authUUID, username, signupData)
	case "extension_officer":
		return s.createExtensionOfficerRecord(authUUID, username, signupData)
	default:
		// Unknown role, don't create specific record
		return nil
	}
}

// createFarmerRecord creates a farmer record
func (s *ProfileService) createFarmerRecord(authUserID uuid.UUID, username string, signupData *SignupData) error {
	// Generate a new ID for the farmer record
	farmerID := generateFarmerID()

	// Set defaults if not provided
	phoneNumber := signupData.PhoneNumber
	cropType := signupData.CropType
	locationID := signupData.LocationID
	language := signupData.Language
	if language == "" {
		language = "en"
	}

	farmer := models.Farmer{
		ID:          farmerID,
		AuthUserID:  &authUserID, // Link to auth user
		Name:        username,
		PhoneNumber: phoneNumber,
		CropType:    cropType,
		LocationID:  locationID,
		Language:    language,
		CreatedAt:   time.Now(),
	}

	var result []models.Farmer
	_, err := s.client.From("farmers").Insert(farmer, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	// If we have crop information, add it to the farmer_crops table
	if cropType != "" {
		return s.addCropToFarmer(farmerID, cropType)
	}

	return nil
}

// addCropToFarmer adds a crop to a farmer's profile
func (s *ProfileService) addCropToFarmer(farmerID int64, cropName string) error {
	// First, find or create the crop
	cropID, err := s.findOrCreateCrop(cropName)
	if err != nil {
		return err
	}

	// Create farmer-crop relationship
	farmerCrop := models.FarmerCrop{
		ID:        uuid.New(),
		FarmerID:  farmerID,
		CropID:    *cropID,
		CreatedAt: time.Now(),
	}

	var result []models.FarmerCrop
	_, err = s.client.From("farmer_crops").Insert(farmerCrop, false, "", "", "").ExecuteTo(&result)
	return err
}

// findOrCreateCrop finds an existing crop or creates a new one
func (s *ProfileService) findOrCreateCrop(cropName string) (*uuid.UUID, error) {
	// First, try to find existing crop
	var existingCrops []models.Crop
	_, err := s.client.From("crops").Select("id", "", false).Eq("name", cropName).ExecuteTo(&existingCrops)
	if err != nil {
		return nil, err
	}

	if len(existingCrops) > 0 {
		return &existingCrops[0].ID, nil
	}

	// Create new crop if not found
	newCrop := models.Crop{
		ID:        uuid.New(),
		Name:      cropName,
		CreatedAt: time.Now(),
	}

	var result []models.Crop
	_, err = s.client.From("crops").Insert(newCrop, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0].ID, nil
	}

	return nil, fmt.Errorf("failed to create crop: %s", cropName)
}

// AddCropsToFarmer adds multiple crops to a farmer
func (s *ProfileService) AddCropsToFarmer(farmerID int64, cropNames []string) error {
	for _, cropName := range cropNames {
		if err := s.addCropToFarmer(farmerID, cropName); err != nil {
			return fmt.Errorf("failed to add crop %s: %w", cropName, err)
		}
	}
	return nil
}

// GetFarmerCrops retrieves all crops for a farmer
func (s *ProfileService) GetFarmerCrops(farmerID int64) ([]models.Crop, error) {
	var farmerCrops []models.FarmerCropWithDetails
	_, err := s.client.From("farmer_crops").
		Select("*, crops(*)", "", false).
		Eq("farmer_id", fmt.Sprintf("%d", farmerID)).
		ExecuteTo(&farmerCrops)
	
	if err != nil {
		return nil, err
	}

	var crops []models.Crop
	for _, fc := range farmerCrops {
		if fc.Crop != nil {
			crops = append(crops, *fc.Crop)
		}
	}

	return crops, nil
}

// createExtensionOfficerRecord creates an extension officer record
func (s *ProfileService) createExtensionOfficerRecord(authUserID uuid.UUID, username string, signupData *SignupData) error {
	// Generate a new ID for the extension officer record
	officerID := generateExtensionOfficerID()

	// Set defaults if not provided
	phoneNumber := signupData.PhoneNumber
	var assignedLocationID *int64
	if signupData.AssignedLocationID != 0 {
		assignedLocationID = &signupData.AssignedLocationID
	}

	officer := models.ExtensionOfficer{
		ID:                 officerID,
		AuthUserID:         &authUserID, // Link to auth user
		Name:               username,
		PhoneNumber:        phoneNumber,
		AssignedLocationID: assignedLocationID,
	}

	var result []models.ExtensionOfficer
	_, err := s.client.From("extension_officers").Insert(officer, false, "", "", "").ExecuteTo(&result)
	return err
}

// generateFarmerID generates a new farmer ID (bigint)
func generateFarmerID() int64 {
	// For now, use timestamp-based ID
	// In production, you might want to use a proper sequence
	return time.Now().UnixNano() / 1000000 // Convert to milliseconds
}

// generateExtensionOfficerID generates a new extension officer ID (bigint)
func generateExtensionOfficerID() int64 {
	// For now, use timestamp-based ID
	// In production, you might want to use a proper sequence
	return time.Now().UnixNano() / 1000000 // Convert to milliseconds
}

// Error definitions
var (
	ErrSupabaseConfigMissing = &ServiceError{Code: "SUPABASE_CONFIG_MISSING", Message: "Supabase configuration missing"}
	ErrProfileCreationFailed = &ServiceError{Code: "PROFILE_CREATION_FAILED", Message: "Failed to create user profile"}
	ErrProfileNotFound       = &ServiceError{Code: "PROFILE_NOT_FOUND", Message: "User profile not found"}
	ErrRoleNotFound          = &ServiceError{Code: "ROLE_NOT_FOUND", Message: "Role not found"}
)

// ServiceError represents a service error
type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
