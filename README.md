<div id="header">

# farMFA

<div class="details">

<span id="author" class="author">Giorgio Azzinnaro</span>  
<span id="email" class="email"><giorgio@azzinna.ro></span>  

</div>

<div id="toc" class="toc">

<div id="toctitle">

Table of Contents

</div>

-   [Concept](#_concept)
-   [Getting started](#_getting_started)
    -   [Split TOTP secret and share](#_split_totp_secret_and_share)
    -   [Generate a TOTP](#_generate_a_totp)
-   [References](#_references)
-   [Glossary](#_glossary)

</div>

</div>

<div id="content">

<div id="preamble">

<div class="sectionbody">

<div class="paragraph">

<span class="image"><a href="LICENSE" class="image"><img
src="https://img.shields.io/github/license/borgoat/farmfa?color=blue&amp;style=flat-square"
alt="License" /></a></span> <span
class="image"><a href="https://goreportcard.com/report/github.com/borgoat/farmfa"
class="image"><img
src="https://goreportcard.com/badge/github.com/borgoat/farmfa"
alt="Go Report Card" /></a></span> <span
class="image"><a href="https://pkg.go.dev/github.com/borgoat/farmfa"
class="image"><img
src="https://pkg.go.dev/badge/github.com/borgoat/farmfa"
alt="PkgGoDev" /></a></span>

</div>

</div>

</div>

<div class="sect1">

## Concept

<div class="sectionbody">

<div class="paragraph">

Multi Factor Authentication is often implemented by using the TOTP
standard [\[RFC6238\]](#RFC6238) from OATH.

</div>

<div class="ulist">

-   The authentication server generates a secret key, stores it, and
    shows it to the user as a QR code or a Base32 string. The user
    stores it usually in a mobile app, or in a password manager.

-   Upon login, the user inputs a One-Time Password. The authenticator
    app or password manager generates this password by applying the TOTP
    algorithm to the secret key and the current time.

-   The authentication server performs the same algorithm on the secret
    key and current time. If the output matches the user’s TOTP - the
    process is successful.

</div>

<div class="paragraph">

The generated One-Time Password, as the name suggests, may only be used
once (or more precisely, within a certain timeframe, depending on the
server implementation).

</div>

<div class="paragraph">

farMFA comes into play in shared environments where access to certain
accounts should be restricted to very special occasions. For example,
access to the root user of an AWS account, especially the root user of
the management account of an AWS Organization, which should only happen
in break-glass scenarios.

</div>

<div class="paragraph">

In this context, we want to restrict access in such a way that multiple
individuals are needed to grant authorisation.

</div>

<div class="paragraph">

First of all, we apply the *Shamir’s Secret Sharing* scheme [\[2\]](#2)
to the original TOTP secret key. This means, for instance, that at least
3 of 5 holders must put together their shares to reconstruct the secret.

</div>

<div class="paragraph">

Additionally, farMFA implements a workflow to reassemble the TOTP secret
in a server, and letting users only access the generated TOTP code. This
way, no single player has to risk accessing and accidentally
leaking/persisting the secret.

</div>

</div>

</div>

<div class="sect1">

## Getting started

<div class="sectionbody">

<div class="paragraph">

The two main workflows are:

</div>

<div class="ulist">

-   getting the TOTP [secret](#secret), splitting it, and sharing it to
    multiple parties ([players](#player)) - this is done locally.

-   putting together the [Tocs](#Toc) and generating the TOTP - done by
    the [server](#server)/[oracle](#oracle).

</div>

<div class="sect2">

### Split TOTP secret and share

<div class="paragraph">

During this phase, farMFA MUST encrypt [Tocs](#Toc) based on the
intended recipient/[player](#player). The current encryption strategy is
based on [age](https://filippo.io/age).

</div>

<div class="paragraph">

Each player generates their age keypair via
[age-keygen](https://github.com/FiloSottile/age#readme):

</div>

<div class="listingblock">

<div class="content">

``` highlight
$ age-keygen -o your_own_key.txt
```

</div>

</div>

<div class="paragraph">

Players then share their public key with the [dealer](#dealer).

</div>

<div class="paragraph">

Now the dealer may start the process. Usually, the dealer is also a
player, so they also provide their own age public key and keep one Toc
for themselves.

</div>

<div class="listingblock">

<div class="content">

``` highlight
$ farmfa dealer \
  --totp-secret HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ \
  -p player_1=age174uc3d7qzaulmm75huxazcaynq5et4s9gdr7ajnau204lqny7asq9ens77 \
  -p player_2=age1pa8m7l4sguj84c5v8qu9gr3mydnmhd8lf633ln2udlp5699uvp2sm2mpzd \
  -p player_3=age1qegd5t5ajlqzewruz38srlrz05w7xgzq4nn2n5gky872up553v8sl64j3j
```

</div>

</div>

<div class="paragraph">

This will yield encrypted Tocs, 1 per player.

</div>

<div class="paragraph">

Each player now receives their Toc. They can decrypt it and inspect it
via the age CLI.

</div>

<div class="listingblock">

<div class="content">

``` highlight
$ age -i your_own_key.txt --decrypt
[paste the encrypted Toc - Ctrl+D]

{
    "group_id": "7GCUCI2Y",
    "group_size": 3,
    "group_threshold": 3,
    "note":"",
    "share":"C2iCgb3pRfxPJw2a7od8p4ShkhrDWAm/Dt6ioQNAVFPZ",
    "toc_id":"5oaAUX9b6aBE"
}
```

</div>

</div>

<div class="paragraph">

Players must now store their Toc securely.

</div>

</div>

<div class="sect2">

### Generate a TOTP

<div class="paragraph">

When a user wants to log in they can start a [session](#session). The
user who wants to log in is called [applicant](#applicant).

</div>

<div class="listingblock">

<div class="title">

Applicant

</div>

<div class="content">

``` highlight
$ http --body POST localhost:8080/sessions toc_zero:='{"group_id":"J7UHQPZK","group_size":5,"group_threshold":2,"share":"5Ovpu-PKEeYXx5ebiQhzU_AT0Z79POf8GGkskDp3its=urkBkVXr-pYjIvTt1ch2YJILCScAoRquLoX_VBxxps4=","toc_id":"TFW52GAK"}'
{
    "complete": false,
    "created_at": "2021-02-24T18:05:53.507396809+01:00",
    "id": "V5K6QD4XUFLRGCZH",
    "kek": "MIotBtYOWrXnQCj6o9rSNIkNeRfIPhNLjEdQtJDDemPRJcKUbme+iq5K2Hc6Ypil6Loi/K9rnN/YrJiKDT/tPi8kFq2WuAY8zl8=",
    "tek": "age1cl5ndmdsq09vs09awlpt8nd4cdu6fpl33lpyyuv75syknqalkpdszwnwyc",
    "toc_group_id": "J7UHQPZK",
    "tocs_in_group": 5,
    "tocs_provided": 1,
    "tocs_threshold": 2
}
```

</div>

</div>

<div class="paragraph">

The [oracle](#oracle) returns:

</div>

<div class="ulist">

-   a session ID.

-   a *Toc Encryption Key (TEK)* - a public key to encrypt individual
    Tocs, so that only the server may use them.

-   a *Key Encryption Key (KEK)* - to decrypt the private part of the
    TEK, so that the server can only decrypt the Tocs when the applicant
    requests it.

</div>

<div class="paragraph">

The applicant shares the *TEK* and session ID with team members who hold
the other [Tocs](#Toc). Those team members who can authorise the
applicant will be named the session’s [constituents](#constituent).

</div>

<div class="paragraph">

Constituents must encrypt and armor their Toc with TEK (once again,
using age).

</div>

<div class="listingblock">

<div class="title">

Constituent

</div>

<div class="content">

``` highlight
$ export ENCTOC=$(echo '{"group_id":"J7UHQPZK","group_size":5,"group_threshold":2,"share":"zxRrozuUaCMgn_u6ajZStlV7RKwhp0keT9aQoXAEruI=nfx2CPJfKiFM32zLmtxHjV94OlZOgBevV1Whrx-lslU=","toc_id":"K5FSSJSV"}' | age -r age1cl5ndmdsq09vs09awlpt8nd4cdu6fpl33lpyyuv75syknqalkpdszwnwyc -a)
```

</div>

</div>

<div class="paragraph">

Constituents upload the encrypted Toc to the [oracle](#oracle),
associating it with the existing session.

</div>

<div class="listingblock">

<div class="title">

Constituent

</div>

<div class="content">

``` highlight
$ http POST localhost:8080/sessions/V5K6QD4XUFLRGCZH/tocs encrypted_toc="$ENCTOC"
HTTP/1.1 200 OK
```

</div>

</div>

<div class="paragraph">

Once the oracle has enough Tocs, the applicant may query the oracle. The
applicant must provide the *KEK* to let the oracle decrypt the Tocs, and
generate the [TOTP](#TOTP).

</div>

<div class="listingblock">

<div class="title">

Constituent

</div>

<div class="content">

``` highlight
$ http --body POST localhost:8080/sessions/V5K6QD4XUFLRGCZH/totp kek="MIotBtYOWrXnQCj6o9rSNIkNeRfIPhNLjEdQtJDDemPRJcKUbme+iq5K2Hc6Ypil6Loi/K9rnN/YrJiKDT/tPi8kFq2WuAY8zl8="
{
    "totp": "824588"
}
```

</div>

</div>

</div>

</div>

</div>

<div class="sect1">

## References

<div class="sectionbody">

<div class="ulist bibliography">

-   <span id="RFC6238"></span>\[RFC6238\] M’Raihi, D., Machani, S., Pei,
    M., and J. Rydell, "TOTP: Time-Based One-Time Password Algorithm",
    RFC 6238, DOI 10.17487/RFC6238, May 2011,
    <a href="https://www.rfc-editor.org/info/rfc6238"
    class="bare">https://www.rfc-editor.org/info/rfc6238</a>.

-   <span id="SSS"></span>\[2\] Adi Shamir. 1979. How to share a secret.
    Commun. ACM 22, 11 (Nov. 1979), 612–613.
    DOI:https://doi.org/10.1145/359168.359176

</div>

</div>

</div>

<div class="sect1">

## Glossary

<div class="sectionbody">

<div class="dlist">

<span id="secret"></span>secret  
A TOTP is a hash generated from a secret. This secret is usually shown
as a QR code and shared between the prover and verifier. In farMFA, the
prover becomes a distributed entity: recipients who share the key
material, and an oracle that actually generates the TOTP.

<span id="Toc"></span>Toc  
The "pieces" in which a TOTP secret gets split.

<span id="deal"></span>deal  
The workflow in which a dealer splits a secret in Tocs and shares them
with multiple players.

<span id="dealer"></span>dealer  
Creates Tocs from a secret, and shares them with players.

<span id="player"></span>player  
During the Tocs creation phase, the individuals who each receive one of
said Tocs.

<span id="session"></span>session  
Describes the workflow in which an applicant requires combining Tocs to
generate a TOTP.

<span id="applicant"></span>applicant  
Initiates a session to request access to a TOTP.

<span id="constituent"></span>constituent  
The individuals who join a session to authorise an applicant to generate
a TOTP, by reaching a quorum/threshold.

<span id="oracle"></span>oracle  
The entity that reconstructs Tocs into TOTP secrets, and generates
one-time passwords. Also called the *prover*, as defined in
[\[RFC6238\]](#RFC6238).

<span id="server"></span>server  
In our context synonym with *oracle*.

<span id="TOTP"></span>TOTP  
As defined in [\[RFC6238\]](#RFC6238): "an extension of the One-Time
Password (OTP) algorithm \[…​\] to support the time-based moving factor".
Used by many applications as a second authentication factor.

</div>

</div>

</div>

</div>

<div id="footer">

<div id="footer-text">

Last updated 2022-07-17 16:01:53 +0200

</div>

</div>
