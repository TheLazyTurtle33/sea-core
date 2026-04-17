package tts

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
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
	MakeTTS("hay daddy", "test")
}

type tts_request struct {
	Text   string `json:"text"`
	Voice  string `json:"voice"`
	Silent int    `json:"-"`
}

var voices = map[string]string{
	"alan":   "en_GB-alan-medium",
	"john":   "en_US-john-medium",
	"lessac": "en_US-lessac-medium",
	"glados": "en_US-glados-high",
	"hal":    "hal",
	"pda":    "pda",
	"trump":  "en_US-trump-high",
	"carlin": "en_US-carlin-high",
}

func MakeTTS(msg, chatter string) {
	logger.Log("MakeTTs msg", "mesg", msg)
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
			if word == "" {
				if i > 0 {
					reqs = append(reqs, req)
				}
				req = tts_request{
					Voice: voices["lessac"],
				}
				continue
			}
			if voice := voices[word]; voice != "" {
				if i > 0 {
					reqs = append(reqs, req)
				}
				req = tts_request{
					Voice: voice,
				}
				continue
			}
		}

		if strings.Contains(word, "/") {
			chars := strings.Split(word, "")
			found := false
			for _, c := range chars {
				if c != "/" {
					found = true
					break
				}
			}
			if !found {
				voice := req.Voice
				if i > 0 {
					reqs = append(reqs, req)
				}
				req = tts_request{
					Silent: len(chars),
				}
				reqs = append(reqs, req)
				req = tts_request{
					Voice: voice,
				}
				continue
			}
		}

		req.Text = req.Text + " " + word
	}
	reqs = append(reqs, req)

	logger.Log("TTS requests", "requests", reqs)

	for _, req := range reqs {
		if req.Silent > 0 {
			wavs = append(wavs, silenceWav(req.Silent))
		} else {
			wavs = append(wavs, makeWav(req))
		}
	}

	wav := mergeWav(wavs...)

	var err error
	if !context.Get().YappyChat {
		err = os.WriteFile(fmt.Sprintf("/app/data/tts/%s_%s.wav", time.Now().Add(2*time.Hour).Format("2006-01-02_15:04:05"), chatter), wav, 0644)
	}
	err = os.WriteFile("/app/data/tts/tts.wav", wav, 0644)
	if err != nil {
		logger.Error("failed to write combined wav file", err)
		return
	}
}

func silenceWav(halfSeconds int) []byte {
	sampleRate := 22050
	channels := 1
	bitDepth := 2 // s16le = 2 bytes per sample

	numSamples := (sampleRate * halfSeconds) / 2
	dataSize := numSamples * channels * bitDepth
	totalSize := 44 + dataSize

	wav := make([]byte, totalSize)

	// RIFF header
	copy(wav[0:4], []byte("RIFF"))
	copy(wav[4:8], intToBytes(uint32(totalSize-8)))
	copy(wav[8:12], []byte("WAVE"))

	// fmt chunk
	copy(wav[12:16], []byte("fmt "))
	copy(wav[16:20], intToBytes(uint32(16))) // chunk size
	copy(wav[20:22], []byte{0x01, 0x00})     // PCM format
	copy(wav[22:24], intToBytes16(uint16(channels)))
	copy(wav[24:28], intToBytes(uint32(sampleRate)))
	copy(wav[28:32], intToBytes(uint32(sampleRate*channels*bitDepth))) // byte rate
	copy(wav[32:34], intToBytes16(uint16(channels*bitDepth)))          // block align
	copy(wav[34:36], intToBytes16(uint16(bitDepth*8)))                 // bits per sample

	// data chunk
	copy(wav[36:40], []byte("data"))
	copy(wav[40:44], intToBytes(uint32(dataSize)))
	// wav[44:] is already zeroed = silence

	return wav
}

func intToBytes16(v uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	return b
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

func mergeWav(wavs ...[]byte) []byte {
	if len(wavs) == 0 {
		logger.Error("no wavs provided to merge", nil)
		return nil
	}
	if len(wavs) == 1 {
		return wavs[0]
	}

	// Write each wav to a temp file
	tmpFiles := []string{}
	for i, wav := range wavs {
		path := fmt.Sprintf("/tmp/tts_part_%d.wav", i)
		if err := os.WriteFile(path, wav, 0644); err != nil {
			logger.Error("failed to write temp wav", err)
			return nil
		}
		tmpFiles = append(tmpFiles, path)
	}
	defer func() {
		for _, f := range tmpFiles {
			os.Remove(f)
		}
	}()

	// Build ffmpeg concat input file
	concatFile := "/tmp/tts_concat.txt"
	concatContent := ""
	for _, f := range tmpFiles {
		concatContent += fmt.Sprintf("file '%s'\n", f)
	}
	if err := os.WriteFile(concatFile, []byte(concatContent), 0644); err != nil {
		logger.Error("failed to write concat file", err)
		return nil
	}
	defer os.Remove(concatFile)

	// ffmpeg concat + normalize to common format
	outFile := "/tmp/tts_merged.wav"
	defer os.Remove(outFile)

	cmd := exec.Command("ffmpeg", "-y",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-ar", "22050", // common sample rate
		"-ac", "1", // mono
		"-sample_fmt", "s16",
		outFile,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		logger.Error("ffmpeg merge failed", err)
		logger.Log("ffmpeg output", "out", string(out))
		return nil
	}

	result, err := os.ReadFile(outFile)
	if err != nil {
		logger.Error("failed to read merged wav", err)
		return nil
	}

	return result
}

func intToBytes(n uint32) []byte {
	return []byte{
		byte(n & 0xFF),
		byte((n >> 8) & 0xFF),
		byte((n >> 16) & 0xFF),
		byte((n >> 24) & 0xFF),
	}
}

type TTS struct {
	action.Action
	Message string
	Name    string
}

func (a TTS) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}

	logger.Log(a.Message)

	if len(v) == 0 {
		flags.Error = fmt.Errorf("Expected Chat Message")
		return flags
	}

	if a.Message == "" {
		if chat, ok := v[0].(datatypes.ChatMessageData); ok {
			parts := strings.SplitN(chat.Message.Text, " ", 2)
			if len(parts) > 1 {
				a.Message = parts[1]
			}
			a.Name = chat.ChatterUserLogin

			if a.Message == "-h" || a.Message == "--help" || a.Message == "help" {
				message := "To use TTS simply type what you want it to say! use [] to add voices (see !tts -v for voices). add `/` for 0.5s of delay ( `//` is 1 sec and so on)"
				flags.AddActions = action.Flag{
					Active: true,
					Actions: []action.Action{
						&actions.ReplyToMessage{Message: message},
					},
					ActionData: []any{chat},
				}
				return flags
			}
			if a.Message == "-v" || a.Message == "--voices" || a.Message == "voices" {
				var message strings.Builder
				message.WriteString("Add [] around a voice to speak in it. voicies:")

				for v := range voices {
					message.WriteString(", " + v)
				}

				flags.AddActions = action.Flag{
					Active: true,
					Actions: []action.Action{
						&actions.ReplyToMessage{Message: message.String()},
					},
					ActionData: []any{chat},
				}
				return flags
			}
		}
	}
	if a.Name == "" {
		a.Name = "INTERNAL"
	}

	MakeTTS(a.Message, a.Name)

	cmd := exec.Command("scp", "-i", "/root/.ssh/tts_key", "-o", "StrictHostKeyChecking=no", "/app/data/tts/tts.wav", "turt@192.168.0.182:/home/turt/StreamShit/audio/tts/tts.wav")
	if err := cmd.Run(); err != nil {
		flags.Error = err
		return flags
	}

	return flags
}

func (a TTS) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
