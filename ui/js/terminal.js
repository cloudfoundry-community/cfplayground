var cfTerminal = function() {
	var event = new Event('terminal_command');
	var self = this;
	self.terminal = $('#terminal').terminal(function(command, term) {

	}, {
		greetings: 'Welcome to Cloud Foundry Playground',
		name: 'cfplayground',
		height: "100%",
		width: "100%",
		enabled: false,
		prompt: "",
		//onFocus: function(term) {
			//workaround of terminal getting focus and enabling input
			//term.focus(false);
		//}
	});

	self.cmdLine = $('#cmdLine').cmd({
		prompt: ' > ',
		width: '100%',
		enabled: false,
		commands: function(command) {
			self.terminal.trigger("terminal_command", command);
		}
	});

	var enabled = function(enable) {
		if (enable) {
			self.cmdLine.enable()
		} else {
			self.terminal.echo('hhi, disabling!!')
			self.terminal.pause()
		}
	}

	var noPrompt = function() {
		//self.terminal.set_prompt("")
	}

	var prompt = function(p) {
		//self.terminal.set_prompt("cloud_foundry > ")
	}

	return {
		echo: this.terminal.echo,
		pause: this.terminal.pause,
		resume: this.terminal.resume,
		prompt: prompt,
		noPrompt: noPrompt,
		enabled: enabled
	}
}
