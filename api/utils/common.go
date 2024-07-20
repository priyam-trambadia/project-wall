package utils

import "os"

type Lable interface {
	GetID() int64
}

func ExtractIDsMakeArray[T Lable](dataArray []T) []int64 {
	ids := make([]int64, 0)

	for _, val := range dataArray {
		ids = append(ids, val.GetID())
	}

	return ids
}

func GetUserActivationLink(token string) string {
	baseURL := os.Getenv("BASE_URL")
	return baseURL + "/user/activate?token=" + token
}

func GetUserPasswordResetLink(token string) string {
	baseURL := os.Getenv("BASE_URL")
	return baseURL + "/user/password/reset?token=" + token
}
