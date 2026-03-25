package v1

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
)

func UploadProcessHandler(c *gin.Context) {
	content, err := os.ReadFile("/otapart/update-flag")
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Failed to read the file: %v\n", err), c)
		return
	}

	// 2. 去除空白字符并转换为字符串
	str := strings.TrimSpace(string(content))

	// 3. 转换为整数
	flag, err := strconv.Atoi(str)
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Conversion failed: %v (content: %s)\n", err, str), c)
		return
	}

	responce.OkWithData(flag, c)
}
