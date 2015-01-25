package view

const play_html = `<html>
<head>
<title>{{.Title}}</title>
<style>
	body           {padding: 0px; margin: 0px; text-align: center; font-family: Arial;}
	h3             {padding: 1px; margin: 2px;}
	#jsRequiredMsg {padding: 3px; margin: 3px; font-weight: bold; background: #faa;}
	#controls      {padding: 2px;}
	#controls *    {margin-left: 3px; margin-right: 3px;}
	#view          {position: relative; padding: 1px;}
	#img           {background: #000; border: 1px solid black;}
	#errMsg        {visibility: hidden; position: absolute; top: 10px; right: 0px; width: 100%; color: #ff3030; font-weight: bold;}
	#footer        {margin-top: 5px; font-size: 90%; font-style: italic;}
</style>
</head>

<body>

<noscript>
	<div id="jsRequiredMsg">
		Your web browser must have JavaScript enabled in order for this page to display correctly!
	</div>
</noscript>

<h3>{{.Title}}</h3>

<div id="controls">
	Quality:
	<select id="quality" title="JPEG image quality.">
		<option value="100">100</option>
		<option value="90">90</option>
		<option value="80">80</option>
		<option value="70" selected>70</option>
		<option value="60">60</option>
		<option value="50">50</option>
		<option value="40">40</option>
		<option value="30">30</option>
		<option value="20">20</option>
		<option value="10">10</option>
		<option value="0">0</option>
	</select>
	
	FPS:
	<select id="fps" title="View refresh rate in Frames Per Second">
		<option value="33">30</option>
		<option value="40">25</option>
		<option value="50" selected>20</option>
		<option value="66">15</option>
		<option value="100">10</option>
		<option value="143">7</option>
		<option value="200">5</option>
		<option value="500">2</option>
		<option value="1000">1</option>
	</select>
	
	<button id="pauseResume" onclick="pauseResume()" {{if not .ShowFreezeBtn}}style="display: none"{{end}} 
		title="Stops refreshing the view but does not pauses the game">Freeze</button>
	
	<button id="newGame" onclick="newGame()">New Game</button>
	
	<a href="/help" target="_blank">Help</a>
	
	<a href="/cheat" target="_blank" title="Get a glimpse of the whole Labyrinth">Cheat</a>
	
	<a href="https://github.com/gophergala/golab" target="_blank" title="Visit Home Page">Home page</a>
</div>

<div id="view">
	<img id="img" width="{{.Width}}" height="{{.Height}}"
		onload="errMsg.style.visibility = 'hidden'; imgLoaded = true;"
		onerror="errMsg.style.visibility = 'visible'; setTimeout('imgLoaded = true;', 1000);"
		onmousedown="imgClicked(event)"/>
	<div id="errMsg">Connection Error or Application Closed!</div>
</div>

<div id="footer">
	Copyright &copy; 2015 Andras Belicza. All rights reserved. <a href="https://github.com/gophergala/golab/blob/master/LICENSE.md" target="_blank">LICENSE</a>
</div>

<script>
	var runId = {{.RunId}};
	var playing = false, imgLoaded = true;
	
	// HTML elements:
	var img            = document.getElementById("img"),
		errMsg         = document.getElementById("errMsg"),
		quality        = document.getElementById("quality"),
		fps            = document.getElementById("fps"),
		pauseResumeBtn = document.getElementById("pauseResume");
	
	// Disable image dragging and right-click context menu:
	img.oncontextmenu = img.ondragstart = function() { return false; }
	
	pauseResume();
	
	// Kick-off:
	refresh();
	setInterval(checkRunId, 10000);
	
	function pauseResume() {
		playing = !playing;
		imgLoaded = true;
		pauseResumeBtn.innerText = playing ? "Freeze" : "Resume";
	}
	
	function refresh() {
		if (playing && imgLoaded) {
			imgLoaded = false;
			img.src = "/img?quality=" + quality.value + "&t=" + new Date().getTime();
			setTimeout(refresh, fps.value);
		}
		else
			setTimeout(refresh, 5);
	}
	
	function imgClicked(e) {
		if (!playing)
			return;
		// Relative mouse coordinates inside image:
		var x, y;
		if (document.all) { // For IE, this is enough (exact):
			x = e.offsetX;
			y = e.offsetY;
		} else {            // For other browsers:
	    	x = e.clientX;
	    	y = e.clientY;
	    	for (var el = img; el; el = el.offsetParent) {
	    		x -= el.offsetLeft - el.scrollLeft + el.clientLeft;
        		y -= el.offsetTop - el.scrollTop + el.clientTop;
	    	}
    	}
    	
		var r = new XMLHttpRequest();
		r.open("GET", "/clicked?x=" + x + "&y=" + y + "&b=" + e.button + "&t=" + new Date().getTime(), true);
		r.send(null);
	}
	
	function checkRunId() {
		if (!playing)
			return;
		var r = new XMLHttpRequest();
		r.open("GET", "/runid?t=" + new Date().getTime(), true);
		r.onreadystatechange = function() {
			if (r.readyState == 4 && r.status == 200 && runId != r.responseText)
				window.location.reload(); // App was restarted, reload page
		};
		r.send(null);
	}
	
	function newGame() {
		var r = new XMLHttpRequest();
		r.open("GET", "/new?t=" + new Date().getTime(), true);
		r.onreadystatechange = function() {
			if (r.readyState == 4 && r.status == 200 && runId != r.responseText)
				// New game was started, reload page
				if (!playing)
					pauseResume();
		};
		r.send(null);
	}
</script>

</body>
</html>
`
