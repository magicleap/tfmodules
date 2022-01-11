# Documentation for Terraform Modules Registry

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DefaultApi* | [**discovery**](Apis/DefaultApi.md#discovery) | **GET** /.well-known/terraform.json | Discovery process
*DefaultApi* | [**download**](Apis/DefaultApi.md#download) | **GET** /{namespace}/{name}/{system}/{version}/archive.tgz | Actually download tarball
*DefaultApi* | [**getDownloadLink**](Apis/DefaultApi.md#getdownloadlink) | **GET** /{namespace}/{name}/{system}/{version}/download | Download module source
*DefaultApi* | [**getLatestVersion**](Apis/DefaultApi.md#getlatestversion) | **GET** /{namespace}/{name}/{system} | Latest Version for a Specific Module Provider
*DefaultApi* | [**listVersions**](Apis/DefaultApi.md#listversions) | **GET** /{namespace}/{name}/{system}/versions | List module versions
*DefaultApi* | [**upload**](Apis/DefaultApi.md#upload) | **POST** /{namespace}/{name}/{system}/{version} | Upload module version


<a name="documentation-for-models"></a>
## Documentation for Models

 - [Module](.//Models/Module.md)
 - [ModuleDetails](.//Models/ModuleDetails.md)
 - [ModuleRegistry](.//Models/ModuleRegistry.md)
 - [ModuleVersion](.//Models/ModuleVersion.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
