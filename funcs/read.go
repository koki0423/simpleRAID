package funcs

import (
	"fmt"
	"io"
)

func (r *SimpleRAID5) Read(size int) ([]byte, error) {
	offset := 0
	var result []byte

	for offset < size {
		parityIndex := r.StripeCnt % len(r.Devices)

		//ストライプからデータを読み込む
		dataStripe, err := r.readStripe(parityIndex)
		if err != nil {
			return nil, err
		}
		if len(dataStripe) > (size - offset) {
			dataStripe = dataStripe[:size-offset]
		}
		result = append(result, dataStripe...)

		offset += len(dataStripe)
		r.StripeCnt++
	}

	return result, nil
}

func (r *SimpleRAID5) readStripe(parityIndex int) ([]byte, error) {
	var stripeData [][]byte
	var parityBlock []byte

	brokenDisk := -1

	for deviceIdx, device := range r.Devices {
		if device == nil {
			if brokenDisk != -1 {
				return nil, fmt.Errorf("more than one disk failed")
			}
			brokenDisk = deviceIdx
			continue
		}

		offset := int64(r.StripeCnt) * int64(r.BlockSize)
		_, err := device.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, err
		}

		block := make([]byte, r.BlockSize)
		_, err = device.Read(block)
		if err != nil {
			return nil, err
		}

		if deviceIdx == parityIndex {
			parityBlock = block
		} else {
			stripeData = append(stripeData, block)
		}
	}

	// 復元作業の開始位置
	var data []byte

	if brokenDisk == -1 {
		// 壊れていない場合、単純にすべて連結
		for i := 0; i < len(stripeData); i++ {
			data = append(data, stripeData[i]...)
		}
		return data, nil
	} else {
		// 復元する場合
		restoreBlock := make([]byte, r.BlockSize)
		copy(restoreBlock, parityBlock)

		for _, block := range stripeData {
			for i := 0; i < r.BlockSize; i++ {
				restoreBlock[i] ^= block[i]
			}
		}

		blockIdx := 0
		for deviceIdx := 0; deviceIdx < len(r.Devices); deviceIdx++ {
			if deviceIdx == parityIndex {
				continue
			}
			if deviceIdx == brokenDisk {
				data = append(data, restoreBlock...)
			} else {
				data = append(data, stripeData[blockIdx]...)
				blockIdx++
			}
		}

		return data, nil
	}
}
