package utils

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
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
