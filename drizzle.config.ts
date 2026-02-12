import { defineConfig } from "drizzle-kit";
import { Config as GatorConfig } from "src/config.ts";

export default defineConfig({
  schema: "db/schema.ts",
  out: "src/db",
  dialect: "postgresql",
  dbCredentials: {
    url: (new GatorConfig(GatorConfig.getDefaultConfigPath())).dbUrl,
  },
});
