var cfWebsocket = function(cfTerm) {
	var wskt = WebSocket;
	var connect = function() {
		try {
			wskt = new WebSocket("ws://" + wsIP + ":" + wsport + "/ws");

			wskt.onopen = function() {
				writeToTerminal("Connected to " + wsIP + ":" + wsport, "system", "");
			};

			wskt.onmessage = function(status) {
				processOutput($.parseJSON(status.data));
			};

			wskt.onclose = function() {
				writeToTerminal('Disconnected from server', "warning", "");
				//$('#progressBar').css("display", "none");
			};
		} catch (err) {
			output('error in connect module: ' + err);
		}

	}

	var send = function(msg) {
		wskt.send(msg);
	}

	var isConnected = function() {
		if (wskt.readyState != 1) {
			alert('Not connected to server, please connect');
			return false;
		} else {
			return true
		}
	}

	var writeToTerminal = function(msg, msgType) {
		if (msgType == "system") {
			cfTerm.echo("[[i;yellow;transparent]" + msg + "\n]");
		} else if (msgType == "warning") {
			cfTerm.echo("[[i;red;transparent]" + msg + "\n]");
		} else if (msgType == "important") {
			cfTerm.echo("[[gi;white;transparent]<< " + msg + " >>\n]");
		} else if (msgType == "input") {
			cfTerm.echo(" > " + msg);
		} else {
			cfTerm.echo(msg);
		}
	}

	function processOutput(data) {
		console.log(data.MsgType)
		if (data.Cmd == "token") {
			cfToken = data.Msg
			$('#fileupload').fileupload("option", "url", "/upload/" + cfToken)
			writeToTerminal("Your session ID is " + cfToken, "important", " ** ")
			$("#divSessionId").html("Your Session Id: " + cfToken)
		} else if (data.Cmd == "echo") {
			writeToTerminal(data.Msg, data.MsgType)
		} else if (data.Cmd == "course") {
			$("#tutorialTitle").html(data.MsgType)
			$("#tutorialStep").html(data.Msg)
		} else {
			writeToTerminal(data.Msg, "")
		}
	}

	return {
		connect: connect,
		send: send,
		isConnected: isConnected,
		writeToTerminal: writeToTerminal
	};
}
