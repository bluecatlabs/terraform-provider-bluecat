# Configuration Record
 This will create or update a BlueCat Configuration in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| name | Required | The Configuration name | Demo |
| properties | Optional | Records properties to be passed | comment=My comments |


## Example of a Configuration resource
    resource "bluecat_configuration" "conf_record" {
      name = "terraform_demo"
      properties = "description=terraform testing config"
    }
