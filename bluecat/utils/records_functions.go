package utils

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetStringFromProperties(propertiesMap map[string]string) string {
	var propertiesStr string
	for prop, value := range propertiesMap {
		propertiesStr += fmt.Sprintf("%s=%s|", prop, value)
	}
	return propertiesStr
}

// RemoveImmutableProperties will remove immutable properties for the record
func RemoveImmutableProperties(properties string, immutableProperties []string) string {
	propertiesMap := entities.GetPropertiesFromString(properties)
	for _, immutableProp := range immutableProperties {
		_, ok := propertiesMap[immutableProp]
		if ok {
			delete(propertiesMap, immutableProp)
		}
	}
	return GetStringFromProperties(propertiesMap)
}

func GetPropertyValue(key, props string) (val string) {
	properties := strings.Split(props, "|")
	for i := 0; i < len(properties); i++ {
		prop := strings.Split(properties[i], "=")
		if prop[0] == key {
			val = prop[1]
			return
		}
	}
	return
}

func ParseDeploymentValue(deploymentString string) (deploy bool) {
	trueValues := []string{"Yes", "yes", "True", "true"}
	return slices.Contains(trueValues, deploymentString)
}

func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}

	// If your client exposes a status code, prefer that:
	type hasStatusCode interface{ StatusCode() int }
	var sc hasStatusCode
	if errors.As(err, &sc) && sc.StatusCode() == http.StatusNotFound {
		return true
	}

	// Fallbacks: exact phrases, then strict 3-digit code extraction
	msg := err.Error()
	return strings.Contains(strings.ToUpper(msg), "404 NOT FOUND")
}

func FilterDataSouceProperties(d *schema.ResourceData, bamProps map[string]string) map[string]string {
	// If user supplied allowed_property_keys, build an allow list and filter
	if v, ok := d.GetOk("allowed_property_keys"); ok {
		allowed := make(map[string]struct{})

		for _, key := range v.(*schema.Set).List() {
			k := strings.TrimSpace(key.(string))
			if k != "" {
				allowed[k] = struct{}{}
			}
		}

		// Build filtered map
		filtered := make(map[string]string)
		for k := range allowed {
			if v, exists := bamProps[k]; exists {
				filtered[k] = v
			}
		}

		return filtered
	}

	// If no allow-list provided, return bamProps unchanged
	return bamProps
}

func FilterProperties(bamProps, cfgProps map[string]string) map[string]string {
	filteredProperties := make(map[string]string, len(cfgProps))

	for key := range cfgProps {
		if value, ok := bamProps[key]; ok {
			filteredProperties[key] = value
		}
	}

	return filteredProperties
}

// parse "a=1|b=2" -> map[string]string{"a":"1","b":"2"}
func ParseProperties(s string) map[string]string {
	out := map[string]string{}
	if s == "" {
		return out
	}
	for _, pair := range strings.Split(s, "|") {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		k := strings.TrimSpace(kv[0])
		v := ""
		if len(kv) == 2 {
			v = kv[1]
		}
		if k != "" {
			out[k] = v
		}
	}
	return out
}

// join map -> "a=1|b=2" with stable key order
func JoinProperties(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, m[k]))
	}
	return strings.Join(parts, "|")
}
