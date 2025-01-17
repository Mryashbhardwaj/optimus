---
id: configuration
title: Configurations
---

# Project Configuration (Client)
See client configuration example on `optimus.sample.yaml`

Optimus project configuration (later on client configuration) can be loaded from file (use `--config` flag), or `optimus.yaml` file in current working directory where the optimus command is executed.

---
**1. Using --config flag**
```sh
$ optimus deploy --config /path/to/config/file.yaml
```
---
**2. Using default optimus.yaml file**
```sh
$ tree
. # current project structure
├── namespace-1
│   └── jobs
│   └── resources
├── namespace-2
│   └── jobs
│   └── resources
└── optimus.yaml # use this file
$ optimus deploy
```

---

If both are exist, then use the file config defined in `--config` flag.

This configuration file should not be checked in version control. 
# Server Configuration
See server configuration example on `config.sample.yaml`

Optimus server configuration can be loaded from file (use `--config` flag), environment variable `OPTIMUS_<CONFIGNAME>`, or `config.yaml` file in executable directory.

---
**1. Using --config flag**
```sh
$ optimus serve --config /path/to/config/file.yaml
```

---
**2. Using environment variable**

All the configs can be passed as environment variables using `OPTIMUS_<CONFIG_NAME>` convention. `<CONFIG_NAME>` is the key name of config using `_` as the path delimiter to concatenate between keys.

For example, to use environment variable, assuming the following configuration layout:

```yaml
version: 1
serve:
  port: 9100
  app_key: randomhash
```

Here is the corresponding environment variable for the above

Configuration key | Environment variable |
------------------|----------------------|
version           | OPTIMUS_VERSION      |
serve.port        | OPTIMUS_PORT         |
serve.app_key     | OPTIMUS_SERVE_APP_KEY|

Set the env variable using export
```sh
$ export OPTIMUS_PORT=9100
```

---
**3. Using default config.yaml from executable binary directory**
```sh
$ which optimus
/usr/local/bin/optimus
```

So the `config.yaml` file can be loaded on `/usr/local/bin/config.yaml`

---

If user specify configuration file using `--config` flag, then any configs defined in env variable and default config.yaml from exec directory will not be loaded.

If user specify the env variable and user also have config.yaml from exec directory, then any key configs from env variable will override the key configs defined in config.yaml from exec directory.

---

App key is used to encrypt credentials and can be randomly generated using
```shell
head -c 50 /dev/random | base64
```
Just take the first 32 characters of the string.