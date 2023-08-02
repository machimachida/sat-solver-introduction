package main

import (
	"fmt"

	"github.com/mitchellh/go-sat"
	"github.com/mitchellh/go-sat/cnf"
)

// n-queen sat solver

const (
	NUMBER_OF_QUEENS = 8
)

var (
	AnswersDoNotWant = [][2]int{
		// 解として不必要な解を列挙する
	}
)

func main() {
	// 変数を作成する
	// 行列の各マスにクイーンが置かれているかどうかを表す変数を作成する
	// 8x8のマス目だと、8x8=64個の変数が必要になる
	variables := make([][]int, NUMBER_OF_QUEENS)
	for i := 0; i < NUMBER_OF_QUEENS; i++ {
		variables[i] = make([]int, NUMBER_OF_QUEENS)
		for j := 0; j < NUMBER_OF_QUEENS; j++ {
			variables[i][j] = generateVariables(i, j)
		}
	}

	// 制約を作成する
	rowCNF := make([][]int, 0, NUMBER_OF_QUEENS*(NUMBER_OF_QUEENS-1)/2+1)
	// 1. 各行には必ず1つのクイーンが置かれている
	for i := 0; i < NUMBER_OF_QUEENS; i++ {
		rowCNF = append(rowCNF, variables[i]) // at least one
		rowCNF = append(rowCNF, MakeCNFAtMostOne(variables[i])...)
	}

	// 2. 各列には必ず1つのクイーンが置かれている
	columnCNF := make([][]int, 0, NUMBER_OF_QUEENS*(NUMBER_OF_QUEENS-1)/2+1)
	for i := 0; i < NUMBER_OF_QUEENS; i++ {
		column := make([]int, NUMBER_OF_QUEENS)
		for j := 0; j < NUMBER_OF_QUEENS; j++ {
			column[j] = variables[j][i]
		}
		columnCNF = append(columnCNF, column) // at least one
		columnCNF = append(columnCNF, MakeCNFAtMostOne(column)...)
	}

	// 3. 左上から右下への斜めには必ず1つのクイーンが置かれている
	diagonalCNF := make([][]int, 0, NUMBER_OF_QUEENS*(NUMBER_OF_QUEENS-1)/2+1)
	for diff := -NUMBER_OF_QUEENS + 2; diff <= NUMBER_OF_QUEENS-2; diff++ {
		start := 0
		if diff < 0 {
			start = -diff
		}
		count := NUMBER_OF_QUEENS - abs(diff)
		line := make([]int, 0, count)
		for i := start; i < start+count; i++ {
			line = append(line, variables[i][i+diff])
		}
		diagonalCNF = append(diagonalCNF, MakeCNFAtMostOne(line)...)
	}

	// 4. 左下から右上への斜めには必ず1つのクイーンが置かれている
	for diff := -NUMBER_OF_QUEENS + 2; diff <= NUMBER_OF_QUEENS-2; diff++ {
		start := 0
		if diff < 0 {
			start = -diff
		}
		count := NUMBER_OF_QUEENS - abs(diff)
		line := make([]int, 0, count)
		for i := start; i < start+count; i++ {
			line = append(line, variables[i][NUMBER_OF_QUEENS-(i+diff+1)])
		}
		diagonalCNF = append(diagonalCNF, MakeCNFAtMostOne(line)...)
	}

	clauses := make([][]int, 0, len(rowCNF)+len(columnCNF)+len(diagonalCNF))
	clauses = append(clauses, rowCNF...)
	clauses = append(clauses, columnCNF...)
	clauses = append(clauses, diagonalCNF...)

	// 解として、AnswersDoNotWantに含まれる解が出ないようにする
	for _, answer := range AnswersDoNotWant {
		clauses = append(clauses, []int{-variables[answer[0]][answer[1]]})
	}

	solver := sat.New()
	solver.AddFormula(cnf.NewFormulaFromInts(clauses))
	sat := solver.Solve()
	if !sat {
		fmt.Println("UNSAT")
		return
	}

	queensPoints := make([][2]int, 0, NUMBER_OF_QUEENS)
	solution := solver.Assignments()
	for i := 0; i < NUMBER_OF_QUEENS; i++ {
		for j := 0; j < NUMBER_OF_QUEENS; j++ {
			if solution[variables[i][j]] {
				fmt.Print("Q")
				queensPoints = append(queensPoints, [2]int{i, j})
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println(queensPoints)
}

func MakeCNFAtMostOne(vals []int) [][]int {
	l := len(vals)
	cnf := make([][]int, 0, l*(l-1)/2)

	for i := 0; i < l-1; i++ {
		for j := i + 1; j < l; j++ {
			cnf = append(cnf, []int{-vals[i], -vals[j]})
		}
	}
	return cnf
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// generateVariables は、行列の各マスにクイーンが置かれているかどうかを表す変数を作成する
// 1行目の1列目の変数は、11とする
// 1行目の2列目の変数は、12とする
// 1行目の3列目の変数は、13とする
// 1行目の4列目の変数は、14とする
func generateVariables(row, column int) int {
	return row*10 + column + 1 // 負の数が必要になるので、i=0,j=0のとき0にならないよう+1する
}
