# wsreplay

Record and playback a websocket session. Playback will attempt to mimic the message timings as closely as possible usually within a fraction of a millisecond depending on the message size and hardware.

## Configuration

Configuration is done with yaml files or CLI flags. The following configuration will record any messages sent to the `target` for `300` seconds. When the duration has elapsed the messages will be written to the path and file noted in `file`.

```yaml
# Recorder settings
target: ws://localhost:8080
# duration in seconds
duration: 300
file: './tapes/session_20221031.gob'
# optional messages to send and when to send them in seconds
sendMessages:
  - at: 0.8 # Send this message 0.8 seconds after connection
    message: '{"action":"auth","params":"S0m3_4uth_T0k3n"}'
  - at: 1.2 # Send this message 1.2 seconds after connection
    message: '{"action":"subscribe","topic":"howBoutThemIggles"}'
# Playback setting
serverAddr: '0.0.0.0:8001'
```

## Recording

To record a session, run the `wsreplay record` command with a path to the configuration file or with flags. Flags are [documented here](./docs/wsreplay_record.md).

```sh
# Uses the configuration file provided.
wsreplay record --config ./path/to/config.yaml

# target, duration, and output file are set with flags
wsreplay record -t ws://localhost:8001 -d 15 -f ./tapes/session1.gob
```

Recording will begin as soon as the connection is made. When the duration has elapsed or `ctrl-c` interrupt is triggered, the file will be written to storage.

## Playback

To playback a session run the `wsreplay playback` command with a config flag or specific flags. Flags are [documented here](./docs/wsreplay_playback.md).

```sh
# Uses the configuration's file value as playback file.
wsreplay playback

# Uses the file at provided path and serves on port 8001
wsreplay playback -f ./tapes/session2.gob -s 0.0.0.0:8001
```

## Configuration Info

Basic command to check the config values of a file without running it. Example of command and output:

```sh
❯ ./wsreplay info
Recorder settings
 - Target: ws://localhost:8001/v1/price/AAPL
 - Duration: 15
 - Output File: ./tapes/derp2.gob

Playback settings
 - Input File: ./tapes/derp2.gob
 - Server Address: localhost:8001
```

---

For more info see the [generated docs](./docs/wsreplay.md).


## Road map

- Output as library for import into other projects.
- Authentication handling.
- Decouple UI and state
- Interactive Mode
