// +build all core data_sources data_user
// +build !exclude_data_sources !exclude_data_user

package graph

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/stretchr/testify/require"
)

// verfies that the data source propagates an error from the API correctly
func TestDataSourceUser_Read_TestDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: nil,
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ListUsers() Failed"))

	resourceData := schema.TestResourceDataRaw(t, DataUser().Schema, nil)
	resourceData.Set("principal_name", "DesireeMCollins@jourrapide.com")
	err := dataUserRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ListUsers() Failed")
}

func TestDataSourceUser_Read_HandlesContinuationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	var calls []*gomock.Call
	calls = append(calls, graphClient.
		EXPECT().
		ListUsers(clients.Ctx, graph.ListUsersArgs{
			SubjectTypes: nil,
		}).
		Return(&graph.PagedGraphUsers{
			GraphUsers:        &usrList1,
			ContinuationToken: &[]string{"2"},
		}, nil).
		Times(1))

	calls = append(calls, graphClient.
		EXPECT().
		ListUsers(clients.Ctx, graph.ListUsersArgs{
			SubjectTypes:      nil,
			ContinuationToken: converter.String("2"),
		}).
		Return(&graph.PagedGraphUsers{
			GraphUsers:        &usrList2,
			ContinuationToken: &[]string{""},
		}, nil).
		Times(1))

	gomock.InOrder(calls...)

	resourceData := schema.TestResourceDataRaw(t, DataUser().Schema, nil)
	resourceData.Set("principal_name", "AllenBMcKinnon@dayrep.com")
	err := dataUserRead(resourceData, clients)
	require.Nil(t, err)
}

func TestDataSourceUser_Read_TestFilterByPricipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: nil,
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUser().Schema, nil)
	idx := 0
	resourceData.Set("principal_name", usrList1[idx].PrincipalName)
	err := dataUserRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, *usrList1[idx].Descriptor, resourceData.Get("descriptor").(string))
	require.Equal(t, *usrList1[idx].PrincipalName, resourceData.Get("principal_name").(string))
	require.Equal(t, *usrList1[idx].DisplayName, resourceData.Get("display_name").(string))
	require.Equal(t, *usrList1[idx].Origin, resourceData.Get("origin").(string))
	require.Equal(t, *usrList1[idx].OriginId, resourceData.Get("origin_id").(string))
	require.Equal(t, *usrList1[idx].MailAddress, resourceData.Get("mail_address").(string))
}

func TestDataSourceUser_Read_TestFilterByOriginOriginId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: nil,
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUser().Schema, nil)
	idx := 3
	resourceData.Set("origin", usrList1[idx].Origin)
	resourceData.Set("origin_id", usrList1[idx].OriginId)
	err := dataUserRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, *usrList1[idx].Descriptor, resourceData.Get("descriptor").(string))
	require.Equal(t, *usrList1[idx].PrincipalName, resourceData.Get("principal_name").(string))
	require.Equal(t, *usrList1[idx].DisplayName, resourceData.Get("display_name").(string))
	require.Equal(t, *usrList1[idx].Origin, resourceData.Get("origin").(string))
	require.Equal(t, *usrList1[idx].OriginId, resourceData.Get("origin_id").(string))
	if usrList1[idx].MailAddress != nil {
		require.Equal(t, *usrList1[idx].MailAddress, resourceData.Get("mail_address").(string))
	} else {
		require.Equal(t, "", resourceData.Get("mail_address").(string))
	}
}

func TestDataSourceUser_Read_TestFilterMailAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: nil,
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUser().Schema, nil)
	idx := 2
	resourceData.Set("mail_address", usrList1[idx].MailAddress)
	err := dataUserRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, *usrList1[idx].Descriptor, resourceData.Get("descriptor").(string))
	require.Equal(t, *usrList1[idx].PrincipalName, resourceData.Get("principal_name").(string))
	require.Equal(t, *usrList1[idx].DisplayName, resourceData.Get("display_name").(string))
	require.Equal(t, *usrList1[idx].Origin, resourceData.Get("origin").(string))
	require.Equal(t, *usrList1[idx].OriginId, resourceData.Get("origin_id").(string))
	require.Equal(t, *usrList1[idx].MailAddress, resourceData.Get("mail_address").(string))
}
