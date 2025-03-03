# IPv4 Netork Record
This will allow creation or update to an IPv4 Network in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the IPv4 Network in the default Configuration if doesn't specify | Demo |
| name | Optional |  The Network name | Server Farm |
| cidr | Optional | The network address in CIDR format. If not provided, the next available network will be created | 10.0.0.0/24 |
| gateway | Optional | Give the IP you want to reserve for gateway, by default the first IP gets reserved for gateway. Can be set only when creating specified network | 10.0.0.1 |
| reserve_ip | Optional | Reserves the number of IP's for later use | 3 |
| template | Optional | IPv4 Template to apply | NetworkTemplateIPv4 |
| parent_block | Optional | The parent block of the network in CIDR format. Required if create next available network | 30.0.0.0/24 |
| size | Optional | The size of the network expressed in the power of 2. Required if create next available network | 256 |
| allocated_id | Optional | The allocated id of the next available network. Required if create next available network | timestamp() |
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
    
    resource "bluecat_ipv4network" "next_available_net_record" {
      configuration = "terraform_demo"
      name = "next available network1"
      reserve_ip = 3
      parent_block = "30.0.0.0/24"
      size = 256
      allocated_id = timestamp()
      properties = ""
      depends_on = [bluecat_ipv4block.block_record]
    }