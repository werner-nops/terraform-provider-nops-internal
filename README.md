# terraform-provider-nops

### Development

To develop:

Create a `dev_overrides` for this provider:

In `~/.terraformrc`
```
provider_installation {

  dev_overrides {
      "registry.terraform.io/nops-io/nops" = "<PATH TO GO BIN>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Then `go install` or `go build`, ensuring the binary produced is placed in the above override directory. This built binary will be used in place of the versioned binary available from the registry.

Once overridden, proceed to the `examples/onboard-nops-resources` directory and do the usual:

You'll need `NOPS_API_TOKEN` exposed either via flags to the apply or via environment with a tool like `direnv`.

`NOPS_HOST` can be set to override the host used for requests to nOps.

```
terraform apply
```

You should see a warning about the override:
```
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - nops-io/nops in <PATH FROM ABOVE>
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
```

### Tests

Run tests, including acceptance tests, with
```
TF_ACC=1 make test
```

### Docs

Regenerate documentation with
```
go generate ./...
```