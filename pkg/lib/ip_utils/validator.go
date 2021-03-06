package ip_utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/syunkitada/goapp/pkg/lib/error_utils"
	"github.com/syunkitada/goapp/pkg/lib/ip_utils/ip_utils_model"
)

func ParseNetwork(subnet string, gateway string, startIp string, endIp string) (*ip_utils_model.Network, error) {
	_, parsedSubnet, subnetErr := net.ParseCIDR(subnet)
	if subnetErr != nil {
		return nil, error_utils.NewInvalidDataError("subnet", subnet, "FailedParse")
	}

	parsedGateway := net.ParseIP(gateway)
	if parsedGateway == nil {
		return nil, error_utils.NewInvalidDataError("gateway", gateway, "FailedParse")
	}

	parsedStartIp := net.ParseIP(startIp)
	if parsedStartIp == nil {
		return nil, error_utils.NewInvalidDataError("startIp", startIp, "FailedParse")
	}

	parsedEndIp := net.ParseIP(endIp)
	if parsedEndIp == nil {
		return nil, error_utils.NewInvalidDataError("endIp", endIp, "FailedParse")
	}

	if !parsedSubnet.Contains(parsedStartIp) {
		return nil, error_utils.NewInvalidDataError("startIp", startIp, "startIp should be countained in subnet")
	}

	if !parsedSubnet.Contains(parsedEndIp) {
		return nil, error_utils.NewInvalidDataError("endIp", endIp, "endIp should be countained in subnet")
	}

	if !parsedSubnet.Contains(parsedGateway) {
		return nil, error_utils.NewInvalidDataError("gateway", gateway, "gateway should be countained in subnet")
	}
	if CompareIp(parsedStartIp, parsedEndIp) != -1 {
		return nil, error_utils.NewInvalidDataError("endIp", endIp, "endIp should be bigger than startIp")
	}
	return &ip_utils_model.Network{
		Subnet:  parsedSubnet,
		Gateway: parsedGateway,
		StartIp: parsedStartIp,
		EndIp:   parsedEndIp,
	}, nil
}

func IncrementIp(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		//only add to the next byte if we overflowed
		if ip[i] != 0 {
			break
		}
	}
}

// CompareIp compare ip1, ip2
// ip1が大きければ1, ip2が大きければ-1, 同じなら0を返す
func CompareIp(ip1 net.IP, ip2 net.IP) int {
	len := len(ip1)
	for i := 0; i < len; i++ {
		if ip1[i] > ip2[i] {
			return 1
		} else if ip1[i] < ip2[i] {
			return -1
		}
	}

	return 0
}

func GenerateRandomMac() (string, error) {
	buf := make([]byte, 5)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	oui := []uint8{0x02}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", oui[0], buf[0], buf[1], buf[2], buf[3], buf[4]), nil
}

func GenerateUniqueRandomMac(macMap map[string]bool, limit int) (string, error) {
	var mac string
	var err error
	for i := 0; i < limit; i++ {
		if mac, err = GenerateRandomMac(); err != nil {
			return "", err
		}

		if _, ok := macMap[mac]; !ok {
			return mac, err
		}
	}

	return "", fmt.Errorf("Failed Generate Mac: Exceeded Limit %d", limit)
}

func AddIntToIp(ip net.IP, value int) net.IP {
	intIp := ip2int(ip)
	intIp += uint32(value)
	newIp := int2ip(intIp)
	return newIp
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
