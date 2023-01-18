package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	_ "github.com/lib/pq"
)
var cred dbCred
var mainMenu string = "\n*****************************\n - Если хотите посмотреть всю таблицу вакансий, наберите \"посмотреть\", \n - Если хотите найти вакансию по названию наберите \"найти\"\n - Если хотите добавить строку - наберите \"добавить\", \n - Если хотите выйти из программы, наберите \"выход\"\n*****************************\n"
var searchQry string = " SELECT vacancies.id, vacancy_name, key_skills, salary, vacancy_desc, job_types.job_type FROM vacancies JOIN job_types ON vacancies.job_type = job_types.id WHERE vacancy_name ILIKE '%"
var db *sql.DB

type dbCred struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

type vacQuery struct {
	ID int
	vacName   string
	keySkills string
	vacDesc   string
	salary    int
	jobCode   int
	jobType string
}

func connectDB(cred dbCred) (*sql.DB, error) {

	// Connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cred.host, cred.port, cred.user, cred.password, cred.dbName)
	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func searchDialog() (string, bool) {
	fmt.Println("Введите название вакансии частично или полностью, либо наберите \"назад\" для выхода в предыдущее меню")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "назад" {
				return "", false
			}
			searchKey := scanner.Text()
			return searchKey, true
		}
	}
}
func loadVacs (qry string) ([] vacQuery, error) {
	result:=[] vacQuery{}
	rows, err := db.Query(qry)
	if err != nil {
		return  result,err
	}
	defer rows.Close()
	for i:=0; rows.Next(); i++ {
		result=append(result, vacQuery{})
		err = rows.Scan(&result[i].ID, &result[i].vacName,&result[i].keySkills, &result[i].salary, &result[i].vacDesc, &result[i].jobType)
		if err != nil {
			return  result,err
		}
		err = rows.Err()
	if err != nil {
		return result, err
	}
}
return  result,err
}

func showVacs(resSlice []vacQuery) error {

	var counter int
	var err error
	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for i,line:= range resSlice {
counter++
			if i == 0 {
			
			_,err= fmt.Fprintln(w, "\tID\tНазвание вакансии\tКлючевые навыки\tОписание вакансии\tЗарплата\tТип работы\t")
			if err != nil {
				return err
			}
			_,err= fmt.Fprintln(w, "\t--\t-----------------\t------------------------------------------\t-----------------------------------------------------------------\t--------\t----------\t")
			if err != nil {
				return err
			}
		}
		

		_,err= fmt.Fprintf(w, "\t%v\t%v\t%v\t%v\t%v\t%v\t\n",line.ID, line.vacName, line.keySkills, line.vacDesc, line.salary, line.jobType)
		if err != nil {
			return err
		}
		_,err= fmt.Fprintln(w, "\t--\t-----------------\t------------------------------------------\t-----------------------------------------------------------------\t--------\t----------\t")
		if err != nil {
			return err
		}
	}
	w.Flush()
	if counter == 0 {
		fmt.Println("\nПохоже, по такому запросу в базе ничего не нашлось. Попробуйте изменить запрос")
		fmt.Println("----------------------------------------")
		return err
	}
	
	return err
}
func insertDialog() (vacQuery, bool) {
	var result vacQuery
	var err error
	fmt.Println("введите соответствующие значения строк, разделяя их знаком \"/\": ")
	fmt.Println("название вакансии, ключевые навыки, описание вакансии, зарплата, и код типа работы: 1 для работы в офисе, 2 для удаленной работы и 3 для гибридной формы работы")
	fmt.Println("Например: \"Охранник/Решать конфликтные ситуации, обращаться с оружием, разгадывать сканворды/Человек, который следит за порядком в офисном здании/50000/1\"")
	fmt.Println("или наберите \"назад\" для выхода в предыдущее меню")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "назад" {
				return vacQuery{}, false
			}

			queryString := strings.Split(scanner.Text(), "/")
			if len(queryString) != 5 {
				fmt.Println("Неверное количество аргументов, повторите ввод")
				continue
			}
			result.vacName = queryString[0]
			result.keySkills = queryString[1]
			result.vacDesc = queryString[2]
			result.salary, err = strconv.Atoi(queryString[3])
			if err != nil {
				fmt.Println("Ошибка ввода данных в поле \"Зарплата\", повторите ввод")
				continue
			}
			result.jobCode, err = strconv.Atoi(queryString[4])
			if err != nil {
				fmt.Println("Ошибка ввода данных в поле \"код типа работы\", повторите ввод")
				continue
			}
			if result.jobCode > 3 || result.jobCode < 1 {
				fmt.Println("Код работы может быть только следующих значений: 1 для работы в офисе, 2 для удаленной работы и 3 для гибридной формы работы. Ввод других значений не допускается")
				continue
			}
		}

		return result, true

	}
}
func insert(q vacQuery) error {
	
	stmt, err := db.Prepare("INSERT INTO vacancies (vacancy_name,key_skills, vacancy_desc ,  salary, job_type) VALUES($1, $2,$3,$4,$5)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(q.vacName, q.keySkills, q.vacDesc, q.salary, q.jobCode)
	if err != nil {
		return err
	}
	stmt.Close()
	fmt.Println("Запись добавлена")
	return nil
}

func mainDialog() error {
	fmt.Println(mainMenu)
	scanner := bufio.NewScanner(os.Stdin)
OuterLoop:
	for {
		// Print a prompt
		fmt.Print("> ")

		// Scan for input
		if scanner.Scan() {
			switch {
			case scanner.Text() == "посмотреть":
				res,err := loadVacs(searchQry+"%'")
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				err=showVacs(res)
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "найти":
				keyWord, proceed := searchDialog()
				if proceed {
					res,err := loadVacs(searchQry+keyWord+"%'")
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				err=showVacs(res)
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "добавить":
				query, proceed := insertDialog()
				if proceed {
					err := insert(query)
					if err != nil {
						fmt.Println("Ошибка внесения данных в таблицу, попробуйте еще раз", err)
						return err
					}
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "выход":
				fmt.Println("Всего хорошего!")
				break OuterLoop
			default:
				fmt.Println("Неверно введена команда, попробуйте еще раз")
			}
		}
	}
	return nil
}
func main() {
	var err error	
	flag.IntVar(&cred.port, "port", 5432, "Port for DB connection")
	flag.StringVar(&cred.host, "h", "localhost", "DB host IP")
	flag.StringVar(&cred.password, "p", "my_awesome_password", "DB connection password")
	flag.StringVar(&cred.user, "u", "postgres", "DB connection user name")
	flag.StringVar(&cred.dbName, "n", "vacancies", "DB name")
flag.Parse()
	db, err = connectDB(cred)

	if err != nil {
		fmt.Println("Ошибка подключения базы данных", err)
		return
	}
	defer db.Close()
	err=mainDialog()
	if err!=nil{
		fmt.Println(err)
	}

}
