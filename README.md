GoLab
===

Introduction
---

**Gopher's Labyrinth** (or just **GoLab**) is a 2-dimensional Labyrinth game where you control [Gopher](http://golang.org/doc/gopher/frontpage.png) (who else) and your goal is to get to the Exit point of the Labyrinth. But beware of the bloodthirsty _Bulldogs_, the ancient enemies of gophers who are endlessly roaming the Labyrinth!

Controlling Gopher is very easy: just click with your _left_ mouse button to where you want him to move (but there must be a free straight line to it). You can even queue multiple target points forming a _path_ on which Gopher will move along. If you click with the _right_ mouse button, the path will be cleared.

<img src="https://github.com/gophergala/golab/blob/master/golab-screenshot.png" alt="GoLab Screenshot" title="alt="GoLab Screenshot">

GoLab is written completely in [Go](http://golang.org/), but there is a thin HTML layer because the User Interface (UI) of the game is an HTML page (web page). GoLab doesn't use any platform dependent or native code, so you can start the application on any platforms supported by a Go compiler (including Windows, Linux and MAC OS-X). Since the UI is a simple HTML page, you can play the game from any browsers on any platforms, even from mobile phones and tablets (no HTML5 capable browser is required). Also the device you play from doesn't need to be the same computer where you start the application, so for example you can start the game on your desktop computer and connect to it and play the game from your smart phone. The solution used (web UI server) provides multi-player support out-of-the-box, although this Labyrinth game doesn't make use of it (the same Gopher can be controlled by all clients). Everything is stored in the (Go) application, you can close the browser and reopen it (even on a different device) and nothing will be lost.

How to get it or install it
---

Of course in the _"Go"_ way using `"go get"`:

`go get github.com/gophergala/golab`

The executable binary `golab` (produced by `"go install"`) is _self-contained_: it contains all resources embedded (e.g. images, html templates), nothing else is required for it to run. On startup by default the application opens the UI web page in your default browser.

Configuration and Tweaking
---

GoLab can be configured and tweaked through command line parameters or flags. Execute `golab -h` to see the available command line options and their description. For completeness and for those who didn't install GoLab, here is the output:

    Usage of golab:
      -autoOpen=true: Auto-opens the UI web page in the default browser
      -bulldogs=10: the number of Bulldogs in an area of 1,000 Blocks; valid range: 0..50
      -cols=33: the number of columns in the Labyrinth; must be odd; valid range: 9..99
      -loopDelay=50: loop delay of the game engine, in milliseconds; valid range: 10..100
      -port=1234: Port to start the UI web server on; valid range: 0..65535
      -rows=33: the number of rows in the Labyrinth; must be odd; valid range: 9..99
      -v=80: moving speed of Gopher and the Bulldogs in pixel/sec; valid range: 20..200
      -viewHeight=700: height of the view image in pixels in the UI web page; valid range: 150..2000
      -viewWidth=700: width of the view image in pixels in the UI web page; valid range: 150..2000

Used Packages
---

GoLab uses only the standard library that comes with the Official Go distributions. GoLab doesn't rely on any external or 3rd party libraries.

Used packages from the standard library and their utilisation:

- [http/net](http://golang.org/pkg/net/http/) package is used as the UI server
- [image](http://golang.org/pkg/image/) package and its sub-packages ([image/color](http://golang.org/pkg/image/color/) and [image/draw](http://golang.org/pkg/image/draw/)) are used to draw the graphics of GoLab
- [image/png](http://golang.org/pkg/image/png/) is used to read image resources of the game
- [image/jpeg](http://golang.org/pkg/image/jpeg/) is used to generate the view of the game (labyrinth) for HTTP clients (browsers)
- [html/template](http://golang.org/pkg/html/template/) package is used to generate the UI web page
- [encoding/base64](http://golang.org/pkg/encoding/base64/) package is used to generate and decode embedded image resources to/form Base64 strings
- [flag](http://golang.org/pkg/flag/) package is used to enable basic configuration through the command line

Under the Hood (Implementation)
---

**Game Engine / Simulation**

As mentioned earlier, everything is calculated and stored in the (Go) application. As an architectural pattern, I chose [Model-View-Controller (MVC)](http://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller). Although I did not enforce everything but logically this pattern is followed.

The `model` package defines the basic types and data structures of the game. The `view` package is responsible for the UI of the game. The UI is a thin HTML layer, it contains an HTML page with some embedded JavaScript. No external JavaScript libraries are used, everything is "self-made". At the GoLab "side" the `net/http` package is used to serve the HTTP clients (browsers).

The `ctrl` package is the controller or the _engine_ of the game, it implements all the game logic. It runs in an endless loop, and processes events from the UI client(s), performs calculation of moving objects, performs certain checks (like winning and dying) and updates the image / view of the Labyrinth.

Since there might be multiple goroutines running parallel, communication between the `view` and the `ctrl/model` is done via channels. Also to prevent incomplete/flickering images sent to the clients, the engine performs explicit "model" locking while the next phase of the game is being calculated. 

**Communication between the (Go) application and the browser (UI):**

- When GoLab is started, it starts an HTTP(web) server.
- Either GoLab auto-opens the UI web page in the default browser (default) or the player manually opens it.
- The UI web page is served by the web server.
- The UI web page presents the view of the game in the form of an HTML image. This image is then periodically refreshed (by JavaScript code).
- Clicks on the view image is detected by JavaScript code and are sent back to the server via AJAX calls. The server processes them.
- Quality is a parameter which is attached to the image urls when the view is requested.
- The FPS parameter is just used at the client side to time image refreshing.
- New Game requests are also sent via AJAX calls.
- The Cheat link opens a new browser tab directed to a URL whose handler sends a snapshot image of the whole Labyrinth.
- The web page constantly monitors the application, and if the application is closed or network error occurs, proper notification/error messages are displayed to the user. The web page automatically "reconnects" if the application becomes available again.
- The web page also automatically detects if the application is restarted, and in this case will reload itself. 

Usefulness
---

Since GoLab is a game, its usefulness might be questioned. GoLab's usefulness is that it is an example solution and a reference implementation that you can create portable games or applications with graphics in Go with an implicit portable UI with just using the standard library of Go. GoLab doesn't rely on any external or 3rd party libraries.

LICENSE
---

See [LICENSE](https://github.com/gophergala/golab/blob/master/LICENSE.md)

GoLab's Gopher is a derivative work based on the Go gopher which was designed by Renee French. ([http://reneefrench.blogspot.com/](http://reneefrench.blogspot.com/)). Licensed under the Creative Commons 3.0 Attributions license.

The source of other images can be found in the [resources/source.txt](https://github.com/gophergala/golab/blob/master/resources/source.txt) file.
