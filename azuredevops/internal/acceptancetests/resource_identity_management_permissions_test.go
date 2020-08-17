// +build all permissions resource_identity_management_permissions
// +build !exclude_permissions !resource_identity_management_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func hclIdentityManagementPermissions(projectName string, createPrincipalPermssions bool, permissions map[string]string) string {
	projectResource := testutils.HclProjectResource(projectName)
	szPermissions := datahelper.JoinMap(permissions, "=", "\n")

	identity := ""
	if createPrincipalPermssions {
		identity = "identity = data.azuredevops_group.project-readers.id"
	}
	return fmt.Sprintf(`
%s

data "azuredevops_group" "project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

data "azuredevops_group" "project-contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}

resource "azuredevops_identity_management_permissions" "test" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.project-contributors.id
	%s
	permissions = {
		%s
	}
  }

`, projectResource, identity, szPermissions)
}

func TestAccIdentityManagementPermissions_SetProjectPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclIdentityManagementPermissions(projectName, false, map[string]string{
		"Read":             "Allow",
		"Write":            "Deny",
		"Delete":           "NotSet",
		"ManageMembership": "NotSet",
	})

	tfNode := "azuredevops_identity_management_permissions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManageMembership", "notset"),
				),
			},
		},
	})
}

func TestAccIdentityManagementPermissions_UpdateProjectPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclIdentityManagementPermissions(projectName, false, map[string]string{
		"Read":             "Allow",
		"Write":            "Deny",
		"Delete":           "NotSet",
		"ManageMembership": "NotSet",
	})
	config2 := hclIdentityManagementPermissions(projectName, false, map[string]string{
		"Read":             "Allow",
		"Write":            "Allow",
		"Delete":           "Deny",
		"ManageMembership": "Allow",
	})

	tfNode := "azuredevops_identity_management_permissions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManageMembership", "notset"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManageMembership", "allow"),
				),
			},
		},
	})
}

func TestAccIdentityManagementPermissions_SetPrincipalPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclIdentityManagementPermissions(projectName, true, map[string]string{
		"Read":             "Allow",
		"Write":            "Deny",
		"Delete":           "NotSet",
		"ManageMembership": "NotSet",
	})

	tfNode := "azuredevops_identity_management_permissions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttrSet(tfNode, "identity"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManageMembership", "notset"),
				),
			},
		},
	})
}

func TestAccIdentityManagementPermissions_UpdatePrincipalPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclIdentityManagementPermissions(projectName, true, map[string]string{
		"Read":             "Allow",
		"Write":            "Deny",
		"Delete":           "NotSet",
		"ManageMembership": "NotSet",
	})
	config2 := hclIdentityManagementPermissions(projectName, true, map[string]string{
		"Read":             "Allow",
		"Write":            "Allow",
		"Delete":           "Deny",
		"ManageMembership": "Allow",
	})

	tfNode := "azuredevops_identity_management_permissions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttrSet(tfNode, "identity"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManageMembership", "notset"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttrSet(tfNode, "identity"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManageMembership", "allow"),
				),
			},
		},
	})
}
