dashing-go
==========

A [Go][1] port of [shopify/dashing][2], built upon [Martini][3].

Still under heavy construction!

### Current Status

* The `/widget/:id` endpoint is done. You can now post JSON data to individual widgets.
* The `/events` endpoint (which emits Server-Sent Events) is done. Registered jobs can now transmit data to widget identifiers.
* For an example of how to write jobs in dashing-go, please refer to the [demo dashboard][4].

Credits
-------

Much of the code is referenced from [golang-sse-todo][5] by @rwynn.

[1]: http://golang.org/
[2]: http://shopify.github.io/dashing/
[3]: http://martini.codegangsta.io/
[4]: https://github.com/gigablah/dashing-go-demo
[5]: https://github.com/rwynn/golang-sse-todo
