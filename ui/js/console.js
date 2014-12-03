/* DEPRECATED
 * ==========
 *
 * this file is not being used anymore
 *
*/

var cfTerminal = function(node) {

	var console = $("<div>", {
		"class": "cli_console",
		"html": "Welcome to Cloud Foundry Playground<br><br>"
	}).appendTo(node);

	//add caption text "Welcome to Cloud Foundry Playground"

	var msgBody = $("<div>", {
		"class": "cli_body",
	}).appendTo(console);

	var lastLine = $("<div>", {
		"class": "cli_lastLine"
	}).appendTo(console);

	var cursor = $("<span>", {
		"class": "cli_cursor",
		"html": ">&nbsp;"
	}).appendTo(lastLine);


	var inputSpan = $("<span>", {
		"class": "cli_inputSpan",
		"html": ""
	}).appendTo(lastLine);


	//add blinker, put it at the bottom of a div, with input textbox position on top

	var blinker = $("<span>", {
		"class": "cli_blinker",
		"html": "_"
	}).appendTo(lastLine);

	// var inputBox = $("<input>", {
	// 	"class": "cli_inputBox",
	// 	"type": "text"
	// }).appendTo(inputDiv);

	var cfMessage = function(msg) {
		$("<div>", {
			"html": "<pre class='pre_no_bootstrap'>" + msg + "</pre>"
		}).appendTo(msgBody);
	}

	// inputBox.keypress(function() {
	// 	inputBox.width((5 + (11 * inputBox.val().length)) + "pt")
	// })

	console.keypress(function(d) {

	})

	return {
		cfMessage: cfMessage
	}
}
