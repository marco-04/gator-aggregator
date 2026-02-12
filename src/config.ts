import fs from "node:fs"
import os from "node:os";
import path from "node:path";

export interface ConfigInterface {
  db_url: string;
  current_user_name: string;
}

export class Config {
  private cfg: ConfigInterface;
  private path: string;

  constructor(path: string) {
    this.path = path;
    const stat = fs.statSync(path, { throwIfNoEntry: false });
    if (!stat) {
      this.cfg = {
        db_url: "",
        current_user_name: ""
      }
      this.write();
    }
    const cfg = fs.readFileSync(path, "utf8");
    this.cfg = JSON.parse(cfg);
  }

  get dbUrl() {
    return this.cfg.db_url;
  }

  get currentUserName() {
    return this.cfg.current_user_name;
  }

  get configObject() {
    return this.cfg;
  }

  set dbUrl(db_url: string) {
    this.cfg.db_url = db_url;
  }

  set currentUserName(current_user_name: string) {
    this.cfg.current_user_name = current_user_name;
  }

  write() {
    fs.writeFileSync(this.path, JSON.stringify(this.cfg));
  }

  static getDefaultConfigPath() {
    return path.join(os.homedir(), ".gatorconfig.json");
  }
}

