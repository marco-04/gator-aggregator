package cli

var cmds = map[string]cmd {
	"login": {
		description: "Set current user",
		argNum: 1,
		usageStr: "<user>",
		callback: login,
	},
	"dburl": {
		description: "Set db_url in config file",
		argNum: 1,
		usageStr: "<db_url>",
		callback: setDBURL,
	},
	"register": {
		description: "Register a user in the database",
		argNum: 1,
		usageStr: "<user>",
		callback: register,
	},
	"reset": {
		description: "[DEBUG] Reset database to allow for easier testing",
		argNum: 0,
		usageStr: "",
		callback: reset,
	},
	"users": {
		description: "List users in the database",
		argNum: 0,
		usageStr: "",
		callback: listUsers,
	},
	"agg": {
		description: "Fetch feeds",
		argNum: 1,
		usageStr: "<time_between_reqs>",
		callback: agg,
	},
	"addfeed": {
		description: "Add a feed",
		argNum: 2,
		usageStr: "<name> <url>",
		callback: doCallbackLoggedIn(addFeed),
	},
	"feeds": {
		description: "List all feeds for the current user",
		argNum: 0,
		usageStr: "",
		callback: listFeeds,
	},
	"follow": {
		description: "Follow a feed",
		argNum: 1,
		usageStr: "<url>",
		callback: doCallbackLoggedIn(followFeed),
	},
	"unfollow": {
		description: "Unfollow a feed",
		argNum: 1,
		usageStr: "<url>",
		callback: doCallbackLoggedIn(unfollowFeed),
	},
	"following": {
		description: "List all feeds the current user is following",
		argNum: 0,
		usageStr: "",
		callback: listFollowFeeds,
	},
	"help": {
		description: "Print usage text",
		argNum: 0,
		usageStr: "",
		callback: printHelp,
	},
}

