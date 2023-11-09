# (Old) Accrescent Developer Console

**NOTE: This project is no longer maintained. If you're looking for Accrescent's
current developer console, see https://github.com/accrescent/parcelo.**

The Accrescent developer console - a web application for developers to upload
and manage their apps in the Accrescent app store.

## Development/testing

To set up the development/testing environment for the developer console, follow
these steps:

1. Create an OAuth app from the developer settings of your GitHub account or
   organization. Set the homepage URL to `https://localhost:8080` and the
   authorization callback URL to `https://localhost:8080/auth/github/callback`.
2. Generate a new client secret and store it in `devconsole/.env` as
   `GH_CLIENT_SECRET`. Store the app's client ID as `GH_CLIENT_ID`. Store the
   authorization callback URL as `OAUTH2_REDIRECT_URL`.
3. Set `SIGNER_GH_ID` to the value of the `id` field from
   `https://api.github.com/users/<username>`.
4. Set `REPO_URL` to `http://repo:8080`.
5. Set `API_KEY` to the same string in both `devconsole/.env` and
   `reposerver/.env`.
6. Set `PUBLISH_DIR` in `reposerver/.env` to a folder name such as `/apps`. This
   directory is internal to the container.
7. Generate a TLS certificate & key and store them as `certs/cert.pem` &
   `certs/key.pem` respectively.
8. Start the application by running `docker compose up`
9. The web application is now accessible at `https://localhost:8080`
