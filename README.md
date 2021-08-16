# Requirements
---
- Go version: 1.14
- Terraform version: 0.13.1

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
  password = "api_password"
} 
```

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
resource "bluecat_ipv4block" "block_record" {
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
resource "bluecat_ipv4network" "net_record" {
  configuration = "terraform_demo"
  name = "network1"
  cidr = "30.0.0.0/24"
  gateway = "30.0.0.12"
  reserve_ip = 3
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
  ip4_address = "30.0.0.22"
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
  ip4_address = "30.0.0.22"
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
  ip4_address = "30.0.0.22"
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
  ip4_address = "30.0.0.124"
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
  network = "30.0.0.0/24"
  ip4_address = "30.0.0.30"
  ttl = 1
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

## 3. Executing the provider:
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

### Adding/updating the resource as the plan
`terraform apply`

### Removing the resource
`terraform destroy`
