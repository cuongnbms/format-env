# Format env
format-env is a Go command-line utility to generate and format environment files (.env) using a template. This tool allows you to specify the environment (e.g., dev, testing, staging, etc.) and a template file to generate the corresponding environment configuration file.

# Features
- Reads key-value pairs from an existing .env file.
- Applies the values to a template for formatting.
- Supports multiple stages (dev, testing, unstable, staging).
- Allows specifying the template path and stage dynamically through command-line arguments.

# Installation
## Option 1: Donwload binary file
```sh
# for macos arm64 
wget https://github.com/cuongnbms/format-env/releases/download/v1.1.0/fenv_darwin_arm64 -O fenv
sudo mv fenv /usr/local/bin/
sudo chmod +x /usr/local/bin/fenv
```

## Option 2: Compile from source
To build and install the utility, make sure you have Go installed, and then run:
```
go build -o fenv fenv.go
```

Move the binary to your $PATH for easy access:
```
sudo mv fenv /usr/local/bin/
sudo chmod +x /usr/local/bin/fenv
```
# Usage

To use the utility, run it from the command line by specifying the template file and the environment stage you want to generate.
```
fenv <env_dir> <stages>
```
- env_dir: Path to the env dir (e.g., `env`). The template file name `_template.env` must have in this folder
- stages: The environment stage, separate by comma (e.g. `dev,testing,staging`).

Example
```
| env
|--- _template.env
|--- dev.env
|--- staging.env
|--- prod.env

# run format
fenv env dev,staging,prod
```

# Pre-commit

This repository can be used as a remote hook with the
[`pre-commit`](https://pre-commit.com) framework. Add the following config to
the project whose environment files should be formatted:

```yaml
repos:
  - repo: https://github.com/cuongnbms/format-env
    rev: v1.1.0
    hooks:
      - id: format-env
        args:
          - env
          - dev,test,prod
```

The first argument is the directory containing `_template.env`. The second is
the comma-separated list of stages to format. Install and verify the hook with:

```sh
pre-commit install
pre-commit run --all-files
```

When the hook changes an environment file, `pre-commit` stops the commit. Review
the changes, stage the files again, and retry the commit.

Use multiple hook entries when a project has more than one environment
directory:

```yaml
repos:
  - repo: https://github.com/cuongnbms/format-env
    rev: v1.1.0
    hooks:
      - id: format-env
        name: format service environment files
        args: [services/api/env, dev,test,prod]
      - id: format-env
        name: format worker environment files
        args: [services/worker/env, dev,prod]
```

# Template Syntax
The template file should use Go’s text/template syntax. For example:
```
# General
# =======================================================================
STAGE={{ .STAGE }}
PORT={{ df .PORT "8000" }}

# Database
# =======================================================================
DB_HOST={{ .DB_HOST }}
DB_PORT={{ .DB_PORT }}
DB_USER={{ .DB_USER }}
DB_PASS={{ .DB_PASS }}
DB_NAME={{ df .DB_NAME "example" }}

# Redis
# =======================================================================
REDIS_ENABLE_SSL={{ df .REDIS_ENABLE_SSL "false" }}
REDIS_URL={{ .REDIS_URL }}
```
Note: This template will render the environment variable `REDIS_ENABLE_SSL` and use a default value `false` if it is not provided.
