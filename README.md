# Sea AUCA authentication service for AU Cloud infrastructure

Version 0.0.1

The main service of the AU Cloud. It handles user management and serves as authorization authority.

The API is split into 2 groups: User API and Auth API.

## Auth API

Authentication is done via JWT tokens, using the Ed25519 curve signature. Service supports OAuth 2.0 and OpenID (partially), with custom changes (removed some fields and added extended permission support).

## Environment setup

### Install rel migration cli app

```bash
go install github.com/go-rel/rel/cmd/rel@latest
```

Do not forget to export the `$GOPATH/bin` to `$PATH` for rel to be executable.

### Database

Make sure the Postgres database is setup on your machine: the development username is `postgres` and the password is same. Check documentation on how to edit `pg_hba.conf` file to setup passwords.

The databse cluster on your machine has to have a database called `auc_auth_dev`. If you don't have it - create. It is used to run tests and perform local development.

### Local Email service

To test emails in development envirionment install [MailHog](https://github.com/mailhog/MailHog).
```sh
$ go install github.com/mailhog/MailHog
```

Configure your postifx configuration by editing following fields in `/etc/postfix/main.cf`
```
myhostname=localhost
relayhost=[127.0.0.1]:1025
```

Make sure that `$GOPATH/bin` is included in your `$PATH`.

Run following command to get the hash for credentials setup (instead of password you can put any string):
```sh
$ MailHog bcrypt password
```
Then create file `secret/mailhog.conf` and set ist content as follows:
```
user:< MailHog bcrypt output>
```

After that setup is complete and you should run the mail server with this command: 
```sh
$ MailHog -auth-file=secret/mailhog.conf
```