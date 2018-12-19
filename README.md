# Get Secrets

This is a tool for reading various config files from S3 and exposing them to a sub-process as Env-Vars

## Raison d'Etre

The main idea of this is to avoid shipping config files that contain secrets with the
app, or for copying them into place as part of the deployment process; they can be
"deployed" at the last moment before the app is started.

In most cases, other systems can manage this for you (e.g. Puppet or Kubernetes
Secrets) but in cases such as running non-k8s Docker containers, this becomes more
complicated partly because the idealised smaller image size means configuration
management tooling is not really appropriate.

This is a middle-road option that avoids needing to add config management code,
mount special file-systems, or include secrets-delivery client code (e.g. Hashicorp
Vault).

# How it works

Rather than running your app / binary directly, you would instead run `get-secrets`
and pass the path/file of your app plus any command-line arguments _your_ app needs, as
additional arguments for `get-secrets`.

NB: If you want to run `get-secrets` directly, it will print out a series of `export VAR=VAL` lines
that can probably be eval'd directly by most shells; escaping various provided env vals
correctly can be challenging, however.

It will then use certain configurations (see below) to find the secrets files, parse them, and
pass those settings via env-vars to a new app / binary:

- The (current) S3 version walks through all the files found in the S3 bucket + prefix, as specified
by `s3.path` (below)
- Each file is read and parsed according to it's filename extension.
-- e.g. `*.env` is parsed as `VAR=VAL` key pairs, as per 12-factor "dot-env" format
-- `*.json` are JSON-formatted files.
-- Other supported formats are specified by [Viper](https://github.com/spf13/viper#what-is-viper)
- The contents of the files are merged together into a hash-map, in alpha order, with earlier in
  the order getting precedence over later.
- The app / binary is then exec'd, and the above hash-map is passed into it as an env-var list.

## Docker

Whilst primarily designed for Docker containers, `get-secrets` is agnostic to the environment it
runs in.

# Configuration

`get-secrets` uses the [Viper](https://github.com/spf13/viper) library to configure itself.

This is set-up to look for a `.secrets.toml` file in three possible places (in order of
preference):
- in the `$HOME` dir
- in the directory that `$SECRETS_BASE` points to
- in the directory that the `get-secrets` binary runs from

Within that file, it can contain something like the following:

```yaml
## This sets debug-mode for `get-secrets` (only):
## Can also use "$SECRETS_DEBUG"
debug = true
## The runtime and local directory:
## Can also use "$SECRETS_BASE"
base = "RUNTIME_DIR"

[application]
## Really only used for logging. Allows you to give a "friendly" name in logging for _your_ app:
## Can also use "$APPLICATION_NAME"
name = "APP-NAME"
## Really only used for logging. Allows you to give a "friendly" name in logging for the type
## of environment your app will be running in -- e.g. "production" or "development":
## Can also use "$ENVIRONMENT"
environment = "ENV-NAME"

[dotenv]
## Whether `get-secrets` should skip straight to running the new app, rather than downloading
## and processing the secrets files from S3:
## Can also use "$SKIP_SECRETS"
skip = false

[s3]
## Base path on S3 that contains the secrets files:
## Can also use "$SECRETS_S3_DOTENV_PATH"
path = "s3://BUCKET/DIR-PATH"

[logging]
## This is the logging output-format.  Either "text" or "json":
## Can also use "$SECRETS_LOGGING_FORMAT"
format = "text"

[logging.sentry]
## This allows sending errors or failures within `get-secrets` to https://sentry.io
## Can also use "$SECRETS_LOGGING_SENTRY_DSN"
dsn = "https://ID:TOKEN@sentry.io/PROJECT"
```

Note that these all map to optional env vars, this means you can override the above settings
by passing in an env var to `get-secrets` as mentioned in the code comments, above.

# Gotchas

## Root CA Certificates

Root CA Certificates get updated at various times.  If the ones in the Localhost or Docker container
are out-of-date, `get-secrets` may get a `x509: failed to load system roots and no roots provided`
error.

This usually means you need to update the Root CA Certs.

## AWS Authentication

This is around how you pass AWS authentication:  https://docs.aws.amazon.com/sdk-for-go/api/aws/session/

e.g.
- EC2 and ECS expose this access via an http://169.x.x.x address.
- You could use the `~/.aws/credentials` on a local machine
- The `AWS_ACCESS_KEY_ID` and `$AWS_SECRET_ACCESS_KEY` env vars.

In a Docker container, you won't be able to use `~/.aws/credentials` (etc) without specifically
creating or mounting that dir into the container, and on a non-ECS container, you won't have
access to the http://169.x.x.x service.

# Testing your config

There are two ways to test `get-secrets`. One involves building locally, the other involves
building a Docker container.

The first requires Go to be installed locally, and is a little faster, but the latter allows
better control against unexpected differences in your machine c.f. the target machinewhere you
would run the command.

If you run one of the following (changing the appropriate entries, e.g. {ChangeMe}):

## Local Testing

```bash

make clean test all

(
  export \
    AWS_ACCESS_KEY_ID='{ChangeMe}' AWS_SECRET_ACCESS_KEY='{ChangeMe}' \
    AWS_DEFAULT_REGION={ChangeMe} \
    SECRETS_S3_DOTENV_PATH='s3://{BUCKET}/{PREFIX}' \
    APPLICATION_NAME=test -e ENVIRONMENT=envtest \
    SECRETS_DEBUG=1;
  ./bin/github.com/mexisme/get-secrets
)
```

## Docker Testing

```bash
docker build --tag get-secrets .

docker run \
  -e 'AWS_ACCESS_KEY_ID={ChangeMe}' -e 'AWS_SECRET_ACCESS_KEY={ChangeMe}' \
  -e AWS_DEFAULT_REGION={ChangeMe} \
  -e SECRETS_S3_DOTENV_PATH='s3://{BUCKET}/{PREFIX}' \
  -e APPLICATION_NAME=test -e ENVIRONMENT=envtest \
  -e SECRETS_DEBUG=1 \
  -it --rm --entrypoint sh get-secrets:latest

./get-secrets
```

... the `get-secrets` app should print out a list of environment variables, in the format they'll
be sent via the `execve(2)` stdlib function.

Note that `execve(2)` uses a much simpler format for parsing env-vars provided to it:
- it splits each line at the first `=` character
- the portion to the left of the split is the env-var name
- the portion to the right is the contents
-- shell-style quoting and escaping rules *do not* apply when passing from `get-secrets` to
   the app / binary, but they *do* apply when parsing `.env` formatted files.

# Alternatives

- Kubernetes Secrets
- Include code that can read directly from KMS, Hashicorp Vault (et al)
- Inject a file into system at run-time
- Mount a read-only volume containing secrets into your system
- Mount a named-pipe from host into the container that delivers one-shot-access
  secret data (containerised only)
