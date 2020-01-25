# CA Key and Self-Signed Cert

Using the `openssl` command we can generate a CA key file `ca.key` and the self-signed CA x509 cert file `ca.crt` in one command. Some guides may have these two commands as separate operations.

```console
$ openssl req -new -x509 \
    -nodes \
    -newkey rsa:4096 \
    -keyout ca.key \
    -sha256 \
    -days 365 \
    -out ca.crt \
    -subj "/C=US/ST=MI/L=AnnArbor/CN=ssh.ca/emailAddress=kgruber1@emich.edu"
```

* The `ca.key` is a root key, and would need to protected in a production deployment.
* The `ca.crt` which was generated is the [root certificate](https://en.wikipedia.org/wiki/Root_certificate) that needs to be distributed to each of the computers we want to establish trust with.

## SSH Server Key and Self-Signed Cert

Now we need to generate keys and certifcates for each for our SSH servers we wish to deploy.

```console
$ openssl req -new \
    -nodes \
    -newkey rsa:4096 \
    -keyout server.name.key \
    -sha256 \
    -days 365 \
    -subj "/C=US/ST=MI/O=PicatLabs/CN=ssh.server.name" \
    -reqexts SAN \
    -config <(cat /etc/ssl/openssl.cnf \
        <(printf "\n[SAN]\nsubjectAltName=IP:127.0.0.1")) \
    -out server.name.csr
```

We can then verify the content of the CSR:

```console
$ openssl req -in server.name.csr -noout -text
Certificate Request:
    Data:
        Version: 0 (0x0)
        Subject: C=US, ST=MI, O=PicatLabs, CN=ssh.server.name
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (4096 bit)
                Modulus:
                    00:bc:66:25:0a:55:9d:04:ab:1b:51:2e:58:6f:4d:
                    7d:fc:bd:5b:7c:1f:2d:6e:b1:66:8b:da:cf:37:47:
                    f4:48:2e:23:bc:03:99:3f:fa:64:eb:5b:93:97:59:
                    99:dc:3b:19:09:f3:df:c7:ae:f4:da:b6:f9:01:f7:
                    1d:20:81:d4:be:30:bc:80:d7:74:ca:55:d3:e5:bd:
                    fc:a2:a4:80:8b:de:27:60:85:84:e8:db:b0:12:fc:
                    a2:91:3c:dc:f3:83:9c:84:12:56:14:2e:32:1a:87:
                    1d:b5:4c:a8:3c:e7:77:07:1d:a4:e1:4b:c5:43:9f:
                    86:bc:48:8c:f4:2b:53:f0:e6:7b:6d:32:7d:5d:7a:
                    f7:c2:37:7b:e1:2a:b7:a7:ac:1f:06:3a:99:21:84:
                    85:57:43:b0:2e:c9:b2:62:ad:bb:26:49:5c:53:aa:
                    f1:e6:9b:fc:e0:7a:0b:cc:f4:5d:55:f5:b2:54:39:
                    5c:b7:d7:29:d8:ea:e0:8c:2a:f2:0f:db:66:6d:2e:
                    f7:ea:63:28:5b:a2:c0:0c:4d:9a:e5:11:72:d3:e1:
                    20:a4:16:a4:f2:ed:a8:b4:76:62:dc:90:22:8f:e4:
                    be:75:f0:94:b6:66:5b:cf:ea:5d:a5:ca:d6:a2:ac:
                    e0:bd:ed:b9:e4:bd:68:05:df:ce:d1:e3:f8:77:4a:
                    9a:54:04:2b:d1:86:00:60:9d:75:4f:4e:8f:71:e7:
                    ed:cd:8d:99:aa:d1:82:e2:f1:b6:93:8c:b4:a1:01:
                    81:74:ed:1f:e5:6d:54:e7:a5:49:63:a0:24:e8:83:
                    73:fa:e7:bc:25:ef:ca:72:7d:c2:cd:ae:3e:50:a9:
                    f2:4e:a3:8e:cf:d6:80:b1:ed:14:f0:cb:93:8f:6d:
                    0d:1f:1b:d5:64:3d:6a:3c:e6:13:f2:15:c9:ca:46:
                    36:2a:4c:55:0a:bf:eb:b7:d1:96:bc:84:b3:be:b0:
                    34:e0:12:01:c5:be:f3:de:01:c9:ac:33:62:34:ed:
                    7a:ad:77:56:c0:16:d3:9f:20:4c:6b:f3:6c:e8:1f:
                    87:d3:f4:4b:a8:1e:7c:58:2d:c5:ad:c6:e4:8e:a5:
                    c1:b8:36:32:f6:4a:c9:26:16:24:19:6d:a6:db:85:
                    e9:be:b7:d1:cd:ed:d3:d1:e6:9e:80:85:0f:a5:40:
                    5d:af:34:1f:19:ad:ac:77:f9:87:b5:fe:99:e0:2a:
                    1a:ad:1b:b4:68:13:9a:20:c4:49:e3:c9:33:2f:c4:
                    4f:e8:20:46:9a:8b:c7:c7:69:66:fa:15:54:3b:07:
                    3f:0e:2f:1a:7b:3c:9e:13:f1:df:b8:d6:1a:26:73:
                    a4:e2:84:b7:6e:55:78:34:4d:dc:0e:79:52:a3:34:
                    33:17:31
                Exponent: 65537 (0x10001)
        Attributes:
        Requested Extensions:
            X509v3 Subject Alternative Name:
                DNS:ssh.server.name
    Signature Algorithm: sha256WithRSAEncryption
         ab:00:fe:03:53:51:cb:c8:7a:9b:9b:5f:ad:23:9e:a4:c2:3f:
         7b:d8:a2:71:6d:7d:0c:c1:86:38:c9:97:1d:b8:bb:03:8a:26:
         28:81:c5:cc:68:fc:90:21:a9:32:2a:ec:4e:59:12:99:d9:79:
         51:92:30:27:2b:fc:a2:14:ed:63:23:f9:e1:6a:1c:3c:8b:98:
         b9:f3:0c:5c:0f:bb:57:a4:ea:ff:b2:a6:d9:85:8c:e9:92:b9:
         26:61:1b:65:88:48:d7:60:29:9f:98:14:85:b5:34:f7:25:7a:
         2c:23:ee:44:3b:8b:55:86:15:2e:33:eb:f4:4e:a2:14:5b:c5:
         a1:88:ac:28:a2:05:c3:18:77:9a:47:97:5c:d1:42:1a:6c:fd:
         16:2a:9e:d0:a9:ee:fb:ea:0f:0f:1d:22:12:e5:e0:39:f9:6b:
         a1:90:d1:e1:a8:db:81:a4:5c:25:bf:bb:b7:7d:47:53:15:c5:
         92:cf:28:e1:b2:a5:e7:65:70:bc:02:01:1a:9b:15:f1:56:0a:
         18:08:0e:d1:97:42:2e:82:03:ba:f1:74:5a:45:cd:9c:43:e7:
         23:bc:55:7d:16:7e:bc:f7:26:2f:f9:88:4e:9a:2f:c8:3c:ad:
         a9:72:33:d3:48:25:b8:d0:fc:a5:96:fb:89:f6:88:36:c8:fe:
         8a:6a:76:6a:0e:6d:cf:f0:f2:a1:62:f3:df:02:0b:44:ed:eb:
         ca:b8:60:d9:a1:ee:e3:d2:f5:90:f0:47:dc:f5:05:1c:71:ed:
         b0:41:ae:0d:8b:60:8d:1b:73:4a:59:7f:ba:c9:b8:32:ba:98:
         57:39:f2:26:6d:80:cb:cc:75:1f:e1:14:99:c3:15:03:c0:57:
         dd:6b:36:72:f7:2a:05:e8:9d:c5:9e:80:8a:0d:de:71:0e:c4:
         7d:69:33:0f:93:f6:f1:88:c4:27:4d:9d:4e:77:f5:3a:b9:87:
         62:bf:a5:d8:e3:c0:41:7e:0a:09:04:43:f7:20:fd:0d:86:98:
         fc:c4:cb:5a:9d:1e:64:b1:d3:bd:ba:52:0f:18:07:41:ab:61:
         ed:9a:32:84:cf:df:12:44:00:84:b4:c2:23:05:16:ee:3d:3f:
         3d:c0:f8:aa:64:e0:5f:c5:e0:a4:6e:38:37:14:2d:80:44:94:
         73:dd:36:b8:d5:ef:8e:90:17:a0:c6:ae:9b:9e:fc:c2:63:3e:
         db:bb:39:57:14:37:31:ac:7e:1a:15:58:40:4b:85:21:80:69:
         37:24:b6:65:d4:53:e9:44:06:99:6b:3f:53:b6:c3:40:05:6b:
         99:de:5a:c1:5d:fc:3d:51:bf:37:f7:5f:5f:76:1a:95:54:85:
         aa:09:a4:6e:6d:c3:f1:85
```

We then use the root `ca.key`/`ca.crt` and the `server.name.csr` to generate the server's `server.name.crt` file:

```console
$ openssl x509 -req \
    -in server.name.csr \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out server.name.crt \
    -days 500 \
    -sha256
Signature ok
subject=/C=US/ST=MI/O=PicatLabs/CN=ssh.server.name
Getting CA Private Key
```

Then verify the certificate's content:

```console
$ openssl x509 -in server.name.crt -text -noout
Certificate:
    Data:
        Version: 1 (0x0)
        Serial Number: 16033129432805395214 (0xde811e45a039cb0e)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=US, ST=MI, L=AnnArbor, CN=ssh.ca/emailAddress=kgruber1@emich.edu
        Validity
            Not Before: Jan  1 02:20:46 2020 GMT
            Not After : May 15 02:20:46 2021 GMT
        Subject: C=US, ST=MI, O=PicatLabs, CN=ssh.server.name
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (4096 bit)
                Modulus:
                    00:bc:66:25:0a:55:9d:04:ab:1b:51:2e:58:6f:4d:
                    7d:fc:bd:5b:7c:1f:2d:6e:b1:66:8b:da:cf:37:47:
                    f4:48:2e:23:bc:03:99:3f:fa:64:eb:5b:93:97:59:
                    99:dc:3b:19:09:f3:df:c7:ae:f4:da:b6:f9:01:f7:
                    1d:20:81:d4:be:30:bc:80:d7:74:ca:55:d3:e5:bd:
                    fc:a2:a4:80:8b:de:27:60:85:84:e8:db:b0:12:fc:
                    a2:91:3c:dc:f3:83:9c:84:12:56:14:2e:32:1a:87:
                    1d:b5:4c:a8:3c:e7:77:07:1d:a4:e1:4b:c5:43:9f:
                    86:bc:48:8c:f4:2b:53:f0:e6:7b:6d:32:7d:5d:7a:
                    f7:c2:37:7b:e1:2a:b7:a7:ac:1f:06:3a:99:21:84:
                    85:57:43:b0:2e:c9:b2:62:ad:bb:26:49:5c:53:aa:
                    f1:e6:9b:fc:e0:7a:0b:cc:f4:5d:55:f5:b2:54:39:
                    5c:b7:d7:29:d8:ea:e0:8c:2a:f2:0f:db:66:6d:2e:
                    f7:ea:63:28:5b:a2:c0:0c:4d:9a:e5:11:72:d3:e1:
                    20:a4:16:a4:f2:ed:a8:b4:76:62:dc:90:22:8f:e4:
                    be:75:f0:94:b6:66:5b:cf:ea:5d:a5:ca:d6:a2:ac:
                    e0:bd:ed:b9:e4:bd:68:05:df:ce:d1:e3:f8:77:4a:
                    9a:54:04:2b:d1:86:00:60:9d:75:4f:4e:8f:71:e7:
                    ed:cd:8d:99:aa:d1:82:e2:f1:b6:93:8c:b4:a1:01:
                    81:74:ed:1f:e5:6d:54:e7:a5:49:63:a0:24:e8:83:
                    73:fa:e7:bc:25:ef:ca:72:7d:c2:cd:ae:3e:50:a9:
                    f2:4e:a3:8e:cf:d6:80:b1:ed:14:f0:cb:93:8f:6d:
                    0d:1f:1b:d5:64:3d:6a:3c:e6:13:f2:15:c9:ca:46:
                    36:2a:4c:55:0a:bf:eb:b7:d1:96:bc:84:b3:be:b0:
                    34:e0:12:01:c5:be:f3:de:01:c9:ac:33:62:34:ed:
                    7a:ad:77:56:c0:16:d3:9f:20:4c:6b:f3:6c:e8:1f:
                    87:d3:f4:4b:a8:1e:7c:58:2d:c5:ad:c6:e4:8e:a5:
                    c1:b8:36:32:f6:4a:c9:26:16:24:19:6d:a6:db:85:
                    e9:be:b7:d1:cd:ed:d3:d1:e6:9e:80:85:0f:a5:40:
                    5d:af:34:1f:19:ad:ac:77:f9:87:b5:fe:99:e0:2a:
                    1a:ad:1b:b4:68:13:9a:20:c4:49:e3:c9:33:2f:c4:
                    4f:e8:20:46:9a:8b:c7:c7:69:66:fa:15:54:3b:07:
                    3f:0e:2f:1a:7b:3c:9e:13:f1:df:b8:d6:1a:26:73:
                    a4:e2:84:b7:6e:55:78:34:4d:dc:0e:79:52:a3:34:
                    33:17:31
                Exponent: 65537 (0x10001)
    Signature Algorithm: sha256WithRSAEncryption
         55:69:ed:16:5b:a6:2b:c2:d7:78:be:1b:75:b6:3a:40:54:f0:
         3c:d3:3f:f1:f2:94:6e:41:f0:01:cf:d2:48:8e:1c:29:38:b7:
         60:4e:3d:7b:60:fd:6e:d8:6f:e2:a5:51:15:50:f0:37:a8:03:
         9d:8a:57:88:8d:f7:c0:d1:7a:09:cb:c8:4a:f7:29:3f:37:b7:
         d3:c8:de:2c:b9:e1:e3:db:e9:24:cf:84:cd:bf:35:5e:30:44:
         f4:3f:f2:37:de:f6:f6:d4:0f:ec:f2:9c:c5:b6:9d:cb:ab:d7:
         55:25:e1:51:f2:7e:fd:dd:e6:ea:fc:b6:d5:c3:5f:f0:88:c2:
         fd:05:18:7f:5d:4a:4d:c2:c7:30:f6:4b:7a:c9:16:76:de:23:
         8e:df:7f:f1:7d:c8:9e:ff:d5:2c:01:85:4c:b8:d0:66:b4:dd:
         a1:93:18:39:34:fe:8a:e2:c4:ab:42:c4:65:40:4c:3a:df:76:
         17:77:29:e7:57:13:98:0b:1d:2f:70:62:9a:24:02:63:9d:08:
         f2:3e:20:6e:ca:66:4e:5d:0f:3b:f7:fa:0b:3d:6b:51:f9:10:
         ce:dd:86:a4:fd:c4:6a:1c:3c:ef:9d:3e:c0:d5:a0:23:36:6b:
         53:c3:b5:01:9a:7e:b8:f4:ab:87:89:0a:b3:87:6c:6c:8d:cc:
         2b:15:fd:20:1f:87:de:0f:d2:26:e3:91:9d:37:48:51:a2:79:
         9a:2b:86:0e:e8:c4:10:fa:23:52:9d:c8:ad:c3:1b:2d:44:78:
         d4:1a:d8:26:44:b4:e8:82:47:b1:73:8f:67:16:9f:97:74:7c:
         4b:85:c0:bc:0b:77:b2:cd:42:8b:47:6d:21:b1:a5:dd:85:74:
         ac:f2:03:ef:f7:30:7b:d7:57:1b:27:b3:af:57:8d:88:27:93:
         d5:ac:1b:a0:89:75:e9:10:76:0f:0c:04:35:67:b8:56:79:e9:
         8d:09:1d:01:3f:f6:ee:6d:3f:ba:02:b5:23:6b:4b:9d:a4:5e:
         4d:75:a3:ed:6d:c1:92:d9:c7:af:c1:9c:f2:db:35:21:f6:7e:
         76:72:20:27:d9:f2:6a:d2:2c:0c:3a:08:c5:1f:e2:1f:f6:c3:
         63:c8:e9:bf:04:a0:de:ee:db:de:55:c5:da:87:c8:a4:85:eb:
         23:5b:30:16:38:7a:3a:44:ae:57:e8:46:31:93:84:91:d4:13:
         c3:ec:89:86:04:b2:33:08:35:db:8f:77:7f:f5:fd:fc:3d:c3:
         7c:eb:f7:d5:14:c0:8f:4f:57:c8:98:30:38:7c:0b:56:dd:97:
         5d:d6:e8:08:c0:d6:f4:9b:b8:08:ef:2b:0e:5c:a2:76:8a:ff:
         a4:65:d1:d7:ee:b7:9d:3c
```

* `server.name.crt`
* `server.name.csr`
* `server.name.key`

## SSH Client Key and Self-Signed Cert

```console
$ openssl req -new \
    -nodes \
    -newkey rsa:4096 \
    -keyout client.name.key \
    -sha256 \
    -days 365 \
    -subj "/C=US/ST=MI/O=PicatLabs/CN=ssh.client.name" \
    -reqexts SAN \
    -config <(cat /etc/ssl/openssl.cnf \
        <(printf "\n[SAN]\nsubjectAltName=IP:127.0.0.1")) \
    -out client.name.csr
```

```console
$ openssl x509 -req \
    -in client.name.csr \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out client.name.crt \
    -days 500 \
    -sha256
Signature ok
subject=/C=US/ST=MI/O=PicatLabs/CN=ssh.client.name
Getting CA Private Keyy
```
