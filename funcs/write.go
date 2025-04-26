package funcs

import ()

func (r *SimpleRAID5) Write(data []byte) error {
	offset := 0
	dataLen := len(data)

	for offset < dataLen {
		end := offset + (r.BlockSize * (len(r.Devices) - 1))
		if end > dataLen {
			end = dataLen
		}

		stripeData := data[offset:end]
		parityIndex := r.StripeCnt % len(r.Devices)

		err := r.writeStripe(stripeData, parityIndex)
		if err != nil {
			return err
		}

		offset += (r.BlockSize * (len(r.Devices) - 1))
		r.StripeCnt++
	}

	return nil
}

func (r *SimpleRAID5) writeStripe(data []byte, parityIndex int) error {
	//4096バイト単位でデータを分割
	var blocks [][]byte
	offset := 0
	dataLen := len(data)

	for offset < dataLen {
		end := offset + r.BlockSize
		if end > dataLen {
			end = dataLen
		}
		block := make([]byte, r.BlockSize)
		copy(block, data[offset:end])
		blocks = append(blocks, block)
		offset += r.BlockSize
	}

	//パリティの計算
	parity := make([]byte, r.BlockSize)
	for i := 0; i < r.BlockSize; i++ {
		for j := 0; j < len(blocks); j++ {
			parity[i] ^= blocks[j][i] //XORの結合法則により，論理の計算順序が変わっても大丈夫
		}
	}

	//データ書き込み
	blockIdx := 0
	for deviceIdx := range r.Devices {
		if r.Devices[deviceIdx] == nil {
			continue
		}

		if parityIndex == deviceIdx {
			// パリティデバイスにはパリティを書き込む
			_, err := r.Devices[deviceIdx].Write(parity)
			if err != nil {
				return err
			}
		} else {
			// データデバイスにはデータを書き込む
			_, err := r.Devices[deviceIdx].Write(blocks[blockIdx])
			if err != nil {
				return err
			}
			blockIdx++ // データブロックのインデックスを進める
		}
	}

	return nil
}
