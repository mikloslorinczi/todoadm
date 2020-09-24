# TODOADM
or TODO Admin

Is a CLI tool to manage the TODO Service. It is written in Go, using spf13's Cobra and Viper libraries. The HTTP client used to talk with the TODO API leverages h2non's Gentleman lib.

## Usage

In order to communicate with the API you need 3 variables to be set:
- TODO_HOST (the URL of the TODO Service)
- ADMIN_HEADER (the key in the HTTP Admin Header)
- ADMIN_SECRET (the value in the HTTP Admin Header)

You may set these values with either a config file (**todoadm_conf.yaml**), environment variables (prefixed with **TODOADM_**) or command line flags.

to test the connectivity ping the service with:\
```$ todoadm ping```
