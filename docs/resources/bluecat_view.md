# Views
This will allow creation or update of a View in Address Manager. The attributes are:

| Attribute     | Required/optional | Description                                                                          | Example             |
|---------------| --- |--------------------------------------------------------------------------------------|---------------------|
| configuration | Optional | The Configuration. Creating the View in the default Configuration if doesn't specify | Demo                |
| name          | Required | The name of view                                                                     | InternalView        |
| properties    | Optional | View's properties to be passed                                                       | comment=My comments |


## Example of a View resource

    resource "bluecat_view" "view_name" {
      configuration = "terraform_demo"
      name = "InternalView"
      properties = ""
    }
