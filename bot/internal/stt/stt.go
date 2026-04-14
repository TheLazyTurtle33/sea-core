package stt

import (
	"io"
	"os"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/gorilla/websocket"
)

func Test() {
	uri := "ws://vosk:2700"
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		logger.Error("Failed to connect to Vosk", err)
		return
	}
	defer conn.Close()

	file, err := os.Open("/app/data/tts_test_fixed.wav")
	if err != nil {
		logger.Error("Could not open test file", err)
		return
	}
	defer file.Close()

	// skip the 44-byte WAV header
	_, err = file.Seek(44, 0)
	if err != nil {
		logger.Error("Seek error", err)
	}

	buffer := make([]byte, 64000) // 0.5s of audio at 16kHz
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("Error reading file", err)
			break
		}

		logger.Debug("Sending audio chunk", "size", n)

		// Send audio chunk
		err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			logger.Error("Error sending to websocket", err)
			break
		}

		// Read back the partial or full result
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Error reading response", err)
			break
		}

		logger.Log("Vosk Output", "msg", string(message))
	}

	// Tell Vosk we are finished to get the final transcription
	conn.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
	_, finalMsg, _ := conn.ReadMessage()
	logger.Log("Final Result", "msg", string(finalMsg))
}
