package main

func main() {

	var terminal TerminalI = Terminal{currentLDB: ""}
	var commander CommanderI = Commander{}

	go terminal.run(&commander)

	// var commander CommanderI = Commander{}

	// commander.CreateLDB("cc_coding_challenge")

	// commander.Set("cc_coding_challenge", "kai1", 1000)
	// commander.Set("cc_coding_challenge", "kai2", 1100)
	// commander.Set("cc_coding_challenge", "kai3", 1200)
	// commander.Set("cc_coding_challenge", "kai4", 1300)
	// commander.Set("cc_coding_challenge", "kai5", 1400)

	// fmt.Println("5 :", commander.GetRank("cc_coding_challenge", "kai5"))
	// fmt.Println("4 : ", commander.GetRank("cc_coding_challenge", "kai4"))
	// fmt.Println("3 : ", commander.GetRank("cc_coding_challenge", "kai3"))
	// fmt.Println("2 : ", commander.GetRank("cc_coding_challenge", "kai2"))
	// fmt.Println("1 : ", commander.GetRank("cc_coding_challenge", "kai1"))

	// commander.CreateLDB("cc2")
	// commander.Set("cc2", "kai1", 1400)
	// commander.Set("cc2", "kai2", 1300)
	// commander.Set("cc2", "kai3", 1200)
	// commander.Set("cc2", "kai4", 1100)
	// commander.Set("cc2", "kai5", 1000)

	// fmt.Println("1 :", commander.GetRank("cc2", "kai5"))
	// fmt.Println("2 : ", commander.GetRank("cc2", "kai4"))
	// fmt.Println("3 : ", commander.GetRank("cc2", "kai3"))
	// fmt.Println("4 : ", commander.GetRank("cc2", "kai2"))
	// fmt.Println("5 : ", commander.GetRank("cc2", "kai1"))

}
