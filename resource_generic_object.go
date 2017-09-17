package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/satori/go.uuid"
)

func resourceGenericObject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"yaml": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceGenericObjectCreate,
		Read:   resourceGenericObjectRead,
		Update: resourceGenericObjectUpdate,
		Delete: resourceGenericObjectDelete,
	}
}

/*
- Generates and sets `id` globally unique id
- Calls kubectl apply
*/
func resourceGenericObjectCreate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	id := uuid.NewV4().String()
	d.SetId(id)
	return p.kubectl.Apply(d.Get("yaml").(string), id)
}

/**
- noop for now
- pretty sure this is unnecessary for min functionality because `apply` handles all of the state transition logic for us
*/
func resourceGenericObjectRead(d *schema.ResourceData, provider interface{}) error {
	return nil
}

// Calls kubectl apply
func resourceGenericObjectUpdate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	d.Partial(true)
	if d.HasChange("yaml") {
		err := p.kubectl.Apply(d.Get("yaml").(string), d.Id())
		if err != nil {
			return err
		}
		d.SetPartial("yaml")
	}
	d.Partial(false)
	return nil
}

// Calls kubectl delete
func resourceGenericObjectDelete(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	yamlStr := d.Get("yaml").(string)
	return p.kubectl.Delete(yamlStr)
}
