package permissions

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceIdentityManagementPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityManagementCreateOrUpdate,
		Read:   resourceIdentityManagementRead,
		Update: resourceIdentityManagementCreateOrUpdate,
		Delete: resourceIdentityManagementDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"identity": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ForceNew:     true,
				Optional:     true,
			},
		}),
	}
}

func resourceIdentityManagementCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Identity, createIdentityManagementToken)
	if err != nil {
		return err
	}

	if err = securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceIdentityManagementRead(d, m)
}

func resourceIdentityManagementRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Identity, createIdentityManagementToken)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn)
	if err != nil {
		return err
	}
	if principalPermissions == nil {
		d.SetId("")
		log.Printf("[INFO] Permissions for ACL token %q not found. Removing from state", sn.GetToken())
		return nil
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceIdentityManagementDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Identity, createIdentityManagementToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func createIdentityManagementToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID := d.Get("project_id").(string)

	aclToken := fmt.Sprintf("%s", projectID)
	if v, ok := d.GetOk("identity"); ok {

		idlist, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: converter.String(v.(string)),
		})

		if err != nil {
			return "", err
		}
		if idlist == nil || len(*idlist) != 1 {
			return "", fmt.Errorf("Failed to load identity information for defined principals [%s]", v.(string))
		}
		aclToken += "\\" + (*idlist)[0].Id.String()
	}

	return aclToken, nil
}
