package netcli

import (
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/adzsx/gwire/internal/utils"
)

func Subnet() string {
	cidr, _ := net.InterfaceAddrs()

	return fmt.Sprint(cidr[1])
}

func CalcAddr(cidr string) (string, string) {
	mask := utils.FilterChar(cidr, "/", false)
	maskN, err := strconv.Atoi(mask)

	utils.Err(err, true)

	var subnetmask []string

	for i := 0; i < int(math.Floor(float64(maskN)/8)); i++ {
		subnetmask = append(subnetmask, "255")
	}

	maskN -= len(subnetmask) * 8

	if maskN != 0 {
		converted := 0
		for i := 0; i <= maskN; i++ {
			converted += 256 >> maskN
		}
		log.Println(converted)
		subnetmask = append(subnetmask, strconv.Itoa(converted))
	}

	rest := 4 - len(subnetmask)
	for i := 0; i < rest; i++ {
		subnetmask = append(subnetmask, "0")
	}

	return utils.FilterChar(cidr, "/", true), strings.Join(subnetmask, ".")
}

func GetHosts(cidr string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ipList []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {

		if ip.String()[len(ip.String())-3:] != "255" && ip.String()[len(ip.String())-2:] != ".0" {
			ipList = append(ipList, ip.String())
		}
	}

	return ipList, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Info() (string, string, string, string) {
	ip, mask := CalcAddr(Subnet())
	list, err := GetHosts(Subnet())
	if err != nil {
		utils.Err(err, true)
	}
	nHosts := strconv.Itoa(len(list))

	return ip, mask, nHosts, ""
}
