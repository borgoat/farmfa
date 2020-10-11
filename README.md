# farMFA

[![Go Report Card](https://goreportcard.com/badge/github.com/giorgioazzinnaro/farmfa)](https://goreportcard.com/report/github.com/giorgioazzinnaro/farmfa)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/giorgioazzinnaro/farmfa)](https://pkg.go.dev/github.com/giorgioazzinnaro/farmfa)

## Concept

Multi Factor Authentication is usually implemented by using the TOTP standard from OATH.

* A secret key is shared upon activation and stored by the user (usually in an app such as Authenticator) and by the authentication server.  
* Upon login, the user, after providing the credentials, will input a One-Time Password.
  This password is generated applying the TOTP algorithm to the secret key and to the current time.
* The server will generate the same password, and if they match, the user will be able to go through.
  The secret key is never shared again by the user or the server.
  
The generated One-Time Password, as the name suggests, may only be used once
(or more precisely, within a timeframe of around 30 seconds to 90 seconds, depending on the server implementation).

farMFA comes into play in enterprise environments where access to certain accounts should be restricted to very special occasions
(for example, access to the root user of an AWS master account).  
In this context, we can secure the access so that, after the credentials to the account are retrieved,
the second level of authorisation must come from multiple individuals.  
First of all, we apply __Shamir's Secret Sharing__ to the original TOTP secret key, so that at least 3 of 5 holders are needed to reconstruct it.
Additionally, the TOTP secret key is only ever reconstructed in farMFA's server memory, so that no single player ever has to risk losing it.
After having reconstructed the secret, farMFA will then generate one or more OTPs for the user, until the session expires.

## Getting started

farMFA is a client-server application.
Some operations are stateless (the creation of the shares), while managing authentication sessions relies on the server memory.
Nothing is ever persisted on disk for higher security.

```sh
$ farmfa 
Far away MFA

Usage:
  farmfa [command]

Available Commands:
  dealer      Dealers need help from players to retrieve secrets, they initiate sessions
  help        Help about any command
  server      Start the server
  shares      Commands to manage TOTP shares

Flags:
  -a, --address string   The endpoint to the API (default "http://localhost:8080")
  -h, --help             help for farmfa

Additional help topics:
  farmfa player Players are those holding shares and helping a dealer retrieve a secret

```

### Bootstrap

First we need the TOTP secret key, this is only used once to generate the shares.  
This action is stateless so farMFA can be shut afterwards.

```sh
curl -s localhost:8081/shares -d secret_key=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ -d shares=5 -d threshold=3 | jq
```
```json
{
  "shares": [
    "MXNXOP3C$0$TsRiP9SL8opNLwvhy2EbcideGxXvej32NT2TEIy3L7w=58odXHDe6982kApdDvJoDcTb11TdNEa45ZS8ZLQNEvs=",
    "MXNXOP3C$1$xuzV0o0Xzcv3_jJY4fapSeS5w7trBhlNihmo6l2iuV8=lmdRw700HxtHD4teQq9IqvDHGjAvHaAYdOsSe0U60Fo=",
    "MXNXOP3C$2$xNdHwK2pyE1vaWF84RYuJZIOH5PevCdEiFXiYgDODZY=Uzg1QwQCST1nYUcRiVeSAQ_SyrfZsR17WSDEM09khN8=",
    "MXNXOP3C$3$-tM9f1_4IFcUAOSwCQ6tD9KaOBg2ffScLOYOZjmX54Y=xHn9v7_7WO2viNw92gAAxURU3miHES945brpK-ge9iQ=",
    "MXNXOP3C$4$pwDSw_aw5_xczj5PeZBC2aS5dGNbc29EoovIuNI3xpA=36DfwL7yhETAv5fDVKhtGIpdI5qaWBhcCc1PZyZVPj8="
  ]
}
```

###### TODO
- Should obviously return GPG-encrypted shares for the different share holders

### Authentication session

Here's what happens when a user wants to log in.

```sh
curl -s localhost:8081/sessions -d shares=5 -d threshold=3 -d first_share='MXNXOP3C$0$TsRiP9SL8opNLwvhy2EbcideGxXvej32NT2TEIy3L7w=58odXHDe6982kApdDvJoDcTb11TdNEa45ZS8ZLQNEvs='
```
```json
{
  "id": "4UZL34ZJMJWIPFGRUP7SAR6EIMT2U7XQ4DU3CBHV",
  "private": "X5WSDZXDORITYHR5XFMH65FBMMLZLHCWXGT5MXS5"
}
```

Only the ID is shared with the other share holders, private is kept for later.

---

The following two commands would be executed by two separate individuals having access to their own key.  

```sh
curl -s localhost:8081/sessions/4UZL34ZJMJWIPFGRUP7SAR6EIMT2U7XQ4DU3CBHV/shares -d share='MXNXOP3C$1$xuzV0o0Xzcv3_jJY4fapSeS5w7trBhlNihmo6l2iuV8=lmdRw700HxtHD4teQq9IqvDHGjAvHaAYdOsSe0U60Fo='
share joined
```

```sh
curl -s localhost:8081/sessions/4UZL34ZJMJWIPFGRUP7SAR6EIMT2U7XQ4DU3CBHV/shares -d share='MXNXOP3C$2$xNdHwK2pyE1vaWF84RYuJZIOH5PevCdEiFXiYgDODZY=Uzg1QwQCST1nYUcRiVeSAQ_SyrfZsR17WSDEM09khN8='
share joined
```

We can now check that enough shares have been provided:

```sh
curl -s localhost:8081/sessions/4UZL34ZJMJWIPFGRUP7SAR6EIMT2U7XQ4DU3CBHV | jq
```
```json
{
  "complete": true,
  "prefix": "MXNXOP3C"
}
```

---

The original user may now retrieve the TOTP code, this will only be valid for 30 seconds.
After that, this link will no longer be valid and the user will not be able to authenticate anymore without repeating the procedure.
(The expiration of the link could be longer so that fresh TOTP may be retrieved e.g. for 1 hour)

```sh
curl -s localhost:8081/sessions/4UZL34ZJMJWIPFGRUP7SAR6EIMT2U7XQ4DU3CBHV/totp -d private='X5WSDZXDORITYHR5XFMH65FBMMLZLHCWXGT5MXS5'
```
```json
{
  "totp":"862888"
}
```