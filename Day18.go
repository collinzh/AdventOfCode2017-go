package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

type Instruction struct {
	instruction string
	param1      string
	param2      string
}

func main() {

	if file, err := os.Open("Day18_t.txt"); err == nil {
		fmt.Println(Day18P1(file))
		defer file.Close()
	}
	if file, err := os.Open("Day18.txt"); err == nil {
		fmt.Println(Day18P1(file))
		defer file.Close()
	}
	if file, err := os.Open("Day18.txt"); err == nil {
		Day18P2(file)
		defer file.Close()
	}
}

func Day18P2(input *os.File) int {
	instructions, numInstructions := ParseInstruction(input)
	zeroToOne := make(chan int64, 500)
	oneToZero := make(chan int64, 500)
	monitor0 := make(chan bool)
	monitor1 := make(chan bool)
	doneZero := make(chan bool)
	doneOne := make(chan bool)
	go func() {
		recv0 := false
		recv1 := false
		for {
			select {
			case res := <-monitor0:
				recv0 = res
			case res := <-monitor1:
				recv1 = res
			}
			if recv0 && recv1 && len(zeroToOne) == 0 && len(oneToZero) == 0 {
				close(zeroToOne)
				close(oneToZero)
				break
			}
		}
	}()

	go InstructionRunner(0, instructions, numInstructions, zeroToOne, oneToZero, monitor0, doneZero)
	go InstructionRunner(1, instructions, numInstructions, oneToZero, zeroToOne, monitor1, doneOne)

	<-doneZero
	<-doneOne

	return 0
}

func InstructionRunner(pid int, instructions []Instruction, numIns int, recv <-chan int64, send chan<- int64, monitor chan bool, doneChannel chan<- bool) {
	registers := make(map[string]int64, 26)
	registers["p"] = int64(pid)

	pcounter := 0
	msgSent := 0
	channelClosed := false

	for pcounter < numIns {
		nextInst := instructions[pcounter]
		pcounter ++
		var num int64
		var err error
		if num, err = strconv.ParseInt(nextInst.param2, 10, 0); err != nil {
			num = registers[nextInst.param2]
		}

		if nextInst.instruction == "set" {
			registers[nextInst.param1] = num

		} else if nextInst.instruction == "add" {
			registers[nextInst.param1] += num

		} else if nextInst.instruction == "mul" {
			registers[nextInst.param1] *= num

		} else if nextInst.instruction == "mod" {
			registers[nextInst.param1] %= num

		} else if nextInst.instruction == "snd" {
			var toSend int64
			if toSend, err = strconv.ParseInt(nextInst.param1, 10, 0); err != nil {
				toSend = registers[nextInst.param1]
			}
			msgSent ++
			fmt.Printf("PID %d sending msg %d - %d\n", pid, msgSent, toSend)
			send <- toSend

		} else if nextInst.instruction == "rcv" {
			monitor <- true
			fmt.Printf("PID %d waiting\n", pid)
			recvd, hasMore := <-recv
			if !hasMore {
				channelClosed = true
				break
			}
			fmt.Printf("PID %d received %d\n", pid, recvd)
			registers[nextInst.param1] = recvd
			monitor <- false

		} else if nextInst.instruction == "jgz" {
			var val int64
			if val, err = strconv.ParseInt(nextInst.param1, 10, 0); err != nil {
				val = registers[nextInst.param1]
			}

			if val > 0 {
				pcounter --
				pcounter += int(num)
			}
		}
	}
	if !channelClosed {
		monitor <- true
	}

	fmt.Printf("PID %d sent %d messages\n", pid, msgSent)

	doneChannel <- true
}

func ParseInstruction(input *os.File) ([]Instruction, int) {
	instructions := make([]Instruction, 1)
	scanner := bufio.NewScanner(input)
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		cmd := words[0]
		reg := words[1]
		var target string
		if len(words) > 2 {
			target = words[2]
		} else {
			target = ""
		}

		if counter == 0 {
			instructions[0] = Instruction{cmd, reg, target}
		} else {
			instructions = append(instructions, Instruction{cmd, reg, target})
		}

		counter ++
	}
	return instructions, counter
}

func Day18P1(input *os.File) int {
	registers := make(map[string]int, 26)
	var lastPlayed int
	//var lastRecovered int
	instructions, numIns := ParseInstruction(input)
	pcounter := 0

	for pcounter < numIns {
		nextIns := instructions[pcounter]
		pcounter ++
		var num int
		var err error
		if num, err = strconv.Atoi(nextIns.param2); err != nil {
			num = registers[nextIns.param2]
		}

		if nextIns.instruction == "set" {
			registers[nextIns.param1] = num

		} else if nextIns.instruction == "add" {
			registers[nextIns.param1] += num

		} else if nextIns.instruction == "mul" {
			registers[nextIns.param1] *= num

		} else if nextIns.instruction == "mod" {
			registers[nextIns.param1] %= num

		} else if nextIns.instruction == "snd" {
			lastPlayed = registers[nextIns.param1]

		} else if nextIns.instruction == "rcv" {
			if registers[nextIns.param1] != 0 {
				return lastPlayed
			}

		} else if nextIns.instruction == "jgz" {
			if registers[nextIns.param1] > 0 {
				pcounter --
				pcounter += num
			}
		}
	}

	return lastPlayed
}
