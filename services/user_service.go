package services

import (
	"context"
	"database/sql"
	"encoding/json"

	db "github.com/kawai-network/veridium/internal/database/generated"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

// User operations

func (s *UserService) GetUser(ctx context.Context, userID string) (db.User, error) {
	return s.queries.GetUser(ctx, userID)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return s.queries.GetUserByEmail(ctx, stringToNullString(email))
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return s.queries.GetUserByUsername(ctx, stringToNullString(username))
}

type CreateUserParams struct {
	ID              string
	Username        *string
	Email           *string
	Avatar          *string
	Phone           *string
	FirstName       *string
	LastName        *string
	IsOnboarded     bool
	ClerkCreatedAt  *int64
	EmailVerifiedAt *int64
	Preference      map[string]interface{}
}

func (s *UserService) CreateUser(ctx context.Context, params CreateUserParams) (db.User, error) {
	now := currentTimestampMs()

	return s.queries.CreateUser(ctx, db.CreateUserParams{
		ID:              params.ID,
		Username:        toNullString(params.Username),
		Email:           toNullString(params.Email),
		Avatar:          toNullString(params.Avatar),
		Phone:           toNullString(params.Phone),
		FirstName:       toNullString(params.FirstName),
		LastName:        toNullString(params.LastName),
		IsOnboarded:     boolToInt(params.IsOnboarded),
		ClerkCreatedAt:  toNullInt64(params.ClerkCreatedAt),
		EmailVerifiedAt: toNullInt64(params.EmailVerifiedAt),
		Preference:      toNullJSON(params.Preference),
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

type UpdateUserParams struct {
	ID         string
	Username   *string
	Email      *string
	Avatar     *string
	Phone      *string
	FirstName  *string
	LastName   *string
	Preference map[string]interface{}
}

func (s *UserService) UpdateUser(ctx context.Context, params UpdateUserParams) (db.User, error) {
	preferenceJSON := "{}"
	if params.Preference != nil {
		if data, err := json.Marshal(params.Preference); err == nil {
			preferenceJSON = string(data)
		}
	}

	return s.queries.UpdateUser(ctx, db.UpdateUserParams{
		Username:   toNullString(params.Username),
		Email:      toNullString(params.Email),
		Avatar:     toNullString(params.Avatar),
		Phone:      toNullString(params.Phone),
		FirstName:  toNullString(params.FirstName),
		LastName:   toNullString(params.LastName),
		Preference: sql.NullString{String: preferenceJSON, Valid: true},
		UpdatedAt:  currentTimestampMs(),
		ID:         params.ID,
	})
}

func (s *UserService) UpdateUserOnboarding(ctx context.Context, userID string, isOnboarded bool) error {
	return s.queries.UpdateUserOnboarding(ctx, db.UpdateUserOnboardingParams{
		IsOnboarded: boolToInt(isOnboarded),
		UpdatedAt:   currentTimestampMs(),
		ID:          userID,
	})
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.queries.DeleteUser(ctx, userID)
}

// User Settings operations

func (s *UserService) GetUserSettings(ctx context.Context, userID string) (db.UserSetting, error) {
	return s.queries.GetUserSettings(ctx, userID)
}

type UpsertUserSettingsParams struct {
	UserID        string
	TTS           map[string]interface{}
	Hotkey        map[string]interface{}
	KeyVaults     *string
	General       map[string]interface{}
	LanguageModel map[string]interface{}
	SystemAgent   map[string]interface{}
	DefaultAgent  map[string]interface{}
	Tool          map[string]interface{}
	Image         map[string]interface{}
}

func (s *UserService) UpsertUserSettings(ctx context.Context, params UpsertUserSettingsParams) (db.UserSetting, error) {
	return s.queries.UpsertUserSettings(ctx, db.UpsertUserSettingsParams{
		ID:            params.UserID,
		Tts:           toNullJSON(params.TTS),
		Hotkey:        toNullJSON(params.Hotkey),
		KeyVaults:     toNullString(params.KeyVaults),
		General:       toNullJSON(params.General),
		LanguageModel: toNullJSON(params.LanguageModel),
		SystemAgent:   toNullJSON(params.SystemAgent),
		DefaultAgent:  toNullJSON(params.DefaultAgent),
		Tool:          toNullJSON(params.Tool),
		Image:         toNullJSON(params.Image),
	})
}

// User Plugins operations

func (s *UserService) ListUserPlugins(ctx context.Context, userID string) ([]db.UserInstalledPlugin, error) {
	return s.queries.ListUserPlugins(ctx, userID)
}

func (s *UserService) GetUserPlugin(ctx context.Context, userID, identifier string) (db.UserInstalledPlugin, error) {
	return s.queries.GetUserPlugin(ctx, db.GetUserPluginParams{
		UserID:     userID,
		Identifier: identifier,
	})
}

type InstallUserPluginParams struct {
	UserID       string
	Identifier   string
	Type         string
	Manifest     map[string]interface{}
	Settings     map[string]interface{}
	CustomParams map[string]interface{}
}

func (s *UserService) InstallUserPlugin(ctx context.Context, params InstallUserPluginParams) (db.UserInstalledPlugin, error) {
	now := currentTimestampMs()

	return s.queries.InstallUserPlugin(ctx, db.InstallUserPluginParams{
		UserID:       params.UserID,
		Identifier:   params.Identifier,
		Type:         params.Type,
		Manifest:     toNullJSON(params.Manifest),
		Settings:     toNullJSON(params.Settings),
		CustomParams: toNullJSON(params.CustomParams),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (s *UserService) UninstallUserPlugin(ctx context.Context, userID, identifier string) error {
	return s.queries.UninstallUserPlugin(ctx, db.UninstallUserPluginParams{
		UserID:     userID,
		Identifier: identifier,
	})
}

// Helper functions
