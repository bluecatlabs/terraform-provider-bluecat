# CNAME Record
This data source allows to retrieve the following information
(attributes) for a CNAME record (Alias) in BlueCat Address Manager:


| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. If not passed, the CNAME record will be queried in the default Configuration | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be queried under default view | Internal |
| zone | Optional | The Zone in which you want to get the CNAME record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| canonical | Required | The name of the CNAME record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| linked_record | Required | The record name that's linked to the CNAME record | server1.bluecatnetworks.com |
| ttl | Optional | The TTL value. | 300 |
| allowed_property_keys | Optional | The list of properties that should be returned from BAM | ["property_name1", "property_name2"] |

## Example of CNAME Record dataset

    data "bluecat_cname_record" "aliasname" {
      configuration="terraform_demo"
      view="internal"
      zone="bluecatnetworks.com"
      canonical="aliasname"
      linked_record="webserver"
    }
    
    output "aliasname_data" {
      value = data.bluecat_cname_record.aliasname
    }

    output "aliasname_id" {
      value = data.bluecat_cname_record.aliasname.id
    }

    output "aliasname_props" {
      value = data.bluecat_cname_record.aliasname.properties
    }