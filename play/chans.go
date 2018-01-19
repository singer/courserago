package chans
import (
	"fmt"
)

func main(){
 c := make(chan interface{},1)
 c <- 10
 fmt.Println(c) 
 fmt.Println("hi")
 var d int
 d = (<-c).(int)
 fmt.Println(d)
 
}
