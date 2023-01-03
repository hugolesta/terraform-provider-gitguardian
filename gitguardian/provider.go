package gitguardian

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hugolesta/terraform-provider-gitguardian/api/client"
)
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("GITGUARDIAN_TOKEN", nil),
				Description: "GitGuardian API token",
			},
			"url": {
				Type: schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("GITGUARDIAN_URL", " https://api.gitguardian.com/v1/"),
				Description: "GitGuardian API url",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"create_team": resourceCreateTeam(),
		},
		
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	
	token := d.Get("token").(string)
	url := d.Get("url").(string)

	return client.NewClient(token, url)
}
