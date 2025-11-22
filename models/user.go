package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID                      primitive.ObjectID `bson:"_id"`
	Username                string             `bson:"username"`
	Password                string             `bson:"password_hash"`
	PushNotificationEnabled bool               `bson:"pushNotificationEnabled"`
	PushToken               []string           `bson:"pushToken"`
}
type UserResponse struct {
	ID                      string   `json:"id"`
	Username                string   `json:"username"`
	Token                   string   `json:"token"`
	PushNotificationEnabled bool     `bson:"pushNotificationEnabled"`
	PushToken               []string `bson:"pushToken"`
}

func UserToUserResponse(user *User, token string) UserResponse {
	return UserResponse{
		ID:                      user.ID.Hex(),
		Username:                user.Username,
		Token:                   token,
		PushNotificationEnabled: user.PushNotificationEnabled,
		PushToken:               user.PushToken,
	}
}
