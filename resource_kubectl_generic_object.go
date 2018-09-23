package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"errors"
	"fmt"
	"strings"
)

func resourceKubectlGenericObject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"yaml": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_applied_configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"api_group": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
		Create:        resourceKubectlGenericObjectCreate,
		Read:          resourceKubectlGenericObjectRead,
		Update:        resourceKubectlGenericObjectUpdate,
		Delete:        resourceKubectlGenericObjectDelete,
		CustomizeDiff: resourceKubectlGenericObjectDiff,
		Exists:        resourceKubectlGenericObjectExists,
		Importer: &schema.ResourceImporter{
			State: resourceKubectlGenericObjectImportState,
		},

	}
}

func resourceKubectlGenericObjectExists(d *schema.ResourceData, provider interface{}) (bool, error) {
	k := provider.(*providerConfig)
	return k.ObjectExists(d.Id())
}

/**
- recalculate the resource name/namespace/kind(/apigroup?)
- update the 'last_applied_configuration'
*/
func resourceKubectlGenericObjectRead(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	obj, err := p.GetObject(d.Id())
	if err != nil {
		return err
	}
	if d.Get("last_applied_configuration").(string) != obj.LastAppliedConfigurationHash() {
		d.Set("last_applied_configuration", obj.LastAppliedConfigurationHash())
	}
	return nil
}

func resourceKubectlGenericObjectDiff(d *schema.ResourceDiff, provider interface{}) error {
	p := provider.(*providerConfig)
	if d.HasChange("yaml") {
		// figure out if immutable things have changed
		obj, err := p.NewObject(d.Get("yaml").(string))
		if err != nil {
			return err
		}
		if obj.Metadata.Name != d.Get("name").(string) {
			d.SetNew("name", obj.Metadata.Name)
		}
		if obj.Metadata.Namespace != d.Get("namespace").(string) {
			d.SetNew("namespace", obj.Metadata.Namespace)
		}
		if obj.Kind != d.Get("kind").(string) {
			d.SetNew("kind", obj.Kind)
		}
		if obj.ApiGroup() != d.Get("api_group").(string) {
			d.SetNew("api_group", obj.ApiGroup())
		}
	}
	return nil
}

/*
- Generates and sets `id` globally unique id
- Calls kubectl apply
*/
func resourceKubectlGenericObjectCreate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	obj, err := p.NewObject(d.Get("yaml").(string))
	if err != nil {
		return err
	}

	// Temporarily ignoring existing objects to preserve existing behavior
	/*
	exists, err := p.ObjectExists(obj.ResourceId())
	if err != nil {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("Object already exists. Try importing it with: terraform import <resource_name> %s", obj.ResourceId()))
	}
	// */
	err = p.Apply(obj)
	if err != nil {
		return err
	}
	d.Set("last_applied_configuration", obj.LastAppliedConfigurationHash())
	d.Set("name", obj.Metadata.Name)
	d.Set("namespace", obj.Metadata.Namespace)
	d.Set("kind", obj.Kind)
	d.Set("api_group", obj.ApiGroup())
	d.SetId(obj.ResourceId())

	// todo: detect when we create a CRD and add it to the provider's "apiResources" (so usages will know if the resource is namespaced or not)

	return nil
}

// Calls kubectl apply
func resourceKubectlGenericObjectUpdate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)

	if d.HasChange("yaml") || d.HasChange("last_applied_configuration") {
		obj, err := p.NewObject(d.Get("yaml").(string))
		if err != nil {
			return err
		}
		err = p.Apply(obj)
		if err != nil {
			return err
		}
		d.Set("last_applied_configuration", obj.LastAppliedConfigurationHash())
	}
	return nil
}

// Calls kubectl delete
func resourceKubectlGenericObjectDelete(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	namespace := d.Get("namespace").(string)
	kind := d.Get("kind").(string)
	name := d.Get("name").(string)
	return p.Delete(namespace, kind, name)
}

func resourceIdParts(resourceId string) (namespace, kind, name string, err error) {
	parts := strings.Split(resourceId, "/")
	if len(parts) != 3 {
		err = errors.New(fmt.Sprintf("ResourceId should have exactly 3 parts: '%s'", resourceId))
	}
	namespace = parts[0]
	kind = parts[1]
	name = parts[2]
	return
}


func resourceKubectlGenericObjectImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	p := meta.(*providerConfig)

	results := make([]*schema.ResourceData, 1)
	results[0] = d

	obj, err := p.GetObject(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("last_applied_configuration", obj.LastAppliedConfigurationHash())
	d.Set("name", obj.Metadata.Name)
	d.Set("namespace", obj.Metadata.Namespace)
	d.Set("kind", obj.Kind)
	d.Set("api_group", obj.ApiGroup())

	return results, nil
}