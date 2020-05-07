# yapi

`yapi` is a very simple API client (a la Postman) which uses a yaml file for configuration.

Start with a yapi.yml file:

```yaml
# routes.yapi.yml
---
root: http://localhost:3000
session: my_project
headers:
  Content-Type: application/json
  Accept: application/json
output:
  headers: true

# Routes

users.index:
  path: /users
  params:
    page: 2
    per: 100

users.show:
  path: /users/1

users.create:
  path: /users
  method: POST
  body: |
    {
      "name": "Test User",
      "email": "testuser@example.com",
    }
```

Specify a routes file and line number and the nearest request specification will be run:

```sh
$ yapi routes.yapi.yml:20
```

or run a route by name:

```sh
$ yapi --name users.show routes.yapi.yml
```


## Installation

You'll need Go, then clone the git repo and run:

```sh
make
make install
```

## Details

Cookies will be used if a `session` name is provided in the settings.

### Available Settings

`root`: The root path for requests in this project.

`headers`: Headers to apply to every request. Requests can overwrite these when needed.

`session`: A name to store cookies for. Session cookies are saved to `~/.cache/yapi/session_name.jar.json`

`output`: Toggles whether to print various items in the output. Currently only `headers: true|false` is supported.

### Available Request Options

`path`: Path to append to `root` path defined in project settings.

`method`: HTTP method type.

`headers`: Define additional headers here, they will be merged with any headers defined on the project level.

`params`: Add query params to the request.

`body`: Add a body to the request.

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/dbridges/yapi.
