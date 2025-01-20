> This repository is part of an [article series](https://dev.to/kuroski/building-a-webhook-payload-delivery-service-in-go-31bg)

## Introduction

This repository contains two applications:

- A [Webhook payload delivery service](https://dev.to/kuroski/building-a-webhook-payload-delivery-service-in-go-31bg) for local development
- A Telegram bot that logs GitHub Actions workflows

I will write the second part of the article once I have more time, but feel free to read the [first one](https://dev.to/kuroski/building-a-webhook-payload-delivery-service-in-go-31bg).

https://github.com/user-attachments/assets/ca3580e0-61ba-4b11-bd61-f5f7e8e07204

## Setup

- Add a `.env` by `cp .env.sample .env`

``` dotenv 
# this can be found when creating your bot through BotFather
# https://core.telegram.org/bots/tutorial#obtain-your-bot-token
TELEGRAM_BOT_TOKEN=
# this can either be found through the browser URL, it should have a format liks "-1234567890"
# or you can use https://telegram.me/rawdatabot
TELEGRAM_CHAT_ID=

# this is for dev environment + it is explained on the article how it works
DEV_CLI_SOURCE_URL=https://your-service.com # e.g. https://smee.io/
DEV_CLI_TARGET_URL=http://web:3000/webhook

# you must have a github app created with the Webhook URL pointing to your server
# this will be the URL of the deployed webserver
# https://your-server.com/webhooks - for prod
# https://your-server.com/channel/<any-wildcard> - for dev
# more info on how to generate the key below
GITHUB_APP_PRIVATE_KEY=

# the APP_ID and APP_INSTALLATION_ID can be found directly in the Github App page
GITHUB_APP_ID=
GITHUB_APP_INSTALLATION_ID=
```

#### Generating github app private key

- Go to github app page
- Click on "generate a private key" in the **Private keys** section
- Save the `.pem` file appropriately
- Generate a base64 representation of the item using `base64 -w 0 -i my.pem > encoded-private-key.txt` or `pbcopy < base64 -w 0 -i my.pem`
- Copy the content into the `GITHUB_APP_PRIVATE_KEY` env variable

## Commands
- `make dev` - runs docker environment with CLI and Web
  - CLI forwards `DEV_CLI_SOURCE_URL` to `DEV_CLI_TARGET_URL
  - Web defaults to `http://localhost:3000` -- `web` alias to docker env
  - Default webhook endpoint is `http://localhost:3000/webhook`
- `make cli-run source=http://my-source target=http://localhost:3000/webhook` - runs only CLI environment, don't forget the `source` and `target` params
- `make web-run` - runs webserver defaulting to `http://localhost:3000`
- `make web-run-docker` - runs production version of webserver - defaults to `http://localhost:8080`
- `kamal deploy` or `kamal setup` to deploy things

For most cases, just run 

```bash 
make dev
```

If you want to debug server, just run cli separately

```bash
make cli-run source=http://my-source target=http://localhost:3000/webhook
```

## Deploying the server

You can find a guide on the [first article](https://dev.to/kuroski/building-a-webhook-payload-delivery-service-in-go-31bg)

## Ansible

> All credits for the playbooks are from
> https://github.com/guillaumebriday/kamal-ansible-manager

When provisioning a new droplet, just configure ansible inside infra folder and run it

Copy the inventory example file:
```bash
$ cp infra/hosts.ini.example infra/hosts.ini
```

Update the `<host1>` with your server's IP address (you can have multiple servers):
```bash
$ vim hosts.ini
```

Install the requirements:
```bash
$ ansible-galaxy collection install -r infra/requirements.yml
$ ansible-galaxy role install -r infra/requirements.yml
```

### Configuring vars

Variables can be configured in the `playbook.yml` file.
Also, you can override default variables provided in [geerlingguy/ansible-role-swap](https://github.com/geerlingguy/ansible-role-swap/blob/master/defaults/main.yml) to adjust the swap settings.

For instance:
```yml
  vars:
    security_autoupdate_reboot: "true"
    security_autoupdate_reboot_time: "03:00"
    swap_file_size_mb: '1024'
```

### Running the playbook

Run the playbook:
```bash
$ ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i infra/hosts.ini infra/playbook.yml
```
