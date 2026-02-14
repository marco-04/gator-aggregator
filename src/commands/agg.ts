import { fetchFeed } from "src/rss";
import { State } from "../state.js";

export async function commandAgg(_: State) {
  console.log(JSON.stringify(await fetchFeed("https://www.wagslane.dev/index.xml")));
}
