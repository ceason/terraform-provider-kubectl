package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"errors"
	"fmt"
)

func resourceGenericObject() *schema.Resource {
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
		},
		Create:        resourceGenericObjectCreate,
		Read:          resourceGenericObjectRead,
		Update:        resourceGenericObjectUpdate,
		Delete:        resourceGenericObjectDelete,
		CustomizeDiff: resourceGenericObjectDiff,
		Exists:        resourceGenericObjectExists,
	}
}

func resourceGenericObjectExists(d *schema.ResourceData, provider interface{}) (bool, error) {
	k := provider.(*providerConfig)
	return k.ObjectExists(d.Id())
}

/**
- recalculate the resource name/namespace/kind(/apigroup?)
- update the 'last_applied_configuration'
*/
func resourceGenericObjectRead(d *schema.ResourceData, provider interface{}) error {
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

func resourceGenericObjectDiff(d *schema.ResourceDiff, provider interface{}) error {
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
	}
	return nil
}

/*
- Generates and sets `id` globally unique id
- Calls kubectl apply
*/
func resourceGenericObjectCreate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	obj, err := p.NewObject(d.Get("yaml").(string))
	if err != nil {
		return err
	}
	exists, err := p.ObjectExists(obj.ResourceId())
	if err != nil {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("Object already exists. Try importing it with: terraform import <resource_name> %s", obj.ResourceId()))
	}
	err = p.Apply(obj)
	if err != nil {
		return err
	}
	d.Set("last_applied_configuration", obj.LastAppliedConfigurationHash())
	d.Set("name", obj.Metadata.Name)
	d.Set("namespace", obj.Metadata.Namespace)
	d.Set("kind", obj.Kind)
	d.SetId(obj.ResourceId())

	// todo: detect when we create a CRD and add it to the provider's "apiResources" (so usages will know if the resource is namespaced or not)

	return nil
}

// Calls kubectl apply
func resourceGenericObjectUpdate(d *schema.ResourceData, provider interface{}) error {
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
func resourceGenericObjectDelete(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	yamlStr := d.Get("yaml").(string)
	return p.Delete(yamlStr)
}
