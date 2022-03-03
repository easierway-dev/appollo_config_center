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
		if consulValue, ok := consulKV[k]; ok {
			if apolloValue != consulValue {
				notEqualKey = append(notEqualKey, k)
			} else {
				continue
			}
		} else {
			notExistKey = append(notExistKey, k)
		}
	}
	fmt.Println("notEqualKey=", notEqualKey)
	fmt.Println("notExistKey=", notExistKey)
	return nil
}
