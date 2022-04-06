package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
)

func limitingDegeriGirisi(limit_deger uint8) uint8 {
	fmt.Println("->LIMITING DEGERI GIRISI:(2,4,6,8,16,32) 0->EXIT")
	fmt.Scanln(&limit_deger)
	return limit_deger
}

func ReadFile(filename string) (*[]byte, error) { //Dosya okuma islemi.
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return &b, err
}

func Create(newFileName string) { //Dosya olusturma islemi.

	f, err := os.Create(newFileName) //Istek dosya olusturuldu.
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func Write(b *[]byte, newFileName string) { //Dosyaya yazma islemi.
	err_ := os.WriteFile(newFileName, *b, 0666) //b []byte turunden olan dizi <dosyaya> yazildi.
	if err_ != nil {
		log.Fatal(err_)
	}

}

func makeDataArray(taintedLoveArray_b *[]byte, wh *WaveHeader, ss int) *[]byte {
	//array := make([]byte, 44+((len(*taintedLoveArray_b)-44)/ss))

	var iterasyon int = 8
	if ss == 2 || ss == 4 || ss == 6 || ss == 8 {
		iterasyon = 8
	} else {
		iterasyon = 32
	}

	var dataCount uint16 = 0
	var kalan int = (len((*taintedLoveArray_b)) - 44) % (iterasyon)
	var diziSonuSorgusu int = ((len((*taintedLoveArray_b)) - 44) - kalan)

	array := []byte{}

	array = append(array, wh.ChunkID...)
	//wh.ChunkSize = ((len(*taintedLoveArray_b) - 44) / ss) + 44 - 8 - 36
	//wh.ChunkSize = diziSonuSorgusu/ss + 44
	wh.ChunkSize = (diziSonuSorgusu/iterasyon)*(iterasyon/ss) + 44
	//fmt.Println("wh.ChunkSize", wh.ChunkSize)
	array = append(array, int32ToBytes(wh.ChunkSize)...)
	array = append(array, []byte(wh.Format)...)

	array = append(array, wh.Subchunk1ID...)
	array = append(array, int32ToBytes(wh.Subchunk1Size)...)
	array = append(array, int16ToBytes(wh.AudioFormat)...)

	wh.NumChannels = wh.NumChannels * wh.BitsPerSample / 8
	array = append(array, int16ToBytes(wh.NumChannels)...) // NumChannels * BitsPerSample / 8 (number of bytes per sample)

	wh.SampleRate = wh.SampleRate / (ss * 2) //"/ ss" sildim.!
	array = append(array, int32ToBytes(wh.SampleRate)...)

	wh.ByteRate = wh.SampleRate * wh.NumChannels * wh.BitsPerSample / 8
	array = append(array, int32ToBytes(wh.ByteRate)...) // SampleRate * NumChannels * BitsPerSample / 8

	wh.BlockAlign = wh.NumChannels * wh.BitsPerSample / 8
	array = append(array, int16ToBytes(wh.BlockAlign)...) // NumChannels * BitsPerSample / 8 (number of bytes per sample)

	array = append(array, int16ToBytes(wh.BitsPerSample)...)
	//fmt.Println("wh.BitsPerSample", wh.BitsPerSample)

	array = append(array, wh.Subchunk2ID...)

	//wh.Subchunk2Size = (len(*taintedLoveArray_b) - 44) / ss
	//wh.Subchunk2Size = diziSonuSorgusu /
	wh.Subchunk2Size = (diziSonuSorgusu / iterasyon) * (iterasyon / ss)
	//fmt.Println("wh.Subchunk2Size :", wh.Subchunk2Size)
	array = append(array, int32ToBytes(wh.Subchunk2Size)...)

	for i := 44; i < diziSonuSorgusu; i += iterasyon {
		if diziSonuSorgusu == i {
			break
		}
		for k := 0; k < (iterasyon / ss); k++ {
			if diziSonuSorgusu == i {
				break
			}
			array = append(array, (*taintedLoveArray_b)[i+k])
			dataCount++
		}

	}

	/*
		for i := 44; i < diziSonuSorgusu; i += iterasyon {
			if diziSonuSorgusu == i {
				break
			}
			for k := 0; k < (iterasyon / ss); k++ {
				if diziSonuSorgusu == i {
					break
				}
				array = append(array, (*taintedLoveArray_b)[i+k])
				dataCount++
			}

		}

			for i := 44; i < (len(*taintedLoveArray_b)); i += iterasyon {
				if diziSonuSorgusu == i {
					break
				}
				for k := 0; k < (iterasyon / ss); k++ {
					if diziSonuSorgusu == i {
						break
					}
					array = append(array, (*taintedLoveArray_b)[i+k])
					dataCount++
				}

			}
	*/

	v := reflect.ValueOf(*wh)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	fmt.Println(values)

	//fmt.Println(array)

	return &array
}

func makeHeaderArray(taintedLoveArray_b *[]byte) *[]byte {
	array := make([]byte, 44)

	for i := 0; i < len(array); i++ {
		array[i] = (*taintedLoveArray_b)[i]
	}

	//fmt.Println(array)

	return &array
}

// byte array int formatina cevirildi.
func bits32ToInt(b *[]byte) int {

	var payload uint32
	buf := bytes.NewReader(*b)
	err := binary.Read(buf, binary.LittleEndian, &payload)
	if err != nil {
		panic(err)
	}
	return int(payload) // easier to work with ints
}

func bits16ToInt(b *[]byte) int {
	if len(*b) != 2 {
		panic("Expected size 4!")
	}
	var payload int16
	buf := bytes.NewReader(*b)
	err := binary.Read(buf, binary.LittleEndian, &payload)
	if err != nil {
		// TODO: make safe
		panic(err)
	}
	return int(payload) // easier to work with ints
}

func readHeader(b []byte) WaveHeader {
	hdr := WaveHeader{}
	fmt.Println("\n  ----ANA SES DOSYASI (TAINTEDLOVE)----")
	//----------Chunk id----------
	chunkID := b[0:4]
	//fmt.Println("ChunkID :")
	hdr.ChunkID = chunkID
	fmt.Println("hdr.chunkID->", string(hdr.ChunkID))
	if string(hdr.ChunkID) != "RIFF" {
		// Validation of the header file
		panic("Invalid file")
	}
	//----------------------------

	//----------ChunkSize---------
	chunkSize := b[4:8]
	hdr.ChunkSize = bits32ToInt(&chunkSize) // easier to work with ints
	fmt.Println("chunkSize :", hdr.ChunkSize)
	//fmt.Println("chunkSize b[4:8]:", chunkSize)

	//----------------------------

	//-----------Format-----------
	format := b[8:12]
	fmt.Println("format :", string(format))
	if string(format) != "WAVE" {
		panic("Format should be WAVE")
	}
	hdr.Format = string(format)
	//----------------------------

	//--------Subchunk1ID---------
	Subchunk1ID := b[12:16]
	hdr.Subchunk1ID = Subchunk1ID
	fmt.Println("SubChunk1Id :", string(Subchunk1ID))
	//----------------------------

	//--------Subchunk1Size-------
	Subchunk1Size := b[16:20]
	hdr.Subchunk1Size = bits32ToInt(&Subchunk1Size)
	fmt.Println("Subchunk1Size :", hdr.Subchunk1Size)
	//----------------------------

	//---------AudioFormat--------
	AudioFormat := 1
	hdr.AudioFormat = AudioFormat
	fmt.Println("AudioFormat :", hdr.AudioFormat)
	//----------------------------

	//---------NumChannels--------
	NumChannels := b[22:24]
	hdr.NumChannels = bits16ToInt(&NumChannels)
	println("NumChannels :", hdr.NumChannels)
	//----------------------------

	//---------SampleRate---------
	SampleRate := b[24:28]
	hdr.SampleRate = bits32ToInt(&SampleRate)
	fmt.Println("SampleRate :", hdr.SampleRate)
	//----------------------------

	//---------ByteRate-----------
	// SampleRate * NumChannels * BitsPerSample / 8
	ByteRate := b[28:32]
	hdr.ByteRate = bits32ToInt(&ByteRate)
	fmt.Println("ByteRate :", hdr.ByteRate)
	//----------------------------

	//--------BlockAlign----------
	// NumChannels * BitsPerSample / 8 (number of bytes per sample)
	BlockAlign := b[32:34]
	hdr.BlockAlign = bits16ToInt(&BlockAlign)
	fmt.Println("BlockAlign :", hdr.BlockAlign)
	//----------------------------

	//---------BitsPerSample------
	BitsPerSample := b[34:36]
	hdr.BitsPerSample = bits16ToInt(&BitsPerSample)
	fmt.Println("BitsPerSample", hdr.BitsPerSample)
	//----------------------------

	//-------ExtraParamSize-------
	Subchunk2ID := b[36:40]
	hdr.Subchunk2ID = (Subchunk2ID)
	fmt.Println("Subchunk2ID", string(hdr.Subchunk2ID))
	//----------------------------

	//--------ExtraParams---------
	Subchunk2Size := b[40:44]
	hdr.Subchunk2Size = bits32ToInt(&Subchunk2Size)
	fmt.Println("SubchunkSize2", hdr.Subchunk2Size)
	//----------------------------

	return hdr
}

func int32ToBytes(i int) []byte {
	b := make([]byte, 4)
	in := uint32(i)
	binary.LittleEndian.PutUint32(b, in)
	return b
}

func int16ToBytes(i int) []byte {
	b := make([]byte, 2)
	in := uint16(i)
	binary.LittleEndian.PutUint16(b, in)
	return b
}

func Kur(newFileName *string, taintedLove_arr **[]byte, ss uint8) {

	Create(*newFileName)                     //Yeni dosya olusturuldu.
	arr := makeHeaderArray(*taintedLove_arr) //Yeni dosya icin set edilecek yeni bir dizi olusturuldu.
	oldHeader := readHeader(*arr)
	wArray := makeDataArray(*taintedLove_arr, &oldHeader, int(ss))
	Write(wArray, *newFileName) //Yeni dosyanin icine, yeni dizi geçirildi.
	fmt.Println("\n Dosya yazildi...")
}

func main() {

	taintedLove_arr, err := ReadFile("./taintedLove.wav") //sarki []byte array olarak okundu.
	if err != nil {
		panic(err)
	}
	newFileName := "./asdasd.wav" //Yeni dosyanin ismi.
	var limitingDegeri uint8
	limitingDegeri = 1

	for {
		if limitingDegeri != 0 {
			switch limitingDegeri {
			case 2:
				Kur(&newFileName, &taintedLove_arr, limitingDegeri)
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)

			case 4:
				Kur(&newFileName, &taintedLove_arr, limitingDegeri)
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)

			case 6:
				Kur(&newFileName, &taintedLove_arr, limitingDegeri)
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)

			case 8:
				Kur(&newFileName, &taintedLove_arr, limitingDegeri)
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)

			case 16:
				Kur(&newFileName, &taintedLove_arr, limitingDegeri)
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)

			case 32:
				Kur(&newFileName, &taintedLove_arr, limitingDegeri)
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)

			default:
				fmt.Println("\n IZIN VERILEN DEGERLER -> (2,4,6,8,16,32)")
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)
			}

		} else {
			break
		}

	}

}
