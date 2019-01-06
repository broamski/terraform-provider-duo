Terraform Provider for Duo Security
==================

- Website: https://www.terraform.io
- Documentation: *TBD*
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by:

* [Brian Nuszkowski](https://github.com/broamski)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11+
-	[Go](https://golang.org/doc/install) 1.11.0 or higher


Using the provider
----------------------

```
provider "duo" {
    ikey = "DIWJ8X6AEYOR5OMC6TQ1"
    api_host = "api-XXXXXXXX.duosecurity.com."
}

resource "duo_admin" "test_user" {
    email = "sir.brian@email.com"
    name = "SIR BRIAN NUSZKOWSKI"
    phone = "+18005551234"
}

resource "duo_admin_auth_factors" "default_admin_factors" {
    mobile_otp_enabled = true
    push_enabled = true
}

resource "duo_integration" "1pass" {
    name = "Family 1Password"
    type = "1password"
}
```

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/broamski/terraform-provider-duo`

```sh
$ mkdir -p $GOPATH/src/github.com/broamski; cd $GOPATH/src/github.com/broamski
$ git clone git@github.com:broamski/terraform-provider-duo
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/broamski/terraform-provider-duo
$ make build
```

Developing the provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-duo
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. `DUO_IKEY`, `DUO_SKEY`, and `DUO_API_HOST` environment variables must be set in order to successfully run acceptance tests

```sh
$ make testacc
```
