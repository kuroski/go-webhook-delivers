- Add a `.env` with
``` dotenv 
TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
TELEGRAM_CHAT_ID=<DESIRED_CHAT_ID> # don't forget the `-` in front of the ID
DEV_CLI_SOURCE_URL=<YOUR_SOURCE_ID>
DEV_CLI_TARGET_URL=http://web:3000/webhook # if using docker, use `web` instead of `localhost`
```

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

## Ansible

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

## Generating github app private key

- Go to github app page
- Click on "generate a private key" in the **Private keys** section
- Save the `.pem` file appropriately
- Generate a base64 representation of the item using `base64 -w 0 -i my.pem > encoded-private-key.txt` or `pbcopy < base64 -w 0 -i my.pem`
- Copy the content into the `GITHUB_APP_PRIVATE_KEY` env variable