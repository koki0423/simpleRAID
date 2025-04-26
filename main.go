package main

import (
	"fmt"
	"os"
	"sd-raid/funcs"
)

const diskCount = 4 // ディスクの数

func main() {
	paths := []string{"disk/disk1.img", "disk/disk2.img", "disk/disk3.img", "disk/disk4.img"}

	// デバイスオープン、存在しない場合はnilにする
	var files []*os.File
	for _, path := range paths {
		file, err := os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Warning: %s is missing, treating as failed disk\n", path)
			files = append(files, nil)
			continue
		}
		files = append(files, file)
	}

	// RAID初期化
	raid := &funcs.SimpleRAID5{
		Devices:    files,
		BlockSize:  4096,
		StripeSize: (len(files) - 1) * 4096,
		StripeCnt:  0,
	}
	defer raid.Close()

	//ディスク故障時の読み出しテストのときは1をコメントアウトして実行する
	// 1. RAID5にデータを書き込む
	if err := writeToRaid(raid); err != nil {
		fmt.Println("Write error:", err)
		return
	}

	// 2. 書き込んだデータサイズ取得（元データと同じ）
	inputInfo, err := os.Stat("testData/test_input.png")
	if err != nil {
		fmt.Println("Error getting input file info:", err)
		return
	}

	size := int(inputInfo.Size())

	// 3. RAID5からデータを読み出す
	raid.StripeCnt = 0 // ストライプカウントをリセット
	recoveredData, err := raid.Read(size)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	// 4. 読み出したデータをファイルに保存
	err = os.WriteFile("testData/recovered_output.png", recoveredData, 0644)
	if err != nil {
		fmt.Println("Error writing recovered file:", err)
		return
	}

	fmt.Println("Recovery completed. Check 'testData/recovered_output.png'")
}

func writeToRaid(raid *funcs.SimpleRAID5) error {
	// 入力データ読み込み
	inputData, err := os.ReadFile("testData/test_input.png")
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// RAID5に書き込み
	err = raid.Write(inputData)
	if err != nil {
		return fmt.Errorf("failed to write to RAID5: %w", err)
	}

	fmt.Println("Write to RAID5 completed successfully.")
	return nil
}
