package platform

import "os"

func GetIpInfoToken() string {
	return os.Getenv("IPINFO_TOKEN")
}
