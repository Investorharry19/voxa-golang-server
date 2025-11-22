package config

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var Cloud *cloudinary.Cloudinary

func InitCloudinary() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatal("Cloudinary init error:", err)
	}

	Cloud = cld

	// Debug
	fmt.Println("Cloudinary initialized:")
	fmt.Println("Cloud name:", os.Getenv("CLOUDINARY_CLOUD_NAME"))
}
