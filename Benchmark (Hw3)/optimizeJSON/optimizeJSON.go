package optimizejson

// Тип данных для которого будет сгенерирован код анмаршалинга
// с помощью easyjson

// команда генерации кода
// easyjson -all optimizeJSON.go
type User struct {
	Browsers []string `json:"browsers"`
	// Company  string   `-"`
	// Country  string   `-"`
	Email string `json:"email"`
	// Job      string   `-`
	Name string `json:"name"`
	// Phone    string   `-`
}
