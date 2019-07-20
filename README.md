# Poll Bot

Matrix bot to do polls. What more do you need?

![image](https://user-images.githubusercontent.com/5547783/60209177-029b8e80-9852-11e9-8aee-c91d7ccaaec1.png)
## Build

```bash
git clone https://github.com/babolivier/matrix-poll-bot.git
cd matrix-poll-bot
go build
```

## Run

```
./matrix-poll-bot --config /path/to/config.yaml
```

## Configure

See [config.sample.yaml](/config.sample.yaml).

## Docker

```bash
git clone https://github.com/babolivier/matrix-poll-bot.git
cd matrix-poll-bot
docker build -t matrix-poll-bot .
docker run -v /path/to/config/dir:/data matrix-poll-bot
# Add a file named 'config.yaml' to /path/to/config/dir
```
