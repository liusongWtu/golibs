package ip

import (
	"fmt"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String(), nil
}
