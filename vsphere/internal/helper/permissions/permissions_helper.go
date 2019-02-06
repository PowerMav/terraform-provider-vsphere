package permission

import (
	"errors"
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// ByID check if a permissions exist, and return that permissions
func ByID(client *govmomi.Client, id string) (*types.Permission, error) {
	entityID, entityType, principal, err := SplitID(id)
	if err != nil {
		return nil, err
	}
	m := object.NewAuthorizationManager(client.Client)
	finder := find.NewFinder(client.Client, true)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	ref := types.ManagedObjectReference{
		Type:  entityType,
		Value: entityID,
	}

	entity, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}

	permissions, err := m.RetrieveEntityPermissions(ctx, entity.Reference(), true)
	if err != nil {
		return nil, err
	}

	for _, permission := range permissions {
		if permission.Principal == principal {
			return &permission, nil
		}
	}

	return nil, errors.New("There is no prinicipal with name " + principal)
}

// Create Entity Permission
func Create(client *govmomi.Client, entityID string, principal string, entityType string, roleID int, group bool, propagate bool) error {
	m := object.NewAuthorizationManager(client.Client)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	ref := types.ManagedObjectReference{
		Type:  entityType,
		Value: entityID,
	}

	perms := []types.Permission{types.Permission{
		Principal: principal,
		RoleId:    int32(roleID),
		Group:     group,
		Propagate: propagate,
	}}

	return m.SetEntityPermissions(ctx, ref, perms)
}

// Remove Entity Permission
func Remove(client *govmomi.Client, permission *types.Permission) error {
	m := object.NewAuthorizationManager(client.Client)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	return m.RemoveEntityPermission(ctx, permission.Entity.Reference(), permission.Principal, permission.Group)
}

// SplitID takes the permission's ID and splits it into the folder and principal.
func SplitID(id string) (string, string, string, error) {
	s := strings.Split(id, ":")
	if len(s) != 3 {
		return "", "", "", fmt.Errorf("role ID does not contain principal, entity type, and entity ID")
	}
	return s[0], s[1], s[2], nil
}

// ConcatID takes a permission's folder and principal and generates an ID.
func ConcatID(id, entityType, principal string) string {
	return fmt.Sprintf("%s:%s:%s", id, entityType, principal)
}
