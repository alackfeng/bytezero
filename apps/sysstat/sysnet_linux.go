package sysstat

import (
	"net"
	"syscall"
)

// GetNetworkMaxSpeed - 获取网卡最大速率.
func GetNetworkMaxSpeed(networkInterface net.Interface) (uint32, error) {
	cmd := exec.Command("ethtool", networkInterface.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	reg := regexp.MustCompile("Speed:(.*)Mb/s")
	res := reg.FindStringSubmatch(test)
	if len(res) != 2 {
		return 0, fmt.Errorf("regexp speed error")
	}
	speed, _ := strconv.ParseInt(strings.Trim(res[1], " "), 10, 64)
	return speed, nil
}
