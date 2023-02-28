# BlueCat Provider for Terraform

The Terraform provider uses BlueCat's REST API version 23.2.1 or above.  

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
-   IPv4 Block (bluecat_ipv4block)
-   IPv4 Network (bluecat_ipv4network)
-   IPv4 DHCP Range (bluecat_dhcp_range)
-   IPv4 IP Address (bluecat_ip_allocation, bluecat_ip_association)
-   Host Record (bluecat_host_record)
-   PTR Record (bluecat_ptr_record)
-   CNAME Record (bluecat_cname_record)
-   TXT Record (bluecat_txt_record)
-   Generic Record (bluecat_generic_record)

## Data Sources

Below are the available BlueCat data sources:

-   IPv4 Block (bluecat_ipv4block)
-   IPv4 Network (bluecat_ipv4network)
-   Host Record (bluecat_host_record)
-   CNAME Record (bluecat_cname_record)

For the latest updates, please see the BlueCat [Terraform Plugin Administration Guide](https://docs.bluecatnetworks.com/r/en-US/BlueCat-Terraform-Plugin-Administration-Guide/21.10.1)
