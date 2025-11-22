package controllers

import (
	"fmt"

	"github.com/Investorharry19/voxa-golang-server/database"
	"github.com/Investorharry19/voxa-golang-server/models"
	"github.com/Investorharry19/voxa-golang-server/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"status": "Okay"})
}

// RegisterUser godoc
// @Summary Register User
// @Description Register a new user
// @Tags Account
// @Accept json
// @Produce json
// @Param registerData body models.UserRequestDTO true "User registration data"
// @Success 201 {object} map[string]string "User created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 461 {object} map[string]string "Username already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /account/register [post]
func RegisterUser(c *fiber.Ctx) error {
	println("Register")
	registerData := models.UserRequestDTO{}
	if err := c.BodyParser(&registerData); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid Json")

	}

	if registerData.Password == "" || registerData.Username == "" {
		return utils.ErrorResponse(c, 400, "Username and password are required")

	}
	hashed, err := utils.HashPasswordSecure(registerData.Password)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Internal server Error")
	}
	userCollection := database.GetCollection("users")
	saveUser := models.User{
		Username:  registerData.Username,
		Password:  hashed,
		PushToken: make([]string, 0),
		ID:        primitive.NewObjectID(),
	}
	res, err := userCollection.InsertOne(c.Context(), saveUser)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range we.WriteErrors {
				if writeErr.Code == 11000 {
					return utils.ErrorResponse(c, 461, "Username already exists")
				}
			}
		}
		return utils.ErrorResponse(c, 500, err.Error())
	}
	return utils.SuccessResponse(c, 201, "User Created", res)

}

// LoginUser godoc
// @Summary Login User
// @Description Authenticate a user and return a JWT token
// @Tags Account
// @Accept json
// @Produce json
// @Param loginData body models.UserRequestDTO true "User login data"
// @Success 201 {object} map[string]interface{} "Logged In"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Invalid Credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /account/login [post]
func LoginUser(c *fiber.Ctx) error {
	loginData := models.UserRequestDTO{}
	if err := c.BodyParser(&loginData); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid Json")
	}
	if loginData.Username == "" || loginData.Password == "" {
		return utils.ErrorResponse(c, 400, "Username and password are required")
	}

	userCollection := database.GetCollection("users")
	userdata := models.User{}
	if err := userCollection.FindOne(c.Context(), bson.M{"username": loginData.Username}).Decode(&userdata); err != nil {
		return utils.ErrorResponse(c, 404, "Invalid Credentials")
	}

	passwordMatch, err := utils.VerifyPasswordSecure(userdata.Password, loginData.Password)
	fmt.Println(loginData.Password, userdata.Password, passwordMatch)
	if err != nil {
		println(err)
		return utils.ErrorResponse(c, 500, "Internal server error")
	}
	println(76)
	if !passwordMatch {
		return utils.ErrorResponse(c, 404, "Invalid Credentials")
	}

	token, _ := utils.GenerateJWT("", string(userdata.ID.Hex()), 10)
	loginresponse := models.UserToUserResponse(&userdata, token)

	return utils.SuccessResponse(c, 201, "Logged In", loginresponse)
}

// GetCurrentUser godoc
// @Summary Get Current User
// @Description Retrieve the authenticated user's information
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]interface{} "User authenticated"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /account/current-user [get]
func GetCurrentUser(c *fiber.Ctx) error {
	id := fmt.Sprintf("%v", c.Locals("userId"))
	fmt.Println(id)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, 401, "Bad request")
	}
	userResponse := models.User{}
	userCollection := database.GetCollection("users")

	if err := userCollection.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&userResponse); err != nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, 400, "Bad request")
	}

	token, _ := utils.GenerateJWT("", userResponse.ID.Hex(), 60*24*7)
	return utils.SuccessResponse(
		c, 201, "User authenticated", models.UserToUserResponse(&userResponse, token))
}

/*




 */
