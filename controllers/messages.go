package controllers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Investorharry19/voxa-golang-server/config"
	"github.com/Investorharry19/voxa-golang-server/database"
	"github.com/Investorharry19/voxa-golang-server/models"
	"github.com/Investorharry19/voxa-golang-server/utils"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddTextMessage godoc
// @Summary Add Text Message
// @Description Send a text message from a user
// @Tags MessageRoutes
// @Accept json
// @Produce json
// @Param messageData body models.TextMessageRequestDTO true "Text message data"
// @Success 201 {object} map[string]interface{} "Message sent"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "User does not exist"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /message/send/text-message [post]
func AddTextMessage(c *fiber.Ctx) error {

	requestData := models.TextMessageRequestDTO{}
	if err := c.BodyParser(&requestData); err != nil {
		return utils.ErrorResponse(c, 400, "INvalid Json")
	}
	if requestData.OwnerUsername == "" || requestData.MessageText == "" {
		return utils.ErrorResponse(c, 400, "ownerusername and message text are required")
	}
	requestData.ID = primitive.NewObjectID()
	requestData.Type = "text"
	requestData.CreatedAt = time.Now()

	userCollection := database.GetCollection("users")
	user := models.User{}
	if err := userCollection.FindOne(c.Context(), bson.M{"username": requestData.OwnerUsername}).Decode(&user); err != nil {
		return utils.ErrorResponse(c, 404, "User does not exist")
	}
	messageCollection := database.GetCollection("messages")
	res, err := messageCollection.InsertOne(c.Context(), requestData)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Internal server error")
	}

	return utils.SuccessResponse(c, 201, "message sent", res)
}

// SendAudioMessage godoc
// @Summary Send Audio Message
// @Description Upload and send an audio message for a user
// @Tags MessageRoutes
// @Accept mpfd
// @Produce json
// @Param ownerUsername formData string true "Owner username"
// @Param voice formData string true "Voice filter option"
// @Param file formData file true "Audio file"
// @Success 200 {object} map[string]string "Audio message uploaded successfully"
// @Failure 400 {object} map[string]string "Invalid request or file upload error"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /message/send/audio-message [post]
func SendAudioMessage(c *fiber.Ctx) error {

	// Get form data
	ownerUsername := c.FormValue("ownerUsername")
	voice := c.FormValue("voice")

	userCollection := database.GetCollection("users")
	// Find user
	user := models.User{}
	err := userCollection.FindOne(c.Context(), bson.M{"username": ownerUsername}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "No user with this username"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Database error"})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "No file uploaded"})
	}

	// Get filter settings
	filters := utils.GetFilterSetting(voice)
	if filters == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid voice option"})
	}

	// Create temp files
	tempDir := os.TempDir()
	timestamp := time.Now().UnixNano()
	tempInputPath := filepath.Join(tempDir, fmt.Sprintf("input_%d.mp3", timestamp))
	tempOutputPath := filepath.Join(tempDir, fmt.Sprintf("output_%d.mp3", timestamp))

	// Cleanup
	defer os.Remove(tempInputPath)
	defer os.Remove(tempOutputPath)

	// Save uploaded file to temp
	if err := c.SaveFile(file, tempInputPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save uploaded file"})
	}

	// Run FFmpeg
	args := []string{
		"-i", tempInputPath,
		"-af", filters,
		"-c:a", "libmp3lame",
		"-b:a", "128k",
		"-ac", "1",
		"-f", "mp3",
		"-y",
		tempOutputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("FFmpeg error: %s\n", string(output))
		return c.Status(500).JSON(fiber.Map{"message": "Error processing audio"})
	}

	// Upload to Cloudinary
	uploadResult, err := config.Cloud.Upload.Upload(c.Context(), tempOutputPath, uploader.UploadParams{
		ResourceType: "video",
	})
	if err != nil {
		fmt.Printf("Cloudinary upload error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "Error uploading to Cloudinary"})
	}

	// Save message to database
	newMessage := models.AudioMessageRequestDTO{
		ID:            primitive.NewObjectID(),
		OwnerUsername: ownerUsername,
		AudioUrl:      uploadResult.SecureURL,
		PublicId:      uploadResult.PublicID,
		CreatedAt:     time.Now(),
		Type:          "audio",
	}

	messageCollection := database.GetCollection("messages")
	_, err = messageCollection.InsertOne(c.Context(), newMessage)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save message"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save media reference"})
	}

	return c.Status(200).JSON(fiber.Map{
		"cloudinaryUrl": uploadResult.SecureURL,
	})
}

// get all messages
// GetAllMessages godoc
// @Summary Get All Messages
// @Description Retrieve all messages for the authenticated user
// @Tags MessageRoutes
// @Accept json
// @Produce json
// @Success 200 {array} models.Message "List of messages"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "No messages found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /message/get-messages [get]
func GetAllMessages(c *fiber.Ctx) error {
	id := fmt.Sprintf("%v", c.Locals("userId"))
	fmt.Println(id)
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, 400, "Bad Request")
	}
	userCollection := database.GetCollection("users")

	user := models.User{}
	err = userCollection.FindOne(c.Context(), bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return utils.ErrorResponse(c, 400, "Bad Request")
	}
	messageCollection := database.GetCollection("messages")
	messages := make([]models.Message, 0)

	fmt.Println(user.Username)
	cursor, err := messageCollection.Find(c.Context(), bson.M{"ownerusername": user.Username})
	if err != nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, 404, "No user with this username")
	}
	for cursor.Next(c.Context()) {
		message := models.Message{}
		if err := cursor.Decode(&message); err != nil {
			return utils.ErrorResponse(c, 500, "INternal server error")
		}
		messages = append(messages, message)
	}

	return utils.SuccessResponse(c, 201, "", messages)
}

// Mark as read
// MarkAsRead godoc
// @Summary Mark Message as Read
// @Description Mark a specific message as read by the authenticated user
// @Tags MessageRoutes
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} map[string]interface{} "Message updated"
// @Failure 400 {object} map[string]string "Invalid message ID or user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /message/mark-as-read/{id} [patch]
func MarkAsRead(c *fiber.Ctx) error {
	id := fmt.Sprintf("%v", c.Locals("userId"))
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrorResponse(c, 400, "")
	}
	messageId := c.Params("id")
	messageObjectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid user id")
	}
	messageCollection := database.GetCollection("messages")
	update := bson.M{
		"$set": bson.M{"isopened": true},
	}
	result := messageCollection.FindOneAndUpdate(
		c.Context(),
		bson.M{"_id": messageObjectId},
		update,
	)

	return utils.SuccessResponse(c, 201, "Message updated", result)

}

// Star message
// StarMessage godoc
// @Summary Star or Unstar a Message
// @Description Mark a specific message as starred or unstarred by the authenticated user
// @Tags MessageRoutes
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param state body models.MessageMarkAsRead true "Star state"
// @Success 200 {object} map[string]interface{} "Message updated"
// @Failure 400 {object} map[string]string "Invalid message ID or request body"
// @Failure 404 {object} map[string]string "Message not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /message/star-message/{id} [patch]
func StarMessage(c *fiber.Ctx) error {
	id := fmt.Sprintf("%v", c.Locals("userId"))
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrorResponse(c, 400, "")
	}
	messageId := c.Params("id")
	messageObjectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid user id")
	}
	state := models.MessageMarkAsRead{}
	if err := c.BodyParser(&state); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid JSON")
	}

	if state.State == nil {
		return utils.ErrorResponse(c, 400, "starred state required")
	}
	messageCollection := database.GetCollection("messages")
	update := bson.M{
		"$set": bson.M{"isstarred": state.State},
	}

	message := models.Message{}
	messageCollection.FindOne(c.Context(), bson.M{"_id": messageObjectId}).Decode(&message)

	if message.OwnerUsername == "" {
		return utils.ErrorResponse(c, 404, "This message does not exits")

	}
	result := messageCollection.FindOneAndUpdate(
		c.Context(),
		bson.M{"_id": messageObjectId},
		update,
	)

	return utils.SuccessResponse(c, 201, "Message updated", result)

}

// delete one message
// DeleteOneMessage godoc
// @Summary Delete a Message
// @Description Delete a specific message by the authenticated user
// @Tags MessageRoutes
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} map[string]interface{} "Message deleted"
// @Failure 400 {object} map[string]string "Invalid message ID or user ID"
// @Failure 404 {object} map[string]string "Message not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /message/delete-message{id} [delete]
func DeleteOneMessage(c *fiber.Ctx) error {
	id := fmt.Sprintf("%v", c.Locals("userId"))
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrorResponse(c, 400, "")
	}
	messageId := c.Params("id")
	messageObjectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid user id")
	}

	messageCollection := database.GetCollection("messages")

	result := messageCollection.FindOneAndDelete(
		c.Context(),
		bson.M{"_id": messageObjectId},
	)

	return utils.SuccessResponse(c, 201, "Message deleted", result)

}

// delete all message
// DeleteAllMessages godoc
// @Summary Delete All Messages
// @Description Delete all messages of the authenticated user
// @Tags MessageRoutes
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "All messages deleted"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /message/delete-all-messages [delete]
func DeleteAllMessages(c *fiber.Ctx) error {
	id := fmt.Sprintf("%v", c.Locals("userId"))
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrorResponse(c, 400, "")
	}

	messageCollection := database.GetCollection("messages")

	result, err := messageCollection.DeleteMany(
		c.Context(),
		bson.M{"ownerusername": objectId},
	)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Internal server error")
	}

	return utils.SuccessResponse(c, 201, "Message deleted", result)

}

// HandleVideoBuffer godoc
// @Summary Convert Audio to Video
// @Description Generate a video from an audio URL and a background image
// @Tags AudioProcessing
// @Accept json
// @Produce mp4
// @Param audioUrl query string true "URL of the audio file"
// @Success 200 {file} binary "Video file returned"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /convert [get]
func HandleVideoBuffer(c *fiber.Ctx) error {
	audioURL := c.Query("audioUrl")
	imagePath := "./bg.jpg"

	videoBuffer, err := utils.ConvertAudioToVideoBuffer(audioURL, imagePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "video/mp4")
	return c.Send(videoBuffer)
}

// ProcessAudioMessage godoc
// @Summary Process audio with voice filter
// @Description Applies voice filter to uploaded audio file
// @Tags AudioProcessing
// @Accept multipart/form-data
// @Produce audio/mpeg
// @Param file formData file true "Audio file"
// @Param voice formData string true "Voice filter (1-4)"
// @Success 200 {file} binary "Processed audio"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /process [post]
func ProcessAudioMessage(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "No file uploaded"})
	}

	// Get voice parameter
	voice := c.FormValue("voice")

	// Get filter settings
	filters := utils.GetFilterSetting(voice)
	if filters == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid voice option"})
	}

	// Create temp files
	tempDir := os.TempDir()
	timestamp := time.Now().UnixNano()
	tempInputPath := filepath.Join(tempDir, fmt.Sprintf("input_%d.mp3", timestamp))
	tempOutputPath := filepath.Join(tempDir, fmt.Sprintf("output_%d.mp3", timestamp))

	// Cleanup
	defer os.Remove(tempInputPath)
	defer os.Remove(tempOutputPath)

	// Save uploaded file to temp
	if err := c.SaveFile(file, tempInputPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save uploaded file"})
	}

	// Run FFmpeg using exec (more control over complex filters)
	args := []string{
		"-i", tempInputPath,
		"-af", filters,
		"-c:a", "libmp3lame",
		"-b:a", "128k",
		"-ac", "1",
		"-f", "mp3",
		"-y",
		tempOutputPath,
	}

	// Debug: print the full command
	fmt.Printf("FFmpeg command: ffmpeg %v\n", args)

	cmd := exec.Command("ffmpeg", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("FFmpeg error: %s\n", string(output))
		return c.Status(500).JSON(fiber.Map{"message": "Error processing audio"})
	}

	fmt.Println("FFmpeg processing finished.")

	// Read output file
	audioData, err := os.ReadFile(tempOutputPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to read processed audio"})
	}

	// Set headers and return
	c.Set("Content-Type", "audio/mpeg")
	c.Set("Content-Disposition", `attachment; filename="processed.mp3"`)

	return c.Send(audioData)
}

/*



 */
