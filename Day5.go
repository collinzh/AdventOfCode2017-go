package main

import (
	"os"
	"bufio"
	"strconv"
	list2 "container/list"
	"fmt"
	"log"
)

func main() {
	part1()
	part2()
}

func part2() {
	f, err := os.Open("Day5.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer f.Close()
	list := list2.New()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		num, _ := strconv.Atoi(string(line))
		list.PushBack(num)
	}

	array := make([]int, list.Len())
	for next := list.Front(); next != nil; next = next.Next() {
		array = append(array, next.Value.(int))
	}

	position := 0
	count := 0
	for {
		if position < 0 || position >= len(array) {
			break
		}
		newPos := array[position] + position
		if array[position] >= 3 {
			array[position] --
		} else {
			array[position] ++
		}
		position = newPos
		count ++
	}

	fmt.Println(count)
}

func part1() {
	f, err := os.Open("Day5.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer f.Close()
	list := list2.New()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		num, _ := strconv.Atoi(string(line))
		list.PushBack(num)
	}

	array := make([]int, list.Len())
	for next := list.Front(); next != nil; next = next.Next() {
		array = append(array, next.Value.(int))
	}

	position := 0
	count := 0
	for {
		if position < 0 || position >= len(array) {
			break
		}
		newPos := array[position] + position
		array[position] ++
		position = newPos
		count ++
	}

	fmt.Println(count)
}
