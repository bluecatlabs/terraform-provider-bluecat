package utils

import (
	"sort"
	"terraform-provider-bluecat/bluecat/entities"
)

const (
	// DeploymentOptionLookupServerID tells the API to resolve the option when the
	// assigned server is not known ahead of time.
	DeploymentOptionLookupServerID = -1
	// DeploymentOptionAllServersID targets the "All Servers" assignment used by
	// creates and deletes in this provider.
	DeploymentOptionAllServersID = 0
)

// CreateDeploymentOptions creates each configured deployment option on the given
// target object using the provider's "All Servers" assignment convention.
func CreateDeploymentOptions(objMgr *ObjectManager, target entities.DeploymentOption, deploymentOptions map[string]string) error {
	for optionName, optionValue := range deploymentOptions {
		target.Name = optionName
		target.Value = optionValue
		target.ServerID = DeploymentOptionAllServersID
		if _, err := objMgr.CreateDeploymentOption(target); err != nil {
			return err
		}
	}
	return nil
}

// ReadDeploymentOptions reads only the option names already present in config or
// state, because the BlueCat API exposes item lookups rather than list reads.
func ReadDeploymentOptions(objMgr *ObjectManager, target entities.DeploymentOption, configured map[string]string) (map[string]string, error) {
	optionNames := GetSortedMapKeys(configured)
	deploymentOptions := make(map[string]string, len(optionNames))
	for _, optionName := range optionNames {
		target.Name = optionName
		target.ServerID = DeploymentOptionLookupServerID
		option, err := objMgr.GetDeploymentOption(target)
		if err != nil {
			return nil, err
		}
		deploymentOptions[optionName] = option.Value
	}
	return deploymentOptions, nil
}

// UpdateDeploymentOptionsForTarget diffs the old and new Terraform maps and
// applies the necessary create, replace, and delete calls for one target object.
func UpdateDeploymentOptionsForTarget(objMgr *ObjectManager, target entities.DeploymentOption, currentRaw interface{}, newRaw interface{}) error {
	newDeploymentOptions := ExpandStringMap(newRaw)
	if newDeploymentOptions == nil {
		newDeploymentOptions = make(map[string]string)
	}
	currentDeploymentOptions := ExpandStringMap(currentRaw)
	if currentDeploymentOptions == nil {
		currentDeploymentOptions = make(map[string]string)
	}

	for optionName, optionValue := range currentDeploymentOptions {
		if _, ok := newDeploymentOptions[optionName]; !ok {
			target.Name = optionName
			target.Value = optionValue
			target.ServerID = DeploymentOptionAllServersID
			if _, err := objMgr.DeleteDeploymentOption(target); err != nil {
				return err
			}
		}
	}

	for optionName, optionValue := range newDeploymentOptions {
		currentValue, ok := currentDeploymentOptions[optionName]
		target.Name = optionName
		target.Value = optionValue
		if ok {
			if currentValue == optionValue {
				continue
			}
			target.ServerID = DeploymentOptionAllServersID
			target.Value = currentValue
			if _, err := objMgr.DeleteDeploymentOption(target); err != nil {
				return err
			}
			target.ServerID = DeploymentOptionAllServersID
			target.Value = optionValue
			if _, err := objMgr.CreateDeploymentOption(target); err != nil {
				return err
			}
			continue
		}
		target.ServerID = DeploymentOptionAllServersID
		if _, err := objMgr.CreateDeploymentOption(target); err != nil {
			return err
		}
	}

	return nil
}

// GetSortedMapKeys returns the keys of the given map in stable lexical order.
func GetSortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
