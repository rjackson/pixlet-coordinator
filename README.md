# Pixlet Coordinator

The goal of this project is to be able to manage & schedule multiple Pixlet apps to be shown on an LED matrix display, and directly drive that display from this code running on a Raspberry Pi.

This began as an ambition to build a DIY Tidbyt, but the community maintained firmware's architecture was quite lacking by having to:

1. Flash hardcoded network credentials and a URL that'll render a webp image for the Tidbyt to serve
2. Only being able to serve a single app at a time with `pixlet serve`
3. Not being able to configure the apps (e.g. set timezone) when serving (config can be set in the URL, but that's not helpful when that's flashed to a device)
4. Pixlet apps having quite short-lived lifetimes, by having to be encoded into an image and separately pushed to the display.

The current state of the project solves issues 1 to 3.

Todo list for future improvements:

- We currently build the images for an app and then push it to the display. This means the display will be blank for a period of time in-between apps, as the next app's images are built.
  - Move to a queue? Prebuild images to populate queue, and then render?
  - Might mess up time-sensitive data, if any apps assume they'll be running near-realtime
- Support genuine realtime rendering, allowing for long-running processes/apps. I.e. a game of snake that isn't just an animation on loop.
- Web interface for managing apps, configuration, and scheduling

## How to run

- SSH onto a Pi
- Pull down Go dependencies
- Build rpi-rgb-led-matrix library (C stuff)
- Install libwebp something or other

```
go build .

sudo run pixlet-coordinator
```
