package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const IPIFY_API_URL = "https://api64.ipify.org"

func dataSourceIp() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Data source in the Terraform provider ip to retrieve ipv4 & ipv6.",

		Schema: map[string]*schema.Schema{
			"nw_interface": {
				// This description is used by the documentation generator and the language server.
				Description: "The network interface to fetch the ip if the machine has multiple NICs.",
				Type:        schema.TypeString,
				Default:     "",
				Optional:    true,
			},
			"public_ip": {
				Description: "Whether to fetch public ip address",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"ip_v4": {
				Description: "v4 IP address",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ip_v6": {
				Description: "v6 IP address",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_ip_v4": {
				Description: "Public v4 IP address",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_ip_v6": {
				Description: "Public v6 IP address",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Read: dataSourceIpRead,
	}
}

func dataSourceIpRead(d *schema.ResourceData, meta interface{}) error {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	// First we fetch through all the interfaces available on the host
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// retrieve all ifaceNames to check whether it matches with the one passed as an argument
	ifaceNames := make([]string, len(ifaces))

	for i := range ifaces {
		ifaceNames[i] = ifaces[i].Name
	}

	present := false

	// If nw_interface is not empty & doesn't equate to interfaces present on the host, we throw an error
	for i := range ifaceNames {
		if ifaceNames[i] == d.Get("nw_interface") || d.Get("nw_interface") == "" {
			present = true
			break
		}
	}

	if !present {
		return fmt.Errorf("passed %s network interface doesn't exist", d.Get("nw_interface"))
	}

	// We iterate through each interface to get the interface the user passed.
	// If not, we use the last non-loopback ip address of last interface
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		nwIfaceName := iface.Name

		fmt.Printf("[DEBUG] Passed network interface is %s", nwIfaceName)

		if nwIfaceName == d.Get("nw_interface") {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						if err := d.Set("ip_v4", ipnet.IP.String()); err != nil {
							return err
						}
						if err := d.Set("nw_interface", nwIfaceName); err != nil {
							return err
						}
						ipHash := sha256.Sum256([]byte(ipnet.IP.String()))
						d.SetId(hex.EncodeToString(ipHash[:]))
					}
					if ipnet.IP.To16() != nil {
						if err := d.Set("ip_v6", ipnet.IP.String()); err != nil {
							return err
						}
					}
				}
				// we return as soon as we find that interface passed is present in the list
				return nil
			}
		} else {
			for _, addr := range addrs {
				// check the address type and if it is not a loopback the display it
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						if err := d.Set("ip_v4", ipnet.IP.String()); err != nil {
							return err
						}
						if err := d.Set("nw_interface", iface.Name); err != nil {
							return err
						}
						ipHash := sha256.Sum256([]byte(ipnet.IP.String()))
						d.SetId(hex.EncodeToString(ipHash[:]))
					}
					if ipnet.IP.To16() != nil {
						if err := d.Set("ip_v6", ipnet.IP.String()); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// fetching public ip
	if d.Get("public_ip") == true {
		res, err := http.Get(IPIFY_API_URL)
		if err != nil {
			return err
		}
		ip, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if strings.Contains(string(ip), ":") {
			if err := d.Set("public_ip_v6", string(ip)); err != nil {
				return err
			}
		} else {
			if err := d.Set("public_ip_v4", string(ip)); err != nil {
				return err
			}
		}
	}

	return nil
}
