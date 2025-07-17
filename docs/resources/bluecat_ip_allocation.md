# IP Allocation Record
This resource will allow the allocation of an IP Address (or next available) from a network while creating a host record. The attributes are:

| Attribute     | Required/optional | Description                                                                                                 | Example                    |
|---------------| --- |-------------------------------------------------------------------------------------------------------------|----------------------------|
| configuration | Optional | The Configuration. Creating the record in the default Configuration if doesn't specify                      | Demo                       |
| view          | Optional | The view which contains the details of the zone. If not provided, record will be created under default view | Internal                   |
| zone          | Optional | The Zone in which you want to update the record. If not provided, the absolute name must be FQDN ones       | bluecatnetworks.com        |
| name          | Required | The name of the IP record. Must be FQDN if the Zone is not provided                                         | webapp.bluecatnetworks.com |
| network       | Required | The Network address in CIDR format                                                                          | 10.0.0.0/24                |
| ip_address    | Optional | The IPv4/IPv6 IP Address. If this is not passed, you will get next available IP Address from the network    | 10.0.0.12                  |
| ip_version    | Optional | Options are ivp4 and ipv6. If left blank, ipv4 will be used                                                  | ipv4                       |
| mac_address   | Optional | The MAC address                                                                                             | 11:22:33:44:55:66          |
| action        | Optional | Desired IP4 address state: MAKE_STATIC / MAKE_RESERVED / MAKE_DHCP_RESERVED                                 | MAKE_STATIC                |
| template      | Optional | IPv4 Template which you want to assign                                                                      | ipTemplateIPv4             |
| properties    | Optional | Records properties to be passed                                                                             | comment=My comments        |
| to_deploy | Optional | Whether or not to deploy the resource to the BDDS, acceptable true values are yes/Yes true/True | yes |

## Example of an IP Allocation resource

    resource "bluecat_ip_allocation" "host_allocate" {
      configuration = "terraform_demo"
      view = "gg"
      zone = "gateway.com"
      name = "ipalloc"
      network = "30.0.0.0/24"
      ip_address = "30.0.0.22"
      ip_version = "ipv4"
      mac_address = "223344556688"
      properties = ""
      depends_on = [bluecat_ip4network.net_record]
    }
