# Onboarding Payer account with AWS resources

In order to test the `nOps` provider alongside the official providers locally we need to do a couple of extra steps:
- Run `go mod tidy` to install the required dependencies
- Make local changes
- Run `go install .` to build the provider executable file, the location of this file varies but it will be on you `go` bin. Ex: `/Users/your_user/go/bin/`
- If the build succeeded without errors, you should see a file with the name `terraform-provider-nops` in that directory
- If you don't already have it, create a TF plugins directory. Usually this folder is in this location `~/.terraform.d/plugins`
- Execute the following command from this directory `terraform providers mirror ~/.terraform.d/plugins`, this will create a local copy of all required providers to that directory.
- Create or update you `.terraformrc` file to contain the following:
```
provider_installation {
  filesystem_mirror {
    path = "/Users/your_user/.terraform.d/plugins"
  }
  direct {
    exclude = ["terraform.local/*/*"]
  }
}
```
- Now create the following directory structure inside you plugins directory, replace the platform with your own. You can use the mirrored providers as reference
```
./terraform.local/custom/nops/1.0.0/darwin_arm64
```
- Now copy the executable file generated earlier from your `go` bin directory over to this folder, adding as suffix the version
```
./terraform.local/custom/nops/1.0.0/darwin_arm64/terraform-provider-nops_v1.0.0
```
- Now we can run `terraform init` from this directory, this will use the providers we mirrored and the versions we set, important that the version in the required providers config is the same
as the one we have mirrored
- Do your testing