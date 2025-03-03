# Host Record

This data source allows to retrieve the following information
(attributes) for a Host Record in BlueCat Address Manager:

| Attribute | Required/optional | Description | Example                    |
| --- | --- | --- |----------------------------|
| configuration | Optional | The Configuration. Getting the Host record in the default Configuration if doesn't specify | Demo                       |
| view | Optional | The view which contains the details of the zone. If not provided, record will be queried under default view | Internal                   |
| zone | Optional | The Zone in which the Host record resides. If not provided, the absolute name must be FQDN  | bluecatnetworks.com        |
| fqdn | Required | The name of the Host record. Must be FQDN if the Zone is not provided | webapp.bluecatnetworks.com |
| ip_address | Required | The IP address assigned to the Host record | 10.0.0.12 or 2003:1000:10  |
| ttl | Optional | The TTL value of the host record | 300                        |
| properties | Optional | The properties of the Host Record | attribute=value            |

## Example of a Host Record dataset

    data "bluecat_host_record" "webserver_host" {
      configuration="terraform_demo"
      view="internal"
      zone="bluecatnetworks.com"
      fqdn="webserver"
      ip_address="10.0.0.101"
    }
    
    output "webserver_host_data" {
      value = data.bluecat_host_record.webserver_host
    }

    output "webserver_host_id" {
      value = data.bluecat_host_record.webserver_host.id
    }
    
    output "webserver_host_properties" {
      value = data.bluecat_host_record.webserver_host.properties
    }
