![Release with goreleaser](https://github.com/greymatter-io/gmcat/workflows/Release%20with%20goreleaser/badge.svg)

# Meshtool in golang

This tool is to be used with Grey Matter Service Mesh

## Prerequisites

- Golang 1.14
- [Grey Matter](https://greymatter.io/grey-matter) Service Mesh
- [Quickstart Certs](./certs/Readme.md)
- make
- [Gorelease](https://goreleaser.com)

## Install

`brew install greymatter-io/greymatter/gmcat`

## Configure

This tool is setup to use either a `config.yaml` in the same directory as the binary or environment variables.

- Configure by placing [certs](certs/Readme.md) or pointing to the certs
- Configure edge communication to the edge url of your mesh and specify catalog's endpoint
- Specify the `CATALOG_FILE_NAME` (current default is `06.catalog.json`)
- Point json config path to a directory using `$JSON_CONFIG_PATH/<config directory>/$CATALOG_FILE_NAME`

## Usage

After installing the binary:

- **create**: `gmcat create -f <project directory name>`
- **Delete**: `gmcat delete`
  - `-f <config directory>` will use the "clusterName" to delete what is in the mesh
- **Search**: `gmcat search`
  - `-f <config directory>` will search for the clusterName and present a diff between the input file and what is in the mesh
- **Update**: `gmcat update`
  - `-f <config directory>` will search for the clusterName and interactively use the provided file as a starting point (uses vi to edit).  This creates a backup in the same directory.
- **Version**: `gmcat version`

### Docker Usage

The `gmcat` binary can be used through docker by volume mounting the gm-config-json files, certs, and configuring it via envar or config.yaml placed in `/app/` directory.

The quickstart certs are baked in to help speed up interfacing with quickstart based meshs.  

ex: This example uses version but any can be used

```console
docker run -it -v $(pwd)/example-greymatter:/app/example-greymatter -v /Users/kyleg/.ssh/moosecanswim.p12:/app/moosecanswim.p12 gmcat-dev:0.1.0 version
```

## Makefile

| target      | comment                                                    |
| ----------- | ---------------------------------------------------------- |
| build-dev   | builds dev binary with version set to"dev"                 |
| build-prod  | builds dev binary with version set to value in `./version` |
| docker-dev  | packages alpine specific binary into `gmcat-dev:<version>` |
| docker-prod | packages alpine specific binary into `gmcat:<version>`     |
