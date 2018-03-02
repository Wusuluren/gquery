package gquery

type StrStrMap map[string]string

func (ss StrStrMap) Get(key string) string {
	defaultValue := ""
	if value, ok := ss[key]; ok {
		return value
	}
	return defaultValue
}
