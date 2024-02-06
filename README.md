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

## What's missing?

Several things are still to be done

* Bugs
  * Order of items is not always preserved
* UX
  * Add indicators to all operations
  * Disable elements while an operation is in progress
* TodoMVC
  * Routing All/Active/Completed
  * Clear completed
  * Toggle all
  * Persistence

## Credit

Created by [Matteo Vaccari](https://matteo.vaccari.name)
