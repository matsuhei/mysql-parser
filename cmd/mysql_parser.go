/*
Copyright © 2022 matsuhei

*/
package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"mysql-parser/pkg"
)

// mysqlParserCmd represents the mysqlParser command
var mysqlParserCmd = &cobra.Command{
	Use:   "mysql_parser",
	Short: "mysql parser command for creating table command.",
	Long: `💻 mysql parser command for creating table command. 💻 you can get table_name & column format that you want.
args 0 table_name`,
	Run: func(cmd *cobra.Command, args []string) {
		showMysqlTables(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(mysqlParserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mysqlParserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mysqlParserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func showMysqlTables(cmd *cobra.Command, args []string) {
	fmt.Printf("command exexute : %s\n", cmd.Use)

	db, err := getDB()
	if err != nil {
		fmt.Printf("failed... err: %s\n", err)
		return
	}

	exec, err := db.Query(fmt.Sprintf("SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"%s\" order by ordinal_position asc;", args[0]))

	if err != nil {
		fmt.Printf("failed... err: %s\n", err)
		return
	}

	tables := pkg.Table{Name: args[0], Columns: []pkg.Column{}}
	defer func(exec *sql.Rows) {
		err := exec.Close()
		if err != nil {
			fmt.Printf("failed... err: %s\n", err.Error())
		}
	}(exec)
	for exec.Next() {
		var cn, ck, dt, in, cc interface{}
		if err := exec.Scan(&cn, &ck, &dt, &in, &cc); err != nil {
			log.Fatal(err)
			break
		}

		tables.Columns = append(tables.Columns, pkg.Column{
			ColumnName:    fmt.Sprintf("%s", cn),
			ColumnKey:     fmt.Sprintf("%s", ck),
			DataType:      fmt.Sprintf("%s", dt),
			IsNullable:    fmt.Sprintf("%s", in),
			ColumnComment: fmt.Sprintf("%s", cc),
		})
	}

	ft := fmt.Sprintf("%s", viper.Get("format-type"))
	switch ft {
	case "figjam":
		figjam := pkg.TableForFigjamDatabase{TableName: args[0], Color: "#1e1e1e", Columns: []pkg.ColumnForFigjam{}}
		for _, v := range tables.Columns {
			figjam.Columns = append(figjam.Columns, pkg.ColumnForFigjam{
				Name:     v.ColumnName,
				Type:     v.DataType,
				KeyType:  pkg.GetKeyType(v.ColumnKey),
				Nullable: pkg.GetNullable(v.IsNullable),
			})
		}

		if jStr, err := json.Marshal(figjam); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%s\n", jStr)
		}

	default:
		fmt.Println("format type not found.")
	}

}

func getDB() (*sql.DB, error) {
	ds := viper.Get("datasource")
	dbName := viper.Get("dbname")
	user := viper.Get("db-user")
	password := viper.Get("password")
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, ds, dbName))
}
