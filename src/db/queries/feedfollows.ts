import { State } from "../../state.js";
import { feed_follows, feeds, users } from "../schema.js";
import { getUserID } from "./users.js";
import { getFeed } from "./feeds.js";
import { eq } from "drizzle-orm";

export async function createFeedFollow(state: State, username: string, feedURL: string) {
  const userID = await getUserID(state, username);
  const feed = await getFeed(state, feedURL);
  const [result] = await state.db.insert(feed_follows).values({
    userId: userID,
    feedId: feed.id
  }).returning();
  return result;
}

export async function getFollowsForUser(state: State, username: string) {
  const userID = await getUserID(state, username);
  const result = await state.db.select()
    .from(feed_follows)
    .innerJoin(users, eq(users.id, feed_follows.userId))
    .innerJoin(feeds, eq(feeds.id, feed_follows.feedId))
    .where(eq(users.id, userID));
  return result;
}

