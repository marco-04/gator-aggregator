import { initCLI } from "./commands/commands.js";
import { Config } from "./config.js";

async function main() {
  try {
    const cfg = new Config(Config.getDefaultConfigPath());

    await initCLI(cfg, ...process.argv.slice(2));
  } catch(err) {
    if (err instanceof SyntaxError) {
      console.log(`Invalid config file:\n${err.message}`);
    } else {
      console.log(`Config error: ${(err as Error).message}`);
    }
    process.exit(1);
  }
  process.exit(0);
}

await main();
