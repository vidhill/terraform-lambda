# terraform-lambda

[![go workflow](https://github.com/vidhill/terraform-lambda/actions/workflows/go.yml/badge.svg)](https://github.com/vidhill/terraform-lambda/actions)

Test lambda service, written in `go`, provisioning all required AWS resources using terraform.

```mermaid
sequenceDiagram
    S3 Source Bucket->>Lambda:  "s3:ObjectCreated" event
    Lambda->>S3 Destination Bucket:  Resized image
```
