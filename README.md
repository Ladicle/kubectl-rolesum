# kubectl-bindrole

Finding Kubernetes Roles bound to a specified ServiceAccount, Group or User.

![screenshot](./img/screenshot.png)

## Installation and Usage

## for macOS

bindrole supports [homebrew](https://docs.brew.sh/Installation) :beer:

```
brew tap Ladicle/kubectl-bindrole
brew install kubectl-bindrole
```

## for other devices

The easiest way is to download binary from the [release page](https://github.com/Ladicle/kubectl-bindrole/releases).
You can also download this repository and install it using Makefile.

```
$ kubectl-bindrole -h  # or kubectl bindrole -h

Usage of kubectl-bindrole:
      --as string                      Username to impersonate for the operation
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --cache-dir string               Default HTTP cache directory (default "/home/ladicle/.kube/http-cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
  -k, --subject-kind string            The Kind of subject which is bound Roles. (default "ServiceAccount")
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
  -v, --version                        Print command version
```

This command works both as a kubectl plugin and as a standalone.
