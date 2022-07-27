# Zones and Sub zones
This data source allows to retrieve the following information
(attributes) for a zone and sub zone in Address Manager:

| Attribute | Required/optional | Description | Example |
| --- | --- | --- | --- |
| configuration | Optional | The Configuration. Getting the Zone in the default Configuration if doesn't specify | Demo |
| view | Optional |  The view which contains the details of the zone. If not provided, zone will be got under default view | Internal |
| zone | Required | The absolute name of zone or sub zone | example.com |
| deployable | Optional |  Zone's deployable property | True |
| server_roles | Optional |  The list of server roles. The format of each server role will be 'role type, server fqdn' | ["primary, server1", "secondary, server2"] |
| properties | Optional | Zone's properties | comment=My comments |


## Example of a Zone and Sub zone dataset

    data "bluecat_zone" "sub_zone" {
      configuration="terraform_demo"
      view="Internal"
      zone="example.com"
    }

    output "sub_zone_data" {
      value = data.bluecat_zone.sub_zone
    }

    output "id" {
      value = data.bluecat_zone.sub_zone.id
    }

    output "deployable" {
      value = data.bluecat_zone.sub_zone.deployable
    }

    output "server_roles" {
      value = data.bluecat_zone.sub_zone.server_roles
    }
