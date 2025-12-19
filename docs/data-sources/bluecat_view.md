# View
This data source allows to retrieve the following information
(attributes) for a view in Address Manager:

| Attribute     | Required/optional | Description                                                                         | Example             |
|---------------| --- |-------------------------------------------------------------------------------------|---------------------|
| configuration | Optional | The Configuration. Getting the Zone in the default Configuration if doesn't specify | Demo                |
| view          | Required | The name of view                                                                    | InternalView        |
| deployable          | Optional | If the view is to be deployable                                                                    | true                |
| server_roles          | Optional | The list of server roles. The format of each server role will be 'role type, server fqdn                                                                    | ["primary, server1","secondary, server2"]        |
| allowed_property_keys | Optional | The list of properties that should be returned from BAM | ["property_name1", "property_name2"] |


## Example of a View dataset

    data "bluecat_view" "view_name" {
      configuration="terraform_demo"
      name="InternalView"
    }

    output "view_data" {
      value = data.bluecat_view.view_name
    }

    output "id" {
      value = data.bluecat_view.view_name.id
    }
