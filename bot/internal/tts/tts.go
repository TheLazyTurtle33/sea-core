package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

func Test() {
	/*
		wav1, wav2 := []byte{}, []byte{}
		{

			body := strings.NewReader(`{"text": "voice 1"}`)

			resp, err := http.Post("http://piper:5000/", "application/json", body)
			if err != nil {
				logger.Error("failed to reach piper", err)
				return
			}
			defer resp.Body.Close()

			logger.Log("piper response", "status", resp.StatusCode, "content-type", resp.Header.Get("Content-Type"))

			wav1, err = io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("failed to read piper response", err)
				return
			}

			logger.Log("TTS read done", "bytes", len(wav1), "first_bytes", fmt.Sprintf("%x", wav1[:min(16, len(wav1))]))

			err = os.WriteFile("/app/data/tts_test.wav", wav1, 0644)
			if err != nil {
				logger.Error("failed to write wav file", err)
				return
			}
		}
		{
			body := strings.NewReader(`
			{
				"text": "voice 2",
				"voice": "en_GB-alan-medium"
			}
			`)

			resp, err := http.Post("http://piper:5000/", "application/json", body)
			if err != nil {
				logger.Error("failed to reach piper", err)
				return
			}
			defer resp.Body.Close()

			logger.Log("piper response", "status", resp.StatusCode, "content-type", resp.Header.Get("Content-Type"))

			wav2, err = io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("failed to read piper response", err)
				return
			}

			logger.Log("TTS read done", "bytes", len(wav2), "first_bytes", fmt.Sprintf("%x", wav2[:min(16, len(wav2))]))

			err = os.WriteFile("/app/data/tts_test2.wav", wav2, 0644)
			if err != nil {
				logger.Error("failed to write wav file", err)
				return
			}
		}
		wav := murgeWav(wav1, wav2)
		err := os.WriteFile("/app/data/tts_combined.wav", wav, 0644)
		if err != nil {
			logger.Error("failed to write combined wav file", err)
			return
		}
		pcm := wav[44:]
		conn, err := net.Dial("tcp", "192.168.0.182:9999")
		if err != nil {
			logger.Error("faild to connect to lan Mic", err)
			return
		}
		defer conn.Close()

		_, err = io.Copy(conn, bytes.NewReader(pcm))
		if err != nil {
			logger.Error("faild to send audio", err)
			return
		}

		time.Sleep(500 * time.Millisecond)
		logger.Log("TTS sent", "bytes", len(pcm))
	*/
	MakeTTS("hi [john] hello", "test")
}

type tts_request struct {
	Text  string `json:"text"`
	Voice string `json:"voice"`
}

var voices = []string{
	"alan",
	"john",
	"lessac",
}

func MakeTTS(msg, chatter string) {
	reqs := []tts_request{}
	wavs := [][]byte{}

	words := strings.Split(msg, " ")

	req := tts_request{
		Voice: "lessac",
	}

	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		if word[0] == '[' && word[len(word)-1] == ']' {
			word = word[1 : len(word)-1]
			if slices.Contains(voices, word) {
				if i > 0 {
					reqs = append(reqs, req)
				}
				req = tts_request{}
				req.Voice = word
				continue
			}
		}

		req.Text = req.Text + " " + word
	}
	reqs = append(reqs, req)

	logger.Log("TTS requests", "requests", reqs)

	for _, req := range reqs {
		wavs = append(wavs, makeWav(req))
	}

	wav := murgeWav(wavs...)

	err := os.WriteFile(fmt.Sprintf("/app/data/tts/%s-%s.wav", chatter, time.Now().Add(2*time.Hour).Format("2006-01-02 15:04:05")), wav, 0644)
	if err != nil {
		logger.Error("failed to write combined wav file", err)
		return
	}
}

func makeWav(req tts_request) []byte {

	var wav []byte

	data, err := json.Marshal(req)
	if err != nil {

	}
	body := bytes.NewReader(data)

	resp, err := http.Post("http://piper:5000/", "application/json", body)
	if err != nil {
		logger.Error("failed to reach piper", err)
		return nil
	}
	defer resp.Body.Close()

	logger.Log("piper response", "status", resp.StatusCode, "content-type", resp.Header.Get("Content-Type"))

	wav, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read piper response", err)
		return nil
	}

	logger.Log("TTS generated")

	return wav
}

func murgeWav(wavs ...[]byte) []byte {
	if len(wavs) == 0 {
		logger.Error("no wavs provided to murge", nil)
		return nil
	}
	if len(wavs) == 1 {
		return wavs[0]
	}
	header := wavs[0][:min(44, len(wavs[0]))]

	wavdata := []byte{}
	for _, wav := range wavs {
		data := wav[44:]
		wavdata = append(wavdata, data...)
	}

	wav := make([]byte, 44+len(wavdata))
	copy(wav[:44], header)
	copy(wav[4:8], intToBytes(uint32(36+len(wavdata))))
	copy(wav[40:44], intToBytes(uint32(len(wavdata))))
	copy(wav[44:], wavdata)

	return wav

}

func intToBytes(n uint32) []byte {
	return []byte{
		byte(n & 0xFF),
		byte((n >> 8) & 0xFF),
		byte((n >> 16) & 0xFF),
		byte((n >> 24) & 0xFF),
	}
}
