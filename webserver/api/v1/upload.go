package v1

import (
	"fmt"
	"io"
	"strconv"

	"github.com/edgehook/ithings/common/utils"
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func UploadChunkHandler(c *gin.Context) {
	// 1. Parse form parameters
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		responce.FailWithMessage("Unable to get chunk file", c)
		return
	}
	defer file.Close()

	// Read file data
	chunkData, err := io.ReadAll(file)
	if err != nil {
		responce.FailWithMessage("Failed to read chunk data", c)
		return
	}

	// Get form parameters
	fileName := c.PostForm("name")
	targetMD5 := c.PostForm("md5")
	fileSize, _ := strconv.ParseInt(c.PostForm("size"), 10, 64)
	totalChunks, _ := strconv.Atoi(c.PostForm("chunks"))
	chunkIndex, _ := strconv.Atoi(c.PostForm("chunk"))

	// Parameter validation
	if fileName == "" || targetMD5 == "" || fileSize <= 0 || totalChunks <= 0 {
		responce.FailWithMessage("Missing required parameters", c)
		return
	}

	// 2. Process or create upload session
	upload := utils.GetOrCreateUpload(targetMD5, fileName, targetMD5, fileSize, totalChunks)

	// Check if chunk has already been uploaded
	if upload.Received[chunkIndex] {
		responce.Ok(c)
		return
	}

	// 3. Save chunk
	chunkPath, err := utils.SaveChunk(upload.FileID, chunkIndex, chunkData)
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Failed to save chunk err %v", err.Error()), c)
		return
	}

	// 4. Update upload status
	utils.UpdateUploadStatus(upload, chunkIndex, chunkPath)

	// 5. Check if all chunks have been uploaded
	if chunkIndex == upload.TotalChunks-1 {
		klog.Infof("All chunks received for file %s, starting merge", upload.FileName)
		// Auto merge chunks
		mergedPath, fileMD5, err := utils.MergeChunks(upload)
		if err != nil {
			responce.FailWithMessage(fmt.Sprintf("Failed to merge chunks err: %v", err.Error()), c)
			return
		}

		// Verify MD5
		verified := utils.VerifyMD5(fileMD5, upload.TargetMD5)
		if !verified {
			responce.FailWithMessage("MD5 verification failed after merging chunks", c)
			return
		}

		// Update final status
		utils.CompleteUpload(upload, mergedPath, fileMD5, verified)

		// Clean up temporary files
		utils.CleanupTempFiles(upload.FileID)
		responce.Ok(c)
		return
	}

	responce.Ok(c)
}

func RebootHandler(c *gin.Context) {
	go func() {
		if err := utils.SysReboot(); err != nil {
			klog.Errorf("system reboot failed: %v", err)
		}
	}()
	responce.Ok(c)
}
