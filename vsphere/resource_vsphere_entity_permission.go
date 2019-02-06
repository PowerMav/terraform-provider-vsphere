package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/permissions"
)

func resourceVSphereEntityPermission() *schema.Resource {
	return &schema.Resource{
		Read:   resourceVSphereEntityPermissionRead,
		Create: resourceVSphereEntityPermissionCreate,
		Delete: resourceVSphereEntityPermissionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereEntityPermissionImport,
		},

		Schema: map[string]*schema.Schema{
			"principal": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"entity_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"propagate": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVSphereEntityPermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	entityID, entityType, principal, err := permission.SplitID(d.Id())
	if err != nil {
		return err
	}
	p, err := permission.ByID(client, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if err = d.Set("propagate", p.Propagate); err != nil {
		return err
	}
	if err = d.Set("role_id", fmt.Sprint(p.RoleId)); err != nil {
		return err
	}
	if err = d.Set("group", p.Group); err != nil {
		return err
	}
	if err = d.Set("principal", principal); err != nil {
		return err
	}
	if err = d.Set("entity_id", entityID); err != nil {
		return err
	}
	if err = d.Set("entity_type", entityType); err != nil {
		return err
	}
	return nil
}

func resourceVSphereEntityPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	principal := d.Get("principal").(string)
	entityID := d.Get("entity_id").(string)
	entityType := d.Get("entity_type").(string)
	group := d.Get("group").(bool)
	roleID := d.Get("role_id").(int)
	propagate := d.Get("propagate").(bool)
	err := permission.Create(client, entityType, entityID, principal, roleID, group, propagate)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(permission.ConcatID(entityID, entityType, principal))
	return resourceVSphereEntityPermissionRead(d, meta)
}

func resourceVSphereEntityPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	p, err := permission.ByID(client, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	err = permission.Remove(client, p)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceVSphereEntityPermissionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
