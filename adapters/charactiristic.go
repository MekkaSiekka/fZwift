package adapters

import (
	"fmt"

	"github.com/MekkaSiekka/fZwift/defines"
)

type PeripheralCharactiristic interface {
	ParseCharBuffer(buf []byte) map[int]int
}

type FitnessMachineChar struct{}

/*
 */
func GetNthBitAsInt(buf []byte, idx int) int {
	byteIdx := idx / 8
	idxInByte := idx % 8
	bt := buf[byteIdx]
	for i := 0; i < 8; i++ {
		bit := (bt >> uint(i)) & 1
		fmt.Print(bit)
	}
	fmt.Println()
	fmt.Printf("idxInByte %v \n", idxInByte)
	bit := bt >> uint(idxInByte) & 1
	fmt.Printf("result bit %v \n", bit)
	return int(bit)
}
func (fm *FitnessMachineChar) ParseCharBuffer(buf []byte) map[int]int {
	ret := make(map[int]int)
	ret[defines.CandenceMeasureEnalbe] = GetNthBitAsInt(buf, defines.CandenceMeasureEnalbe)
	ret[defines.PowerMeasureEnabled] = GetNthBitAsInt(buf, defines.PowerMeasureEnabled)
	fmt.Printf("Return map %v \n", ret)
	return ret
}
