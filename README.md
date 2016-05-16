# Imterra Forge
Forge is an universal, extensible, and distributed build-tool. It provides
simple to configure project build management, including incremental build
support. On top of this, it also provides support for compilation using servers
or any computers running the server-side program. This is done in order to
speed up the compilation process.

In this file, we describe basic installation, configuration and usage.

## Installation

### Dependencies

Since Forge is written in Go, in order to install from source, Go compiler is
required. Forge has been developed and tested with official Go compiler,
version 1.4.2+. You can get Go at the [official website](https://golang.org).

In order to parse configuration files, Forge also requires installation of
[simpleyaml library](https://github.com/smallfish/simpleyaml).

You can download source code for Forge from
[Forge GitHub](https://github.com/imterra/forge).

### Client-side program (Master)

Client-side program is located in `client/` folder. Building it is simple, just
run `go build` from the directory. Building the project creates executable file
named `client`, which should be moved to user's `PATH` variable and renamed to
`forge`.

### Server-side program (Slave, Worker)

The server-side program source code resides in `worker/` directory. Compilation
is the same as in client, just run `go build` in the directory. After that,
the slave executable must be moved to user's `PATH~ and renamed to
`forge-server`.

## Configuration

### Server-side

Server-side configuration is done when running the server, using command-line
flags:

 - **jobs** - maximal number of simultaneously running compilation jobs on the
   current server, equivalent to make's `-j` flag (default: number of CPU cores
   available)
 - **port** - TCP port to listen on for incoming RPC requests (default: 1103)
 - **root** - root folder of Forge projects hierarchy on the server, this
   setting overrides `FORGE_ROOT` environmental variable

### Client-side

The client can use a textual configuration file, by default `~/.forge.yaml` and
`/etc/forge.yaml`. These are written in YAML with simple format of
`flag: value`. The following flags are defined:

 - **jobs** - number of jobs of local worker, this number is passed on to
   locally-started worker, unless set to 0 (default: number of CPU cores)
 - **root** - root folder of Forge project hierarchy overrides value set by
   `FORGE_ROOT`, the value is sent to locally-started worker
 - **worker** - a list of `host:port` pairs, can be either comma-separated or
   provided through multiple use of the flag, this specifies the workers that
   are used (outside of the locally-started one) for compilation

**Note:** In case of multiple different settings, settings from
`/etc/forge.yaml` are overriden by `~/.forge.yaml~, which is overriden by
(in case of **root**) `FORGE_ROOT`, which are all overrides by job flags.

## Usage

Forge expects a list of (at least one) targets to be built. These targets can
be either FQTN (Full-Qualified Target Names) for any project within forge root,
as well as names of local targets, if the current working directory is
within the forge tree.

### FQTN

Full-Qualified Target Name (FQTN) describes location of target within the whole
Forge hieararchy. It starts with double forward slash (`//`) and then continues
with path from the root directory to the build.yaml file, which describes the
target. After the last slash, the name of the target is written.

For example, for target `foo` defined in file `bar/baz/foo/build.yaml`, the
FQTN would be `//bar/baz/foo/foo`.

### BUILD files

BUILD files describe the project's targets. They are written in YAML format and
are named `build.yaml`. The document is an object, with each field being the
name of the target.

Every target definition is an object, containing the following fields:

 - **type** - type of a target, describing action generation, currently only
   `app_c` and `lib_c` are supported.
 - **sources** - list of names (in either relativne form, or FQTN form) of
   target's source files
 - **resources** - list of resource file names (in either relativne to BUILD
   file or in FQTN form), resource files are used to define header files for
   `app_c` and `lib_c` target types
 - **dependencies** - list of target names (either relative or FQTN), on which
   the successful completion of this target is dependent

#### Recommendations

Forge hierarchy and `FORGE_ROOT` were created in order to avoid target name
clashes between projets. They provide a single hierarchical place, in which
caution must be exercised to avoid conflicts. That's why, when developing a
project, it should be placed under the company's name internal directory
structure (i.e. //imterra/forge, for Imterra's own Forge project). For
projects developed on GitHub, or similar code sharing platform, we recommend
using the site's name as the top-level folder, with the usual path being the
further hierarchy. Thus `https://github.com/USERNAME/PROJECTNAME` would be
hosted under `//github/USERNAME/PROJECTNAME` hierarchy.
