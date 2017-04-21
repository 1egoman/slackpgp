# Slack PGP

A slackbot to pgp encrypt and armor text sent via a slack slash command.

## Installation

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

After deployment, click on `View`. You'll be brought to a setup guide that'll walk you through
setting up slack to work with this slackbot.

## Usage

First, register your public key with the service by running `/pgp init`:

```
> /pgp init
Click here to configure your public key: http://localhost:8000/onboard/nan5a090N-q8nIKs23ZIxfTAkWfb5pthbQyyZMOjQbs=
```

Click on that link ad paste in your public key. We'll encrypt messages sent to you with this key.

Then, send an encrypted message: `/pgp @foo Hey! This message is secret!`

The bot will then encrypt that message with your public key and post it in slack for you:

```
Hey @foo, here's a message from @bar:
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


### Security
There is one database table that only contains a user and their respective public key. No other
personal information is collected or stored.
