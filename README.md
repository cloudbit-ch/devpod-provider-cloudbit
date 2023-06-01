# Cloudbit Provider for DevPod

## Getting started

The provider is available for auto-installation using

```sh
devpod provider add cloudbit-ch/devpod-provider-cloudbit
devpod provider use cloudbit-ch/devpod-provider-cloudbit
```

Follow the on-screen instructions to complete the setup.

Needed variables will be:

- CLOUDBIT_TOKEN - Cloudbit Application Token

### Creating your first devpod env with cloudbit

After the initial setup, just use:

```sh
devpod up .
```

You'll need to wait for the machine and environment setup.

### Customize the VM Instance

This provides has the seguent options

| NAME              | REQUIRED | DESCRIPTION                             | DEFAULT                |
|-------------------|----------|-----------------------------------------|------------------------|
| CLOUDBIT_LOCATION | true     | The cloudbit location to use.           | BIT1                   |
| CLOUDBIT_IMAGE    | true     | The cloudbit image to use.              | linux-ubuntu-20.04-lts |
| CLOUDBIT_PRODUCT  | true     | The cloudbit compute VM product to use. | b1.1x1                 |

Options can either be set in `env` or using for example:

```sh
devpod provider set-options -o CLOUDBIT_PRODUCT=b1.2x2
```