# kubectl-clusternet

A `kubectl` plugin for interacting with [Clusternet](https://github.com/clusternet/clusternet).

![License](https://img.shields.io/github/license/clusternet/kubectl-clusternet)
![Release](https://img.shields.io/github/v/release/clusternet/kubectl-clusternet)
![build](https://github.com/clusternet/kubectl-clusternet/actions/workflows/ci.yml/badge.svg)
![Downloads](https://img.shields.io/github/downloads/clusternet/kubectl-clusternet/total?color=green)

## Installation

### Install With Krew

`kubectl-clusternet` can be installed using [Krew](https://github.com/kubernetes-sigs/krew),
please [install Krew with this guide](https://krew.sigs.k8s.io/docs/user-guide/setup/install/) first.

Then you can install `Clusternet` kubectl plugin with,

```bash
$ kubectl krew install clusternet
```

### Download Binary

Alternatively, `kubectl-clusternet` can be directly downloaded
from [released packages](https://github.com/clusternet/clusternet/releases).

Download a tar file matching your OS/Arch, and extract `kubectl-clusternet` binary from it.

Then copy `./kubectl-clusternet` to a directory in your executable `$PATH`.

### Build on Your Own

Clone this repo and run `make bin`

```bash
$ git clone https://github.com/clusternet/kubectl-clusternet
$ make bin
```

Then copy `./dist/kubectl-clusternet` to a directory in your executable `$PATH`.

## How it works

```bash
$ kubectl clusternet -h
Usage:
  clusternet [flags]
  clusternet [command]

Available Commands:
  apply       Apply a configuration to a resource by filename or stdin
  create      Create a resource from a file or from stdin.
  delete      Delete resources by filenames, stdin, resources and names, or by resources and label selector
  edit        Edit a resource on the server
  get         Display one or many resources
  help        Help about any command
  scale       Set a new size for a Deployment, ReplicaSet or Replication Controller
  version     Print the plugin version information
```
