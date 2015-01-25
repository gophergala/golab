package view

const help_html = `<html>
<head>
<title>{{.Title}}</title>
<style>
	body           {padding: 0px; margin: 0px; margin-left: auto; margin-right: auto; text-align: center; font-family: Arial; width: 600px;}
	h3             {padding: 1px; margin: 2px;}
	#content       {text-align: justify;}
	#footer        {margin-top: 7px; padding-top: 3px; font-size: 90%; font-style: italic; border-top: 1px solid #888;}
</style>
</head>

<body>

<h3>{{.Title}}</h3>

<div id="content">
	<p>
		<b>Gopher's Labyrinth</b> (or just <b>GoLab</b>) is a 2-dimensional Labyrinth game where you control
		<a href="http://golang.org/doc/gopher/frontpage.png" target="_blank">Gopher</a> (who else)
		and your goal is to get to the Exit point of the Labyrinth. But beware of the bloodthirsty <i>Bulldogs</i>,
		the ancient enemies of gophers who are endlessly roaming the Labyrinth!
	</p>
	<p>
		Controlling Gopher is very easy: just click with your <i>left</i> mouse button to where you want him to move
		(but there must be a free straight line to it). You can even queue multiple target points forming a <i>path</i>
		on which Gopher will move along. If you click with the <i>right</i> mouse button, the path will be cleared.
	</p>
</div>

<div id="close">
	<input type="button" value="Close Help" onclick="window.close()">
</div>

<div id="footer">
	Copyright &copy; 2015 Andras Belicza. All rights reserved. <a href="https://github.com/gophergala/golab/blob/master/LICENSE.md" target="_blank">LICENSE</a>
</div>

</body>
</html>
`
