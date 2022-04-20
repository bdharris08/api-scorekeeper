# api-scorekeeper
An example api deployment for the scorekeeper library.
This is a practice project.

### Local testing
There is a local testing stack available via docker compose.
It will stand up the api container and a postgres container, then run db migration scripts against postgres.
```sh
docker compose up -d
```

Some example curl commands:
```sh
curl --location --request GET 'http://localhost:3000/hello'
```

```sh
curl --location --request POST 'localhost:3000/scores/trial' --header 'Content-Type: application/json' --data-raw '{"action": "hop","time": 100}'

curl --location --request POST 'localhost:3000/scores/trial' --header 'Content-Type: application/json' --data-raw '{"action": "skip","time": 100}'

curl --location --request POST 'localhost:3000/scores/trial' --header 'Content-Type: application/json' --data-raw '{"action": "jump","time": 100}'

curl --location --request GET 'localhost:3000/scores/trial/average'
```

You should then see some result like 
```json
[
  {
    "action":"hop", 
    "avg": 100
  },
  {
    "action":"skip", 
    "avg": 100
  },
  {
    "action":"jump", 
    "avg": 100
  },
]
```

### CICD
Code - Github actions
- see [test.yml](./.github/workflows/test.yml)
- On pull requests to `main` branch,
  - make sure the go code passes `vet`, `build`, and `test`
- On pushes to `main`,
  - Run the same checks, then
  - tag this version
  - build and push the image to AWS ECR

Infrastructure - Terraform
- see [terraform.yml](./.github/workflows/terraform.yml)
- uses Terraform Cloud for remote backend and configuration
- `terraform.tfvars` shouldn't be committed, so they are stored as workspace variables in terraform cloud.
- On pull requests to `main` branch,
  - check format and validate terraform
  - output the plan in the PR so it can be seen by reviewers
- On pushes to `main`,
  - Run the same checks, then
  - if plan is successful, run apply


### Infrastructure
Infrastructure is defined in [/terraform](./terraform/).

##### Backend
The backend is provided by Terraform Cloud.
This could easily be swapped out for a state file in a s3 bucket, but for this project I decided to use Terraform Cloud for the ease of setting up CICD.

##### Architecture
It uses https://github.com/turnerlabs/terraform-ecs-fargate-apigateway as a starting point. Full suite of changes necessary can be found in these PRs:
- https://github.com/bdharris08/api-scorekeeper/pull/14
- https://github.com/bdharris08/api-scorekeeper/pull/15

Architecture diagram: TODO




