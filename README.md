# mtls

> 🔒 [mTLS](https://en.wikipedia.org/wiki/Mutual_authentication) server and client library.

## Install

```console
$ go get -u -v github.com/picatz/mtls/...
...
```

## Generate Certs

CA Cert and Key

```golang
caCertPEM, caPrivKeyPEM, err := cert.NewCA(
    cert.WithNewECDSAKey(),
    cert.WithCommonName("ca"),
)
```

Server Cert and Key

```golang
caPemReader := bytes.NewReader(caPEM)
caPrivKeyReader := bytes.NewReader(caPrivKeyPEM)

serverCertPEM, serverPrivKeyPEM, err := cert.NewServerFromCA(
    caPemReader,
    caPrivKeyReader,
    cert.WithNewECDSAKey(),
    cert.WithCommonName("server"),
)
```

Client Cert and Key

```golang
caPemReader := bytes.NewReader(caPEM)
caPrivKeyReader := bytes.NewReader(caPrivKeyPEM)

clientCertPEM, clientPrivKeyPEM, err := cert.NewClientFromCA(
    caPemReader,
    caPrivKeyReader,
    cert.WithNewECDSAKey(),
    cert.WithCommonName("client"),
)
```
