# IBM Technology Zone Command Line Interface (itz)

![build status](https://github.com/lfDev28/itzcli/actions/workflows/build-go.yml/badge.svg) ![release status](https://github.com/lfDev28/itzcli/actions/workflows/release-cli.yml/badge.svg)

The `itz` command line interface is a command line interface that provides CLI access to IBM Technology Zone.

This is a fork of the original itz cli repository. The original repository can be found [here](https://github.com/cloud-native-toolkit/itzcli).

## Introduction

Using `itz`, you can:

- List your existing reservations and get their status.
- List the available pipelines that you can install in a Red Hat OpenShift cluster
  that you reserved in TechZone.
- Install or deploy products outside of TechZone using infrastructure as code in your
  own OpenShift cluster.
- Reserve an IBM Techzone Environment.
- Extend the reservation of an IBM Techzone Environment.

## Quickstart

See the [QUICKSTART](QUICKSTART.md).

For usage documentation, see the documentation in the [docs](docs/itz.md) folder.

## Installing `itz`

Release packages for your OS can be found at https://github.com/lfDev28/itzcli/releases.

### Installing on Linux

To install `itz` on Linux, use the _install.sh_ script as shown here:

```bash
$ curl https://raw.githubusercontent.com/cloud-native-toolkit/itzcli/main/scripts/install.sh | bash -
```

Verify the installation by using the `itz version` command to view the current
version.

By default, the script installs `itz` in `/usr/local/bin`. If you would like this to be in
a different location, you can set the `ITZ_INSTALL_HOME` environment variable, like
this:

```bash
$ export ITZ_INSTALL_HOME=~/bin
$ curl https://raw.githubusercontent.com/lfdev28/itzcli/main/scripts/install.sh | bash -
```

### Installing on Mac

> **_Note: if you have version 1.24 and installed itz with `brew`, you must
> use brew to uninstall itz and then re-install it._**

#### If you have itz already installed

If you have itz already installed and `itz version` outputs _1.24_ (or a lower
version), you must follow these steps first:

1. Use brew to uninstall the existing itz.
   ```bash
   brew uninstall itz
   ```
1. Untap the existing repository.
   ```bash
   brew untap lfdev28/itzcli-maximo
   ```

Once you have uninstalled itz, you can proceed to
"[Installing itz using brew](#installing-itz-using-brew)".

#### Installing itz using brew

To install `itz` using [Homebrew](), follow these steps:

1. Tap the cask.
   ```bash
   brew tap lfdev28/itzcli-maximo
   ```
2. Install ITZ with brew.
   ```bash
   brew install itzcli_maximo
   ```

### Signing on to IBM Technology Zone

Version v0.1.245 and higher of `itz` supports IBM's Single Sign On (SSO) to
authenticate against the TechZone APIs. To log in, type the following:

```bash
$ itz login
```

This command will automatically open a browser. You can log into IBM Verify with
your IBM ID. When you are done, you can close the browser window.

#### Signing on without a browser

For headless VMs or scripts, the `itz login` command also supports authentication
using the token stored in a file. Log in at https://techzone.ibm.com and view your
profile using the **My profile** page. Copy the value from the **API token** field
and store the value in a file (e.g., `~/token.txt`). Then use the command as
shown to store the value:

```bash
$ itz login --from-file ~/token.txt
```
