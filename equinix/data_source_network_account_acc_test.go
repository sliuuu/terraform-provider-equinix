package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccNetworkAccountDataSource(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"resourceName": "tf-account",
		"metro_code":   metro.(string),
		"status":       "active",
	}
	resourceName := fmt.Sprintf("data.equinix_network_account.%s", context["resourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAccountDataSource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "number"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "ucm_id"),
				),
			},
		},
	})
}

func testAccNetworkAccountDataSource(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "%{resourceName}" {
  metro_code = "%{metro_code}"
  status     = "%{status}"
}
`, ctx)
}
