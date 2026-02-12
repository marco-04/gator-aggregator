import { State } from "../state.js";
import { Config } from "../config.js";
import { commandLogin } from "./user.js";

export function initCLI(config: Config, ...args: string[]) {
  const cmd = args[0];
  if (cmd === undefined) {
    console.log("No command was provided");
    process.exit(1);
  }

  if (availableCommands[cmd] === undefined) {
    console.log(`Unknown command "${cmd}"`);
    process.exit(1);
  }

  const state = {
    cfg: config
  };

  const subArgs = args.slice(1);
  try {
    validateArgs(cmd, ...subArgs);
    availableCommands[cmd].callback(state, ...subArgs);
  } catch (err) {
    console.log(`Error: ${(err as Error).message}`);

    process.exit(1);
  }
}

export type CLICommand = {
  description: string,
  argNum: number,
  callback: (state: State, ...args: string[]) => void;
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
}

function commandHelp(_: State) {
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

