# Requirements
---
- Go version: 1.20
- Terraform version: 1.6.4

# Compile the source code as a Terraform's provider
---
- Go to inside of the project directory `cd terraform`
- Compile the project: `go build`

# Running
---
## 1. Preparing the configuration:
---
### Application configuration:
There's the app.yml file that specifying the application configurations such as logging configurations...
Sample configuration file:
```
logging:
  level: warn
  file: provider_bluecat.log
```

### Provider configuration:
Create a file `main.tf` with the following
``` 
provider "bluecat" {
  server = "127.0.0.1"
  api_version = "1"
  transport = "http"
  port = "5000"
  username = "api_user"
  password = "encryption_password/plain_text_password"
  encrypt_password = true/false
} 
```
encrypt_password: Default is false, to indicate if the password is encrypted

## 2. Preparing the resource:
---
Note: The "depends_on" property in each resource to indicate the plan for actions, so that resources are created and destroyed in the correct order
### Resource Configuration:
```
resource "bluecat_configuration" "conf_record" {
  name = "terraform_demo"
  properties = "description=terraform testing config"
}
```
### Resource IPv4 Block:
```
resource "bluecat_block" "block_record" {
  configuration = "terraform_demo"
  name = "block1"
  parent_block = ""
  address = "30.0.0.0"
  cidr = "24"
  properties = "allowDuplicateHost=enable"
  depends_on = [bluecat_configuration.conf_record]
}
```
### Resource IPv4 Network:
```
resource "bluecat_network" "net_record" {
  configuration = "terraform_demo"
  name = "network1"
  cidr = "30.0.0.0/24"
  gateway = "30.0.0.12"
  reserve_ip = 3
  properties = ""
  depends_on = [bluecat_ipv4block.block_record]
}
```
```
resource "bluecat_network" "next_available_net_record" {
  configuration = "terraform_demo"
  name = "next available network1"
  reserve_ip = 3
  parent_block = "30.0.0.0/24"
  size = 256
  allocated_id = timestamp()
  properties = ""
  depends_on = [bluecat_ipv4block.block_record]
}
```
### Resource IPv4 Address Allocation:
```
resource "bluecat_ip_allocation" "host_allocate" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  name = "testhost"
  network = "30.0.0.0/24"
  ip_address = "30.0.0.22"
  mac_address = "223344556688"
  properties = ""
  depends_on = [bluecat_ipv4network.net_record]
}
```
```
resource "bluecat_ip_allocation" "address_allocate" {
  configuration = "terraform_demo"
  view = "gg"
  zone = ""
  name = "testaddress"
  network = "30.0.0.0/24"
  ip_address = "30.0.0.22"
  mac_address = "223344556688"
  properties = ""
  depends_on = [bluecat_ipv4network.net_record]
}
```
### Resource IPv4 Address Association:
```
resource "bluecat_ip_association" "address_associaion" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  name = "testaddress"
  network = "30.0.0.0/24"
  ip_address = "30.0.0.22"
  mac_address = "223344556688"
  properties = ""
  depends_on = [bluecat_ip_allocation.host_allocate]
}
```
### Resource Host Record:
```
resource "bluecat_host_record" "host_record" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  absolute_name = "testhost"
  ip_address = "30.0.0.124"
  ttl = 123
  properties = ""
  depends_on = [bluecat_ipv4network.net_record]
}
```
### Resource PTR Record:
```
resource "bluecat_ptr_record" "ptr_record" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  name = "host30"
  ip_address = "30.0.0.30"
  ttl = 1
  reverse_record = "True"
  properties = ""
  depends_on = [bluecat_ipv4network.net_record]
}
```
### Resource CNAME Record:
```
resource "bluecat_cname_record" "cname_record" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  absolute_name = "cname2"
  linked_record = "host1.gateway.com"
  ttl = 123
  properties = ""
  depends_on = [bluecat_host_record.host_record]
}
```
### Resource TXT Record:
```
resource "bluecat_txt_record" "txt_record" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  absolute_name = "txt"
  text = "text"
  ttl = 123
  properties = ""
}
```
### Resource Generic Record:
```
resource "bluecat_generic_record" "generic_record" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  type = "NS"
  absolute_name = "test_NS"
  data = "text"
  ttl = 123
  properties = ""
}
```
### Resource DHCP Range:
```
resource "bluecat_dhcp_range" "dhcp_range" {
  configuration = "terraform_demo"
  network = "30.0.0.0/24"
  start = "30.0.0.20"
  end = "30.0.0.30"
  properties = ""
  template = "testtemplate"
  depends_on = [bluecat_ipv4network.net_record]
}
```
### Resource Zone and Sub zone:
```
resource "bluecat_zone" "sub_zone" {
  configuration = "terraform_demo"
  view = "Internal"
  zone = "example.com"
  deployable = "True"
  server_roles = [“primary, server1”, “secondary, server2”]
  properties = ""
}
```
## 3. Preparing the datasource:

### Datasource IPv4 Block:
```
data "bluecat_block" "test_ip4block" {
  configuration = "terraform_demo"
  cidr = "20.0.0.0/24"
}

output "output_block" {
  value = data.bluecat_block.test_ip4block
}
```
### Datasource IPv4 Network:
```
data "bluecat_network" "test_ip4network" {
  configuration = "terraform_demo"
  cidr = "20.0.0.0/24"
}

output "output_network" {
  value = data.bluecat_network.test_ip4network
}
```
### Datasource Host Record:
```
data "bluecat_host_record" "test_record" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  fqdn = "host"
}

output "output_host" {
  value = data.bluecat_host_record.test_record
}
```

### Datasource CNAME Record:
```
data "bluecat_cname_record" "test_cname" {
  configuration = "terraform_demo"
  view = "gg"
  zone = "gateway.com"
  linked_record = "host.gateway.com"
  canonical = "cname"
}

output "output_cname" {
  value = data.bluecat_cname_record.test_cname
}
```
### Datasource Zone and Sub zone:
```
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
```

## 4. Executing the provider:
---
### Initialize the provider:

In case of you're using the local build of the provider, you need to prepare the structure to be able to installing the provider as below. Otherwise, just run the `terraform init`.
- Create the directory to store the providers: <HOME_DIR>/providers
- Create the provider structure for the provider under the directory at step 1: <HOSTNAME>/<NAMESPACE>/<TYPE>/<VERSION>/<PLATFORM>/<PROVIDER_BINARY>. For example: test.com/hashicorp/bluecat/1.0.0/windows_amd64/terraform-provider-bluecat.exe
- Add the block of configuration for the provider
    ```
    terraform {
      required_providers {
        <TYPE> = {
          version = ">= <VERSION>"
          source = "<HOSTNAME>/<NAMESPACE>/<TYPE>"
        }
      }
    }
    ```
    For example:
    ```
    terraform {
      required_providers {
        bluecat = {
          version = ">= 1.0.0"
          source = "test.com/hashicorp/bluecat"
        }
      }
    }
    ```
- Install your provider: `terraform init -plugin-dir=<HOME_DIR>/providers`

### Checking out the plan
`terraform plan`

### Adding/updating resources as the plan
If your configuration resources need to be created in order that you have written them in main.tf file:
`terraform apply -parallelism=1`
If the order of resources' creation is not important use:
`terraform apply`
Option for automatic creation without getting the prompt for approving the creation:
`terraform apply -auto-approve`

### Removing resources
`terraform destroy`

### Importing resources
#### 1. Creating imports.tf configuration
Create a new file called imports.tf and define the resource blocks that you want to import.

Resource import definition:
```
import {
  to = <RESOURCE_TYPE>.<RESOURCE_CUSTOM_NAME>
  id = "<RESOURCE_ID>"
}
```
Example of importing two blocks and one network:
```
import {
  to = bluecat_ipv4block.block_record_import_1
  id = "2.0.0.0/8"
}

import {
  to = bluecat_ipv4block.block_record_import_2
  id = "2.2.0.0/16"
}

import {
  to = bluecat_ipv4network.network_record_import_1
  id = "2.2.2.0/24"
}
```
Notice: You need to check if every resource specified in the blocks above exists as a real object in BAM.

#### 2. Generating configuration for resources that you want to import
`terraform plan -generate-config-out=generated_resources.tf` \
Review the generated configuration to see if everything is as it is expected for a certain type of the resource.

#### 3. Apply generated configuration to import your infrastructure
`terraform apply -auto-approve` \
You can use only `apply` without `-auto-approve` if you want to revise once again if everything is ok and manually type `yes` to approve that you want to import resources.

Notice: When you want to create additional resources that you want to import you can:
1. Leave the generated_resources.tf file and generate newly one with different name \
or
2. Use following steps:
   * Copy generated resources from the generated_resources.tf file and copy them to the main config file (main.tf file)
   * Delete the generated_resources.tf file and create the same one with command from Step 2 (`terraform plan -generate-config-out=generated_resources.tf`)

Also, you need to remove imports block if that resources are already imported.
