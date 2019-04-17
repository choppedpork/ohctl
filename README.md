# ohctl

Because sometimes you need to turn the lights off and the volume down from the command line...

## Setup

ohctl wil connect to `http://localhost:80` by default. You can override this by using a config file (`$HOME/.ohctl.yaml`): 

```yaml
host: openhab
port: 8080
```

or environment variables:

```
$ OHCTL_HOST=openhab OHCTL_PORT=8080 ohctl get items
```

## Usage

Fetch a list of all items:

    ohctl get items

Fetch status of an item:

    ohctl get item sofa_light

Execute a command on an item:

    ohctl cmd kitchen_volume increase