# IPv6 Block Record
This will allow creation or update to an IPv6 Block in Address Manager. The attributes are:

| Attribute | Required/optional | Description                                                                                      | Example            |
| --- | --- |--------------------------------------------------------------------------------------------------|--------------------|
| configuration | Optional | The Configuration. Creating the IPv6Block record in the default Configuration if doesn't specify | Demo               |
| name | Required | The Block name                                                                                   | Server Farm        |
| parent_block | Optional | The parent Block. Specified to creating the child Block. The IPv6 Block in CIDR format           |                    |
| address | Required | IPv6 Block's address                                                                             | 2003:1000::      |
| cidr | Required | IPv6 Block's CIDR                                                                                | 65                 |
| ip_version | Required | Options are ivp4 and ipv6. For this creation, ipv6 should be used                                                                           | 65                 |
| properties | Optional | Records properties to be passed                                                                  | comment=My comments |


## Example of a IPv6 Block resource

    resource "bluecat_block" "block_record" {
      configuration = "terraform_demo"
      name = "block1"
      parent_block = ""
      address = "2003:1000::"
      cidr = "65"
      properties = ""
      ip_version = "ipv6"
      depends_on = [bluecat_configuration.conf_record]
    }
