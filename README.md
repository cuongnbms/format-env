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
wget https://github.com/cuongnb14/format-env/releases/download/v1.0.1/fenv_darwin_arm64 -O fenv
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
