# Base Terraform

Creates the foundational infrastructure for the application's infrastructure.
These Terraform files will create a [remote state][state] and a [registry][ecr].
Most other infrastructure pieces will be created within the `environments` directory.


## Included Files

+ `main.tf`  
The main entry point for the Terraform run.

+ `variables.tf`  
Common variables to use in various Terraform files.

+ `ecr.tf`  
Creates an AWS [Elastic Container Registry (ECR)][ecr] for the application.


## Usage

Typically, the base Terraform will only need to be run once, and then should only
need changes very infrequently.

```
# Sets up Terraform to run
$ terraform init

# Executes the Terraform run
$ terraform apply
```


## Variables

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| app | Name of the application. This value should usually match the application tag below. | string |  | yes |
| image_tag_mutability | The tag mutability setting for the repository. | string | IMMUTABLE |  |
| region | The AWS region to use for the bucket and registry; typically `us-east-1`. Other possible values: `us-east-2`, `us-west-1`, or `us-west-2`. <br>Currently, Fargate is only available in `us-east-1`. | string | `us-east-1` |  |
| tags | A map of the tags to apply to various resources. The required tags are: <br>+ `application`, name of the app <br>+ `environment`, the environment being created <br>+ `team`, team responsible for the application <br>+ `contact-email`, contact email for the _team_ <br>+ `customer`, who the application was create for | map | `<map>` | yes |


## Additional Information

+ [Terraform remote state][state]

+ [Terraform providers][provider]

+ [AWS ECR][ecr]



[state]: https://www.terraform.io/docs/state/remote.html
[provider]: https://www.terraform.io/docs/providers/
[ecr]: https://aws.amazon.com/ecr/
