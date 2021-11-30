# TXT Record
This resource will create a TXT record in Address Manager with the specific data in Address and Host Record.  The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the TXT record in the default Configuration if doesn't specify | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Optional | The Zone in which you want to update a TXT record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| absolute_name | Required | The name of the TXT record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| text | Required | The text data | 10.0.0.0/24 |
| ttl | Optional | The TTL value. Default is -1  | 300 |
| properties | Optional | Records properties to be passed | comment=My comments |

## Example of a TXT Record resource

    resource "bluecat_txt_record" "txt_record" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      absolute_name = "txt"
      text = "text"
      ttl = 123
      properties = ""
    }