# Keybase Slackbot

A slackbot to pgp encrypt and armor text sent via a slack slash command.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## Usage


First, register your public key with the service by running `/pgp init`:

```
> /pgp init
Click here to configure your public key: http://localhost:8000/onboard/nan5a090N-q8nIKs23ZIxfTAkWfb5pthbQyyZMOjQbs=
```

That link lets you configure your public key.

Then, send an encrypted message: `/pgp @user Hey! This message is secret!`

The bot will then encrypt that message with your public key and post it in slack for you:

```
-----BEGIN PGP MESSAGE-----
To-Slack-User: user
Sent-By: slackbot

wcBMA98g6SKTAs8NAQgAKX5gw/rLHiaUBtOZlrmVMXHMgBiwv5KXuDbHghHQzMTS
VoCVd9WRrvPqmLPqtM1aceIVFFEz3rlqcw2Nt0leDYVKkOUVJ7jUIpiD5cGFkp76
wR73Sl5dKttMjwTw5cfJADr+PYZib6suut5f0clj0ZpvgFvzymsULgOyXYlrLdq/
M2YBjflGli9+fv6T0kZbQzYLrv/R/sVaL5jcQIT6YUQqQa9O5VO6ZI6Hx1tZk4qs
pF1F7dz2onCn5R6wZamRsygdfPqb+7R4qbnbf4LH+GYe6PO4X+EW9K7aZnhoPQOR
bMuObwgMtJeb7jd7c498pEgBPEh3qlkF8RPInpiBvtLgAeRfAVU9w8deWG+p4M44
zxHY4V+V4DLgluFqRuAf4hxuWEXgaeQYpT1uQH8QS41Xb2EHYFXS4BTkUDkHa2H7
2noyDYqpTAf4R+KOjN1X4fHBAA==
=33OY
-----END PGP MESSAGE-----
```
