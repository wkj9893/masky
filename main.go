package main

import (
	"log"

	"github.com/lucas-clemente/quic-go/http3"
)

// func main() {
// 	count := 0
// 	wg := sync.WaitGroup{}
// 	wg.Add(100)
// 	for i := 0; i < 100; i++ {
// 		go func() {
// 			count++
// 			temp := count
// 			os.Setenv("1", fmt.Sprint(temp))
// 			if os.Getenv("1") != fmt.Sprint(temp) {
// 				fmt.Println("zzzzzzzzz")
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }

func main() {
	log.Fatal(http3.ListenAndServe("", "", "", nil))
}
