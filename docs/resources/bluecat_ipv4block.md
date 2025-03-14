# IPv4 Block Record
This will allow creation or update to an IPv4 Block in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the IPvBlock record in the default Configuration if doesn't specify | Demo |
| name | Required |  The Block name | Server Farm |
| parent_block | Optional | The parent Block. Specified to creating the child Block. THe Block in CIDR format |  |
| address | Required | IPv4 Block's address | 10.0.0.0 |
| cidr | Required | IPv4 Block's CIDR | 24 |
| ip_version    | Optional | Options are ivp4 and ipv6. If left blank, ipv4 will be used                                                  | ipv4                       |
| properties | Optional | Records properties to be passed | comment=My comments |


## Example of a IPv4 Block resource

    resource "bluecat_block" "block_record" {
      configuration = "terraform_demo"
      name = "block1"
      parent_block = ""
      address = "30.0.0.0"
      cidr = "24"
      ip_version = "ipv4"
      properties = "allowDuplicateHost=enable"
      depends_on = [bluecat_configuration.conf_record]
    }
