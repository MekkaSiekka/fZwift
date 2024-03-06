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

	featureBlock := buf[0:4]
	ret[defines.CandenceMeasureEnabled] = GetNthBitAsInt(featureBlock, defines.CandenceMeasureEnabled)
	ret[defines.PowerMeasureEnabled] = GetNthBitAsInt(featureBlock, defines.PowerMeasureEnabled)
	fmt.Printf("Return map 1%v \n", ret)

	targetSettingBlock := buf[4:8]
	ret[defines.IndoorBikeControlParamEnabled] = GetNthBitAsInt(targetSettingBlock, defines.IndoorBikeControlParamEnabled)
	ret[defines.ResistentTargetSettingEnabled] = GetNthBitAsInt(targetSettingBlock, defines.ResistentTargetSettingEnabled)
	ret[defines.PowerTargetSettingEnabled] = GetNthBitAsInt(targetSettingBlock, defines.PowerTargetSettingEnabled)
	fmt.Printf("Return map 2%v \n", ret)

	return ret
}
