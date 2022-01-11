# Design notes

## APIs compatible with renovate

The [protocol specified by hashicorp documentation](https://www.terraform.io/internals/module-registry-protocol) allows to implement a module registry compatible with `terraform` client.

But this API is not sufficient for `renovate`. Actually, `renovate` uses one specific API from public registry https://registry.terraform.io which is the  [`Latest Version for a Specific Module Provider`](https://www.terraform.io/registry/api-docs#latest-version-for-a-specific-module-provider).

This is used by `renovate` to compare the latest version available in registry and the version from the analysed code.

The important fields from the response are :
- `version` : returns the latest version in the registry
- `source` : returns the URL of the code repository where `renovate` can fetch the CHANGELOG of the module and display the difference between the versions

You will find more details on the APIs in [./api/README.md]

## OpenAPI v3

We have chosen to use OpenAPI v3 format to specify the APIs and the responses (`api/modules.yaml`). The [`oapi-codegen`](https://github.com/deepmap/oapi-codegen) tool helps generating the server code based on the API specs.

## Server configuration

We chose to use only envars to configure the server. This choice is due to the fact that we run the image as a pod in kubernetes. The configuration with envar is pretty simple, so no need to complexify the code with config files or command flags.

## Code Structure

### cmd

The `cmd` folder contains the startup of the HTTP server. We use `chi` which is a simple server. Other http servers are availble with [`oapi-codegen`](https://github.com/deepmap/oapi-codegen#registering-handlers)

The `cmd` also configures the storage backend to use. Currently, we have two backends `gcp` for a GCS bucket and `fake` that is only used as a test backend.

### pkg/backends

A `backend` implements an interface that is to list the mdoules, write or read the file from the backend.

### pkg/modules

This package implements the API, ie the functions that are refered as `operationId` in `api/modules.yaml` specification. The `oapi-codegen` generated the boilerplate into this folder.

## Usage of GCS object metadata

We use the [GCS object metadata](https://cloud.google.com/storage/docs/metadata#custom-metadata) to store and fetch the information about a module. Actually, only the download API needs to access the stored file. Almost all the APIs simply return some formatted information.

GCS API offers a way to return only selected metadata, which makes an efficient way to return information about a large number of modules.

The metadata we use is actually quite small : we only need the version of the module (`x-module-version`) used to list the available versions and the URL of the code source (`x-module-source`) used only by `renovate` for the CHANGELOG analysis.

## APIs analysis

See also [api/README.md]

### Discovery API

This API is used by the terraform client. It is required for the HTTP server to be recognized by the terraform client as a valid registry.
See more in https://terraform.io/internals/remote-service-discovery#discovery-process

The discovery API returns the actual location of the registry. For a module registry, it expects the key `modules.v1`. In our case, the discovery API is hosted in the server as the registry itself, so it returns a path (that can be modified with `MODULE_PATH` envar if needed).

### Download a module workflow

When the terraform client detects a module in the code, it will not download the module directly. Actually, it will call the `list versions` API (ie **GET** `/{namespace}/{name}/{system}`) to check if the version exists.

Once the expected version is found in the returned versions, it will call the `download` API (ie **GET** `/{namespace}/{name}/{system}/{version}/download`) which actually does not return the expected tarball. This API returns a link in the header `X-Terraform-Get` that will actually download the file. This header could point at any URL. In our case, we use an API **GET** `/{namespace}/{name}/{system}/{version}/archive.tgz`
The reason why is to benefit from the capability of terrafiorm client to recognize the extension of the downloaded archive as described in https://www.terraform.io/language/modules/sources#fetching-archives-over-http

### Upload a module

The API is not specifically defined in terraform documentation and is more an helper for the CI jobs.

In our case, we have chose n to implement a **POST** API that will store the content of a `tgz` archive in a specific structure (the file will even be renamed). This allows when listing the versions for a given modules to use the path filter capability from GCP API.

The `version` that will be stored as a metadata is retrieved from the API parameter.

The URL of the code repository is optionnally passed as a header `module-source`

We have added an existence check with the envar `OVERWRITE` : if set to `0`, the upload API will fail if the module version already exists.

### Renovate API

Renovate only needs the information about the latest available version. To do so, we use the `x-module-version` object metadata (written by the upload API) to fetch all the available versions (same as list versions API), we sort the list in descending order, so that the first item is the latest version.

If the metadata `x-module-source` is defined, it will be returned in the response, and will be used by renovate to fetch the CHANGELOG.md content.

