# IPv4 Block Record
This data source allows to retrieve the following information
(attributes) for a IPv4 Block in Address Manager:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Getting the IPvBlock record in the default Configuration if doesn't specify | Demo |
| name | Optional |  The Block name | Server Farm |
| parent_block | Optional | The parent Block. Specified to getting the child Block. THe Block in CIDR format |  |
| cidr | Required | IPv4 Block's CIDR | 10.0.0.0/24 |
| properties | Optional | The properties of the IPv4 Block | attribute=value |


## Example of a IPv4 Block dataset

    data "bluecat_ipv4block" "toronto_block" {
      configuration="terraform_demo"
      parent_block="10.0.0.0/8"
      cidr="10.0.0.0/16"
    }

    output "toronto_block_data" {
      value = data.bluecat_ipv4block.toronto_block
    }

    output "toronto_block_data" {
      value = data.bluecat_ipv4block.toronto_block.id
    }