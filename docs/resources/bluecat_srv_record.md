# SRV Record
This resource will create a SRV record (Alias) in Address Manager with the specific name supplied and Host Record. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the SRV record in the default Configuration if doesn't specify | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Optional | The Zone in which you want to update a SRV record. If not provided, the absolute name must be FQDN ones | bluecatnetworks.com |
| absolute_name | Required | The name of the SRV record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| linked_record | Required | The record name that will be linked to the SRV record | server1.bluecatnetworks.com |
| weight | Required | This is the weight, used to determine which server to connect to if multiple servers have the same priority | 10 |
| port | Required | This is the port number on which the service is listening | 8080 |
| ttl | Optional | The TTL value. Default is -1 | 300 |
| priority | Required | The priority of the record, a lower value is a higher priority | 2 |
| properties | Optional | Records properties to be passed | comment=My comments |
| name | Optional | The name that terraform will use to update the fqdn of the record. *Make sure* to update the absolute name to match the newly updated name after using this parameter | webapp2 |
| to_deploy | Optional | Whether or not to deploy the resource to the BDDS, acceptable true values are yes/Yes true/True | yes |

## Example of a SRV Record resource

    resource "bluecat_srv_record" "test_srv_record" {
    configuration = "Demo"
    view = "Internal"
    zone = "example.com"
    absolute_name = "srv.example.com"
    linked_record = "ns1.example.com"
    weight = 10
    priority = 1
    port = 8080
    }
