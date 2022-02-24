package ip

import (
	"crypto/sha256"
	"encoding/hex"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceV4() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Description: "v4 IP address",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Read: resourceBoardRead,
	}
}

func resourceBoardRead(data *schema.ResourceData, meta interface{}) error {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if err := data.Set("ip", ipnet.IP.String()); err != nil {
					return err
				}
				ipHash := sha256.Sum256([]byte(ipnet.IP.String()))
				data.SetId(hex.EncodeToString(ipHash[:]))
			}
		}
	}
	return nil
}
