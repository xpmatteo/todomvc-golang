# GO + Htmx • [TodoMVC](http://todomvc.com)

> Trying to see how the classic single-page-application TodoMVC could be implemented with traditional server-side-rendered
> html, and a sprinkle of Htmx


## Implementation

Almost no JavaScript!  Thanks to [HTMX](https://htmx.org/), we can mimic the weird ways that editing operations are
triggered in the JavaScript SPA using only the attributes of HTMX.  The body of the page is replaced at every request, 
with no full page reloads.

The server is written in [Go](https://go.dev/), a language that I'm in the process of learning.

The html is all rendered server-side.

## What's missing?

Several things are still to be done

* Bugs
  * The POST /edit request triggers 2 POST calls and 2 GET calls ?!?
  * The POST /toggle request triggers 2 POST calls ???
* UX
  * Add indicators to all operations
  * Disable elements while an operation is in progress
* Performance
  * Perhaps the HTML content should be returned by the POST operation; this saves a round trip.
* TodoMVC
  * Routing All/Active/Completed
  * Clear completed
  * Toggle all
  * Persistence

## Credit

Created by [Matteo Vaccari](https://matteo.vaccari.name)
