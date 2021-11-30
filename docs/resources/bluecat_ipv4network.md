# IPv4 Netork Record
This will allow creation or update to an IPv4 Network in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the IPv4 Network in the default Configuration if doesn't specify | Demo |
| name | Optional |  The Network name | Server Farm |
| cidr | Required | The network address in CIDR format | 10.0.0.0/24 |
| gateway | Optional | Give the IP you want to reserve for gateway, by default the first IP gets reserved for gateway | 10.0.0.1 |
| reserve_ip | Optional | Reserves the number of IP's for later use | 3 |
| template | Optional | IPv4 Template to apply | NetworkTemplateIPv4 |
| properties | Optional | Records properties to be passed | comment=My comments |


## Example of a IPv4 Network Record resource

    resource "bluecat_ipv4network" "net_record" {
      configuration = "terraform_demo"
      name = "network1"
      cidr = "30.0.0.0/24"
      gateway = "30.0.0.12"
      reserve_ip = 3
      properties = ""
      depends_on = [bluecat_ipv4block.block_record]
    }