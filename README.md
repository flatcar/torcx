<img align="left" width="70px" src="Documentation/torcx.png" />

# torcx - a boot-time addon manager

[![Apache](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://travis-ci.org/coreos/torcx.svg?branch=master)](https://travis-ci.org/coreos/torcx)

torcx (pronounced _"torks"_) is a boot-time manager for system-wide ephemeral customization of Linux systems.
It has been built specifically to work with an immutable OS such as [Flatcar Container Linux][flatcar-cl] by Kinvolk.

[flatcar-cl]: https://www.flatcar-linux.org/releases/

torcx focuses on:
* providing a way for users to add additional binaries and services, even if not shipped in the base image
* allowing users to pin specific software versions, in a seamless and system-wide way
* supplying human- and machine-friendly interfaces to work with images and profiles

## Getting started

This project provides a very lightweight add-ons manager for otherwise immutable distributions.
It applies collections of addon packages (named, respectively, "profiles" and "images") at boot-time, extracting them on the side of the base OS.

Profiles are simple JSON files, usually stored under `/etc/torcx/profiles/` with a `.json` extension, containing a set of image-references:

```json
{
  "kind": "profile-manifest-v1",
  "value": {
    "images": [
      {
        "name": "foo-addon",
        "reference": "0.1",
        "remote": "com.example.foo"
      }
    ]
  }
}

```

Image archives are looked up in several search paths, called "stores":
 1. Vendor store: usually on a read-only partition, it contains addons distributed together with the OS image
 1. User store: usually on a writable partition, it contains images provided by the user
 1. Runtime store: additional search path specified at runtime

At boot-time, torcx unpacks and propagates the addons defined in the active profile, specified in `/etc/torcx/next-profile`.
Once done, torcx seals the system into its new state and records its own metadata under `/run/metadata/torcx`.

## Contributing

Please see [CONTRIBUTING](https://github.com/kinvolk/contribution/) and the [Kinvolk Code of Conduct](https://github.com/kinvolk/contribution/blob/master/CODE_OF_CONDUCT.md)

## License

torcx is released under the Apache 2.0 license.
See the [LICENSE](LICENSE) file for all details.

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](http://developercertificate.org/) site for details.
