digestor
========

Personal digest mail service

```
~/.digestor.json
```

## crontab task
```
# you can put it into your crontab -e
15 15 * * * /bin/bash -l -c '/path/to/digestor/digestor' >> /path/to/digestor/cron-job.log 2>&1
```
