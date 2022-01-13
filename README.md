# tfmodules

We use a lot of Terraform to provision resource across many different providers. As a consequence, we create a lot of modules, helping reusing the code across many projects so as to reduce code complexity. But it is not always easy to maintain the dependencies up to date when the modules evolve.

To ease depencies maintenance, we use [Renovate](https://github.com/renovatebot/renovate) that is able to detect many types of dependencies, among them Terraform modules. Renovate is able to compare the module versions in code and to open a Pull Request/Merge Request with the CHANGELOG difference to ease the analysis by the developer.

But it detects only the modules from Terraform public registries. We haven't found a private module registry that is compatible with Renovate. Therefore, we implemented ours.

`tfmodules` is a registry for Terraform modules written in `go`, compatible with [Renovate](https://github.com/renovatebot/renovate) and using GCS bucket as a storage backend. It can be deployed as a k8s pod (not only GKE, any provider actually), as a cloudrun service or as a standalone container.

We tried to make the code open enough for ianyone integrating to another backend than GCS bucket.

## How to install the server

The image can run almost anywhere. It needs to access a preexisting GCS bucket and run as an IAM service account that has the permissions to access this bucket. The application uses the GCS API to communicate with the bucket. There is no need to mount the bucket as a volume.

### Important : the server has to be served as HTTPS

Terraform client requires the registry to be served with a valid TLS certificate.

### How to deploy in GKE with workload identity and Config Connector

The example [examples/gke] shows how to deploy the registry using config connector to provision the bucket and the workload identoty to grant permisions.

The example does not show how to serve the application as an ingress or a virtual service.

### Run in Cloudrun

The container can also run as a as a Cloudrun service. A Cloudrun service can cusomized to run with a service account. This service account will have to be granted the permissions to read and write objects in the bucket.

The TLS certificate is automatically created with Cloud Run. However, the image has to be stored in GCR to be pulled.
### Run locally

The container can also run locally (or in any VM) while accessing remote GCS bucket. You can [generate a service account key and expose it to the image with envar GOOGLE_APPLICATION_CREDENTIALS](https://cloud.google.com/docs/authentication/production)

```
export GOOGLE_APPLICATION_CREDENTIALS=path/to/file.json
export GOOGLE_BUCKET=ml-test-modules-registry
make server
```

Running locally won't allow to use it as a terraform registry as not served as HTTPS, but help testing using `curl` for instance :

```
curl -v localhost:8080/test/mymodule/gcp/versions
```

### Server Configuration

The server can be configured with envars.

- `BACKEND` : storage backend to use, `gcs` or `fake`, default `gcs`
- `OVERWRITE` : accepts to overwrite existing modules with same version, default `0` ie prevents from overriding
- `GOOGLE_BUCKET` : name of the GCS bucket to use. Mandatory if backend is `gcp`
- `MODULE_PATH` : path that serves the modules, default, `/`
- `PORT` : port to listen to, default `8080`
- `LISTEN` : accepted IP range, default `0.0.0.0`
- `VERBOSE` : debug logs, default `0`

## API

### Discovery

| Method  | Path                        |
| ------- | --------------------------- |
| **GET** | /.well-known/terraform.json |

Returns a json response to implement discovery process by Terraform
The returned JSON contains the paths for each API versions. Here we implement only the v1.
See more in https://www.terraform.io/internals/remote-service-discovery

#### Parameters

None

#### Response

**status 200** : success

Sample JSON :
```json
{
    "modules.v1":"/"
}
```

### Upload file

| Method   | Path                                   |
| -------- | -------------------------------------- |
| **POST** | /{namespace}/{name}/{system}/{version} |

Uploads a module tarball to the registry.
This API is not defined in Terraform protocol. It is an helper to push new nodules to the registry.
It takes as extra optional parameter the `module-source` which is the code source, that will enable RenovateBot
to detect and scan the CHANGELOG.md (if existing)

#### Parameters

| parameter       | type   | required? | description                                                                                             |
| --------------- | ------ | --------- | ------------------------------------------------------------------------------------------------------- |
| `namespace`     | path   | yes       | unique on a particular hostname, that can contain one or more modules that are somehow related          |
| `name`          | path   | yes       | name of the module                                                                                      |
| `system`        | path   | yes       | name of a remote system that the module is primarily written to target (for example, `gcp`, `aws`, ...) |
| `version`       | path   | yes       | version of the module                                                                                   |
| `module-source` | header | no        | URL of the git repository containing the changelog for renovate                                         |

#### Response

**status 201** : success

**status 418** : failure

**status 403** : module already exists when overwriting is disabled

#### Example

To create the release the module v0.0.2 from local archive file `myfile.tar.gz`:

```
curl -X POST --data-binary "@myfile.tar.gz" localhost:8080/test/mymodule/gcp/0.0.2 -H "module-source: https://whatever.com/wherever.git"
```

### List versions

| Method  | Path                                  |
| ------- | ------------------------------------- |
| **GET** | /{namespace}/{name}/{system}/versions |

Returns the available versions for a given fully-qualified module.
This is required by Terraform client to get the modules.
See more in https://www.terraform.io/internals/module-registry-protocol#list-available-versions-for-a-specific-module

#### Parameters

| parameter   | type | required? | description                                                                                             |
| ----------- | ---- | --------- | ------------------------------------------------------------------------------------------------------- |
| `namespace` | path | yes       | unique on a particular hostname, that can contain one or more modules that are somehow related          |
| `name`      | path | yes       | name of the module                                                                                      |
| `system`    | path | yes       | name of a remote system that the module is primarily written to target (for example, `gcp`, `aws`, ...) |

#### Response

**status 200** : success. Returns JSON list of the available versions

Sample JSON :
```json
{
    "modules": [
        {
            "versions":[
                {"version":"1.2.3"},
                {"version":"0.1.2"},
                {"version":"0.0.1"}
            ]
        }
    ]
}
```

**status 418** : failure

#### Example

```
curl localhost:8080/test/mymodule/gcp/versions
```

### Get latest version

| Method  | Path                         |
| ------- | ---------------------------- |
| **GET** | /{namespace}/{name}/{system} |

Returns the latest version of a module for a single provider.
This API is not part of the Terraform module protocol but is needed for RenovateBot support.
This API is defined in official Terraform registry.
The complete response is not implemented, but only the fields that are needed by RenovateBot
See more in https://www.terraform.io/registry/api-docs#latest-version-for-a-specific-module-provider

#### Parameters

| parameter   | type | required? | description                                                                                             |
| ----------- | ---- | --------- | ------------------------------------------------------------------------------------------------------- |
| `namespace` | path | yes       | unique on a particular hostname, that can contain one or more modules that are somehow related          |
| `name`      | path | yes       | name of the module                                                                                      |
| `system`    | path | yes       | name of a remote system that the module is primarily written to target (for example, `gcp`, `aws`, ...) |

#### Response

**status 200** : success. Returns a JSON object

Sample JSON :
```json
{
    "name":"mymodule",
    "namespace":"test",
    "provider":"gcp",
    "source":"https://my.git/terraform/mymodule",
    "version":"1.2.3",
    "versions":[
        "1.2.3",
        "0.1.2",
        "0.0.1"
    ]
}
```

**status 418** : failure

#### Example
```
curl localhost:8080/test/mymodule/gcp
```

### Get archive link

| Method  | Path                                            |
| ------- | ----------------------------------------------- |
| **GET** | /{namespace}/{name}/{system}/{version}/download |

This does not actually download the module tarball, but sends a link to the tarball.
The tarball could be hosted in another domain tha the API. Not the case here.
See more in https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version

#### Parameters

| parameter   | type | required? | description                                                                                             |
| ----------- | ---- | --------- | ------------------------------------------------------------------------------------------------------- |
| `namespace` | path | yes       | unique on a particular hostname, that can contain one or more modules that are somehow related          |
| `name`      | path | yes       | name of the module                                                                                      |
| `system`    | path | yes       | name of a remote system that the module is primarily written to target (for example, `gcp`, `aws`, ...) |
| `version`   | path | yes       | version of the module                                                                                   |

#### Response

**status 204** : success. No JSON. Returns the archive link as response header `X-Terraform-Get`

```
* Connection state changed (MAX_CONCURRENT_STREAMS == 256)!
< HTTP/2 204
< date: Thu, 13 Jan 2022 10:04:54 GMT
< content-type: application/json
< x-terraform-get: /test/mymodule/gcp/1.2.3/archive.tgz
```

**status 418** : failure

#### Example
```
curl localhost:8080/test/mymodule/gcp/1.2.3/download
```

### Download module

| Method  | Path                                               |
| ------- | -------------------------------------------------- |
| **GET** | /{namespace}/{name}/{system}/{version}/archive.tgz |

Actually download module source. The API contains `.tgz` to force the autodetection from Terraform
See more in https://www.terraform.io/language/modules/sources#fetching-archives-over-http

#### Parameters

| parameter   | type | required? | description                                                                                             |
| ----------- | ---- | --------- | ------------------------------------------------------------------------------------------------------- |
| `namespace` | path | yes       | unique on a particular hostname, that can contain one or more modules that are somehow related          |
| `name`      | path | yes       | name of the module                                                                                      |
| `system`    | path | yes       | name of a remote system that the module is primarily written to target (for example, `gcp`, `aws`, ...) |
| `version`   | path | yes       | version of the module                                                                                   |

#### Response

**status 200** : success.

Download the file as `application/x-gzip`

**status 418** : failure

#### Example
```
curl localhost:8080/test/mymodule/gcp/1.2.3/archive.tgz -o local.tar.gz
```

## Build

### Generate server code from OpenAPI specification

If you need to change the API, you have to install [`oapi-codegen`](https://github.com/deepmap/oapi-codegen) to generate code

```
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
make generate
```

This generates the file `pkg/modules/modules.gen.go`

### Build and push image

We build and push the image using [`ko`](https://github.com/google/ko) from Google
```
go install github.com/google/ko
make push
```

You can change the repository by overriding the variable `KO_DOCKER_REPO`
```
make KO_DOCKER_REPO=wherever.com/whatever build
```

