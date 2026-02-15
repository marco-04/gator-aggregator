import { newState, State } from "../state.js";
import { Config } from "../config.js";
import { commandLogin, commandRegister, commandReset, commandUsers } from "./user.js";
import { commandAddfeed, commandAgg, commandFeedFollow, commandFeedUnfollow, commandFollowing, commandListfeed } from "./agg.js";
import { User } from "../db/schema.js";
import { getUser } from "src/db/queries/users.js";

export async function initCLI(cfg: Config, ...args: string[]) {
  const cmd = args[0];
  if (cmd === undefined) {
    console.log("No command was provided");
    process.exit(1);
  }

  if (availableCommands[cmd] === undefined) {
    console.log(`Unknown command "${cmd}"`);
    process.exit(1);
  }

  const state = newState(cfg);

  const subArgs = args.slice(1);
  try {
    validateArgs(cmd, ...subArgs);
    await availableCommands[cmd].callback(state, ...subArgs);
  } catch (err) {
    console.log(`Error: ${(err as Error).message}`);

    process.exit(1);
  }
}

export type Callback = (state: State, ...args: string[]) => Promise<void>;
export type UserCallback = (user: User, state: State, ...args: string[]) => Promise<void>;

export type CLICommand = {
  description: string,
  argNum: number,
  callback: Callback
}

const availableCommands: Record<string, CLICommand> = {
  help: {
    description: "Print the help text",
    argNum: 0,
    callback: commandHelp
  },
  login: {
    description: "Login to user",
    argNum: 1,
    callback: commandLogin
  },
  register: {
    description: "Register an user",
    argNum: 1,
    callback: commandRegister
  },
  users: {
    description: "Get all registered users",
    argNum: 0,
    callback: commandUsers
  },
  agg: {
    description: "Aggregate feeds",
    argNum: 1,
    callback: commandAgg
  },
  addfeed: {
    description: "Add a feed",
    argNum: 2,
    callback: loggedIn(commandAddfeed)
  },
  feeds: {
    description: "List feeds",
    argNum: 0,
    callback: commandListfeed
  },
  follow: {
    description: "Follow feed",
    argNum: 1,
    callback: loggedIn(commandFeedFollow)
  },
  following: {
    description: "Follow feed",
    argNum: 0,
    callback: loggedIn(commandFollowing)
  },
  unfollow: {
    description: "Unfollow feed",
    argNum: 1,
    callback: loggedIn(commandFeedUnfollow)
  },
  reset: {
    description: "Delete all user records",
    argNum: 0,
    callback: commandReset
  },
}

async function commandHelp(_: State) {
  console.log("Usage:");
  for (const key of Object.keys(availableCommands)) {
    console.log(`${key}: ${availableCommands[key]}`);
  }
}

function validateArgs(cmd: string, ...args: string[]) {
  const expected = availableCommands[cmd]!.argNum;
  const got = args.length;
  if (got !== expected) {
    throw new Error(`Wrong number of arguments: expected ${expected}, got ${got}`);
  }
}

export function loggedIn(handler: UserCallback): Callback {
  return async (state: State, ...args: string[]) => {
    const user = await getUser(state, state.cfg.currentUserName);
    if (!user) {
      throw new Error(`User ${state.cfg.currentUserName} not found`);
    }

    await handler(user, state, ...args);
  }
}

