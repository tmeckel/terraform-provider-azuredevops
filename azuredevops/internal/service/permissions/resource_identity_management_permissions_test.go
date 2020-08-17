// +build all permissions resource_identity_management_permissions
// +build !exclude_permissions !resource_identity_management_permissions

package permissions

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

var identityPermissionsProjectID = "f454422e-57b3-442a-8dde-b1b6b7c40b95"
var identityPermissionsIdentitySubject = "vssgp.Uy0xLTktMTU1MTM3NDI0NS0zNTMwNTgyMDY3LTEwNTQ4MTM1MTYtMjQ1ODAwNjc3Mi0yNTI1NjM2NjkxLTAtMC0wLTAtMw"
var identityPermissionsIdentityID = "38d476af-3062-47c5-a293-979e0a56063e"

func TestIdentityManagementPermissions_CreateToken_ProjectGlobal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients, _ := initClient(ctrl)

	d := getIdentityManagementPermissionsResource(t, identityPermissionsProjectID, "")
	token, err := createIdentityManagementToken(d, clients)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, identityPermissionsProjectID, token)
}

func TestIdentityManagementPermissions_CreateToken_DontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients, identityClient := initClient(ctrl)

	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &identityPermissionsIdentitySubject,
		}).
		Return(nil, fmt.Errorf("@@ReadIdentities@@failed")).
		Times(1)

	d := getIdentityManagementPermissionsResource(t, identityPermissionsProjectID, identityPermissionsIdentitySubject)
	token, err := createIdentityManagementToken(d, clients)
	assert.NotNil(t, err)
	assert.Empty(t, token)
}

func TestIdentityManagementPermissions_CreateToken_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients, identityClient := initClient(ctrl)

	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &identityPermissionsIdentitySubject,
		}).
		Return(&[]identity.Identity{}, nil).
		Times(1)

	d := getIdentityManagementPermissionsResource(t, identityPermissionsProjectID, identityPermissionsIdentitySubject)
	token, err := createIdentityManagementToken(d, clients)
	assert.NotNil(t, err)
	assert.Empty(t, token)
}

func TestIdentityManagementPermissions_CreateToken_MultipleResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients, identityClient := initClient(ctrl)

	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &identityPermissionsIdentitySubject,
		}).
		Return(&[]identity.Identity{
			{
				Id: converter.UUID("c6c8d792-895f-49b2-a263-5d2edc1fd9f5"),
			},
			{
				Id: converter.UUID("ae3c83a2-729f-4914-82a2-ca7262627cba"),
			},
		}, nil).
		Times(1)

	d := getIdentityManagementPermissionsResource(t, identityPermissionsProjectID, identityPermissionsIdentitySubject)
	token, err := createIdentityManagementToken(d, clients)
	assert.NotNil(t, err)
	assert.Empty(t, token)
}

func TestIdentityManagementPermissions_CreateToken_WithIdentity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients, identityClient := initClient(ctrl)

	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &identityPermissionsIdentitySubject,
		}).
		Return(&[]identity.Identity{
			{
				Id: converter.UUID(identityPermissionsIdentityID),
			},
		}, nil).
		Times(1)

	d := getIdentityManagementPermissionsResource(t, identityPermissionsProjectID, identityPermissionsIdentitySubject)
	token, err := createIdentityManagementToken(d, clients)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	ref := fmt.Sprintf("%s\\%s", identityPermissionsProjectID, identityPermissionsIdentityID)
	assert.Equal(t, ref, token)
}

func getIdentityManagementPermissionsResource(t *testing.T, projectID string, identity string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceIdentityManagementPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if identity != "" {
		d.Set("identity", identity)
	}
	return d
}

func initClient(ctrl *gomock.Controller) (*client.AggregatedClient, *azdosdkmocks.MockIdentityClient) {
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	return clients, identityClient
}
