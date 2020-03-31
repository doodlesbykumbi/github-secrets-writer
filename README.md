# github-secrets-writer

**github-secrets-writer** is a command-line tool that simplifies the process of creating or updating [Github secrets](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets) by carrying out the encryption for you, and writing the secrets to Github via the API.
> Secrets  are encrypted environment variables created in a repository and can only be used by GitHub Actions

**Table of Contents**
* [What's all the fuss?](#whats-all-the-fuss)
* [Installation](#installation)
  + [Binaries](#binaries)
  + [Via Go](#via-go)
* [Usage](#usage)
  + [Write secrets](#write-secrets)

## What's all the fuss?

GitHub secrets are encrypted using public-key authenticated encryption and the Poly1305 cipher algorithm. The [Github developer documentation](https://developer.github.com/v3/actions/secrets/#create-or-update-a-secret-for-a-repository) suggests carrying out the encryption using [LibSodium](https://libsodium.gitbook.io/doc/bindings_for_other_languages). 

😞 Without **github-secrets-writer**,  a user would have to piece together bits of code to:

1. Encrypt a secret using LibSodium (requires installing dependencies)
2. Write the secret to Github using the API
3. Repeat (1) and (2) for multiple secrets

🚀 **github-secrets-writer** offers the user some convenience by doing all the above for you, without the need to install additional dependencies. It comes as a binary and it uses Go's supplementary [cryptography libraries](https://go.googlesource.com/crypto) (that are interoperable with LibSodium) to carry out the encryption, then writes the secrets to Github using the [Go library for accessing the GitHub API](https://github.com/google/go-github).

## Installation
### Binaries
For installation instructions from binaries please visit the [Releases Page](https://github.com/doodlesbykumbi/github-secrets-writer/releases).

As an example, if you wanted to install v0.1.0 on OSX you would run the following commands in your terminal:
```
# Download the release archive with the binary
wget https://github.com/doodlesbykumbi/github-secrets-writer/releases/download/v0.1.0/github-secrets-writer_darwin_amd64.tar.gz

# Uncompress the release archive to /usr/local/bin. You might need to run this with `sudo`
tar -C /usr/local/bin -zxvf github-secrets-writer_darwin_amd64.tar.gz
```

### Via Go
```console
$ go get -u github.com/doodlesbykumbi/github-secrets-writer
```

## Usage

**github-secrets-writer** uses flags to specify the Github repository and the secrets to write to it. 

NOTE: An OAuth token **must** be provided via the `GITHUB_TOKEN` environment variable, this is used to authenticate to the Github API.  Access tokens require [`repo` scope](https://developer.github.com/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/#available-scopes) for private repos and [`public_repo` scope](https://developer.github.com/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/#available-scopes) for public repos. GitHub Apps must have the `secrets` permission to use the API. Authenticated users must have collaborator access to a repository to create, update, or read secrets.

```console
$ github-secrets-writer -h
Create or update multiple Github secrets sourced from literal values or files.

Key/value pairs representing a secret name and the source of the secret value are provided via the flags --from-file and --from-literal. Depending on the key/value pairs specified a single invocation may carry out zero or more writes to the Github secrets of the repository.

NOTE: An OAuth token **must** be provided via the 'GITHUB_TOKEN' environment variable, this is used to authenticate to the Github API. Access tokens require 'repo' scope for private repos and 'public_repo' scope for public repos. GitHub Apps must have the 'secrets' permission to use the API. Authenticated users must have collaborator access to a repository to create, update, or read secrets.

Usage:
  github-secrets-writer --owner=owner --repo=repo [--from-literal=secretName1=secretValue1] [--from-file=secretName2=/path/to/secretValue2]

Examples:
# Write a single secret from a literal value
github-secrets-writer --owner=owner --repo=repo --from-literal=secretName1=secretValue1

# Write a single secret from a file
github-secrets-writer --owner=owner --repo=repo --from-file=secretName1=secretFilePath

# Write multiple secrets, one from a literal value and one from a file
github-secrets-writer --owner=owner --repo=repo --from-literal=secretName1=secretValue1 --from-file=secretName2=/path/to/secretValue2

Flags:
      --from-file stringArray      specify secret name and literal value pairs e.g. secretname=somevalue (zero or more)
      --from-literal stringArray   specify secret name and source file pairs e.g. secretname=somefile (zero or more)
  -h, --help                       help for github-secrets-writer
      --owner string               owner of the repository e.g. an organisation or user (required)
      --repo string                name of the repository (required)
```

### Write secrets

To write secrets to a repository you must invoke **github-secrets-writer**  with the relevant flags. Below is an example of writing 3 secrets (2 from literal values and 1 from a source file) to the `example-owner/example-repo` repository.

```console
$ GITHUB_TOKEN=... github-secrets-writer \
    --owner example-owner \
    --repo example-repo \
    --from-literal secretName1=secretValue1 \
    --from-literal secretName2=secretValue2 \
    --from-file secretName3=/path/to/secretValue3
Write results:

secretName1: 204 No Content
secretName2: 204 No Content
secretName3: 204 No Content
```

Following the successful writes, you can [use the encrypted secrets in a workflow](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets#using-encrypted-secrets-in-a-workflow) as shown below.

```yaml
steps:
  - name: Hello world action
    with: # Set the secret as an input
      secretName1: ${{ secrets.secretName1 }}
    env: # Or as an environment variable
      secretName2: ${{ secrets.secretName2 }}
      secretName3: ${{ secrets.secretName3 }}
```

