# Contacting the WebSocket

When sending information over the webSocket, if it
is audio make sure to send it as a binary blob. If
sending over text make sure it's sent as a text blob.

## Starting Speech to Text

We use Vosk for our speech to text service. To start
the docker container used for vosk run the following
command

`docker run -d --rm -p 2700:2700 alphacep/kaldi-en:latest`

If you have a cuda capable gpu and would like to use
it to assist the speech to text process run the following

`docker run -d --rm -p 2700:2700 alphacep/kaldi-en-gpu:latest`

## Contacting Speech to Text

When sending audio over to the go project websocket,
make sure it is downsampled to 16khz as that is what
vosk expects. The project itself handles contacting
vosk directly and sending the audio to it.
