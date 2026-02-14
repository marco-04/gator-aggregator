import { Config } from "./config.js";

import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

import * as schema from "./db/schema.js";

export interface State {
  cfg: Config,
  db: ReturnType<typeof drizzle>
}

export function newState(cfg: Config): State {
  const db = initDB(cfg);

  return {
    cfg,
    db
  };
}

function initDB(cfg: Config): ReturnType<typeof drizzle> {
  const conn = postgres(cfg.dbUrl);
  return drizzle(conn, { schema });
}
