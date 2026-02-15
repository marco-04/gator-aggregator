import { fetchFeed } from "src/rss";
import { State } from "../state.js";
import { feeds, users } from "../db/schema.js";
import { createFeed, listFeeds } from "../db/queries/feeds.js";
import { getUser, getUserFromID } from "../db/queries/users.js";
import { createFeedFollow, getFollowsForUser } from "../db/queries/feedfollows.js";

export type Feed = typeof feeds.$inferSelect;
export type User = typeof users.$inferSelect;

function printFeed(feed: Feed, user: User) {
  console.log(JSON.stringify(feed));
  console.log(JSON.stringify(user));
}

export async function commandAgg(_: State) {
  console.log(JSON.stringify(await fetchFeed("https://www.wagslane.dev/index.xml")));
}

export async function commandAddfeed(state: State, ...args: string[]) {
  const name = args[0];
  const feedURL = args[1];

  const feed = await createFeed(state, name, feedURL);
  await commandFeedFollow(state, feedURL);

  const user = await getUserFromID(state, feed.userId);
  printFeed(feed, user);
}

export async function commandListfeed(state: State) {
  const feeds = await listFeeds(state);
  for (const feed of feeds) {
    console.log(`* "${feed.feeds.name}": ${feed.feeds.url} (${feed.users?.name})`);
  }
}

export async function commandFeedFollow(state: State, ...args: string[]) {
  const feedURL = args[0];

  const feedFollow = await createFeedFollow(state, state.cfg.currentUserName, feedURL);
  console.log(JSON.stringify(feedFollow));
}

export async function commandFollowing(state: State) {
  const feedFollows = await getFollowsForUser(state, state.cfg.currentUserName);
  for (const feed of feedFollows) {
    console.log(`* "${feed.feeds?.name}": ${feed.feeds?.url}`);
  }
}

