package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Person struct {
	XMLName       xml.Name `xml:"row"`
	Id            int      `xml:"id"`
	Guid          string   `xml:"guid"`
	IsActive      bool     `xml:"isActive"`
	Balance       string   `xml:"balance"`
	Picture       string   `xml:"picture"`
	Age           int      `xml:"age"`
	EyeColor      string   `xml:"eyeColor"`
	First_name    string   `xml:"first_name"`
	Last_name     string   `xml:"last_name"`
	Gender        string   `xml:"gender"`
	Company       string   `xml:"company"`
	Email         string   `xml:"email"`
	Phone         string   `xml:"phone"`
	Address       string   `xml:"address"`
	About         string   `xml:"about"`
	Registered    string   `xml:"registered"`
	FavoriteFruit string   `xml:"favoriteFruit"`
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	Rows    []Person `xml:"row"`
}

// SearchServer принимает GET-параметры:
// * `query` - что искать. Ищем по полям записи `Name` и `About` просто подстроку, без регулярок.
// `Name` - это first_name + last_name из xml (вам надо руками пройтись в цикле по записям и сделать такой, автоматом нельзя).
// Если поле пустое - то возвращаем все записи (поиск пустой подстроки всегда возвращает true), т.е. делаем только логику сортировки

// * `order_field` - по какому полю сортировать. Работает по полям `Id`, `Age`, `Name`, если пустой - то сортируем по `Name`,
// если что-то другое - SearchServer ругается ошибкой.

// * `order_by` - направление сортировки (как есть, по убыванию, по возрастанию), в client.go есть соответствующие константы

// * `limit` - сколько записей вернуть

// * `offset` - начиня с какой записи вернуть (сколько пропустить с начала) - нужно для огранизации постраничной навигации

func SearchServer(query string, order_field string, orderBy int, limit int, offset int) (find []User, err error) {
	const filename = `dataset.xml`
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("error open file: %s", err)
		return
	}
	res := Root{}
	err = xml.Unmarshal(data, &res)
	if err != nil {
		log.Printf("error unmarshal xml: %s", err)
		return
	}
	rgx, err := regexp.Compile(fmt.Sprintf(".*%s.*", query))
	// Ищем то, что запросил пользователь
	for _, row := range res.Rows {
		if rgx.MatchString(row.First_name) || rgx.MatchString(row.Last_name) || rgx.MatchString(row.About) {
			find = append(find, User{
				Id:     row.Id,
				Name:   row.First_name + " " + row.Last_name,
				Age:    row.Age,
				About:  row.About,
				Gender: row.Gender})
		}
	}
	// Сортируем данные (работаем по полям `Id`, `Age`, `Name`, если пустой - то сортируем по `Name`,
	// если что-то другое - SearchServer ругается ошибкой)
	switch order_field {
	case "Id":
		switch orderBy {
		case OrderByAsc:
			sort.Slice(find, func(i, j int) bool { return find[i].Id < find[j].Id })
		case OrderByDesc:
			sort.Slice(find, func(i, j int) bool { return find[i].Id > find[j].Id })
		case OrderByAsIs:
		default:
			return find, fmt.Errorf("invalid order setting %d", orderBy)
		}
	case "Age":
		switch orderBy {
		case OrderByAsc:
			sort.Slice(find, func(i, j int) bool { return find[i].Age < find[j].Age })
		case OrderByDesc:
			sort.Slice(find, func(i, j int) bool { return find[i].Age > find[j].Age })
		case OrderByAsIs:
		default:
			return find, fmt.Errorf("invalid order setting %d", orderBy)
		}
	case "Name", "":
		switch orderBy {
		case OrderByAsc:
			sort.Slice(find, func(i, j int) bool {
				return find[i].Name < find[j].Name
			})
		case OrderByDesc:
			sort.Slice(find, func(i, j int) bool {
				return find[i].Name > find[j].Name
			})
		case OrderByAsIs:
		default:
			return find, fmt.Errorf("invalid order setting %d", orderBy)
		}
	default:
		return find, fmt.Errorf("unrecognize field to order: %w", err)
	}
	// Возвращаем итоовый слайс исходя из параметров limit, offset
	return find[min(offset, len(find)):min(offset+limit, len(find))], nil
}

func SearchServerHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	order_field := r.FormValue("order_field")
	orderBy, err := strconv.Atoi(r.FormValue("order_by"))
	if err != nil {
		log.Println("Error parsing Order By", err)
	}
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		log.Println("Error parsing limit", err)
	}
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		log.Println("Error parsing offset", err)
	}
	// fmt.Println(query, order_field, orderBy, limit, offset)
	res, err := SearchServer(query, order_field, orderBy, limit, offset)
	if err != nil {
		log.Println("Error SearchServer", err)
	}
	toSend, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshaling data: %v", err)
	}
	_, err = w.Write(toSend)
	if err != nil {
		log.Printf("error writing data: %v", err)
	}
}

func SearchServerErrors(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	switch query {
	case "StatusUnauthorized":
		w.WriteHeader(http.StatusUnauthorized)
	case "StatusInternalServerError":
		w.WriteHeader(http.StatusInternalServerError)
	case "StatusBadRequest":
		w.WriteHeader(http.StatusBadRequest)
	case "StatusBadRequest_ErrorOrder":
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "ErrorBadOrderField"}`)
	case "StatusBadRequest_Unknown":
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Unknown"}`)
	case "StatusOk_ErrorJSON":
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"error": "Unknown"`)
	case "Timeout":
		time.Sleep(10 * time.Second)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

// Пример запроса
// http://localhost:8080/?query=Boyd&order_field=Id&order_by=1&limit=10&offset=0
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", SearchServerHandler)
	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	server.ListenAndServe()
}
