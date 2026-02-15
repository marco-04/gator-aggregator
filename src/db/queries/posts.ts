import { sql } from "drizzle-orm";
import { State } from "../../state.js";
import { feed_follows, posts } from "../schema.js";
import { RSSItem } from "../../rss.js";

export async function createPost(state: State, feedID: string, post: RSSItem) {
  if (!post.link) {
    throw new Error("Invalid post");
  }
  if (!post.title) {
    post.title = post.link;
  }

  let pubDate: Date | undefined = new Date(post.pubDate);
  if (isNaN(pubDate.getTime())) {
    pubDate = undefined;
  }

  const [result] = await state.db.insert(posts).values({
    title: post.title,
    url: post.link,
    description: post.description,
    pubDate: pubDate,
    feedId: feedID
  }).returning();
  return result;
}

export async function getPostsForUser(state: State, userID: string, numPosts: number) {
  const result = await state.db.execute(
    sql`SELECT * FROM ${posts}
    INNER JOIN ${feed_follows}
    ON ${posts.feedId} = ${feed_follows.feedId}
    WHERE ${feed_follows.userId} = ${userID}
    ORDER BY ${posts.pubDate} DESC NULLS LAST
    LIMIT ${numPosts ? numPosts : 2}`
  );
  return result;
}

