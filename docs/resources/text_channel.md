# Discord Text Channel Resource

A resource to create a text channel

## Example Usage

```hcl-terraform
resource discord_text_channel general {
  name = "general"
  server_id = var.server_id
}
```

## Argument Reference

* `name` (Required) Name of the category
* `server_id` (Required) ID of server this category is in
* `topic` (Optional) Topic of the channel
* `nsfw` (Optional) Whether the channel is NSFW
* `category` (Optional) ID of category to place this channel in
