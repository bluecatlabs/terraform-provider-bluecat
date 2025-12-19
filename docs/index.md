# BlueCat Provider for Terraform

The Terraform provider uses BlueCat's REST API version 25.0.0 or above.  

The following is an example contents of a provider configuration file named main.tf:

```
provider "bluecat" {
    server = "127.0.0.1"
    api_version = "1"
    transport = "http"
    port = "5000"
    username = "api_user"
    encrypt_password = "False"
    password = "api_password"
}
```

Where the fields represent the following:
- **server**: the IP address of the BlueCat REST API image.
- **api_version**: the version of the REST API.
- **transport**: the protocol used to access the REST API.
- **port**: the port used to access the REST API.
- **username**: the username of the API user with the correct permissions to access the REST API.
- **encrypt_password**: (optional) True or false option to use encrypted password in "password" field.
- **password**: When encrypt_password is false or not set, contains the password of the API users with the correct permissions to access the REST API.  If encrypt_password=true, then place the filename of the encrypted password as created in BlueCat Gateway 

**Example**: 

```
encrypt_password = true   
password = ".encrypted_password"
```

To encrypt password, log into Gateway and navigate to Administration > Encrypt Password. 

```
Path: customizations/.encrypted_password
Password: user_password_here
```


Once this is complete, you can use the .encrypted_password value in the BlueCat Provider password field.

## Resources

Below are the available resources for the following objectTypes:

-   Configuration - (bluecat_configuration)
-   Block (bluecat_ipv4block/bluecat_ipv6block)
-   Network (bluecat_ipv4network/bluecat_ipv6network)
-   DHCP Range (bluecat_dhcp_range)
-   IP Address (bluecat_ip_allocation, bluecat_ip_association)
-   Host Record (bluecat_host_record)
-   PTR Record (bluecat_ptr_record)
-   CNAME Record (bluecat_cname_record)
-   TXT Record (bluecat_txt_record)
-   SRV Record (bluecat_srv_record)
-   Generic Record (bluecat_generic_record)
-   DNS Zone (bluecat_zone)
-   View (bluecat_view)

## Data Sources

Below are the available BlueCat data sources:

-   Block - IPv4/IPv6 (bluecat_ipv4block/bluecat_ipv6block)
-   Network - IPv4/IPv6 (bluecat_ipv4network/bluecat_ipv6network)
-   Host Record (bluecat_host_record)
-   CNAME Record (bluecat_cname_record)
-   DNS Zone (bluecat_zone)
-   View (bluecat_view)

To filter out which properties should be used within the Terraform infrastructure, pass the optional field "allowed_property_keys" to the datasource object in the form of "allowed_property_keys = ["property1_name", "property2_name",...]"

## Import Capabilities

You can now import existing BlueCat data into Terraform state. The available BlueCat Objects you can import are:

-  Block
-  Network
-  Zone
-  CNAME
-  Generic Record
-  Host Record
-  TXT Record
-  View

## Bluecat Import Recommendation
If the resource block is not included a "terraform plan -generate-config-out=<yourFileNameHere>.tf" must be run to create a configuration file with the resource that is being imported. This file does not specify optional attributes and uses the default value for required attributes where applicable.

Including the resource block gives the user the ability to specify the exact attributes that should be assigned to the resource as it is imported, overwriting any values that are currently assigned to the resource in BAM

Example
#Block
```
import{
    to=bluecat_ipv4block.import_block
    id="10.2.0.0/16"
}

resource "bluecat_ipv4block" "import_block" {
    configuration = "demo"
    name = "10_2/16 Block (import)"
    parent_block = "10.0.0.0/8"
    address = "10.2.0.0"
    cidr = "16"
    ip_version = "ipv4"
    properties = "allowDuplicateHost=enable"
    depends_on = [bluecat_ipv4block.block_10_record, bluecat_ipv4block.block_10_16_record]
}
```

For the latest updates, please see the BlueCat Product Documents.
