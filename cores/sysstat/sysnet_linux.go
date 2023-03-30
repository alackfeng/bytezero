package sysstat

import (
	"net"
	"strings"
	"fmt"
	"strconv"
	"regexp"
	"os/exec"

)

// GetNetworkMaxSpeed - 获取网卡最大速率.
func GetNetworkMaxSpeed(networkInterface net.Interface) (uint32, error) {
	cmd := exec.Command("ethtool", networkInterface.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	reg := regexp.MustCompile("Speed:(.*)Mb/s")
	res := reg.FindStringSubmatch(string(output))
	fmt.Println("GetNetworkMaxSpeed:", networkInterface.Name, res)
	if len(res) != 2 {
		return 0, nil
	}
	speed, _ := strconv.ParseInt(strings.Trim(res[1], " "), 10, 64)
	return uint32(speed), nil
}
