[![Go](https://github.com/xpmatteo/todomvc-golang/actions/workflows/go.yml/badge.svg)](https://github.com/xpmatteo/todomvc-golang/actions/workflows/go.yml)

# GO + htmx • [TodoMVC](http://todomvc.com)

> Trying to see how the classic single-page-application TodoMVC could be implemented with traditional server-side-rendered
> HTML, and a sprinkle of htmx


## Implementation

Almost no JavaScript!  Thanks to [htmx](https://htmx.org/), we can mimic the weird ways that editing operations are
triggered in the JavaScript SPA using only the attributes of htmx.  The body of the page is replaced at every request,
with no full page reloads.

The server is written in [Go](https://go.dev/), a language that I'm in the process of learning.

The HTML is all rendered server-side.

Interesting features
 * For any request, you can receive the response in Json by using the "accept: application/json" request header
 * Publishes Prometheus metrics
 * Graceful shutdown
 * If you request /active you will receive a full page; if you click on the corresponding button, it will just reload
   the body through an ajax request


## What's missing?

Still to be done:

* UX
  * Add indicators to all operations
  * Disable elements while an operation is in progress
* TodoMVC
  * trim whitespace
  * Clear completed
  * Toggle all
  * Persistence
* Make static assets cacheable by the client
* Configure the app through a config file
* Export health/readiness checks

## Credit

Created by [Matteo Vaccari](https://matteo.vaccari.name)
