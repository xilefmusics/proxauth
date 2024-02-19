# Proxauth

Proxauth is a user-based authentication proxy, which can be used to protect and rewrite different URLs and paths.
It uses JWT tokens stored inside cookies to save the session state and is therefore completely stateless.

## Usage

Since this app is only an proxy it needs an application behind it it can point to.
An example on how to use it with an application can be seen in the [money-app](https://github.com/xilefmusics/money-app/blob/main/docker-compose.yaml).
The docker images for this application are prebuild for all the releases on [DockerHub](https://hub.docker.com/repository/docker/xilefmusics/proxauth).
If you want to test/develop the application you can start it using the following commands:

```bash
cd ./src
go get .
go run .
```

More detailed information can be found in the [Dockerfile](https://github.com/xilefmusics/proxauth/blob/main/Dockerfile).

## Configuration

Proxauth is completely configurable through environment variables.

|Variable|Default|Explanation|
|-|-|-|
|PORT|`8080`|The port inside the container, proxauth listenes at.|
|SERVER_SECRET|`changeMe`|The secret which is used to encrypt the JWT tokens. Change that if you want to kill current sessions.|
|JWT_EXPIRATION_DURATION|`24h`|Time span that defines how long the JWT tokens and therefore the sessions are valid (note that hour is the biggest unit)|
|CONFIG_FILE|`/config/config.yaml`|Location of the config file which get's used if the `CONFIG` variable is not set.|
|CONFIG||Configfile as string. If set the variable `CONFIG_FILE` gets ignored.|

## Config File

The config file can be passed as a file or directly as a string string.
It is in YAML format and consists of two different sections, the users and the rules.

### Users

A user consists of a `username`, a SHA-256 encrypted `password` and a salt `salt`.

### Rules

A rule defines one mapping of an URL to an other.
Optional a protection for this route can be created.

|Param|Default|Explanation|
|-|-|-|
|fromScheme|`http`|The scheme which is used for redirects and logging. It is not responisble for matching the rule.|
|fromHost|`*`|Selects the rule based on the host. Use `*` to match all hosts.|
|fromPath|`/`|The subpath which is used to match the rule.|
|toScheme|`http`|The url scheme which is used to forward the request.|
|toHost|`localhost`|The host the request is forwarded to.|
|toPort||The port the request is forwarded to.|
|toPath|`/`|The path the request is forwarded to.|
|loginPath|`/login`|The path which is used for login. (GET & POST is blocked)|
|loginPath|`/logout`|The path which points to the logout page. (GET is blocked)|
|allowedUsers|`[]`|List of the users that are allowed to access the endpoints. Use `[]` if authentication is disabled and all users can access it.|
|redirectToLogin|`false`|If this is set to true the unauthorized respone is catched and a redirection to the loginpage is made instead.|
|backgroundColor|`#000000`|Color scheme used in html pages surfed by proxauth itself.|
|textColor|`#002200`|Color scheme used in html pages surfed by proxauth itself.|
|primaryColor|`#002200`|Color scheme used in html pages surfed by proxauth itself.|
|title|`Proxauth`|Title of the html pages surfed by proxauth itself.|
|redirects|`empty`|A map of redirects.|

### Example

This config file includes one user and one rule.
This rule matches all of the incomming requests and forwards them to `localhost:8080`.
The route is protected and only the one user has access to it.

```yaml
users:
- username: <username>
  password: sha-256(<password><salt>)
  salt: <salt>

rules:
- toPort: 8080
  allowedUsers: [<username>]
  redirectToLogin: true
```

## License

[![GPL-3.0](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
