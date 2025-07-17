package bluecat

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceImporter(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
	// Here we use a function to parse the import ID (like the example above) to simplify our logic
	address, cidr, err := resourceServiceParseId(d.Id())

	if err != nil {
		return nil, err
	}

	d.Set("address", address)
	d.Set("cidr", cidr)
	d.SetId(fmt.Sprintf("%s/%s", address, cidr))

	return []*schema.ResourceData{d}, nil
}

func recordImporter(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
	// Here we use a function to parse the import ID (like the example above) to simplify our logic
	zoneName, recordName, err := recordParseId(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("zone", zoneName)
	d.Set("absolute_name", recordName)
	absoluteName := fmt.Sprintf("%s.%s", recordName, zoneName)
	d.Set("absoluteName", absoluteName)
	d.SetId(absoluteName)

	return []*schema.ResourceData{d}, nil
}

func zoneImporter(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	zoneName := d.Id()
	d.Set("zone", zoneName)

	return []*schema.ResourceData{d}, nil
}

func viewImporter(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	viewName := d.Id()
	d.Set("zone", viewName)

	return []*schema.ResourceData{d}, nil
}

// Get the properties values and get only the value for the propertyName
func parseRecordPropertyValue(props string, propertyName string) string {
	// ip4_address example:
	// record.Properties = absoluteName=test-host.example.com|addresses=2.2.2.2|addressIds=5340963|reverseRecord=true|
	// linkedRecord example:
	// record.LinkedRecord = ttl=123|absoluteName=test-cname-2.example.com|linkedRecordName=alloc_1.example.com|
	attrs := strings.Split(props, "|")
	for _, attrStr := range attrs {
		attr := strings.Split(attrStr, "=")
		if len(attr) == 2 {
			attrName := attr[0]
			attrVal := attr[1]
			if attrName == propertyName {
				return attrVal
			}
		}
	}
	return ""
}

func getAbsoluteName(d *schema.ResourceData) (string, error) {
	var absoluteName string
	if d.Id() != "" {
		zoneName, recordName, err := recordParseId(d.Id())
		if err != nil {
			log.Debug(err)
			return "", err
		}
		absoluteName = fmt.Sprintf("%s.%s", recordName, zoneName)
	} else {
		absoluteName = d.Get("absolute_name").(string)
	}
	return absoluteName, nil
}
