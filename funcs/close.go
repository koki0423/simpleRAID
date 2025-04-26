package funcs

import "fmt"

func (r *SimpleRAID5) Close() {
	for _, disk := range r.Devices {
		if disk == nil {
			continue
		}
		disk.Close()
		fmt.Printf("Closed disk: %s\n", disk.Name())
	}
}
