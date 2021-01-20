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

## References

Related issues that made this implementation needed:
- https://github.com/aws/aws-cli/issues/5636
- https://github.com/docker/docker-credential-helpers/issues/190
- https://github.com/awslabs/amazon-ecr-credential-helper/issues/229
