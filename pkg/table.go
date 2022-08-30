package pkg

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	ColumnName    string
	ColumnKey     string
	DataType      string
	IsNullable    string
	ColumnComment string
}

type TableForFigjamDatabase struct {
	TableName string            `json:"tableName"`
	Color     string            `json:"color"`
	Columns   []ColumnForFigjam `json:"columns"`
}

type ColumnForFigjam struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	KeyType  KeyType `json:"keyType"`
	Nullable bool    `json:"nullable"`
}

type KeyType string

const (
	Primary KeyType = "primary"
	Foreign KeyType = "foreign" // TODO: 判別方法がないため一旦利用なし
	Unique  KeyType = "unique"
	Index   KeyType = "index"
	None    KeyType = "normal"
)

func GetKeyType(key string) KeyType {
	switch key {
	case "PRI":
		return Primary
	case "UNI":
		return Unique
	case "MUL":
		return Index
	default:
		return None
	}
}

func GetNullable(state string) bool {
	return state == "YES"
}
