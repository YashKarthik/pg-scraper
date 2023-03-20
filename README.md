# [pg-scraper](https://www.yashkarthik.xyz/api/pg-feed)

- A scraped RSS feed of [Paul Graham](http://paulgraham.com/articles.html) 's essays. Get it
[here](https://www.yashkarthik.xyz/api/pg-feed).
- Why? The scraped feed by [Aaron Swartz](http://www.aaronsw.com/) does not include the date; so I
made on that does.
- How do I get date? The date is always (some exceptions) in the first line of the essay itself. So I parse the
essay page's html and extract the date from that. Some exceptions exist, I assign `time.Now()` to them.
