package scheduler

var URLs = []string{}

func PopUrl() string {
	length := len(URLs)
	if length < 1 {
		return ""
	}

	url := URLs[length-1]
	URLs = URLs[:length-1]
	return url
}

func AppendUrl(url string) {
	URLs = append(URLs, url)
}

func CountUrl() int {
	return len(URLs)
}
