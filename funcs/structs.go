package funcs

import (
	"os"
)

type SimpleRAID5 struct {
	Devices    []*os.File  // 使用する仮想デバイス（ファイル）
	BlockSize  int         // 1ブロックのサイズ
	StripeSize int         // 1ストライプにおけるデータサイズ
	StripeCnt  int         // 現在までに処理したストライプ数
}
