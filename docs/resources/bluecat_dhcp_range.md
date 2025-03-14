# DHCP Range Record
This resource will create a DHCP Range for the specified IPv4 Network in Address Manager with the specific name supplied and Host Record. The attributes are:

| Attribute     | Required/optional | Description | Example             |
|---------------| --- | --- |---------------------|
| configuration | Optional | The Configuration. Creating the DHCP Range record in the default Configuration if doesn't specify | Demo                |
| network       | Required |  The network address in CIDR format | 10.0.0.0/24         |
| start         | Optional | Start IP of the DHCP Range | 10.0.0.10           |
| end           | Required | End IP of the DHCP Range | 10.0.0.100          |
| name          | Optional | The name of the DHCP Range | DHCP Floor 1        |
| ip_version    | Optional | Options are ivp4 and ipv6. If left blank, ipv4 will be used. | ipv4                |
| template      | Required | The name of the IPv4 Template to apply to this DHCP Range | DHCP_Template_IPv4  |
| properties    | Optional | Records properties to be passed | comment=My comments |


## Example of a DHCP Range Record resource

    resource "bluecat_dhcp_range" "dhcp_range" {
      configuration = "terraform_demo"
      network = "30.0.0.0/24"
      start = "30.0.0.20"
      end = "30.0.0.30"
      properties = ""
      template = "testtemplate"
      ip_version = "ipv4"
      depends_on = [bluecat_ipv4network.net_record]
    }
