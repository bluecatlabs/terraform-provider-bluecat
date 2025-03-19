# IPv6 Block Record
This data source allows to retrieve the following information
(attributes) for a IPv6 Block in Address Manager:

| Attribute | Required/optional | Description                                                                                     | Example         |
| --- | --- |-------------------------------------------------------------------------------------------------|-----------------|
| configuration | Optional | The Configuration. Getting the IPv6Block record in the default Configuration if doesn't specify | Demo            |
| name | Optional | The Block name                                                                                  | Server Farm     |
| ip_version | Optional | Default is ipv4. Options are ipv4 or ipv6                | ipv6            |
| parent_block | Optional |  The parent block of the IPv4/IPv6 Block. Specify this field to retrieve the child IPv4/IPv6 Block. The parent_block must be in CIDR format | 2000::/3        |
| cidr | Required | IPv6 Block's CIDR                                                                               | 2003:1000::/65  |
| properties | Optional | The properties of the IPv6 Block                                                                | attribute=value |


## Example of a IPv6 Block dataset

    data "bluecat_ipv6block" "toronto_ipv6_block" {
      configuration="terraform_demo"
      ip_version="ipv6"
      cidr="2003:1000::/65"
    }

    output "toronto_ipv6_block_data" {
      value = data.bluecat_ipv6block.toronto_ipv6_block
    }

    output "toronto_ipv6_block_id" {
      value = data.bluecat_ipv6block.toronto_ipv6_block.id
    }
