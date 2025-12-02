# go-github-selfupdate

> fork of https://github.com/rhysd/go-github-selfupdate

Changes:

- Delete APIToken related code because it is not needed for public repositories.
- Add HTTPClient option to Updater config struct to allow custom HTTP clients (e.g., with proxy settings).
- Add context.Context parameters to NewUpdater method to support request cancellation and timeouts.