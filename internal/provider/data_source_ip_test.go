package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIp(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.ip.foo", "nw_interface", regexp.MustCompile("^eth")),
				),
			},
		},
	})
}

func TestAccDataSourcePublicIp(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePublicIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.ip.nat", "public_ip_v4", regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)),
				),
			},
		},
	})
}

const testAccDataSourceIp = `
data "ip" "foo" {
	nw_interface = "eth0"
}
`

const testAccDataSourcePublicIp = `
data "ip" "nat" {
	public_ip = true
}
`
