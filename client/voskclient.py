import pyaudio
import websocket
import json
import threading

# Vosk-Server WebSocket URL
VOSK_SERVER_URL = "ws://localhost:2700"

# Audio Settings
FORMAT = pyaudio.paInt16  # 16-bit PCM
CHANNELS = 1              # Mono audio
RATE = 16000              # 16kHz sample rate (Vosk standard)
CHUNK = 8000              # Buffer size (0.5 seconds of audio)

def send_audio(ws):
    """Capture microphone audio and send to Vosk-Server."""
    audio = pyaudio.PyAudio()
    stream = audio.open(format=FORMAT,
                        channels=CHANNELS,
                        rate=RATE,
                        input=True,
                        frames_per_buffer=CHUNK)

    print("🎤 Listening... Press Ctrl+C to stop.")
    try:
        while True:
            data = stream.read(CHUNK, exception_on_overflow=False)
            ws.send(data, websocket.ABNF.OPCODE_BINARY)
    except KeyboardInterrupt:
        print("\n🛑 Stopping audio stream.")
    finally:
        # Signal end of audio stream
        ws.send(b'')
        stream.stop_stream()
        stream.close()
        audio.terminate()

def on_message(ws, message):
    """Handle messages from Vosk-Server."""
    result = json.loads(message)
    if 'text' in result and result['text']:  # Only print final results
        print("🗣️ Recognized:", result['text'])

def on_open(ws):
    """Send configuration when connection opens."""
    config = json.dumps({"config": {"sample_rate": RATE}})
    ws.send(config)
    # Start audio streaming in a new thread
    threading.Thread(target=send_audio, args=(ws,), daemon=True).start()

def main():
    ws = websocket.WebSocketApp(VOSK_SERVER_URL,
                                on_open=on_open,
                                on_message=on_message)
    ws.run_forever()

if __name__ == "__main__":
    main()
