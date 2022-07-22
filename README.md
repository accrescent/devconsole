# Accrescent Developer Portal

The Accrescent developer portal - a web application for developers to upload and
manage their apps in the Accrescent app repository.

**Note: Accrescent is not yet ready for production usage. Consider all software
and services run by this organization as in a "pre-alpha" stage and fit only for
development and preliminary testing.**

## Development/testing

To set up the development/testing environment for the developer portal, follow
these steps:

1. Create an OAuth app from the developer settings of your GitHub account or
   organization. Set the homepage URL to `https://localhost:8080` and the
   authorization callback URL to `https://localhost:8080/auth/github/callback`.
2. Generate a new client secret and store it in `.env` as `GH_CLIENT_SECRET` at
   the root of the repository. Store the app's client ID as `GH_CLIENT_ID`.
   Store the authorization callback URL as `OAUTH2_REDIRECT_URL`.
4. Generate a TLS certificate & key and store them as `certs/cert.pem` &
   `certs/key.pem` respectively.
5. Start the application by running `docker-compose up`
6. The web application is now accessible at `https://localhost:8080`
