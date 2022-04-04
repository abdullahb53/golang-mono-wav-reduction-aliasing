package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func limitingDegeriGirisi(limit_deger uint8) uint8 {
	fmt.Println("->LIMITING DEGERI GIRISI:(2,4,6,8,16,32) 0->EXIT")
	fmt.Scanln(&limit_deger)
	return limit_deger
}

func ReadFile(filename string) ([]byte, error) { //Dosya okuma islemi.
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return b, err
}

func Create(newFileName string) { //Dosya olusturma islemi.

	f, err := os.Create(newFileName) //Istek dosya olusturuldu.
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func Write(b []byte, newFileName string) { //Dosyaya yazma islemi.
	err_ := os.WriteFile(newFileName, b, 0666) //b []byte turunden olan dizi <dosyaya> yazildi.
	if err_ != nil {
		log.Fatal(err_)
	}

}

func makeNewArray(taintedLoveArray_b *[]byte) []byte {
	array := make([]byte, len(*taintedLoveArray_b)/8)

	for i := 0; i < len(array)-44; i++ {
		array[i+44] = (*taintedLoveArray_b)[i*8]
	}
	for i := 0; i < 44; i++ {
		array[i] = (*taintedLoveArray_b)[i]
	}

	fmt.Println(array)

	return array
}

// turn a 32-bit byte array into an int
func bits32ToInt(b *[]byte) int {

	var payload uint32
	buf := bytes.NewReader(*b)
	err := binary.Read(buf, binary.LittleEndian, &payload)
	if err != nil {
		panic(err)
	}
	return int(payload) // easier to work with ints
}

func main() {

	taintedLove_arr, err := ReadFile("./taintedLove.wav") //sarki []byte array olarak okundu.
	if err != nil {
		panic(err)
	}
	newFileName := "./asdas.wav" //Yeni dosyanin ismi.

	Create(newFileName)                   //Yeni dosya olusturuldu.
	arr := makeNewArray(&taintedLove_arr) //Yeni dosya icin set edilecek yeni bir dizi olusturuldu.
	arr2 := bits32ToInt(&arr)
	fmt.Println(arr2)
	Write(arr, newFileName) //Yeni dosyanin icine, yeni dizi geçirildi.

	/*
		array := make([]byte, len(soundFile)/8)

		for i := 0; i < len(array)-44; i++ {
			array[i+44] = soundFile[i*8]
		}
		for i := 0; i < 44; i++ {
			array[i] = soundFile[i]
		}

		Write(soundArray)
		//Write(array)
	*/

	var limitingDegeri uint8
	limitingDegeri = limitingDegeriGirisi(limitingDegeri)

	for {
		if limitingDegeri != 0 {
			switch limitingDegeri {
			case 2:
				fmt.Println("one")
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)
			case 4:
				fmt.Println("two")
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)
			case 6:
				fmt.Println("three")
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)
			case 8:
				fmt.Println("three")
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)
			case 16:
				fmt.Println("three")
				limitingDegeri = limitingDegeriGirisi(limitingDegeri)
			case 32:
				fmt.Println("three")
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
