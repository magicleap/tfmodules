# tfmodules

`tfmodules` is a registry for Terraform modules, compatible with [RenovateBot](https://github.com/renovatebot/renovate) with GCS bucket as a storage backend.

## Install the server

The image can run almost anywhere. The image needs to access the GCS bucket with the right permissions.

TODO: add deployment examples

### Run in GKE

The container can run as a deployment in GKE using workload identity to bind the service account running the pod to the IAM roles.

### Run in Cloudrun

The container can run as a as a Cloudrun service.

### Run locally

The container can run locally (or in any VM) while accessing remote GCS bucket. You can [generate a service account key and expose it to the image with envar GOOGLE_APPLICATION_CREDENTIALS](https://cloud.google.com/docs/authentication/production)

```
export GOOGLE_APPLICATION_CREDENTIALS=path/to/file.json
export GOOGLE_BUCKET=ml-test-modules-registry
export MODULE_PATH=/
make server
```

In another window :
```
curl -v localhost:8080/test/mymodule/gcp/versions
```

A `fake` server is also available to test without having to connect to GCP. It is used for unit testing

### Server Configuration

We chose to configure the server only with envars, easy to set in a pod

- `BACKEND` : storage backend to use, `gcs` or `fake`, default `gcs`
- `OVERWRITE` : accepts to overwrite existing modules with same version, default `0` ie prevents from overriding
- `GOOGLE_BUCKET` : name of the GCS bucket to use. Mandatory if backend is `gcp`
- `MODULE_PATH` : path that serves the modules, default, `/`
- `PORT` : port to listen to, default `8080`
- `LISTEN` : accepted IP range, default `0.0.0.0`
- `VERBOSE` : debug logs, default `0`

## APIs and curl instructions

### Upload file

To create the module v0.0.2 from local tar.gz
```
curl -X POST --data-binary "@myfile.tar.gz" localhost:8080/test/mymodule/gcp/0.0.2 -H "module-source: https://whatever.com/wherever.git"
```

The `module-source` header is used to pass the URL of the sourcecode. It will be used by renovate to fetch the `CHANGELOG.md` file. The changelog has to follow a standard format.

### List versions

This API is used by `terraform` client
```
curl localhost:8080/test/mymodule/gcp/versions
```

### Get latest version

This specific API is used by renovate to detect the latest available version and compare it to the current one. It will also return the CHANGELOG difference between the two versions.
```
curl localhost:8080/test/mymodule/gcp
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

