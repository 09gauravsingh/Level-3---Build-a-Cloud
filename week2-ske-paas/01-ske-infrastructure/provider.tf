provider "stackit" {
  default_region = var.region

  service_account_key_path = pathexpand(
    "~/.config/stackit/sa-key-2e086d71-a875-4f4c-83d1-7aad85dff500.json"
  )
}
