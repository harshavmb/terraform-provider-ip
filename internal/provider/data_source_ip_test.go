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

const testAccDataSourceIp = `
data "ip" "foo" {
	nw_interface = "eth0"
}
`
