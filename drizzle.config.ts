import { defineConfig } from "drizzle-kit";
import { Config as GatorConfig } from "src/config.ts";

export default defineConfig({
  schema: "src/db/schema.ts",
  out: "drizzle",
  dialect: "postgresql",
  dbCredentials: {
    url: (new GatorConfig(GatorConfig.getDefaultConfigPath())).dbUrl,
  },
});
