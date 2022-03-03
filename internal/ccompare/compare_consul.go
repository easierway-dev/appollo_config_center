package ccompare

import "fmt"

const Application = "application"

var notExistKey []string
var notEqualKey []string

func ApolloCompareWithConsul() error {
	apolloKV, err := GetSingleNameSpaceInfo(Application)
	if err != nil {
		return err
	}
	consulKV, err := GetConsulKV()
	if err != nil {
		return err
	}
	for k, apolloValue := range apolloKV {
		consulValue, ok := consulKV[k]
		if !ok {
			notExistKey = append(notExistKey, k)
		}
		if apolloValue != consulValue {
			notEqualKey = append(notEqualKey, k)
		}
	}
	fmt.Println("notEqualKey=", notEqualKey)
	fmt.Println("notExistKey=", notExistKey)
	return nil
}
