package models

type User struct {
	UserID       string `bson:"_id"           json:"-"` // to Disable UserID from Json Response used "-"
	Name         string `bson:"name"          json:"name"`
	Email        string `bson:"email"         json:"email"`
	Password     string `bson:"password"      json:"password"`
	RefreshToken string `bson:"refreshtoken"  json:"-"`
}

// UserCreateRequest represents the request body for creating a new user
type UserBind struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
