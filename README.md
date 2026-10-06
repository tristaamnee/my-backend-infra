# my-backend-infra

Infrastructure as code for the FlowerShop backend, written with [Pulumi](https://www.pulumi.com/) in Go. It describes a [Google Cloud Run](https://cloud.google.com/run) service that runs the backend container.

This is a small personal practice project (basic level). The application itself lives in [flowershopbe-v4](https://github.com/tristaamnee/flowershopbe-v4).

## What it creates

- A Cloud Run (v2) service `go-backend-service` in `asia-southeast1`
  - Container image pulled from Google Artifact Registry (image reference comes from configuration)
  - 1 vCPU and 512 MiB memory, container port 8080
  - TCP startup probe on port 8080
  - Environment variables for the app (database, cache, mail, payment and JWT settings)
- An IAM binding that lets anyone invoke the service over HTTPS (`roles/run.invoker` for `allUsers`)
- The service URL as a stack output (`uri`)

The Cloud Run service is created with `pulumi.Protect(true)`, so Pulumi refuses to delete or replace it by accident.

## Secrets

No credentials are stored in this repository. The program reads every secret from Pulumi configuration, and the values are kept in a [Pulumi ESC](https://www.pulumi.com/docs/esc/) environment named `flowershop/dev`. The stack file `Pulumi.dev.yaml` only points to that environment.

Configuration keys the program expects (under `pulumiConfig`, prefixed with the project name `my-backend-infra:`):

- `mongodbUri`
- `jwtSecret`
- `redisAddress`
- `redisPassword`
- `gmailEmail`
- `gmailEmailPassword`
- `payosClientID`
- `payosApiKey`
- `payosChecksum`
- `dockerImageLink`

Secret values are marked with `fn::secret` in the ESC environment, so they are masked in CLI output and encrypted in the stack state.

## Usage

Prerequisites: [Pulumi CLI](https://www.pulumi.com/docs/iac/download-install/), Go, a Pulumi Cloud account, and Google Cloud credentials with access to the target project.

1. Create the ESC environment `flowershop/dev` and fill in the keys above (`pulumi env edit flowershop/dev`).
2. Log in and select the stack, then preview the changes:

```
pulumi login
pulumi stack select dev
pulumi preview
```

`pulumi preview` only shows the planned changes. Applying them with `pulumi up` is a separate, deliberate step.

## Related

- Backend application: [flowershopbe-v4](https://github.com/tristaamnee/flowershopbe-v4)
