# Gator: boot.dev rss feed aggregator 🐊
(no, I am not responsible for this pun)

RSS feed aggregator using postgres (why????) as db backend (I know, an sqlite backend would have been more than enough for this thing)

## Usage
TODO

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
*Maybe*. Probably to un-overcomplicate the user system and use sqlite instead of deploying *fine* postgres for an RSS feed aggregator lol
