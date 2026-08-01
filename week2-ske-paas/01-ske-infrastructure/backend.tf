terraform {
  backend "s3" {
    bucket = "storage-tfstate-20a60c06-eu01"
    key    = "week2-ske-paas/ske/terraform.tfstate"
    region = "eu01"

    endpoints = {
      s3 = "https://object.storage.eu01.onstackit.cloud"
    }

    # STACKIT Object Storage is S3-compatible, but it is not AWS.
    # Avoid requests to AWS identity and metadata services.
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_metadata_api_check     = true

    # STACKIT supports path-style object addressing:
    # endpoint/bucket/object-key
    use_path_style = true

    # Useful for S3-compatible implementations.
    skip_s3_checksum = true

    # Request server-side AES-256 encryption.
    encrypt = true

    # Create a .tflock object to prevent concurrent state writes.
    use_lockfile = true
  }
}
