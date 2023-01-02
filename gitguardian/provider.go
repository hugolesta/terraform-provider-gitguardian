package gitguardian

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hugolesta/terraform-provider-gitguardian/gitguardian/api/client"
)
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
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
		ResourcesMap: map[string]*schema.Resource{
			"create_team": resourceCreateTeam(),
		},
		ConfigureFunc: providerConfigure,
		
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)
	url := d.Get("url").(string)
	return client.NewClient(token, url), nil
}
