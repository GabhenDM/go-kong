package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestRBACGroupService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", requiredFeatures{rbac: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	group := &RBACGroup{
		Name:      String("testGroupPleaseIgnore"),
		Comment:   String("testing"),
	}

	createdGroup, err := client.RBACGroups.Create(defaultCtx, group)
	assert.Nil(err)
	assert.NotNil(createdGroup)

	group, err = client.RBACGroups.Get(defaultCtx, createdGroup.ID)
	assert.Nil(err)
	assert.NotNil(group)

	group.Comment = String("new comment")
	group, err = client.RBACGroups.Update(defaultCtx, group)
	assert.Nil(err)
	assert.NotNil(group)
	assert.Equal("new comment", *group.Comment)
	err = client.RBACGroups.Delete(defaultCtx, createdGroup.ID)
	assert.Nil(err)
}

func TestGroupRoles(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", requiredFeatures{rbac: true})
	assert := assert.New(T)
	client, err := NewTestClient(nil, nil)

	assert.Nil(err)
	assert.NotNil(client)
	
	workspace := Workspace{
		Name: String("test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, &workspace)

	assert.Nil(err)
	assert.NotNil(createdWorkspace)

	roleA := &RBACRole{
		Name: String("roleA"),
	}
	roleB := &RBACRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.RBACRoles.CreateWorkspaceRole(defaultCtx, roleA,createdWorkspace)
	assert.Nil(err)
	createdRoleB, err := client.RBACRoles.CreateWorkspaceRole(defaultCtx, roleB,createdWorkspace)
	assert.Nil(err)

	group := &RBACGroup{
		Name:      String("testGroupPleaseIgnore"),
		Comment:   String("testing"),
	}

	createdGroup, err := client.RBACGroups.Create(defaultCtx, group)
	assert.Nil(err)
	assert.NotNil(createdGroup)


	updatedGroup, err := client.RBACGroups.AddRole(defaultCtx, createdGroup.ID, createdRoleA, createdWorkspace )
	assert.Nil(err)
	assert.NotNil(updatedGroup)

	updatedGroup, err = client.RBACGroups.AddRole(defaultCtx, createdGroup.ID, createdRoleB, createdWorkspace )
	assert.Nil(err)
	assert.NotNil(updatedGroup)

	roleList, err := client.RBACGroups.ListRoles(defaultCtx, createdGroup.ID)
	assert.Nil(err)
	assert.NotNil(roleList)
	assert.Equal(2, len(roleList))

	err = client.RBACGroups.DeleteRole(defaultCtx, createdGroup.ID, createdRoleA, createdWorkspace)
	assert.Nil(err)

	// Get Roles after delete
	deletedRoleList, err := client.RBACGroups.ListRoles(defaultCtx, createdGroup.ID)
	assert.Nil(err)
	assert.Equal(1, len(deletedRoleList))
}