Beluga
--------

A fun little chat bot for Discord, inspired from the [Narwhal IRC bot](https://github.com/narwhalirc/narwhal). View all available commands in Discord using `!help`.

[![Report](https://goreportcard.com/badge/github.com/EbonJaeger/beluga)](https://goreportcard.com/report/github.com/EbonJaeger/beluga) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
--------

## Building
Beluga is a regular Go program using Go modules, so build it how you would any normal Go application.

## Installation
Create a Discord bot [here](https://discordapp.com/developers/applications/me). Next, add the bot to your Discord server using this link, replacing the Client ID with your bot's ID:
```
https://discordapp.com/oauth2/authorize?client_id=<CLIENT_ID>&scope=bot
```

## CLI Usage
```
./beluga [OPTIONS]
```
Options:
```
-c, --configDir - Specify the directory to use for configuration files
-h, --help   - Print the help message
```

## License
Copyright © 2020 Evan Maddock (EbonJaeger)

Beluga is available under the terms of the Apache-2.0 license
