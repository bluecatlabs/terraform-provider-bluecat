# IP Association Record
This resource is for mapping an IP Address to a host record in Address Manager. The attributes are:

| Attribute | Required/optional | Description                                                                                                 | Example |
| --- | --- |-------------------------------------------------------------------------------------------------------------| --- |
| configuration | Optional | The Configuration. Creating the record in the default Configuration if doesn't specify                      | Demo |
| view | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Optional | The Zone in which you want to update the record. If not provided, the absolute name must be FQDN ones       | bluecatnetworks.com |
| name | Required | The name of the record. Must be FQDN if the Zone is not provided                                            | webapp.bluecatnetworks.com |
| network | Required | The Network address in CIDR format                                                                          | 10.0.0.0/24 |
| ip_address | Required | The IPv4/IPv6 IP Address                                                                                    | 10.0.0.12 |
| ip_version    | Optional | Options are ivp4 and ipv6. If left blank, ipv4 will be used                                                 | ipv4                       |
| mac_address | Required | The MAC address                                                                                             | 11:22:33:44:55:66 |
| properties | Optional | Records properties to be passed                                                                             | comment=My comments |

## Example of an IP Association resource

    resource "bluecat_ip_association" "address_associaion" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      name = "testaddress"
      network = "30.0.0.0/24"
      ip_address = "30.0.0.22"
      mac_address = "223344556688"
      properties = ""
      depends_on = [bluecat_ip_allocation.host_allocate]
    }
