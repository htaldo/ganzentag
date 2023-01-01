package main

import (
	"regexp"
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"os"
)

type Task struct {
	ID 					int
	Name 				string
	Description 		string
	Predecessor			[]int
	Duration			int
	ES, EF, LS, LF, H	int
}

func main() {
	//get filename from 1st argument
	//filename := os.Args[1]

	//specify a filename directly in the source code
	filename := "input.tsv"

	tasks, n := getTasks(filename)
	
	//print initial tasks struct
	//for _, task := range tasks {
	//	fmt.Printf("%v\n", task)
	//}

	//print adjacency matrix with all paths
	fmt.Println("---ADJACENCY MATRIX---")
	matrix := adjMatrix(tasks)
	fmt.Printf("\n\n\n")

	//compute EF and LF for every task
	for i := 0; i < len(tasks); i++ {
		for j := 0; j < len(tasks); j++ {
			if (matrix[i][j] == 1){
				if tasks[i].EF > tasks[j].ES {
					tasks[j].ES = tasks[i].EF
				}
				tasks[j].EF = tasks[j].ES + tasks[j].Duration
			}
		}
	}
	
	//set END values
	tasks[n].LF = tasks[n].EF
	tasks[n].LS = tasks[n].LF

	//compute LS, LF and H for every task
	for i := n; i >= 0; i-- {
		for j := n; j >= 0; j-- {
			if (matrix[i][j] == 1){
				if (tasks[i].LF == 0 || tasks[j].LS < tasks[i].LF) && tasks[j].LS != 0 {
					tasks[i].LF = tasks[j].LS
				}
				tasks[i].LS = tasks[i].LF - tasks[i].Duration
				tasks[i].H = tasks[i].LF - tasks[i].EF
			}
		}
	}

	//set START values
	tasks[0].LF = 0
	tasks[0].LS = 0
	tasks[0].H = 0

	//print full tasks struct
	fmt.Println("---TASKS INFO---")
	for _, task := range tasks[:] {
		fmt.Printf("%v\n", task)
	}
	fmt.Printf("\n\n\n")

	//print adjacency matrix of the critic path
	fmt.Println("---CRITICAL PATH MATRIX---")
	cpm := cpMatrix(tasks, matrix)
	fmt.Printf("\n\n\n")

	//print critical path
	fmt.Println("---CRITICAL PATH---")
	criticalPath(tasks, cpm)
}

func getTasks(filename string) ([]Task, int) {
	//create slice of tasks
	//account for the project start and end (add 2)
	size := lineCount(filename) + 2
	tasks := make([]Task, size)	

	//define the start of the project
	tasks[0].ID = 0
	tasks[0].Name = "START"
	tasks[0].Duration = 0
	tasks[0].ES = 0

	//define tasks specified in tsv
	f, _ := os.Open(filename)
	tScanner := bufio.NewScanner(f)
	i := 1
	for tScanner.Scan() {
		//fmt.Printf("%v\n", i)
		fields := splitRow(tScanner.Text())
		tasks[i].ID = i
		tasks[i].Name = fields[0]
		tasks[i].Description = fields[1]
		tasks[i].Predecessor = splitIntByComma(fields[2])
		tasks[i].Duration, _ = strconv.Atoi(fields[3])
		i++		
	}
	f.Close()

	//define the end of the project
	n := i
	tasks[n].ID = n
	tasks[n].Name = "END"
	tasks[n].Predecessor = finalTasks(tasks) 
	tasks[n].Duration = 0 

	return tasks, n
}

func splitRow(row string) []string {
	array := regexp.MustCompile("\t+").Split(row, -1)
	return array
}

func splitIntByComma(field string) []int {
	var ints []int
	strInts := strings.Split(field, ",")
	for _, strInt := range strInts[:] {
		newInt, _ := strconv.Atoi(strInt)
		ints = append(ints, newInt)
	}
	return ints
}

func lineCount(filename string) int{
	f, _ := os.Open(filename)
	tScanner := bufio.NewScanner(f)
	lineCount := 0
	for tScanner.Scan() {
		lineCount++
	}
	f.Close()

	return lineCount
}

func adjMatrix(tasks []Task) [][]int {
	matrix := createSquareMatrix(len(tasks))

	for i := 0; i < len(tasks); i++ {
		for j := 0; j < len(tasks); j++ {
			//switching i and j in the next line gives the transpose
			matrix[i][j] = contains(tasks[j].Predecessor, i)
			fmt.Printf("%d ", matrix[i][j])
		}
			fmt.Printf("\n")
	}
	return matrix
}

func contains(list []int, num int) int {
	for _, elem := range list {
		if elem == num {
			return 1
		}
	}
	return 0
}

func cpMatrix(tasks []Task, adjMatrix [][]int) [][]int {
	matrix := createSquareMatrix(len(tasks))
	for i := 0; i < len(tasks); i++ {
		for j := 0; j < len(tasks); j++ {
			if adjMatrix[i][j] == 1 {
				if tasks[i].H==0 && tasks[j].H == 0 {
					matrix[i][j] = 1
				}
			}
			fmt.Printf("%d ", matrix[i][j])
		}
		fmt.Printf("\n")
	}
	return matrix
}

func createSquareMatrix(size int) [][]int {
	matrix := make ([][]int, size)
	for i := 0; i < size; i++{
		matrix[i] = make([]int, size)
	}
	return matrix
}

func finalTasks(tasks []Task) []int {
	var finalTasks []int
	//skip START and END
	for i := 1; i < len(tasks) - 1; i++ {
		test := 0
		for j := 1; j < len(tasks) - 1; j++ {
			if contains(tasks[j].Predecessor,i) == 1 {
				test = 1
			}
		}
		if test == 0{
			finalTasks = append(finalTasks, i)
		}
	}
	return finalTasks
}

func criticalPath(tasks []Task, cpMatrix [][]int) {
	fmt.Println("Nodes\t\tDuration")
	//since cpMatrix is square, the index doesn't matter
	completionTime := 0
	n := len(tasks)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if cpMatrix[i][j] == 1 {
				fmt.Println(tasks[i].Name, "->", tasks[j].Name, "\t", tasks[i].Duration)
				completionTime += tasks[i].Duration
			}
		}
	}
	fmt.Println("Total\t\t", completionTime)
}
