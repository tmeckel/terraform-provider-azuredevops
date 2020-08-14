---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_identity_management_permissions"
description: |-
  Manages permissions for managing identities (users , groups) in Azure Devops
---

# azuredevops_identity_management_permissions

Manages permissions for managing identities (users , groups) in Azure Devops 

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Identity Management within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `identity`.

### Project level

Permissions for all identities (existing and newly created identities) inside a project are specified, if only the argument `project_id` has a value.

#### Example usage

```hcl
resource "azuredevops_identity_management_permissions" "project-permissions" {
  project_id  = azuredevops_project.test.id
  principal   = data.azuredevops_group.project-readers.id
  permissions = {
    Read             = "NotSet"
    Write            = "Deny"
    Delete           = "Deny"
    ManageMembership = "Deny"
  }
}
```

### Identity level

Identity management permissions for a specific identity are specified if the arguments `project_id` and `identity` are set.

#### Example usage

```hcl
resource "azuredevops_identity_management_permissions" "identity-permissions" {
  project_id    = azuredevops_project.test.id
  identity      = data.azuredevops_group.project-readers.id
  principal     = data.azuredevops_group.project-contributors.id
  permissions   = {
    Read             = "NotSet"
    Write            = "Allow"
    Delete           = "Deny"
    ManageMembership = "Allow"
  }
}
```

## Example Usage

```hcl
resource "azuredevops_project" "test" {
  project_name       = "Test Project"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_group" "project-readers" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

data "azuredevops_group" "project-contributors" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

data "azuredevops_group" "project-administrators" {
  project_id = azuredevops_project.test.id
  name       = "Project administrators"
}

resource "azuredevops_identity_management_permissions" "project-permissions" {
  project_id  = azuredevops_project.test.id
  principal   = data.azuredevops_group.project-readers.id
  permissions = {
    Read             = "NotSet"
    Write            = "Deny"
    Delete           = "Deny"
    ManageMembership = "Deny"
  }
}

resource "azuredevops_identity_management_permissions" "identity-permissions" {
  project_id    = azuredevops_project.test.id
  identity      = data.azuredevops_group.project-readers.id
  principal     = data.azuredevops_group.project-contributors.id
  permissions   = {
    Read             = "NotSet"
    Write            = "Allow"
    Delete           = "Deny"
    ManageMembership = "Allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `identity` - (Optional) The subject descriptor of the identity for which the permssions should be set. 
* `principal` - (Required) The **group** principal to assign the permissions.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`
* `permissions` - (Required) the permissions to assign. The follwing permissions are available

| Permission       | Description                 |
|------------------|-----------------------------|
| Read             | View identity information   |
| Write            | Edit identity information   |
| Delete           | Delete identity information |
| ManageMembership | Manage group membership     |
| CreateScope      | Create identity scopes      |
| RestoreScope     | Restore identity scopes     |

~> **Note** To manage members of a group the permissions `Write` **and** `ManageMembership` must be set simultaneously on the target group

## Relevant Links

* [Azure DevOps Service REST API 5.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-5.1)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
