## Digestor
> Your personal digest mail service

> Get your personal social email stream digest every day

> Leave your browser

------
![tweets](https://raw2.github.com/29decibel/digestor/master/resources/tweets-digest.png)
![hackernews](https://raw2.github.com/29decibel/digestor/master/resources/hackernews-digest.png)
![github](https://raw2.github.com/29decibel/digestor/master/resources/github-digest.png)


## Install
```
# create your own config
# make changes by your self
cp digestor.json.example ~/.digestor.json
```

## crontab task
```
# you can put it into your crontab -e
15 15 * * * /bin/bash -l -c '/path/to/digestor/digestor' >> /path/to/digestor/cron-job.log 2>&1
```
