import { XMLParser } from "fast-xml-parser";

export type RSSFeed = {
  channel: {
    title: string;
    link: string;
    description: string;
    item: RSSItem[];
  }
}

export type RSSItem = {
  title: string;
  link: string;
  description: string;
  pubDate: string;
}

export function isRSSFeed(value: any): value is RSSFeed {
  if (value.channel === undefined) {
    return false;
  }
  if (
    typeof value.channel.title !== "string" ||
    typeof value.channel.link !== "string" ||
    typeof value.channel.description !== "string" ||
    !Array.isArray(value.channel.item)
  ) {
    return false;
  }

  return true;
}

export function isRSSItem(value: any): value is RSSItem {
  if (
    typeof value.title !== "string" ||
    typeof value.link !== "string" ||
    typeof value.description !== "string" ||
    typeof value.pubDate !== "string"
  ) {
    return false;
  }

  return true;
}

export async function fetchFeed(feedURL: string) {
  const response = await fetch(feedURL, {
    method: "GET",
    mode: "cors",
    headers: {
      "User-Agent": "gator"
    }
  });

  const feed = await response.text();

  const parser = new XMLParser();
  let feedObj = parser.parse(feed);

  if (!isRSSFeed(feedObj?.rss)) {
    throw new Error("Not an RSS feed");
  }
  feedObj = feedObj.rss;
  const items: RSSItem[] = [];
  for (const item of feedObj.channel.item) {
    if (isRSSItem(item)) {
      items.push(item);
    }
  }
  feedObj.channel.item = items;

  return feedObj;
}

