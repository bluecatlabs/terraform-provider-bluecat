# CNAME Record
This resource will create a CNAME record (Alias) in Address Manager with the specific name supplied and Host Record. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the CNAME record in the default Configuration if doesn't specify | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Optional | The Zone in which you want to update a CNAME record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| absolute_name | Required | The name of the CNAME record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| linked_record | Required | The record name that will be linked to the CNAME record | server1.bluecatnetworks.com |
| ttl | Optional | The TTL value. Default is -1 | 300 |
| properties | Optional | Records properties to be passed | comment=My comments |

## Example of a CNAME Record resource

    resource "bluecat_cname_record" "cname_record" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      absolute_name = "cname2"
      linked_record = "testhost.gateway.com"
      ttl = 123
      properties = ""
      depends_on = [bluecat_host_record.host_record]
    }
