package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"k8s.io/klog"
)

// UploadStatus structure
type UploadStatus struct {
	FileID      string
	FileName    string
	FileSize    int64
	TargetMD5   string
	TotalChunks int
	Received    map[int]bool
	ChunkPaths  map[int]string
	Status      string // pending, uploading, merging, completed, failed
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	FilePath    string
	FileMD5     string
	Verified    bool
}

// Global storage
var (
	uploads      = make(map[string]*UploadStatus)
	uploadsMutex = &sync.RWMutex{}
)

func GetOrCreateUpload(fileID, fileName, targetMD5 string, fileSize int64, totalChunks int) *UploadStatus {
	uploadsMutex.Lock()
	defer uploadsMutex.Unlock()

	// If fileID is provided, try to get existing upload
	if fileID != "" {
		if upload, exists := uploads[fileID]; exists {
			return upload
		}
	}

	// Create new upload session
	if fileID == "" {
		fileID = uuid.New().String()
	}

	upload := &UploadStatus{
		FileID:      fileID,
		FileName:    fileName,
		FileSize:    fileSize,
		TargetMD5:   targetMD5,
		TotalChunks: totalChunks,
		Received:    make(map[int]bool),
		ChunkPaths:  make(map[int]string),
		Status:      "uploading",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	uploads[fileID] = upload
	return upload
}

// Save chunk file
func SaveChunk(fileID string, chunkIndex int, data []byte) (string, error) {
	// Create temporary directory
	tempDir := fmt.Sprintf("./temp/%s", fileID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("Failed to create temporary directory: %v", err)
	}

	// Save chunk file
	chunkPath := fmt.Sprintf("%s/chunk_%d", tempDir, chunkIndex)
	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		return "", fmt.Errorf("Failed to save chunk file: %v", err)
	}

	return chunkPath, nil
}

// Update upload status
func UpdateUploadStatus(upload *UploadStatus, chunkIndex int, chunkPath string) {
	uploadsMutex.Lock()
	defer uploadsMutex.Unlock()

	upload.Received[chunkIndex] = true
	upload.ChunkPaths[chunkIndex] = chunkPath
	upload.UpdatedAt = time.Now()
	upload.Status = "uploading"
}

// Merge chunks
func MergeChunks(upload *UploadStatus) (string, string, error) {
	// Create final file path
	finalPath := fmt.Sprintf("./uploads/%s_%s", upload.FileID, upload.FileName)

	// 确保 uploads 目录存在
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return "", "", fmt.Errorf("Failed to create uploads directory: %v", err)
	}

	finalFile, err := os.Create(finalPath)
	if err != nil {
		return "", "", fmt.Errorf("Failed to create final file: %v", err)
	}
	defer finalFile.Close()
	// Create MD5 calculator
	hasher := md5.New()

	// Merge all chunks in order
	for i := 0; i < upload.TotalChunks; i++ {
		chunkPath, exists := upload.ChunkPaths[i]
		if !exists {
			return "", "", fmt.Errorf("Missing chunk %d", i)
		}

		// Read chunk data
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return "", "", fmt.Errorf("Failed to read chunk %d: %v", i, err)
		}

		// Write to final file
		if _, err := finalFile.Write(chunkData); err != nil {
			return "", "", fmt.Errorf("Failed to write chunk %d: %v", i, err)
		}

		// Update MD5
		hasher.Write(chunkData)
	}

	// Calculate final MD5
	fileMD5 := hex.EncodeToString(hasher.Sum(nil))

	return finalPath, fileMD5, nil
}

// Verify MD5
func VerifyMD5(actualMD5, targetMD5 string) bool {
	return actualMD5 == targetMD5
}

// Complete upload
func CompleteUpload(upload *UploadStatus, filePath, fileMD5 string, verified bool) {
	uploadsMutex.Lock()
	defer uploadsMutex.Unlock()

	completedAt := time.Now()
	upload.Status = "completed"
	upload.FilePath = filePath
	upload.FileMD5 = fileMD5
	upload.Verified = verified
	upload.CompletedAt = &completedAt
	upload.UpdatedAt = completedAt
}

// Update error status
func UpdateUploadError(upload *UploadStatus, errorMsg string) {
	uploadsMutex.Lock()
	defer uploadsMutex.Unlock()

	upload.Status = "failed"
	upload.UpdatedAt = time.Now()
}

// Clean up temporary files
func CleanupTempFiles(fileID string) {
	tempDir := fmt.Sprintf("./temp/%s", fileID)
	os.RemoveAll(tempDir)
}

// Create necessary directories
func CreateDirectories() {
	dirs := []string{
		"./temp",
		"./uploads",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			klog.Errorf("Create %s error: %v", dir, err)
		}
	}
}
