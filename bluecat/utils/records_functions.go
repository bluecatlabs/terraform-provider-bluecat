package utils

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
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
