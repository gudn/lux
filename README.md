# Lux

Setup script that generates NGINX configuration from folder structure

## Usage

Lux is nginx docker image with `lux` binary. It accepts three parameters:
- templates -- path to templates. One of they must name nginx.conf and it is
  main configuration file.
- projects -- path to folder with project folders. All projects must have
  `lux.yaml` files with variables and template name.
- output -- config output directory.

Also `lux` accept `exec` parameter for execve nginx to no-daemon mod after
successful config building. Also you can use multi-stage builds (just copy
config directory and projects dirs into the same directories).

For example, see this repository and [Dockerfile](./Dockerfile)
