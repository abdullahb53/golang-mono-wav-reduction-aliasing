package main

var (
	ChunkID          = []byte{0x52, 0x49, 0x46, 0x46} // RIFF
	BigEndianChunkID = []byte{0x52, 0x49, 0x46, 0x58} // RIFX
	WaveID           = []byte{0x57, 0x41, 0x56, 0x45} // WAVE
	Format           = []byte{0x66, 0x6d, 0x74, 0x20} // FMT
	Subchunk2ID      = []byte{0x64, 0x61, 0x74, 0x61} // DATA
)

type Sample float64

/*
╔════════╤════════════════╤══════╤═══════════════════════════════════════════════════╗
║ Offset │ Field          │ Size │ -- start of header                                ║
╠════════╪════════════════╪══════╪═══════════════════════════════════════════════════╣
║ 0      │ ChunkID        │ 4    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 4      │ ChunkSize      │ 4    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 8      │ Format         │ 8    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ --     │ --             │ --   │ -- start of fmt                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 12     │ SubchunkID     │ 4    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 16     │ SubchunkSize   │ 4    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 20     │ AudioFormat    │ 2    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 22     │ NumChannels    │ 2    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 24     │ SampleRate     │ 4    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 28     │ ByteRate       │ 4    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 32     │ BlockAlign     │ 2    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 34     │ BitsPerSample  │ 2    │                                                   ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ * 36   │ ExtraParamSize │ 2    │ Optional! Only when not PCM                       ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ * 38   │ ExtraParams    │ *    │ Optional! Only when not PCM                       ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ --     │ --             │ --   │ -- start of data, assuming PCM                    ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 36     │ Subchunk2ID    │ 4    │ (offset by extra params of subchunk 1 if not PCM) ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 40     │ SubchunkSize   │ 4    │ (offset by extra params of subchunk 1 if not PCM) ║
╟────────┼────────────────┼──────┼───────────────────────────────────────────────────╢
║ 44     │ Data           │ *    │ (offset by extra params of subchunk 1 if not PCM) ║
╚════════╧════════════════╧══════╧═══════════════════════════════════════════════════╝
*/

// WaveHeader describes the header each WAVE file should start with
type WaveHeader struct {
	ChunkID   []byte // should be RIFF on little-endian or RIFX on big-endian systems..
	ChunkSize int
	Format    string // sanity-check, should be WAVE (//TODO: keep as []byte?)

	Subchunk1ID    []byte // should contain "fmt"
	Subchunk1Size  int    // 16 for PCM
	AudioFormat    int    // PCM = 1 (Linear Quantization), if not 1, compression was used.
	NumChannels    int    // Mono 1, Stereo = 2, ..
	SampleRate     int    // 44100 for CD-Quality, etc..
	ByteRate       int    // SampleRate * NumChannels * BitsPerSample / 8
	BlockAlign     int    // NumChannels * BitsPerSample / 8 (number of bytes per sample)
	BitsPerSample  int    // 8 bits = 8, 16 bits = 16, .. :-)
	ExtraParamSize int    // if not PCM, can contain extra params
	ExtraParams    []byte // the actual extra params.

	Subchunk2ID   []byte // Identifier of subchunk
	Subchunk2Size int    // size of raw sound data

}

type (
	bytesToIntF   func([]byte) int
	bytesToFloatF func([]byte) float64
)
