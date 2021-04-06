**DEPRECATED**: This project has been deprecated in favor of [kubelogin](https://github.com/int128/kubelogin), which does the same thing, only a lot smarter. Check it out!

# k8s-login

A command line tool for authenticating to Kubernetes clusters.

[![Build status](https://travis-ci.org/anton-johansson/k8s-login.svg)](https://travis-ci.org/anton-johansson/k8s-login)


## Credits

Most of the code here is taken from/inspired by other repositories, primarily the [example app from dex](https://github.com/dexidp/dex/tree/master/cmd/example-app). This project is mainly to have a properly released application that can more easily be re-used.


## Building

```shell
$ make
```


## Installation

Linux:

```shell
$ wget https://github.com/anton-johansson/k8s-login/releases/download/v0.0.3/k8s-login-linux-amd64 -O - | sudo tee /usr/local/bin/k8s-login > /dev/null
$ sudo chmod +x /usr/local/bin/k8s-login
```


## Usage

### Requirements

This command line tool requires you to run [dex](https://github.com/dexidp/dex) with the following client configured:

```yaml
    staticClients:
      - id: k8s-login
        name: k8s-login
        secret: lhHN7keNTf4MXEIH3WF4NUL701qITv9Q
        redirectURIs:
          - http://localhost:5555/callback
```

It also requires you to pass in some arguments to `kube-apiserver`, like this:

```
  --oidc-issuer-url=https://dex.svc.example.com \
  --oidc-client-id=k8s-login \
  --oidc-username-claim=email \
  --oidc-groups-claim=groups \
```

### Command line interface

```shell
$ k8s-login auth <server-name>
```


### Assumptions

This CLI currently assumes that your Kubernetes API is reachable on a URL that looks like this:

```
http(s)://k8s.<something>:6443
```

And it expects your Dex to be reachable by replacing `://k8s.` with `://dex.` and removing the port. It also assumes that the OAuth2 client ID and password is as given in the example above. It can always be overriden with the arguments `--client-id` and `--client-secret`, but it makes the CLI a little messier to use.


## License

This project is licensed under Apache 2.0 license. Also refer to the [license of dex](https://github.com/dexidp/dex/blob/master/LICENSE).
