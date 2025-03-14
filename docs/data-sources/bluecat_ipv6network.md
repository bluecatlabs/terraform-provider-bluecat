# IPv6 Network Record
This data source allows to retrieve the following information
(attributes) for a IPv6 Network in Address Manager:

| Attribute | Required/optional | Description                                                                                 | Example |
| --- | --- |---------------------------------------------------------------------------------------------| -- |
| configuration | Optional | The Configuration. Getting the IPv6 Network in the default Configuration if doesn't specify | Demo |
| name | Optional | The Network name                                                                            | Server Farm |
| cidr | Required | The Network address in CIDR format                                                          | 2003:1000::/65 |
| ip_version | Optional |  Default is ipv4, options are ipv4 or ipv6 | ipv6 |
| properties | Optional | The properties of the IPv6 Network                                                          | attribute=value |


## Example of a IPv6 Network Record dataset

    data "bluecat_network" "toronto_ipv6_network" {
      configuration="terraform_demo"
      ip_version="ipv6"
      cidr="2003:1000::/65"
    }

    output "toronto_ipv6_network_id" {
      value = data.bluecat_network.toronto_ipv6_network.id
    }

    output "toronto_ipv6_network_cidr" {
      value = data.bluecat_network.toronto_ipv6_network.cidr
    }
