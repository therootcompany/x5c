# [x5c](https://github.com/therootcompany/x5c)

A pass-by-query-param online x509 certificate decoder, written in Go.

-   Demo
-   Usage
-   Install
-   API
-   Build

## Demo

<https://x5c.bnna.net/>

Note: until this is a paid service, there's no guarantee of uptime.

## Usage

```sh
x5c-server --web-root ./public/ --port 8080
```

-   `--web-root` can be used to override index.html

## Install

```sh
go install github.com/therootcompany/x5c/cmd/x5c-server
```

## API

```
GET /api/x509?cert=<pem|rfc-base64|url-base64|hex>
```

```json
{
    "issuer": "CN=example.com",
    "subject": "CN=example.com",
    "serial_number": "588847758906910965042193239346229823033307348532",
    "valid_from": "2024-09-07T07:50:45Z",
    "valid_to": "2025-09-07T07:50:45Z",
    "sha1_fingerprint": "D3480B4AFC14E72DAFE350173995967CC5A10AD6",
    "sha256_fingerprint": "5832B347806271E1B341EF364A5C73849B8CB30EECB66C8A2EC8AC7D04CAD74F"
}
```

Example certificate:

```pem
-----BEGIN CERTIFICATE-----
MIIBQDCB86ADAgECAhRnJNjDOT2MzyPw0YgyhI2ziwS+NDAFBgMrZXAwFjEUMBIG
A1UEAwwLZXhhbXBsZS5jb20wHhcNMjQwOTA3MDc1MDQ1WhcNMjUwOTA3MDc1MDQ1
WjAWMRQwEgYDVQQDDAtleGFtcGxlLmNvbTAqMAUGAytlcAMhAFe8ERQZwaGP7UCi
HDcKnCGI8EzOlqEcuGa502FzqDzdo1MwUTAdBgNVHQ4EFgQUcOjieH3j0OY4nrtP
BdDO4XN/rLEwHwYDVR0jBBgwFoAUcOjieH3j0OY4nrtPBdDO4XN/rLEwDwYDVR0T
AQH/BAUwAwEB/zAFBgMrZXADQQBFOO8xhfMbZ2iJbS/mgkOyund5/FVorMZsu/j5
jCjcURccZOG6+gh1ahTGk20QmRLmE2Cf/+WTpfrAa5l6x6cN
-----END CERTIFICATE-----
```

## Build

```sh
GOOS=linux GOARCH=amd64 go build -o x5c-server ./cmd/x5c-server/
```
