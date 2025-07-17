# Generic Record
This resource will create a Generic record in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the Generic record in the default Configuration if doesn't specify | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Optional | The Zone in which you want to update a Generic record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| type | Required | The Type in which you want to create type of Generic record | The following generic record types are available: A6', 'AAAA', 'AFSDB', 'APL', 'CERT', 'DHCID', 'DNAME', 'DS', 'IPSECKEY', 'ISDN','KEY', 'KX', 'LOC', 'MB', 'MG', 'MINFO', 'MR', 'NS', 'NSAP', 'PTR', 'PX', 'RP', 'RT','SINK', 'SPF', 'SSHFP', 'WKS', 'X25'. These records contain name, type, and value information |
| absolute_name | Required | The name of the Generic record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| data | Required | The Data of the Generic record | 10.0.0.12 |
| ttl | Optional | The TTL value. Default is -1  | 300 |
| properties | Optional | Records properties to be passed | comment=My comments |
| to_deploy | Optional | Whether or not to deploy the resource to the BDDS, acceptable true values are yes/Yes true/True | yes |

## Example of a Generic Record resource

    resource "bluecat_generic_record" "generic_record" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      type = "NS"
      absolute_name = "test_NS"
      data = "text"
      ttl = 123
      properties = ""
    }