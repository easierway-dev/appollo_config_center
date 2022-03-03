package ccompare

import "fmt"

const Application = "application"
var key []string
func ApolloCompareWithConsul() error{
	apolloKV, err := GetSingleNameSpaceInfo(Application)
	if err != nil{
		return err
	}
	consulKV, err := GetConsulKV()
	if err != nil{
		return err
	}
	for k, apolloValue := range apolloKV {
		if consulValue,ok:= consulKV[k];ok{
			if apolloValue == consulValue{
				continue
			}
		}
		key = append(key,k)
	}
	fmt.Println("key=",key)
	return nil
}