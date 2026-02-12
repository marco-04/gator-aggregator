import { State } from "../state.js";

export function commandLogin(state: State, ...args: string[]) {
  const userName = args[0];

  state.cfg.currentUserName = userName;
  state.cfg.write();
  console.log(`${userName} set as current user`);
}
