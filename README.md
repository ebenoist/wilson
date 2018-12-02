Wilson
---

Meet Wilson, an open source "always on" assistant. Think Alexa, but no Amazon, pluggable speech recognition, and completely open. This should be able to run on a raspberry pi. I'll provide instructions on assembly once I get there.

## Installation
Right now this is tied to the google cloud, but I suspect I'll add support for some other services in the future including an "offline" mode. Currently this uses [snowboy](https://github.com/Kitt-AI/snowboy) for wake word detection and then Google TTS/STT APIs for the heavy lifting. The rest is just regexps and your imagination.

## Hardware
* [Respeaker](https://respeaker.io/4_mic_array/) 4 mic array and LED ring
* [Rapsberry Pi 3](https://www.raspberrypi.org/products/raspberry-pi-3-model-b-plus/)
