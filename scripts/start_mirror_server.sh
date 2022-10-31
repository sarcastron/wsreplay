#!/bin/bash

PORT=8001

# websocat -t --exit-on-eof ws-l:127.0.0.1:$PORT broadcast:mirror:
websocat -s 0.0.0.0:$PORT
