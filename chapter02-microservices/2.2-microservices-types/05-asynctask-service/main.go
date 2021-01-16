// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

func main() {
	ms := NewMicroservice()

	defer ms.Cleanup()

	ms.Start()
}
