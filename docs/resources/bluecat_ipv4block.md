# IPv4 Block Record
This will allow creation or update to an IPv4 Block in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The BAM configuration. If unspecified, the IPvBlock record is created in the default configuration | Demo |
| name | Required |  The name of the Block | Server Farm IPv4 |
| parent_block | Optional | The parent block (CIDR format). Required for creating a child or next available block | 10.1.0.0/16 |
| address | Optional | IPv4 Block address. Required when creating a specified block | 10.0.0.0 |
| cidr | Optional | Block size as a power of 2. Required for next available block | 24 |
| size | Optional | The size of the block expressed in the power of 2. Required for next available block creation | 256 |
| allocated_id | Optional | Allocated ID of the next available block. Recommended for stable retrieval | timestamp() |
| ip_version    | Optional | Options: ipv4 or ipv6. Defaults to ipv4 if unspecified| ipv4 |
| properties | Optional | Record properties to pass | attribute=value |


## Example of a specified IPv4 Block resource

    resource "bluecat_ipv4block" "block_record" {
      configuration = "terraform_demo"
      name = "block1"
      parent_block = ""
      address = "30.0.0.0"
      cidr = "24"
      ip_version = "ipv4"
      properties = "allowDuplicateHost=enable"
      depends_on = [bluecat_configuration.conf_record]
    }

## Example of a next available IPv4 Block resource

    resource "bluecat_ipv4block" "next_available_block_record" {
      configuration = "terraform_demo"
      name = "next available block1"
      parent_block = "30.0.0.0/16"
      size = "256"
      allocated_id = timestamp()
      ip_version = "ipv4"
      properties = "allowDuplicateHost=enable"
      depends_on = [bluecat_configuration.conf_record]
    }
