# IPv4 Network Record
This data source allows to retrieve the following information
(attributes) for a IPv4 Network in Address Manager:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Getting the IPv4 Network in the default Configuration if doesn't specify | Demo |
| name | Optional |  The Network name | Server Farm |
| cidr | Required | The Network address in CIDR format | 10.0.0.0/24 |
| gateway | Optional |  This is the Gateway address for the Network | 10.0.0.1 |
| ip_version | Optional |  Default is ipv4, options are ipv4 or ipv6 | ipv4 |
| properties | Optional | The properties of the IPv4 Network | attribute=value |


## Example of a IPv4 Network Record dataset

    data "bluecat_ipv4network" "toronto_network" {
      configuration="terraform_demo"
      name="Toronto 1st Floor"
      cidr="10.0.0.0/24"
    }

    output "toronto_network_data" {
      value = data.bluecat_ipv4network.toronto_network
    }

    output "toronto_network_cidr" {
      value = data.bluecat_network.toronto_network.cidr
    }
