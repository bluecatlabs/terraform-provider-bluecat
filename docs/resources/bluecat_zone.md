# Zones and Sub zones
This will allow creation or update of a Zone or Sub zone in Address Manager. The attributes are:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Creating the Zone in the default Configuration if doesn't specify | Demo |
| view | Optional |  The view which contains the details of the zone. If not provided, record will be created under default view | Internal |
| zone | Required | The absolute name of zone or sub zone | example.com |
| deployable | Optional | The deployable flag is False by default and is optional. To make the zone deployable, set the deployable flag to True | True |
| server_roles | Optional | The list of server roles. The format of each server role will be 'role type, server fqdn'. Options for this: FORWARDER, PRIMARY, PRIMARY_HIDDEN, NONE, RECURSION, SECONDARY, SECONDARY_STEALTH, STUB| [“primary, bdds1.example.com", “secondary, bdds2.example.com"] |
| properties | Optional | Zone's properties to be passed | comment=My comments |


## Example of a Zone or Sub zone resource

    resource "bluecat_zone" "sub_zone" {
      configuration = "terraform_demo"
      view = "Internal"
      zone = "example.com"
      deployable = "True"
      server_roles = [“primary, server1”, “secondary, server2”]
      properties = ""
    }
