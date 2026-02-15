import { fetchFeed } from "src/rss";
import { State } from "../state.js";
import { Feed, User } from "../db/schema.js";
import { createFeed, listFeeds } from "../db/queries/feeds.js";
import { createFeedFollow, deleteFeedFollow, getFollowsForUser } from "../db/queries/feedfollows.js";

function printFeed(feed: Feed, user: User) {
  console.log(JSON.stringify(feed));
  console.log(JSON.stringify(user));
}

export async function commandAgg(_: State) {
  console.log(JSON.stringify(await fetchFeed("https://www.wagslane.dev/index.xml")));
}

export async function commandAddfeed(user: User, state: State, ...args: string[]) {
  const name = args[0];
  const feedURL = args[1];

  const feed = await createFeed(state, user.id, name, feedURL);
  await commandFeedFollow(user, state, feedURL);

  printFeed(feed, user);
}

export async function commandListfeed(state: State) {
  const feeds = await listFeeds(state);
  for (const feed of feeds) {
    console.log(`* "${feed.feeds.name}": ${feed.feeds.url} (${feed.users?.name})`);
  }
}

export async function commandFeedFollow(user: User, state: State, ...args: string[]) {
  const feedURL = args[0];

  const feedFollow = await createFeedFollow(state, user.id, feedURL);
  console.log(JSON.stringify(feedFollow));
}

export async function commandFeedUnfollow(user: User, state: State, ...args: string[]) {
  const feedURL = args[0];

  await deleteFeedFollow(state, user.id, feedURL);
  console.log(`${feedURL} unfollowed!`);
}

export async function commandFollowing(user: User, state: State) {
  const feedFollows = await getFollowsForUser(state, user.id);
  for (const feed of feedFollows) {
    console.log(`* "${feed.feeds?.name}": ${feed.feeds?.url}`);
  }
}

