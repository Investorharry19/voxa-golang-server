package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ProbeData struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func ConvertAudioToVideoBuffer(audioURL, imagePath string) ([]byte, error) {
	var tempAudioPath, tempOutputPath string

	// Cleanup function
	cleanup := func() {
		if tempAudioPath != "" {
			os.Remove(tempAudioPath)
		}
		if tempOutputPath != "" {
			os.Remove(tempOutputPath)
		}
	}
	defer cleanup()

	// Check if image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", imagePath)
	}

	fmt.Println("Downloading audio...")
	downloadStart := time.Now()

	// Download audio
	resp, err := http.Get(audioURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	// Create temp files
	tempDir := os.TempDir()
	timestamp := time.Now().UnixNano()
	tempAudioPath = filepath.Join(tempDir, fmt.Sprintf("temp_audio_%d.mp3", timestamp))
	tempOutputPath = filepath.Join(tempDir, fmt.Sprintf("temp_output_%d.mp4", timestamp))

	// Write audio to temp file
	if err := os.WriteFile(tempAudioPath, audioData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp audio: %w", err)
	}
	fmt.Println("Audio downloaded and saved temporarily")
	fmt.Printf("⏱️ Download took: %v\n", time.Since(downloadStart))

	// Get audio duration using ffprobe
	probeResult, err := ffmpeg.Probe(tempAudioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to probe audio: %w", err)
	}

	var probeData ProbeData
	if err := json.Unmarshal([]byte(probeResult), &probeData); err != nil {
		return nil, fmt.Errorf("failed to parse probe data: %w", err)
	}

	audioDuration := probeData.Format.Duration
	fmt.Printf("Audio duration: %s seconds\n", audioDuration)

	// Run FFmpeg conversion
	fmt.Println("Starting FFmpeg conversion...")
	conversionStart := time.Now()

	imageInput := ffmpeg.Input(imagePath, ffmpeg.KwArgs{
		"loop":      1,
		"framerate": 5,
		"t":         audioDuration,
	})

	audioInput := ffmpeg.Input(tempAudioPath)

	err = ffmpeg.Output(
		[]*ffmpeg.Stream{imageInput, audioInput},
		tempOutputPath,
		ffmpeg.KwArgs{
			"c:v":     "libx264",
			"c:a":     "aac",
			"pix_fmt": "yuv420p",
			"vf":      "scale=720:-2",
			"r":       5,
			"preset":  "ultrafast",
			"tune":    "stillimage",
			"crf":     30,
			"b:a":     "96k",
			"ar":      22050,
			"ac":      2,
		},
	).OverWriteOutput().Run()

	if err != nil {
		return nil, fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	fmt.Println("Video conversion completed!")
	fmt.Printf("⏱️ Conversion took: %v\n", time.Since(conversionStart))

	// Read the output file
	videoBuffer, err := os.ReadFile(tempOutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output video: %w", err)
	}

	fmt.Printf("Video file size: %d bytes\n", len(videoBuffer))

	return videoBuffer, nil
}
