# IPv6 Network Record
This will allow creation or update to an IPv6 Network in Address Manager. The attributes are:

| Attribute | Required/optional | Description                                                                                                                                     | Example           |
| --- | --- |-------------------------------------------------------------------------------------------------------------------------------------------------|-------------------|
| configuration | Optional | The Configuration. Creating the IPv6 Network in the default Configuration if doesn't specify                                                    | Demo              |
| name | Optional | The Network name                                                                                                                                | Server Farm       |
| cidr | Optional | The network address in CIDR format. If not provided, the next available network will be created                                                 | 2003:1000::/65    |
| template | Optional | IPv4 Template to apply                                                                                                                          | NetworkTemplateIPv6 |
| parent_block | Optional | The parent block of the network in CIDR format. Required if create next available network                                                       | 2003:1000::/64    |
| properties | Optional | Records properties to be passed                                                                                                                 | comment=My comments |
| ip_version | Optional | Options are ivp4 and ipv6. For this creation, ipv6 should be used                                                    | ipv6              |



## Example of a IPv6 Network Record resource

    resource "bluecat_ipv6network" "net_record" {
      configuration = "terraform_demo"
      name = "network1"
      cidr = "2003:1000::/65"
      ip_version = "ipv6"
      properties = ""
      depends_on = [bluecat_ipv6network.block_record]
    }
