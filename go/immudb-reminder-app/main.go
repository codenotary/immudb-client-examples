package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	_ "github.com/codenotary/immudb/pkg/stdlib"

)

type cfg struct {
	IpAddr   string
	Port     string
	Username string
	Password string
	DBName   string
}

func parseConfig() (c cfg) {
	flag.StringVar(&c.IpAddr, "addr", "localhost", "IP address of immudb server")
	flag.StringVar(&c.Port, "port", "3322", "Port number of immudb server")
	flag.StringVar(&c.Username, "user", "immudb", "Username for authenticating to immudb")
	flag.StringVar(&c.Password, "pass", "immudb", "Password for authenticating to immudb")
	flag.StringVar(&c.DBName, "db", "defaultdb", "Name of the database to use")
	flag.Parse()
	return
}

func main() {
	c := parseConfig()
        connstring := fmt.Sprintf("immudb://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.IpAddr, c.Port, c.DBName )
	db, err := sql.Open("immudb", connstring)
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		defer db.Close()
		return
	}

	_, err = db.ExecContext(context.TODO(), "CREATE TABLE IF NOT EXISTS REMINDERS(id INTEGER AUTO_INCREMENT,title VARCHAR, description VARCHAR, alias VARCHAR, PRIMARY KEY id);")
	if err != nil {
		fmt.Println(err)
	}
	for {
		fmt.Println("-> Welcome to Reminders Console App, built using Golang and Immudb")
		fmt.Println("-> Select a numeric option; \n [1] Create a new Reminder \n [2] Get a reminder \n [3] Delete a reminder \n [4] Quit")
		quitBool := ""
		consoleReader := bufio.NewScanner(os.Stdin)
		consoleReader.Scan()
		userChoice := consoleReader.Text()

		switch userChoice {
		case "1":
			{
				var (
					titleInput,
					descriptionInput,
					aliasInput string
				)

				fmt.Println("You are about to create a new reminder. Please provide the following details:")

				fmt.Println("-> What is the title of your reminder?")
				consoleReader.Scan()
				titleInput = consoleReader.Text()

				fmt.Println("-> What is the description of your reminder?")
				consoleReader.Scan()
				descriptionInput = consoleReader.Text()

				fmt.Println("-> What is an alias of your reminder? [ An alias will be used to retrieve your reminder ]")
				consoleReader.Scan()
				aliasInput = consoleReader.Text()

				err := createReminder(titleInput, descriptionInput, aliasInput, db)
				if err != nil {
					return
				}
				break
			}
		case "2":
			{
				fmt.Println("-> Please provide an alias for your reminder:")
				consoleReader.Scan()
				aliasInput := consoleReader.Text()

				data, getErr := retrieveReminder(aliasInput, db)
				if getErr != nil {
					fmt.Println(getErr)
				}

				fmt.Println(data)
				break
			}
		case "3":
			{
				fmt.Println("-> Please provide the alias for the reminder you want to delete:")

				consoleReader.Scan()
				deleteAlias := consoleReader.Text()

				getErr := deleteReminder(deleteAlias, db)
				if getErr != nil {
					fmt.Println(getErr)
				}
				break
			}
		case "4":
			{
				quitBool = "true"
				break
			}
		default:
			fmt.Printf("-> Option: %v is not a valid numeric option. Try 1 , 2 , 3", userChoice)
		}
		if quitBool == "true" {
			break
		} else {
			quitBool = ""
		}
	}
}

func createReminder(titleInput, aliasInput, descriptionInput string, database *sql.DB) error {

	err := database.PingContext(context.Background())
	if err != nil {
		fmt.Printf("Error checking db connection: %v", err)
	}

	queryStatement := fmt.Sprintf("INSERT INTO REMINDERS ( title, description, alias ) VALUES ( '%v', '%v', '%v' );", titleInput, aliasInput, descriptionInput)

	_, queryErr := database.ExecContext(context.TODO(), queryStatement)

	if queryErr != nil {
		fmt.Printf("Query err: %v", queryErr)
	}

	return nil
}

func retrieveReminder(alias string, database *sql.DB) (string, error) {
	sqlStatement := fmt.Sprintf("SELECT title FROM REMINDERS WHERE alias='%v';", alias)

	data, err := database.QueryContext(context.TODO(), sqlStatement)
	if err != nil {
		fmt.Printf("Query err: %v", err)
	}

	title := ""
	for data.Next() {
		//var description, alias, title  string
		//var isCompleted int

		nErr := data.Scan(&title)
		if nErr != nil {
			fmt.Printf("Error: %v", nErr)
		}

	}
	if title == "" {
		title = "no entry found"
	}
	return title, nil
}

func deleteReminder(alias string, database *sql.DB) error {

	err := database.PingContext(context.Background())
	if err != nil {
		fmt.Printf("Error checking db connection: %v", err)
	}

	queryStatement := fmt.Sprintf("DELETE FROM reminders WHERE alias='%v';", alias)

	_, queryErr := database.ExecContext(context.TODO(), queryStatement)

	if queryErr != nil {
		fmt.Printf("Query err: %v", queryErr)
	}

	return nil
}

