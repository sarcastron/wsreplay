# wsreplay

Record and playback a websocket session. Playback will attempt to mimic the message timings as closely as possible usually within a fraction of a millisecond depending on the message size and hardware.

## Configuration

Configuration is done with yaml files. In the future all options will be configurable with flags. The following configuration will record any messages sent to the `target` for `300` seconds. When the duration has elapsed the messages will be written to the path and file noted in `outputTapeFile`.

```yaml
# Recorder settings
target: ws://localhost:8080
# duration in seconds
duration: 300
outputTapeFile: './tapes/session_20221031.gob'
# Playback setting
ServerAddr: '0.0.0.0:8001'
```

## Recording

To record a session, run the `wsreplay record` command with a path to the configuration file.

```sh
# Uses default path of configuration file: ./config.yaml
wsreplay record

# Uses the configuration file provided.
wsreplay record --config ./path/to/config.yaml
```

Recording will begin as soon as the connection is made. When the duration has elapsed or `ctrl-c` interrupt is triggered, the file will be written to storage.

## Playback

To playback a session run the `wsreplay playback` command. If provided, an input file path will be used. Otherwise the configuration file's `outputTapeFile` value will be used.

```sh
# Uses the configuration file's outputTapeFile value as playback file.
wsreplay playback

# Uses the file at provided path
wsreplay playback -f ./tapes/derp2.gob
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
- Interactive Mode
