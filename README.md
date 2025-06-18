# RSS Reader

Go-based RSS reader for fetching and parsing RSS feeds.

---

## Docker

### Build

```bash
docker build -t rss-parser .
```

### Run

Pass the RSS feed URL as an argument
```bash
docker run rss-parser https://feeds.simplecast.com/mKn_QmLS
```

Otherwise, the default URL https://feeds.simplecast.com/qm_9xx0g will be used
```bash
docker run rss-parser
```