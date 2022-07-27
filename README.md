# Terraform Provider Ip (Terraform Plugin SDK)

This tiny provider will fetch the ip address on the host executing Terraform. Typical use cases are to fetch IP address of the host & append to public keys in .authorized_keys file while SSHing to target device.

## Usage Example
```hcl
# 1. Specify the ip provider to use
terraform {
  required_providers {
    ip = {
      source = "harshavmb/ip"      
    }
  }
}

# 2. Configure the ip Provider
provider "ip" {  
}

# 3. Fetch the ip address
data "ip" "example" {  
}

# 4. If interested in IP address of specific NIC
data "ip" "eth" {  
    nw_interface = "eth0"
}
```


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
