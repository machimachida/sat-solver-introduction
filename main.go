package main

import (
	"encoding/csv"
	"os"
)

const (
	NUMBER_OF_TEAMS    = 10
	NUMBER_OF_DAYS_OFF = 20
)

func main() {
	daysoff, err := ReadQuestionCSV("question.csv")
	if err != nil {
		panic(err)
	}

	// 変数を作成する
	// 1日目から20日目までの10チームの出勤状況を表す変数を作成する

}

func ReadQuestionCSV(filename string) (map[string][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	daysoff := make(map[string][]int, NUMBER_OF_TEAMS)
	for _, row := range rows {
		if daysoff[row[0]] == nil {
			daysoff[row[0]] = make([]int, 0, NUMBER_OF_DAYS_OFF)
			daysoff[row[0]] = append(daysoff[row[1]], 0)
		}
	}

	return daysoff, nil
}

func MakeCNFAtLeastOne() {
}
