package ip

import (
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAzureStackClientConfig_basic(t *testing.T) {
	dataSourceName := "data.ip_v4.current"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckArmClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testIPConfigAttr(dataSourceName, "ip", resourceBoardReadTest()),
				),
			},
		},
	})
}

// Wraps resource.TestCheckResourceAttr to prevent leaking values to console
// in case of mismatch
func testIPConfigAttr(name, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(name, key, value)(s)
	}
}

func resourceBoardReadTest() string {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

const testAccCheckArmClientConfig_basic = `
data "ip_v4" "current" {}
`
