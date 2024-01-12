package main

import "expense-tracker/delivery"

func main() {
	delivery.NewServer().Run()
}
