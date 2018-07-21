# matrix-appservice-iot-proxy

[![#iotproxy:t2bot.io](https://img.shields.io/badge/matrix-%23iotproxy:t2bot.io-brightgreen.svg)](https://matrix.to/#/#iotproxy:t2bot.io)
[![TravisCI badge](https://travis-ci.org/turt2live/matrix-appservice-iot-proxy.svg?branch=master)](https://travis-ci.org/turt2live/matrix-appservice-iot-proxy)

A proxy for the matrix client/server API to help offload some requirements from the client. This is written with IoT in mind where the client device is intended to have the appservice token and a user ID to use. The IoT device will *not* receive appservice traffic, however the device may safely make calls to the proxy without registering the user - the proxy will ensure the desired user is created.


# Installing

Assuming Go 1.9 is already installed on your PATH:
```bash
# Get it
git clone https://github.com/turt2live/matrix-appservice-iot-proxy
cd matrix-appservice-iot-proxy

# Set up the build tools
currentDir=$(pwd)
export GOPATH="$currentDir/vendor/src:$currentDir/vendor:$currentDir:"$GOPATH
go get github.com/constabulary/gb/...
export PATH=$PATH":$currentDir/vendor/bin:$currentDir/vendor/src/bin"

# Build it
gb vendor restore
gb build

# Configure it (edit iot-proxy.yaml to meet your needs)
cp config.sample.yaml iot-proxy.yaml

# Run it
bin/matrix_iot_proxy
```

# Example IoT Code

This project allows IoT devices to simplify their code, making sure that virtual users exist and have joined the room. For instance, consider a simple sensor (that happens to be powered by cURL) data logging device:

```bash
curl -X PUT -H "Content-Type: application/json" --data-binary '{"temperature":22,"units":"Celsius"}' 'http://my.iot.proxy.com:4232/_matrix/client/r0/rooms/!myroom:domain.com/send/com.custom.temperature/_txn_?access_token=YourAppserviceToken&user_id=@.sensor.temperature:domain.com'
```

Note how the IoT device does not need to register `@.sensor.temperature:domain.com` and it does not need to join `!myroom:domain.com` - both are handled automatically by this proxy. Please note that your IoT devices will need to have the appservice token, and therefore you'll need to register the appservice yourself.

The IoT Proxy will additionally replace `_txn_` with a unique value when sending events so your devices do not need to track state. If you don't want to use this functionality, simply don't use the `_txn_` keyword in your request.
