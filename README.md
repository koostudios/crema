# Crema

**Current Version:** 0.1

The essential collection of tools for making your static sites more awesome.

### Crema Forms

A simple way to send your forms on static sites to an email address.

* Point any static form to the server and it will start sending
* Send to any email address, no sign up required


## Changelog

### Version 0.1

Sets up the following endpoints:

* `GET /` Hello World! - Sends a Hello World Message.
* `GET /{email}` Reminds you to please `POST` to this endpoint.
* `POST /{email}` Sends all form content to your email in a Go Template
