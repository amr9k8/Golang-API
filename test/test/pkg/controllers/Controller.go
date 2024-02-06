package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"test/pkg/db/models"
	"test/pkg/db/repository/organizations"
	"test/pkg/db/repository/users"
	"test/pkg/utils"
)

var redisClient *redis.Client

type TokenInfo struct {
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userID"`
	Blacklisted  bool   `json:"blacklisted"`
}

func RedisInit() {
	// Initialize the Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Your Redis server address
		Password: "",           // No password
		DB:       0,            // Default DB
	})
}
func addTokenInfoRedis(tokenInfo TokenInfo) error {

	data, err := json.Marshal(tokenInfo)
	if err != nil {
		return err
	}

	err = redisClient.Set(context.Background(), tokenInfo.RefreshToken, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func verifyRefreshTokenRedis(refreshToken string) (bool, error) {
	val, err := redisClient.Get(context.Background(), refreshToken).Result()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist, refresh token is not valid
			return false, nil
		}
		return false, err
	}

	var tokenInfo TokenInfo
	err = json.Unmarshal([]byte(val), &tokenInfo)
	if err != nil {
		return false, err
	}

	return !tokenInfo.Blacklisted, nil
}

// Users Controllers
func GetUsersController() ([]models.User, error) {
	return users.GetAllUsers()
}
func GetUserController(UserID string) (*models.User, error) {
	return users.GetUserById(UserID)
}
func DeleteUserController(UserID string) error {
	return users.DeleteUserById(UserID)
}
func AddUserController(newuser models.UserBind) error {
	return users.InsertUser(newuser)
}
func LoginUserController(email string, password string) (map[string]string, error) {
	user, err := users.LoginUser(email, password)
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"userid": user.UserID,
		"name":   user.Name,
		"email":  user.Email,
	}
	AccessToken, err := utils.JWTGenerateToken(claims)
	if err != nil {
		return nil, err
	}
	RefreshToken, err := utils.GenerateRefreshToken(user.UserID)
	myMap := make(map[string]string)
	myMap["access_token"] = AccessToken
	myMap["refresh_token"] = RefreshToken
	// TODO  UNCOMMENT TO USE WITH REDIS CONTAINER
	//// Set the refresh token in Redis with the user ID as the key
	//tokenInfo := TokenInfo{
	//	RefreshToken: RefreshToken,
	//	UserID:       user.UserID,
	//	Blacklisted:  false,
	//}
	//err = addTokenInfoRedis(tokenInfo)
	//if err != nil {
	//	fmt.Println("Error adding token:", err)
	//}

	return myMap, err
}
func RefreshTokenController(refershToken string) (map[string]string, error) {
	// Example: Verifying a token
	// TODO  UNCOMMENT TO USE WITH REDIS CONTAINER
	//isValidRedis, _ := verifyRefreshTokenRedis(refershToken)
	//if isValidRedis {

	Claims, IsValid, _ := utils.JWTDecodeToken(refershToken)
	if IsValid {
		fmt.Println("Claims:")
		for key, value := range Claims {
			fmt.Printf("%s: %v\n", key, value)
		}
		UserID := Claims["userid"].(string)
		user, err := users.GetUserById(UserID)
		if err != nil {
			return nil, err
		}
		claims := jwt.MapClaims{
			"userid": user.UserID,
			"name":   user.Name,
			"email":  user.Email,
		}
		AccessToken, err := utils.JWTGenerateToken(claims)
		if err != nil {
			return nil, err
		}
		myMap := make(map[string]string)
		myMap["Message"] = "New AccessToken Generated Successfully!"
		myMap["access_token"] = AccessToken
		myMap["refresh_token"] = refershToken

		return myMap, nil
	}
	//}
	return nil, errors.New("Refresh Token Is Invalid")

}
func UpdateUserController(UserID string, newuser models.User) error {
	return users.UpdateUser(UserID, newuser)
}

// Organizations Controllers
func GetOrganizationsController() ([]models.Organization, error) {
	return organizations.GetAllOrganizations()
}
func GetOrganizationController(OrgID string) (*models.Organization, error) {
	return organizations.GetOrganizationById(OrgID)
}
func AddOrganizationController(NewOrg models.OrganizationBind, UserID string, AccessLevel string) (string, error) {
	return organizations.InsertOrganization(NewOrg, UserID, AccessLevel)

}
func InviteOrganizationController(OrgID string, email string) error {
	return organizations.InviteMemberToOrganization(OrgID, email)
}
func UpdateOrganizationController(OrgID string, updatedOrg models.OrganizationOnly, UserID string) error {

	// 1) get organization
	org, _ := organizations.GetOrganizationById(OrgID)
	// 2) get members
	for i := range org.Members {
		AccessLevel := org.Members[i].AccessLevel
		MemberID := org.Members[i].UserID
		if AccessLevel == "fullaccess" && MemberID == UserID {
			// 3) find if user have full access
			return organizations.UpdateOrganization(OrgID, updatedOrg)
		}

	}
	return errors.New("Insufficient Privilege")

}
func DeleteOrganizationController(OrgID string, UserID string) error {
	// 1) get organization
	org, _ := organizations.GetOrganizationById(OrgID)
	// 2) get members
	for i := range org.Members {
		AccessLevel := org.Members[i].AccessLevel
		MemberID := org.Members[i].UserID
		if AccessLevel == "fullaccess" && MemberID == UserID {
			// 3) find if user have full access
			return organizations.DeleteOrganizationById(OrgID)
		}

	}
	return errors.New("Insufficient Privilege")

}
func TestController() error {
	AccessLevel := "fullaccess"
	UserID := "65b8a2e738f20fe0f0476148"
	return organizations.InsertMemberIntoOrganization("65b9029a4c5f266f78039a9f", UserID, AccessLevel)
}
