package main

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"github.com/ebfe/scard"
	"github.com/gogetth/sscard"
)

func exampleThaiIDCard() {

	// Establish a PC/SC context

	fmt.Println("Establish a PC/SC context")
	context, err := scard.EstablishContext()
	if err != nil {
		fmt.Println("Error EstablishContext:", err)
		return
	}

	// Release the PC/SC context (when needed)
	defer context.Release()

	// List available readers
	readers, err := context.ListReaders()
	if err != nil {
		fmt.Println("Error ListReaders:", err)
		return
	}

	// Use the first reader
	reader := readers[0]
	fmt.Println("Using reader:", reader)

	// Connect to the card
	card, err := context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		fmt.Println("Error Connect:", err)
		return
	}

	// Disconnect (when needed)
	defer card.Disconnect(scard.LeaveCard)

	// Send select APDU
	selectRsp, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardSelect)
	if err != nil {
		fmt.Println("Error Transmit:", err)
		return
	}
	fmt.Println("resp sscard.APDUThaiIDCardSelect: ", selectRsp)

	cid, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardCID)

	if err != nil {
		fmt.Println("Error APDUGetRsp: ", err)
		return
	}
	fmt.Printf("cid: _%s_\n", string(cid))

	fullnameEN, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardFullnameEn)
	if err != nil {
		fmt.Println("Error APDUGetRsp: ", err)
		return
	}
	fmt.Printf("fullnameEN: _%s_\n", ConvertTIS620toUTF8(string(fullnameEN)))

	fullnameTH, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardFullnameTh)
	if err != nil {
		fmt.Println("Error APDUGetRsp: ", err)
		return
	}
	fmt.Printf("fullnameTH: _%s_\n", ConvertTIS620toUTF8(string(fullnameTH)))

	birth, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardBirth)
	if err != nil {
		fmt.Println("Error APDUGetRsp: ", err)
		return
	}
	fmt.Printf("birth: _%s_\n", ConvertTIS620toUTF8(string(birth)))

	cardPhotoJpg, err := sscard.APDUGetBlockRsp(card, sscard.APDUThaiIDCardPhoto, sscard.APDUThaiIDCardPhotoRsp)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	//Write Image
	n2, err := sscard.WriteBlockToFile(cardPhotoJpg, "./idcPhoto.jpg")
	if err != nil {
		fmt.Println("Error WriteBlockToFile: ", err)
		return
	}
	fmt.Printf("Img wrote %d bytes\n", n2)

}

func ConvertTIS620toUTF8(tis620 string) (valueUTF8 string) {
	tis620Valure := []byte(tis620)
	dec := charmap.Windows874.NewDecoder()
	makeUTF := make([]byte, len(tis620Valure)*3)
	n, _, err := dec.Transform(makeUTF, tis620Valure, false)
	if err != nil {
		return "This Filed can't Convert Tis620 to UTF8"
	}
	valueUTF8 = string(makeUTF[:n])
	valueUTF8 = strings.Trim(valueUTF8, "\u0000")
	valueUTF8 = strings.ReplaceAll(valueUTF8, "#", " ")

	valueUTF8 = strings.TrimSpace(valueUTF8)
	return
}
