package search

import (
	"context"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
)

//Result describes one search result
type Result struct {
	//Фраза которую искали 
	Phrase  string 
	//Целиком вся строка в которой нашли вхождение 
	Line    string 
	//Номер позиции (начиая с 1)
	LineNum int64 
	//Номер позиции (начиая с 1)
	ColNum  int64 
}

//All ....
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	regex, _ := regexp.Compile(phrase)

	for _, file := range files {
		wg.Add(1)
		go func(found chan<- []Result, text string) {
			defer wg.Done()
			content := make([]byte, 0)
			content, rerr := ioutil.ReadFile(text)
			if rerr != nil {
				log.Printf("File: %v, Error: %v", text, rerr)
			}
			con := string(content)
			str := strings.Split(con, "\n")
			res := Result{Phrase: phrase}
			result := make([]Result, 0)

			
			for lineNum, line := range str {
				if line == "" {
					continue
				}
				if !strings.Contains(line, phrase) {
					continue
				}
				indexes := regex.FindAllStringIndex(line, -1)

				res.Line = line
				res.LineNum = int64(lineNum + 1)
				for _, index := range indexes {
					res.ColNum = int64(index[0] + 1)
					result = append(result, res)				
				}
			
			}
			if len(result) > 0 {
				found <- result
			}		
		
		}(ch, file)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

//Any....
func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	regex, _ := regexp.Compile(phrase)
	ctx,cancel:= context.WithCancel(ctx)
	for i, file := range files {
		wg.Add(1)
		go func(found chan<- []Result, text string, i int64) {
			defer wg.Done()
			content, rerr := ioutil.ReadFile(text)
			if rerr != nil {
				log.Printf("File: %v, Error: %v", text, rerr)
			}
			con := string(content)
			str := strings.Split(con, "\n")
			res := Result{Phrase: phrase}
			for lineNum, line := range str {
			select{
			case <-ctx.Done():
				return
			default:
			}
			if line == "" {
				continue
			}
			if !strings.Contains(line, phrase) {
				continue
			}
			indexes := regex.FindAllStringIndex(line, -1)

			res.Line = line
			res.LineNum = int64(lineNum + 1)
			for _, index := range indexes {
				res.ColNum = int64(index[0] + 1)
			}
			
		}
	
	
	}(file,i)
}
wg.Wait()
defer close(ch)
cancel()
return ch
}