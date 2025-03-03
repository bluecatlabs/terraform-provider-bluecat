# PTR Record
This resource will create a PTR record (reverse record) in Address Manager with the specific IP Address and Host Record.  The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the PTR record in the default Configuration if doesn't specify | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Required | The Zone in which you want to update a PTR record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| name | Required | The name of the host record | webapp |
| ip_address | Required | The IP address that will be created the PTR record for | 10.0.0.12 |
| reverse_record | Required | To create a reverse record for the pass host | True/False |
| ttl | Optional | The TTL value. Default is -1  | 300 |

## Example of a PTR Record resource

    resource "bluecat_ptr_record" "ptr_record" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      name = "host30"
      ip_address = "30.0.0.30"
      reverse_record = True
      ttl = 1
    }
