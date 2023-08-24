package main

func main() {
	tasks := []string{"task1", "task2", "task3"}

	ch := make(chan string, 3)

	for _, task := range tasks {
		ch <- task
	}
	//Sending task4 to channel.
	ch <- "task4" // deadlock occurs here

	go worker(ch)
}

// G2
func worker(ch chan string) {
	for {
		t := <-ch
		println(t)
	}
}
