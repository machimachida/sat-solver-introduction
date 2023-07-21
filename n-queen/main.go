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

func main() {
	// 変数を作成する
	// 行列の各マスにクイーンが置かれているかどうかを表す変数を作成する
	// 8x8のマス目なので、8x8=64個の変数が必要になる
	// 1行目の1列目の変数は、11とする
	// 1行目の2列目の変数は、12とする
	// 1行目の3列目の変数は、13とする
	// 1行目の4列目の変数は、14とする
	variables := make([][]int, NUMBER_OF_QUEENS)
	for i := 0; i < NUMBER_OF_QUEENS; i++ {
		variables[i] = make([]int, NUMBER_OF_QUEENS)
		for j := 0; j < NUMBER_OF_QUEENS; j++ {
			variables[i][j] = i*10 + j + 1 // 負の数が必要になるので、i=0,j=0のとき0にならないよう+1する
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

	// 3. 斜めには必ず1つのクイーンが置かれている
	diagonalCNF := make([][]int, 0, NUMBER_OF_QUEENS*(NUMBER_OF_QUEENS-1)/2+1)
	for i := -NUMBER_OF_QUEENS; i < NUMBER_OF_QUEENS; i++ {
		for j := 0; j < NUMBER_OF_QUEENS; j++ {
			diagonal := make([]int, 0, NUMBER_OF_QUEENS)
			for k := 0; k < NUMBER_OF_QUEENS; k++ {
				if 0 <= j+i+k && j+i+k < NUMBER_OF_QUEENS {
					diagonal = append(diagonal, variables[k][j+i+k])
				}
			}
			if len(diagonal) == 0 {
				continue
			}

			diagonalCNF = append(diagonalCNF, diagonal) // at least one
			diagonalCNF = append(diagonalCNF, MakeCNFAtMostOne(diagonal)...)
		}
	}
	fmt.Println(diagonalCNF)

	clauses := make([][]int, 0, len(rowCNF)+len(columnCNF)+len(diagonalCNF))
	clauses = append(clauses, rowCNF...)
	clauses = append(clauses, columnCNF...)
	clauses = append(clauses, diagonalCNF...)

	solver := sat.New()
	solver.AddFormula(cnf.NewFormulaFromInts(clauses))
	sat := solver.Solve()
	if !sat {
		println("UNSAT")
		return
	}
	solution := solver.Assignments()
	for i := 0; i < NUMBER_OF_QUEENS; i++ {
		for j := 0; j < NUMBER_OF_QUEENS; j++ {
			if solution[variables[i][j]] {
				print("Q")
			} else {
				print(".")
			}
		}
		println()
	}
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
