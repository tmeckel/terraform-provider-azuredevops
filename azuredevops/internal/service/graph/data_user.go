package graph

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataUsers schema and implementation for users data source
func DataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,

		//https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"principal_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"origin", "origin_id", "mail_address"},
				ExactlyOneOf:  []string{"principal_name", "mail_address", "origin"},
			},
			"mail_address": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"origin", "origin_id", "principal_name"},
				ExactlyOneOf:  []string{"principal_name", "mail_address", "origin"},
			},
			"origin": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"principal_name", "mail_address"},
				RequiredWith:  []string{"origin_id"},
				ExactlyOneOf:  []string{"principal_name", "mail_address", "origin"},
			},
			"origin_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				RequiredWith: []string{"origin"},
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataUserRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	principalName, byPrincipalName := d.GetOk("principal_name")
	mailAddress, byMailAddress := d.GetOk("mail_address")
	origin, byOrigin := d.GetOk("origin")
	originID := d.Get("origin_id").(string)

	if !byPrincipalName && !byMailAddress && !byOrigin {
		return fmt.Errorf("At least one attribute of principal_name, mail_address or the combination of origin and origin_id must be specified for searching")
	}
	var currentToken string
	var puser *graph.GraphUser = nil

found:
	for hasMore := true; hasMore; {
		newUsers, latestToken, err := getUsersWithContinuationToken(clients, nil, currentToken)
		currentToken = latestToken
		hasMore = currentToken != ""
		if err != nil {
			return err
		}
		for _, user := range newUsers {
			if byPrincipalName && user.PrincipalName != nil && strings.EqualFold(principalName.(string), *user.PrincipalName) {
				d.SetId("user#" + *user.PrincipalName)

				d.Set("mail_address", user.MailAddress)
				d.Set("origin", user.Origin)
				d.Set("origin_id", user.OriginId)

				puser = &user
				break found
			} else if byMailAddress && user.MailAddress != nil && strings.EqualFold(mailAddress.(string), *user.MailAddress) {
				d.SetId("user#" + *user.MailAddress)

				d.Set("principal_name", user.PrincipalName)
				d.Set("origin", user.Origin)
				d.Set("origin_id", user.OriginId)

				puser = &user
				break found
			} else if byOrigin && user.Origin != nil && user.OriginId != nil && strings.EqualFold(origin.(string), *user.Origin) && strings.EqualFold(originID, *user.OriginId) {
				d.SetId("user#" + *user.Origin + ":" + *user.OriginId)

				d.Set("mail_address", user.MailAddress)
				d.Set("principal_name", user.PrincipalName)

				puser = &user
				break found
			}
		}
	}

	if puser == nil {
		errMsg := "unable to find user by "
		if byPrincipalName {
			errMsg += "principal name " + principalName.(string)
		} else if byMailAddress {
			errMsg += "mail address " + mailAddress.(string)
		} else if byOrigin {
			errMsg += "origin " + origin.(string) + " and originID " + originID
		}
		return fmt.Errorf(errMsg)
	}

	d.Set("descriptor", puser.Descriptor)
	d.Set("display_name", puser.DisplayName)
	return nil
}
