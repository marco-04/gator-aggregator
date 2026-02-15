import { fetchFeed } from "src/rss";
import { State } from "../state.js";
import { Feed, User } from "../db/schema.js";
import { createFeed, getNextToFetch, listFeeds, markFeedFetched } from "../db/queries/feeds.js";
import { createFeedFollow, deleteFeedFollow, getFollowsForUser } from "../db/queries/feedfollows.js";

function printFeed(feed: Feed, user: User) {
  console.log(JSON.stringify(feed));
  console.log(JSON.stringify(user));
}

async function scrapeFeed(state: State) {
  const nextFeed = await getNextToFetch(state);
  const feed = await fetchFeed(nextFeed.url);
  await markFeedFetched(state, nextFeed.id);
  console.log(`>> Fetched ${nextFeed.url}`);
  for (const item of feed.channel.item) {
    console.log(`* ${item.title}`);
  }
}

function parseDuration(durationStr: string): number | undefined {
  const regex = /^(\d+)(ms|s|m|h)$/;
  const match = durationStr.match(regex);
  if (match === null) {
    return undefined;
  }

  let multiplier: number;
  switch(match[2]) {
    case "ms":
      multiplier = 1;
      break;
    case "s":
      multiplier = 1000;
      break;
    case "m":
      multiplier = 60 * 1000;
      break;
    case "h":
      multiplier = 60 * 60 * 1000;
      break;
    default:
      return undefined;
  }

  return Number(match[1]) * multiplier;
}

function logError(err: Error) {
  console.log(`== ${(err as Error).message} ==`);
}

export async function commandAgg(state: State, ...args: string[]) {
  const durationStr = args[0].toLowerCase();
  const interval = parseDuration(durationStr);
  if (interval === undefined) {
    throw new Error(`Invalid duration ${durationStr}`);
  }

  console.log(`== Collecting feeds every ${durationStr} ==`);

  scrapeFeed(state).catch(logError);
  const scrapeInterval = setInterval(() => {
    scrapeFeed(state).catch(logError);
  }, interval);

  // Kill the program with CTRL+C
  await new Promise<void>((resolve) => {
    process.on("SIGINT", () => {
      console.log("== Shutting down feed aggregator... ==");
      clearInterval(scrapeInterval);
      resolve();
    });
  });
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

