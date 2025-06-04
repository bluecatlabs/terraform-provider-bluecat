# External Host Record
This resource will create an external host record in Address Manager with a list of linked addresses. The attributes are:

| Attribute     | Required/optional | Description | Example                    |
|---------------| --- | --- |----------------------------|
| configuration | Optional | The Configuration. Creating the External Host record in the default Configuration if doesn't specify | Demo                       |
| view          | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal                   |
| absolute_name | Required | The name of the External Host record. Must be an FQDN. | webapp.bluecatnetworks.com |
| addresses    | Required | A list of IP Addresses to link to the external host record. NOTE: Respective "bluecat_ip_allocation"-s need to be created using terraform to keep data consistency | 45.0.0.4,45.0.0.6 |
| properties    | Optional | Records properties to be passed | comment=My comments        |

## Example of a External Host Record resource

    resource "bluecat_external_host_record" "external_host_record" {
      configuration = "terraform_demo"
      view = "gg"
      absolute_name = "testhost.testy.com"
      addresses = "45.0.0.4,45.0.0.6"
      properties = ""
    }