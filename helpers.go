package tgbot

import "strings"

func ClearSlice(s []string) []string {
	news := []string{}
	temp := ""
	for _, i := range s {
		temp = strings.TrimSpace(i)
		if temp != "" {
			news = append(news, temp)
		}
	}
	return news
}
