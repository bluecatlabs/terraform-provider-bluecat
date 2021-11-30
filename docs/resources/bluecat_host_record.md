# Host Record
This resource will create a host record in Address Manager with a specific IP Address. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the Host record in the default Configuration if doesn't specify | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Optional | The Zone in which you want to update a Host record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| absolute_name | Required | The name of the Host record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| ip4_address | Required | The IP address that will be linked to the Host record | 10.0.0.12 |
| ttl | Optional | The TTL value. Default is -1  | 300 |
| properties | Optional | Records properties to be passed | comment=My comments |

## Example of a Host Record resource

    resource "bluecat_host_record" "host_record" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      absolute_name = "testhost"
      ip4_address = "30.0.0.124"
      ttl = 123
      properties = ""
      depends_on = [bluecat_ipv4network.net_record]
    }