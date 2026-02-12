# Gator: boot.dev rss feed aggregator 🐊
(no, I am not responsible for this pun)

RSS feed aggregator using postgres (why????) as db backend (I know, an sqlite backend would have been more than enough for this thing)

## Usage
```
Usage: gator-aggregator <cmd> [args]
Commands:
 - dburl: Set db_url in config file
   usage: dburl <db_url>

 - register: Register a user in the database
   usage: register <user>

 - login: Set current user
   usage: login <user>

 - users: List users in the database
   usage: users

 - addfeed: Add a feed
   usage: addfeed <name> <url>

 - feeds: List all feeds for the current user
   usage: feeds

 - unfollow: Unfollow a feed
   usage: unfollow <url>

 - follow: Follow a feed
   usage: follow <url>

 - following: List all feeds the current user is following
   usage: following

 - browse: Browse posts from your followed feeds (defaults to 2 posts)
   usage: browse [limit]

 - agg: Fetch feeds
   usage: agg <time_between_reqs>

 - help: Print usage text
   usage: help

 - reset: [DEBUG] Reset database to allow for easier testing
   usage: reset
```

You'll need to spawn a Postgres instance, then register a user, add some feeds or follow feeds from other users, then run the `agg` command in the background to continuously fetch posts, that you can then "see" (very much _not prettified_) with the `browse` command (_che al mercato mio padre comprò_)

## Build
**IMPORTANT**: to be able to compile the project, [sqlc](https://sqlc.dev) needs to generate some code, so be sure to run the following command ***before*** trying to compile the project:
```
sqlc generate
```

Database migrations are handled with [goose](https://pressly.github.io/goose), and to initialize the database with the correct schema you need to run beforehand:
```
goose postgres <dburl> up -dir sql/schema
```

## Future updates?
~~*Maybe*. Probably to un-overcomplicate the user system and use sqlite instead of deploying *fine* postgres for an RSS feed aggregator lol~~
I have changed my mind lol, I won't touch this ever again

---
Typescript version on the [`typescript`](https://github.com/marco-04/gator-aggregator/tree/typescript) branch
