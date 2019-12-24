# SlackHell C2

Slack Web Shell C2 for Fun and Profit

## Introduction

Slackhell is a simple tool for generating and controlling web shell backdoor using Slack Bot.

## Features

* Slack C2
* Users management (ACL)
* Generate web shell backdoor
* Web shell backdoor password protected
* Shell session

## Getting Started

These instructions will get you a copy of the project up and running on your local machine.

### Prerequisites

* Before using Slackhell make sure you have created Slack Bot, by following this guide, [create slack app](https://api.slack.com/start/overview).
* Install [golang](https://golang.org/doc/install)

### Build from Source

```bash
$ git clone https://github.com/herwonowr/slackhell.git
$ cd slackhell
$ go build -o slackhell cmd/main.go
```

### Change the Configuration

```bash
$ vim ./data/config/slackhell.toml
```

```toml
[account]
    # Slackhell initial admin account
    # id e.g: UR*******
    # realname e.g: Vulncode
    id = "UR*******"
    realname = "Vulncode"
[slack]
    token = "xoxb-yourslackbottoken"
[database]
    path = "./data/db/slackhell.db"
[log]
    debug = false
```

Configuration:

* **Account** - Initial admin account for Slackhell
  * id - your slack id
  * realname - your slack realname
* **Slack** - Slack bot token
  * token - your slack bot token (xoxb-*)
* **Database** - Slackhell database path
  * path - database path
* **Log** - Log debug
  * debug - set verbose log

### Start Slackhell

```bash
$ slackhell run
```

### Build Docker Image

```bash
$ git clone https://github.com/herwonowr/slackhell.git
$ cd slackhell
$ docker build -t reponame/slackhell:version .
```

### Start Slackhell Docker Image

```bash
$ docker run --rm -v $(pwd)/config/slackhell.toml:/slackhell/config/slackhell.toml -v $(pwd)/db:/slackhell/db reponame/slackhell:version
```

### Start Slackhell Pre Build Docker Image

```bash
$ docker run --rm -v $(pwd)/config/slackhell.toml:/slackhell/config/slackhell.toml -v $(pwd)/db:/slackhell/db herwonowr/slackhell:v1.0.0
```

## Contributing

Please read [CODE OF CONDUCT](CODE_OF_CONDUCT.md) and [CONTRIBUTING](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning.

## Authors

* **Herwono W. Wijaya** - *Initial work* - [Slackhell](https://github.com/herwonowr/slackhell)

## License

This project is licensed under the GNU GPLv3 License - see the [LICENSE](LICENSE) file for details

## Disclaimer

THIS TOOL IS BEING PROVIDED FOR EDUCATIONAL PURPOSES ONLY, WITH THE INTENT FOR RESEARCH PURPOSES ONLY.

You may not use this software for any illegal or unethical purpose; including activities which would give rise to criminal or civil liability.

USE ON YOUR OWN RISK. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER OR CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES.

## Acknowledgments

* [Go Slack API](https://github.com/nlopes/slack)
* [Slacker](https://github.com/shomali11/slacker)