# docker-credential-aws-sso-ecr
Docker credential store for AWS SSO 

Simple docker credential store for AWS ECR when using AWS SSO profiles.
It *does not store aws credentials* (keys and tokens) *on the machine* but it returns them dinamically to docker every time the auth is needed 

## Requirement

- awscli v2
- windows x64

## Instructions

1. Clone this repository and add the `/bin` folder in the `PATH` environment variable to enable docker to discover and run it.
   - Alternatively, download the binary `/bin/docker-credential-aws-sso-ecr.exe` and place it in a folder present in the `PATH` environment variable
   - If, for security reasons, don't trust running the docker-credential-aws-sso-ecr.exe created follow the build instruction here to create it from source
2. Update `~/.docker/config.json` by adding the `docker-credential-aws-sso-ecr` as credStore for the specific registry like:
   ```json
    {
        "credHelpers": {
        // Important bit
            "<ACCOUNT>.dkr.ecr.<REGION>.amazonaws.com": "aws-sso-ecr"
        // Important bit
        },

        "credStore": "desktop",
        "stackOrchestrator": "swarm"
    }
   ```
3. Start pulling and pushing.

### Multiple profiles for same account

:warning:
In case you have multiple profiles for the same account but with different roles you can specify a default role to be used to select the profile by adding an environment variable `DOCKER_CREDSTORE_AWS_SSO_ECR = ROLENAME` in your machine.

If there are not profile with the default `ROLE_NAME`, the first profile matching the account and region will be selected.


## Build

### Requirements
 - [Golang](https://golang.org/)

### Compile command

From the root of the repo.

```bash
# Windows
go build -o bin/docker-credential-aws-sso-ecr.exe src/docker-credential-aws-sso-ecr.go
```

## References

Related issues that made this implementation needed:
- https://github.com/aws/aws-cli/issues/5636
- https://github.com/docker/docker-credential-helpers/issues/190
- https://github.com/awslabs/amazon-ecr-credential-helper/issues/229
