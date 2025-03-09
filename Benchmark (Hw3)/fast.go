package main

import (
	"bufio"
	"fmt"
	optimizejson "hw3/optimizeJSON"
	"io"
	"os"
	"strings"

	"github.com/mailru/easyjson"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	/*
		!!! !!! !!!
		обратите внимание - в задании обязательно нужен отчет
		делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировалм их
		так же обратите внимание на команду в параметром -http
		перечитайте еще раз задание
		!!! !!! !!!
	*/
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	seenBrowsers := map[string]bool{}
	uniqueBrowsers := 0

	fmt.Fprintln(out, "found users:")
	writer := bufio.NewWriter(out)
	for i := 0; scanner.Scan(); i++ {
		user := optimizejson.User{}
		err := easyjson.Unmarshal(scanner.Bytes(), &user)
		// err := json.Unmarshal([]byte(scanner.Text()), &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {
			if ok := strings.Contains(browserRaw, "Android"); ok {
				isAndroid = true
				notSeenBefore := true
				if seenBrowsers[browserRaw] {
					notSeenBefore = false
				}
				if notSeenBefore {
					seenBrowsers[browserRaw] = true
					uniqueBrowsers++
				}
			}

			if ok := strings.Contains(browserRaw, "MSIE"); ok {
				isMSIE = true
				notSeenBefore := true
				if seenBrowsers[browserRaw] {
					notSeenBefore = false
				}
				if notSeenBefore {
					seenBrowsers[browserRaw] = true
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		writer.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
		if i%100 == 0 {
			writer.Flush()
		}
	}
	writer.Flush()
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
