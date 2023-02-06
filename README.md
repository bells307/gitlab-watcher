# gitlab-watcher
**gitlab-watcher** is the service that processes hooks from Gitlab and sends notifications to the corresponding services (for example, *telegram*)

# Launch
To launch the service use the following command:
```
gitlab-watcher -config=<CONFIG_PATH>
```
`-config` is an optional field, the default value is `./config.yml`

## docker-compose
To start the service using `docker-compose` you need to create a file `.env` with environment variables like:
```bash
# http listen address
GH_HTTP_LISTEN=0.0.0.0:8888
# config file path
GH_CONFIG_FILE=./config.yml
# messages template directory for telegram
GH_TELEGRAM_TEMPLATES_DIR=./templates
```

`docker-compose` configuration:
```yaml
version: "3.8"

services:
  bot:
    image: ghcr.io/bells307/gitlab-watcher
    command: gitlab-watcher -config=/opt/gitlab-watcher/config.yml
    volumes:
      - ${GH_CONFIG_FILE}:/opt/gitlab-watcher/config.yml
      - ${GH_TELEGRAM_TEMPLATES_DIR}:/opt/gitlab-watcher/templates/telegram
    ports:
      - "${GH_HTTP_LISTEN}:${GH_HTTP_LISTEN}"
```

# Configuration
The service configuration is described using a `yaml` configuration file. The default configuration file name is `config.yml`
```yaml
# http server endpoint
listen: 0.0.0.0:8888
# send notifications on template processing errors
send_template_errors: true
# telegram bot configuration
telegram:
  # chats in which bot will send messages
  # (optional, that field will fill when bot adds/removes from chat)
  chats:
    - -123456789
  # mapping telegram users to gitlab users (for mapUser() template function)
  users:
    # telegram name - gitlab name
    tguser: gitlabuser
  # show debug messages
  debug: true
  # telegram message parse mode (HTML, MARKDOWN)
  parse_mode: HTML
  # templates directory
  templates: ./templates
  # bot token
  token: YOUR_TOKEN_HERE
```

# Message templates
`gitlab-watcher` uses [golang text templates](https://pkg.go.dev/text/template) for sending messages to the corresponding services. Templates must be in specific directory, described in the configuration for each sender service
```yaml
telegram:
  # ...
  templates: ./templates
  # ...
```
The template text file must have a specific name in a directory for every Gitlab hook type:
- Merge Request: `merge-request.txt`
- Pipeline: `pipeline.txt`

Template files use the same structure as [Gitlab Webhook API](https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html) objects.

You don't need to restart the service after editing the template.

Text template example:
```
{{ if and (eq .ObjectAttributes.State "opened") (eq .ObjectAttributes.Action "open") (not (hasPrefix .ObjectAttributes.Title "DRAFT:")) }}
🔍 {{ .User.Name }} created Merge Request "{{ .ObjectAttributes.Title }}" for the project {{ .Project.Name }}:
{{ .ObjectAttributes.URL }}
{{ else if and (eq .ObjectAttributes.State "merged") (eq .ObjectAttributes.Action "merge") }}
✅ {{ .User.Name }} merged "{{ .ObjectAttributes.Title }}" for the project {{ .Project.Name }}:
{{ .ObjectAttributes.URL }}
{{ end }}
```
Another examples can be found in the repository directory `./templates`