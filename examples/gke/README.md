# Deploy in GKE using Config Connectior and Workload Identity

GCP allows to bind a k8s service account to an IAM service account thanks to [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity).

Config Conector is an operator allowing to deploy GCP assets as k8s manifests. It can be installed as a [GKE addon](https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall).

This example shows installation using Config Connector to provision the bucket, the IAM service account and to configure the Workload Identity.

## Project and Config Connector configuration

The GCP project is named `myproject`.

We consider that a namespace `myproject-cc` is already created with the config connector context configured for the hosting GCP project `myproject`. The service account used to run Config Connector must at least have the permissions to GCS and IAM.

The application will be deployed in the namespace `tfmodules`

## sa.yaml

An IAM service account `sa-for-bucket` will be created in the project using the CRD [`IAMServiceAccount`](https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iamserviceaccount). It will be bound to a k8s service account `tfmodules` in namespace `tfmodules`. This is done using the config connector CRD [`IAMPolicyMember`](https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iampolicymember).

## bucket.yaml

We provision a bucket named `tfmodules-registry` using CRD [`StorageBucket`](https://cloud.google.com/config-connector/docs/reference/resource-docs/storage/storagebucket). The CRD [`StorageBucketAccessControl`](https://cloud.google.com/config-connector/docs/reference/resource-docs/storage/storagebucketaccesscontrol) allows to grant write permissions to the `sa-for-bucket` service account.

## deployment.yaml

We deploy the application as a k8s `deployment`, running as the k8s service account `tfmodules` that has been granted the permissions to read and write in the bucket. We pass the created bucket as an envar `GOOGLE_BUCKET`

