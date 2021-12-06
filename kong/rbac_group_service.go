package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractRBACGroupService handles Groups in Kong.
type AbstractRBACGroupService interface {
	// Create creates an RBAC group in Kong.
	Create(ctx context.Context, group *RBACGroup) (*RBACGroup, error)
	// Get fetches a group in Kong.
	Get(ctx context.Context, nameOrID *string) (*RBACGroup, error)
	// Update updates a group in Kong.
	Update(ctx context.Context, group *RBACGroup) (*RBACGroup, error)
	// Delete deletes a group in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of groups in Kong.
	// opt can be used to control pagination.
	List(ctx context.Context, opt *ListOpt) ([]*RBACGroup, *ListOpt, error)
	// ListAll fetches all groups in Kong.
	ListAll(ctx context.Context) ([]*RBACGroup, error)
	// AddRoles adds a comma separated list of roles to a group.
	AddRole(ctx context.Context, nameOrID *string, role *RBACRole, workspace *Workspace) (*RBACGroupRole, error)
	// DeleteRoles deletes roles associated with a group
	DeleteRole(ctx context.Context, nameOrID *string, role *RBACRole, workspace *Workspace) error
	// ListRoles returns a slice of Kong RBAC roles associated with a group.
	ListRoles(ctx context.Context, nameOrID *string) ([]*RBACGroupRole, error)
}

// RBACGroupService handles Groups in Kong.
type RBACGroupService service


// Create creates an RBAC Group in Kong.
func (s *RBACGroupService) Create(ctx context.Context,
	group *RBACGroup) (*RBACGroup, error) {

	if group == nil {
		return nil, fmt.Errorf("cannot create a nil user")
	}

	endpoint := "/groups"
	method := "POST"

	req, err := s.client.NewRequest(method, endpoint, nil, group)
	if err != nil {
		return nil, err
	}

	var createdGroup RBACGroup
	_, err = s.client.Do(ctx, req, &createdGroup)
	if err != nil {
		return nil, err
	}
	return &createdGroup, nil
}



// Get fetches a Group in Kong.
func (s *RBACGroupService) Get(ctx context.Context,
	nameOrID *string) (*RBACGroup, error) {

	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/groups/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var RBACGroup RBACGroup
	_, err = s.client.Do(ctx, req, &RBACGroup)
	if err != nil {
		return nil, err
	}
	return &RBACGroup, nil
}

// Update updates a Group in Kong.
func (s *RBACGroupService) Update(ctx context.Context,
	group *RBACGroup) (*RBACGroup, error) {

	if group == nil {
		return nil, fmt.Errorf("cannot update a nil Group")
	}

	if isEmptyString(group.ID) && isEmptyString(group.Name) {
		return nil, fmt.Errorf("ID and Name cannot both be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/groups/%v", *group.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, group)
	if err != nil {
		return nil, err
	}

	var updatedGroup RBACGroup
	_, err = s.client.Do(ctx, req, &updatedGroup)
	if err != nil {
		return nil, err
	}
	return &updatedGroup, nil
}

// Delete deletes a User in Kong
func (s *RBACGroupService) Delete(ctx context.Context,
	groupOrID *string) error {

	if isEmptyString(groupOrID) {
		return fmt.Errorf("groupOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/groups/%v", *groupOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Groups in Kong.
// opt can be used to control pagination.
func (s *RBACGroupService) List(ctx context.Context,
	opt *ListOpt) ([]*RBACGroup, *ListOpt, error) {

	data, next, err := s.client.list(ctx, "/groups/", opt)
	if err != nil {
		return nil, nil, err
	}
	var groups []*RBACGroup
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var group RBACGroup
		err = json.Unmarshal(b, &group)
		if err != nil {
			return nil, nil, err
		}
		groups = append(groups, &group)
	}

	return groups, next, nil
}

// ListAll fetches all groups in Kong.
func (s *RBACGroupService) ListAll(ctx context.Context) ([]*RBACGroup, error) {
	var groups, data []*RBACGroup
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, data...)
	}

	return groups, nil
}

// AddRole adds a single role to a Group.
func (s *RBACGroupService) AddRole(ctx context.Context,
	nameOrID *string, role *RBACRole, workspace *Workspace) (*RBACGroupRole, error) {
	
	var roleData struct {
		RoleID *string `json:"rbac_role_id,omitempty" yaml:"rbac_role_id,omitempty"`
		WorkspaceID    *string `json:"workspace_id,omitempty" yaml:"workspace_id,omitempty"`
	}
	roleData.RoleID = role.ID
	roleData.WorkspaceID = workspace.ID
	endpoint := fmt.Sprintf("/groups/%v/roles", *nameOrID)
	req, err := s.client.NewRequest("POST", endpoint, nil, roleData)
	if err != nil {
		return nil, err
	}

	var addedRole RBACGroupRole
	_, err = s.client.Do(ctx, req, &addedRole)
	if err != nil {
		return nil, fmt.Errorf("error updating role: %v", err)
	}
	return &addedRole, nil
}

// DeleteRoles deletes roles associated with a Group
func (s *RBACGroupService) DeleteRole(ctx context.Context,
	nameOrID *string, role *RBACRole, workspace *Workspace) error {
		var roleData struct {
			RoleID *string `json:"rbac_role_id,omitempty" yaml:"rbac_role_id,omitempty"`
			WorkspaceID    *string `json:"workspace_id,omitempty" yaml:"workspace_id,omitempty"`
		}
		roleData.RoleID = role.ID
		roleData.WorkspaceID = workspace.ID
		endpoint := fmt.Sprintf("/groups/%v/roles", *nameOrID)
		req, err := s.client.NewRequest("DELETE", endpoint, nil, roleData)
		if err != nil {
			return err
		}
	
		var deletedRole *RBACGroup
		_, err = s.client.Do(ctx, req, &deletedRole)

		// Weird EOF error when deleting role, so we'll just ignore it
		if err.Error() == "EOF" {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error delete role: %v", err)
		}
		return nil
}

// ListRoles returns a slice of Kong RBAC roles associated with a Group.
func (s *RBACGroupService) ListRoles(ctx context.Context,
	nameOrID *string) ([]*RBACGroupRole, error) {

	endpoint := fmt.Sprintf("/groups/%v/roles", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var listRolesResponse struct {
		Roles []*RBACGroupRole `json:"data,omitempty" yaml:"data,omitempty"`
	}
	_, err = s.client.Do(ctx, req, &listRolesResponse)
	if err != nil {
		return nil, fmt.Errorf("error retrieving list of roles: %v", err)
	}

	return listRolesResponse.Roles, nil
}
