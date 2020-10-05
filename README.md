# NewsAggregator

NewsAggregator consumes multiple RSS feeds and exposes the data over a single endpoint.
The API is designed to be consumed by the NewsAggregatorAPI.

## Usage
The application can be ran directly via go. Go 1.15+ is advised. The below assumes you are running from the root of the repository.
If you are not, you will need to specify where the feed configuration is located.
Any changes to the port should be reflected in the NewsAggregatorAPI configuration.
```bash
   go run main.go
```

The options below can be supplied as flags. 
```bash
  -feedConfiguration string
        path to feed configuration file (default "./config/local.json")
  -port int
        port to start the application on (default 8888)
  -tick duration
        duration between feed polls (default 1h0m0s)
```
### Feed Configuration
Feed configuration is driven by a JSON file in the `/config` directory.
Each entry represents an individual channel. 'Categories' is an array of default categories to be applied to every article in that channel.

Sample Configuration
```JSON[
         {
           "URL": "http://feeds.bbci.co.uk/news/uk/rss.xml",
           "Categories": ["General"]
         },
         {
           "URL": "http://feeds.bbci.co.uk/news/technology/rss.xml",
           "Categories": ["Tech"]
         },
         {
           "URL": "http://feeds.skynews.com/feeds/rss/uk.xml",
           "Categories": ["General"]
         },
         {
           "URL": "http://feeds.skynews.com/feeds/rss/technology.xml",
           "Categories": ["Tech"]
         }
       ]
```

### API
```GET``` `/latest`
Gets the latest set of data available in the aggregator. This will be the complete set of data.
As such any data previously retrieved should be discarded.

E.G `"http://localhost:8888/latest"`

There will be one Entry for each individual channel. Followed by an array of articles.
##### Sample Response
```JSON
[
    {
        "Articles": [
            {
                "Categories": [
                    "General"
                ],
                "ID": "47febb37f0bd154f1d7badec664ea8fc",
                "Provider": "BBCNews",
                "PublishDate": "2020-10-03T17:17:53+01:00",
                "Snippet": "From monkeys to spiders - Prince George, Princess Charlotte and Prince Louis ask about the natural world.",
                "ThumbnailLink": "TODO",
                "Title": "'Do you like spiders?' - Royal children quiz Sir David Attenborough"
            },
            {
                "Categories": [
                    "General"
                ],
                "ID": "5e5fc91a2deeb7a75934fffa920d5266",
                "Provider": "BBCNews",
                "PublishDate": "2020-10-05T16:04:12+01:00",
                "Snippet": "The firm will open a quarter of its theatres from Friday to Sunday, but declined to comment on jobs.",
                "ThumbnailLink": "TODO",
                "Title": "Odeon to open weekends-only at some cinemas"
            },
            ...
            "Channel": "http://feeds.bbci.co.uk/news/uk/rss.xml"
        ],
        ...
    }
}
```