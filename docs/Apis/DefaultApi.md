# DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**discovery**](DefaultApi.md#discovery) | **GET** /.well-known/terraform.json | Discovery process
[**download**](DefaultApi.md#download) | **GET** /{namespace}/{name}/{system}/{version}/archive.tgz | Actually download tarball
[**getDownloadLink**](DefaultApi.md#getDownloadLink) | **GET** /{namespace}/{name}/{system}/{version}/download | Download module source
[**getLatestVersion**](DefaultApi.md#getLatestVersion) | **GET** /{namespace}/{name}/{system} | Latest Version for a Specific Module Provider
[**listVersions**](DefaultApi.md#listVersions) | **GET** /{namespace}/{name}/{system}/versions | List module versions
[**upload**](DefaultApi.md#upload) | **POST** /{namespace}/{name}/{system}/{version} | Upload module version


<a name="discovery"></a>
# **discovery**
> Object discovery()

Discovery process

    Returns a json response to implement discovery process by Terraform The returned JSON contains the paths for each API versions. Here we implement only the v1. See more in https://www.terraform.io/internals/remote-service-discovery 

### Parameters
This endpoint does not need any parameter.

### Return type

[**Object**](..//Models/object.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="download"></a>
# **download**
> File download(namespace, name, system, version)

Actually download tarball

    Actually download module source. The API contains &#x60;.tgz&#x60; to force the autodetection from Terraform See more in https://www.terraform.io/language/modules/sources#fetching-archives-over-http 

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **String**| unique on a particular hostname, that can contain one or more modules that are somehow related | [default to null]
 **name** | **String**| module name | [default to null]
 **system** | **String**| remote system that the module is primarily written to target (aws, gcp, ...) | [default to null]
 **version** | **String**| version of the module | [default to null]

### Return type

[**File**](..//Models/file.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/x-gzip

<a name="getDownloadLink"></a>
# **getDownloadLink**
> getDownloadLink(namespace, name, system, version)

Download module source

    This does not actually download the module tarball, but sends a link to the tarball. The tarball could be hosted in another domain tha the API. Not the case here. See more in https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version 

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **String**| unique on a particular hostname, that can contain one or more modules that are somehow related | [default to null]
 **name** | **String**| module name | [default to null]
 **system** | **String**| remote system that the module is primarily written to target (aws, gcp, ...) | [default to null]
 **version** | **String**| version of the module | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getLatestVersion"></a>
# **getLatestVersion**
> ModuleDetails getLatestVersion(namespace, name, system)

Latest Version for a Specific Module Provider

    Returns the latest version of a module for a single provider. This API is not part of the Terraform module protocol but is needed for RenovateBot support. This API is defined in official Terraform registry. The complete response is not implemented, but only the fields that are needed by RenovateBot See more in https://www.terraform.io/registry/api-docs#latest-version-for-a-specific-module-provider 

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **String**| unique on a particular hostname, that can contain one or more modules that are somehow related | [default to null]
 **name** | **String**| module name | [default to null]
 **system** | **String**| remote system that the module is primarily written to target (aws, gcp, ...) | [default to null]

### Return type

[**ModuleDetails**](..//Models/ModuleDetails.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listVersions"></a>
# **listVersions**
> ModuleRegistry listVersions(namespace, name, system)

List module versions

    Returns the available versions for a given fully-qualified module. This is required by Terraform client to get the modules. See more in https://www.terraform.io/internals/module-registry-protocol#list-available-versions-for-a-specific-module 

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **String**| unique on a particular hostname, that can contain one or more modules that are somehow related | [default to null]
 **name** | **String**| module name | [default to null]
 **system** | **String**| remote system that the module is primarily written to target (aws, gcp, ...) | [default to null]

### Return type

[**ModuleRegistry**](..//Models/ModuleRegistry.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="upload"></a>
# **upload**
> upload(namespace, name, system, version, moduleSource, body)

Upload module version

    Uploads a module tarball to the registry. This API is not defined in Terraform protocol. It is an helper to push new nodules to the registry. It takes as extra optional parameter the &#x60;source&#x60; which is the code source, that will enable RenovateBot to detect an scan the CHANGELOG.md (if existing) 

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **String**| unique on a particular hostname, that can contain one or more modules that are somehow related | [default to null]
 **name** | **String**| module name | [default to null]
 **system** | **String**| remote system that the module is primarily written to target (aws, gcp, ...) | [default to null]
 **version** | **String**| version of the module | [default to null]
 **moduleSource** | **String**| code URL of the module | [optional] [default to null]
 **body** | **File**|  | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/octet-stream
- **Accept**: Not defined

