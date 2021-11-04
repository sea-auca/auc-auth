# Sea AUCA authentication service for AU Cloud infrastructure

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
