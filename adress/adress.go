package adress

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func GetConfig() (error, net.Addr) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return err, nil
	}

	for _, interf := range interfaces {
		if isWirelessInterface(interf.Name) {
			addrs, err := interf.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						fmt.Printf("  IPv4-adress: %s\n", ipNet.IP.String())
						return nil, ipNet
					}
				}
			}
		}
	}
	return errors.New("Can't find wireless interface"), nil
}

func isWirelessInterface(name string) bool {
	return strings.HasPrefix(name, "wlan") ||
		strings.HasPrefix(name, "wl") ||
		strings.Contains(name, "Wi-Fi") ||
		name == "Беспроводная сеть"
}
