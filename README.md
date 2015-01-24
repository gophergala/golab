GoLab
===

Introduction
---

**Gopher's Labyrinth** (or just **GoLab**) is a 2-dimensional Labyrinth game where you control [Gopher](http://golang.org/doc/gopher/frontpage.png) (who else) and your goal is to get to the Exit point of the Labyrinth. Controlling Gopher is very easy: just click with your mouse to where you want to move it (but there must be a free straight line to where you want to move it).

GoLab is written completely in [Go](http://golang.org/), but there is a thin HTML layer because the User Interface (UI) of the game is an HTML page (webpage). GoLab doesn't use any platform dependent or native code, so you can start the application on any platforms supported by a Go compiler (including Windows, Linux and MAC OS-X). Since the UI is a simple HTML page, you can play the game from any browsers on any platforms, even from mobile phones and tablets (no HTML5 capable browser is required). Also the device you play from doesn't need to be the same computer where you start the application, so for example you can start the game on your desktop computer and connect to it and play the game from your smart phone. Everything is stored in the (Go) application, you can close the browser and reopen it (even on a different device) and nothing will be lost.

How to get it or install it
---

Of course in the _"Go"_ way using `"go get"`:

`go get github.com/gophergala/golab`

Implementation
---

GoLab uses only the standard library that comes with the Official Go distributions. GoLab doesn't rely on any external or 3rd party libraries.

Used packages from the standard library and their utilization:

- [http/net](http://golang.org/pkg/net/http/) package is used as the UI server
- [image](http://golang.org/pkg/image/) package and its sub-packages ([image/color](http://golang.org/pkg/image/color/) and [image/draw](http://golang.org/pkg/image/draw/)) are used to draw the graphics of GoLab
- [image/png](http://golang.org/pkg/image/png/) is used to read image resources of the game
- [image/jpeg](http://golang.org/pkg/image/jpeg/) is used to generate the view of the game for HTTP clients (browsers)
- [html/template](http://golang.org/pkg/html/template/) package is used to generate the UI web page


Usefulness
---

Since GoLab is a game, its usefulness might be questioned. GoLab usefulness is that it is a proof-of-concept and a reference implementation that you can create portable games or applications with graphics in Go with an implicit portable UI with just using the standard library of Go. GoLab doesn't rely on any external or 3rd party libraries.

LICENSE
---

See [LICENSE](https://github.com/gophergala/golab/blob/master/LICENSE.md)

