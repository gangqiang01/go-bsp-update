package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

func NewUUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

// return ms
func GetNowTimeStamp() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}
func Md5V(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

/*
* GetLocalMACs
* get the local host's macaddress.
 */
func GetLocalMACs() []string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil
	}

	macAddrs := make([]string, 0)
	for _, ni := range netInterfaces {
		if !strings.HasPrefix(ni.Name, "e") &&
			!strings.HasPrefix(ni.Name, "w") &&
			!strings.HasPrefix(ni.Name, "p") {
			continue
		}
		macAddr := ni.HardwareAddr.String()
		if len(macAddr) > 6 {
			macAddrs = append(macAddrs, strings.ToUpper(macAddr))
		}
	}

	return macAddrs
}
func GetOsType() string {
	return runtime.GOOS
}

func Execute(command string) (string, error) {
	if strings.Contains(GetOsType(), "windows") {
		cmd := exec.Command("cmd", "/C", command)
		cmd.Env = os.Environ()
		output, err := cmd.CombinedOutput()
		if err != nil {
			klog.Errorf("output: %s, err: %v", string(output), err)
			return string(output), err
		}

		return string(output), nil
	}

	// default for unix.
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Env = os.Environ()
	// cmd.Dir = "/usr/local"
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Setpgid: true,
	// }

	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("output: %s, err: %v", string(output), err)
		return string(output), err
	}

	return string(output), nil
}

func Execute1(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("output: %s, err: %v", string(output), err)
		return string(output), err
	}

	return string(output), nil
}

func SysReboot() error {
	if strings.Contains(GetOsType(), "windows") {
		_, err := Execute1("cmd", "/C", "shutdown", "/r", "/t", "0")
		return err
	}

	//default for linux.
	_, err := Execute("reboot")
	return err
}

func AuthUser(username, password string) error {
	// 1. Select command based on system architecture
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("GetRootPath filepath Abs with err: %v", err)
	}

	command := filepath.Join(dir, "linuxPackage", "authPassword")
	if runtime.GOARCH == "arm64" {
		command = filepath.Join(dir, "linuxPackage", "authPassword-ARM64")
	}

	// 2. Check if command file exists
	if _, err := os.Stat(command); os.IsNotExist(err) {
		return fmt.Errorf("Command file not found: %s (architecture: %s)", command, runtime.GOARCH)
	}

	// 3. Set execution permission
	if err := os.Chmod(command, 0755); err != nil {
		// Continue execution even if setting permission fails
		return fmt.Errorf("Warning: cannot set execute permission: %v\n", err)
	}

	// 4. Set 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5. Build and execute command
	cmd := exec.CommandContext(ctx, command, "-username", username, "-password", password)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	// 6. Handle timeout
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("Command execution timeout")
	}

	// 7. Handle execution error
	if err != nil {
		// Try to extract information from error output
		if outputStr == "" {
			outputStr = "no output"
		}
		return fmt.Errorf("Command execution failed: %v", err)
	}

	// 8. Parse output
	// Try to convert directly to integer
	code, err := strconv.Atoi(outputStr)
	if err != nil {
		return fmt.Errorf("Cannot parse command output: %s", outputStr)
	}

	// 9. Determine result based on output code
	switch code {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("User verification failed")
	default:
		return fmt.Errorf("Unknown error")
	}
}
func CopyFile(source, destination string) error {
	// Check if the source file exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		// Original file does not exist, use default file
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return fmt.Errorf("GetRootPath filepath Abs with err: %v", err)
		}

		defaultSource := filepath.Join(dir, "linuxPackage", "advupdate.txt.example")
		fmt.Printf("Original file %s does not exist, trying default file: %s\n", source, defaultSource)

		// Check if the default file exists
		if _, err := os.Stat(defaultSource); os.IsNotExist(err) {
			return fmt.Errorf("original file %s does not exist, and default file %s also does not exist", source, defaultSource)
		}

		// Use the default file
		source = defaultSource
	} else if err != nil {
		// Other errors (e.g., permission issues)
		return fmt.Errorf("failed to check source file %s: %v", source, err)
	}

	// Open the source file
	srcFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", source, err)
	}
	defer srcFile.Close()

	// Get source file info for later verification
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %v", err)
	}

	// Create the destination file
	dstFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", destination, err)
	}
	defer dstFile.Close()

	// Copy the content
	written, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	// Verify the copy size
	if written != srcInfo.Size() {
		return fmt.Errorf("copy size mismatch: source file %d bytes, actually copied %d bytes", srcInfo.Size(), written)
	}

	fmt.Printf("File copied successfully: %s -> %s (%d bytes)\n", source, destination, written)
	return nil
}
