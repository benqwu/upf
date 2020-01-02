package elem

import (
	"bytes"
	"encoding/binary"
	"log"
)

type CreateQER struct {
	EType   IEType
	ELength uint16
	GateStatus
	MBR
	QERID
	QFI
	//TODO::
}

func DecodeCreateQER(data []byte, len uint16) *CreateQER {
	createQER := CreateQER{
		EType:   IETypeCreateQER,
		ELength: len,
	}
	var cursor uint16
	buf := bytes.NewBuffer(data)
	for cursor < createQER.ELength {
		var (
			eType IEType
			eLen  uint16
		)
		if err := binary.Read(buf, binary.BigEndian, &eType); err != nil {
			log.Println(err) //TODO::
		}
		if err := binary.Read(buf, binary.BigEndian, &eLen); err != nil {
			log.Println(err) //TODO::
		}
		eValue := make([]byte, eLen)
		if err := binary.Read(buf, binary.BigEndian, &eValue); err != nil {
			log.Println(err) //TODO::
		}
		switch eType {
		case IETypeGateStatus:
			createQER.GateStatus = *DecodeGateStatus(eValue, eLen)
		case IETypeMBR:
			createQER.MBR = *DecodeMBR(eValue, eLen)
		case IETypeQERID:
			createQER.QERID = *DecodeQERID(eValue, eLen)
		case IETypeQFI:
			createQER.QFI = *DecodeQFI(eValue, eLen)
		default:
			log.Println("err: unknown tlv type", eType) //TODO::
		}
		cursor = cursor + eLen + 4
	}
	return &createQER
}

func EncodeCreateQER(createQER CreateQER) []byte {
	ret := setValue(createQER.EType, createQER.ELength, createQER.GateStatus, createQER.QERID) //GateStatus QERID为M信元
	if HasMBR(createQER.MBR) {
		ret = setValue(ret, createQER.MBR)
	}
	if HasQFI(createQER.QFI) {
		ret = setValue(ret, createQER.QFI)
	}
	return ret
}

//判断是否含有CreateQER
func HasCreateQER(createQER CreateQER) bool {
	if createQER.EType == 0 {
		return false
	}
	return true
}