package gitguardian

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Gaardsholt/go-gitguardian/client"
	"github.com/Gaardsholt/go-gitguardian/teams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hugolesta/terraform-provider-gitguardian/gitguardian"
)

func resourceCreateTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateTeamCreate,
		Read:   resourceCreateTeamRead,
		// Update: resourceCreateTeamUpdate,
		// Delete: resourceCreateTeamDelete,

		Schema: map[string]*schema.Schema{
			"team_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCreateTeamCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(client.ClientOption)
	teamsOps := m.(teams.TeamsCreateOptions)
	teamsOps.Name = d.Get("team_name").(string)
	team, err := teams.NewClient(client)
	if err != nil {
		return err
	}

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		if _, err := team.Create(teamsOps); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return resourceCreateTeamRead(d, m)
}

type Team struct {
	Id int `json:"id"`
	MemberId int `json:"member_id"`
	TeamId int `json:"team_id"`
	TeamPermission string `json:"team_permission"`
	IncidentPermission string `json:"incident_permission"`
}

func buildTeamStruct(d *schema.ResourceData) *teams.TeamGetResult {

	team := &gitguardian.Team{
		Id:  strings.TrimSpace(d.Get("team_id").(string)),
		MemberId:  strings.TrimSpace(d.Get("member_id").(string)),
		TeamId:  strings.TrimSpace(d.Get("team_id").(string)),
		TeamPermission:  strings.TrimSpace(d.Get("team_permission").(string)),
		IncidentPermission:  strings.TrimSpace(d.Get("incident_permission").(string)),

	}

	if attr, ok := d.GetOk("id"); ok {
		team.Id = attr.(int)
	}

	if attr, ok := d.GetOk("member_id"); ok {
		team.MemberId = attr.(int)
	}

	if attr, ok := d.GetOk("team_id"); ok {
		team.TeamId = attr.(int)
	}

	if attr, ok := d.GetOk("team_permission"); ok {
		team.TeamPermissions = attr.(string)
	}

	if attr, ok := d.GetOk("incident_permission"); ok {
		team.IncidentPermission = attr.(string)
	}
	return team
}

func resourceCreateTeamRead(d *schema.ResourceData, m interface{}) error {
	client := m.(client.ClientOption)
	team, err := teams.NewClient(client)
	listMembershipOpt := m.(teams.ListMembershipsOptions)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading GitGuardian Team %s", d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		team_id := d.Get("member_id").(int)

		team,_, err := team.ListMemberships(team_id,listMembershipOpt)
		
		if err != nil {
			return resource.NonRetryableError(
				fmt.Errorf("error Reading teams: %s", err),
			)
		}
		
		d.Set("id", team.Result.ID)
		d.Set("member_id", team.Result.MemberID)
		d.Set("team_id", team.Result.TeamID)
		d.Set("team_permission", team.Result.TeamPermission)
		d.Set("incident_permission", team.Result.IncidentPermission)

		return nil
	})
}
