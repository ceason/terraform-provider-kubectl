package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"errors"
	"fmt"
	"strings"
	"strconv"
)

const (
	PropertymapFieldname = "z"
)

func resourceKubectlGenericObject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"yaml": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
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
			PropertymapFieldname: {
				Type:        schema.TypeMap,
				Description: "Parsed object fields to make 'terraform plan' output more meaningful",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
- update the object's "property map", filtered by which keys are already set
*/
func resourceKubectlGenericObjectRead(d *schema.ResourceData, provider interface{}) error {
	///*
	p := provider.(*providerConfig)
	obj, err := p.GetObject(d.Id())
	if err != nil {
		return err
	}
	actual := obj.Properties()
	cached := d.Get(PropertymapFieldname).(map[string]interface{})
	out := make(map[string]string)
	for path, v := range cached {
		if value, ok := actual[path]; ok {
			out[path] = value
		} else {
			out[path] = v.(string)
		}
	}
	d.Set(PropertymapFieldname, out)
	//*/
	return nil
}

func resourceKubectlGenericObjectDiff(d *schema.ResourceDiff, provider interface{}) error {
	p := provider.(*providerConfig)

	if !d.NewValueKnown("yaml") {
		return errors.New("Could not determine value for 'yaml' field at plan-time. A common cause of this error is an interpolation which depends on another resource's output (recommendation is to avoid this type of interpolation) ")
	}

	// figure out if calculated fields have changed
	if d.HasChange("yaml") || d.Id() != "" {
		cfg, err := p.NewObjectConfig(d.Get("yaml").(string))
		if err != nil {
			return err
		}
		if cfg.ApiVersion != d.Get("api_version").(string) {
			d.SetNew("api_version", cfg.ApiVersion)
		}
		if cfg.Kind != d.Get("kind").(string) {
			d.SetNew("kind", cfg.FullKind())
		}
		if cfg.Metadata.Name != d.Get("name").(string) {
			d.SetNew("name", cfg.Metadata.Name)
		}
		if cfg.Metadata.Namespace != d.Get("namespace").(string) {
			d.SetNew("namespace", cfg.Metadata.Namespace)
		}
		d.SetNew(PropertymapFieldname, cfg.Properties())
	}

	return nil
}

/*
- Generates and sets `id` globally unique id
- Calls kubectl apply
*/
func resourceKubectlGenericObjectCreate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	cfg, err := p.NewObjectConfig(d.Get("yaml").(string))
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
	obj, err := p.Apply(cfg, false)
	if err != nil {
		return err
	}
	d.Set("name", obj.Metadata.Name)
	d.Set("namespace", obj.Metadata.Namespace)
	d.Set("kind", obj.Kind)
	d.Set("api_version", obj.ApiVersion)
	d.Set(PropertymapFieldname, cfg.Properties())
	d.SetId(cfg.ResourceId())

	// todo: detect when we create a CRD and add it to the provider's "apiResources" (so usages will know if the resource is namespaced or not)

	return nil
}

// Calls kubectl apply
func resourceKubectlGenericObjectUpdate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)

	// we check just this field because it's our primary way of communicating the diff to the user (ie we
	// only want to do changes if we've actually told the user what's changing)
	if d.HasChange(PropertymapFieldname) {
		obj, err := p.NewObjectConfig(d.Get("yaml").(string))
		if err != nil {
			return err
		}
		_, err = p.Apply(obj, false)
		if err != nil {
			return err
		}
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

	d.Set("name", obj.Metadata.Name)
	d.Set("namespace", obj.Metadata.Namespace)
	d.Set("kind", obj.Kind)
	d.Set("api_version", obj.ApiVersion)
	d.Set(PropertymapFieldname, obj.Properties())

	return results, nil
}

// Find all leaf values in the provided object, returning a map of those values and their paths
func leafValues(pathPrefix string, obj interface{}) map[string]string {
	out := make(map[string]string)

	// exclude paths with these prefixes from the output
	excludedPrefixes := []string{
		"apiVersion",
		"kind",
		"status.",
		"metadata.name",
		"metadata.namespace",
		"metadata.annotations[kubectl.kubernetes.io/last-applied-configuration]",
		"metadata.generation",
		"metadata.ownerReferences",
		"metadata.resourceVersion",
	}
	isExcluded := func(path string) bool {
		for _, prefix := range excludedPrefixes {
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
		return false
	}

	// delimit path parts with '.' (unless the path part contains '.')
	pathFmt := func(parent, current string) string {
		var suffix string
		if strings.Contains(current, ".") {
			suffix = fmt.Sprintf("[%s]", current)
		} else if parent == "" {
			suffix = current
		} else {
			suffix = "." + current
		}
		return parent + suffix
	}

	// traverse: lists, maps,
	// append: scalars
	var walkNode func(string, interface{})
	walkNode = func(path string, node interface{}) {
		if isExcluded(path) {
			// if path is excluded we don't even traverse tree
			return
		}
		switch v := node.(type) {
		case int:
			out[path] = strconv.Itoa(v)
		case bool:
			out[path] = strconv.FormatBool(v)
		case string:
			out[path] = v
		case map[interface{}]interface{}:
			// walk all child paths
			for key, value := range v {
				walkNode(pathFmt(path, key.(string)), value)
			}
		case []interface{}:
			for idx, value := range v {
				walkNode(pathFmt(path, strconv.Itoa(idx)), value)
			}
		case nil:
			// ignore nil values
		default:
			panic(fmt.Sprintf("I don't know about type %T!\n", v))
		}
	}

	walkNode(pathPrefix, obj)

	return out
}
