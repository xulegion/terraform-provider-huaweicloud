---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud_cce_namespace

Manages a CCE namespace resource within HuaweiCloud.

## Example Usage

### Basic

```hcl
variable "cluster_id" {}

resource "huaweicloud_cce_namespace" "test" {
  cluster_id = var.cluster_id
  name       = "test-namespace"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the namespace resource.
  If omitted, the provider-level region will be used. Changing this will create a new namespace resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID to which the CCE namespace belongs.
  Changing this will create a new namespace resource.

* `name` - (Optional, String, ForceNew) Specifies the name of the namespace. Must be unique. This parameter can
  contain a maximum of 63 characters, which may consist of lowercase letters, digits and hyphens (-), and must
  start and end with lowercase letters and digits. Changing this will create a new namespace resource.
  Exactly one of `name` or `prefix` must be provided.

* `prefix` - (Optional, String, ForceNew) Specifies A prefix used by the server to generate a unique name.
  This parameter can contain a maximum of 63 characters, which may consist of lowercase letters, digits and
  hyphens (-), and must start and end with lowercase letters and digits.
  Changing this will create a new namespace resource. Exactly one of `name` or `prefix` must be provided.

* `annotations` - (Optional, Map, ForceNew) An unstructured key value map for external parameters. Changing this
  will create a new namespace resource.

* `labels` - (Optional, Map, ForceNew) Map of string keys and values for labels. Changing this
  will create a new namespace resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The namespace ID in UUID format.

* `creation_timestamp` - The server time when namespace was created.

* `status` - The current phase of the namespace.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 5 minute.

## Import

CCE namespace can be imported using the cluster ID and namespace name separated by a slash, e.g.:

```
$ terraform import huaweicloud_cce_namespace.test bb6923e4-b16e-11eb-b0cd-0255ac101da1/test-namespace
```
