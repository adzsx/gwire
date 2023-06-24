package utils

import "strings"

func FilterIp(ip string) string {
	var final string

	for _, element := range ip {
		if string(element) != ":" {
			final += string(element)
		} else {
			return final
		}
	}

	return final
}

func FilterPort(ip string) string {
	var final []string

	for index, element := range ip {
		if string(element) == ":" {
			final = append(final, ip[index+1:])
		}
	}

	return strings.Join(final, "")
}

func Remove(slice [][]string, index int) [][]string {
	var final [][]string

	for rn, element := range slice {
		if rn != index {
			final = append(final, element)
		}
	}

	return final
}
