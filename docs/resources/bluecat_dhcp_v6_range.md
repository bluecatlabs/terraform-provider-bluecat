# DHCPv6 Range Record
This resource will create a DHCPv6 Range for the specified IPv6 Network in Address Manager with the specific name supplied and Host Record. The attributes are:

| Attribute | Required/optional | Description | Example        |
| --- | --- | --- |----------------|
| configuration | Optional | The Configuration. Creating the DHCP Range record in the default Configuration if doesn't specify | Demo           |
| network | Required |  The network address in CIDR format | 2003:1000::/64 |
| start | Required | Start IP of the DHCP Range | 2003:1000::1   |
| end | Required | End IP of the DHCP Range | 2003:1000::100 |
| properties | Optional | Records properties to be passed | key=value      |
| ip_version | Optional | Options are ivp4 and ipv6. For this creation, ipv6 should be used                                                    | ipv6              |


## Example of a DHCPv6 Range Record resource

    resource "bluecat_dhcp_range" "dhcp_v6_range" {
      configuration = "terraform_demo"
      network = "2003:1000::/64"
      start = "2003:1000::1"
      end = "2003:1000::100"
      ip_version = "ipv6"
      properties = ""
      depends_on = [bluecat_ipv6network.net_record]
    }
