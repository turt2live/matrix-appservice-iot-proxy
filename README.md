# matrix-appservice-iot-proxy

A proxy for the matrix client/server API to help offload some requirements from the client. This is written with IoT in mind where the client device is intended to have the appservice token and a user ID to use. The IoT device will *not* receive appservice traffic, however the device may safely make calls to the proxy without registering the user - the proxy will ensure the desired user is created.
