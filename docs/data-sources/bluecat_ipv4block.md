# IPv4 Block Record
This data source allows to retrieve the following information
(attributes) for a IPv4 Block in Address Manager:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Getting the IPvBlock record in the default Configuration if doesn't specify | Demo |
| name | Optional |  The Block name | Server Farm |
| ip_version | Optional | If not provided, this will default to ipv4. Options are ipv4 or ipv6|  |
| cidr | Required | IPv4 Block's CIDR | 10.0.0.0/24 |
| properties | Optional | The properties of the IPv4 Block | attribute=value |


## Example of a IPv4 Block dataset

    data "bluecat_ipv4block" "toronto_block" {
      configuration="terraform_demo"
      ip_version="ipv4"
      cidr="10.0.0.0/16"
    }

    output "toronto_block_data" {
      value = data.bluecat_ipv4block.toronto_block
    }

    output "toronto_block_id" {
      value = data.bluecat_ipv4block.toronto_block.id
    }