package taskpage
import(
	"logs"
	"queue"
	"time"
)

func handleResponse(){
	for{
		if PageQueue.Empty() == false {
			task := PageQueue.PopOneTask()
			
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func Start(){
	go handleResponse()
}