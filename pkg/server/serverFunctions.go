package serveDB

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/mrkovshik/grpc_vacancy_database/grpc/proto"
)

type dbCred struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}
var cred  = dbCred {
	host: "localhost",
	port:     5432,
	user:     "postgres",
	password: "17pasHres19!",
	dbName:   "vacancies",
	}

func connectDB() (*sql.DB, error) {
	// Connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cred.host, cred.port, cred.user, cred.password, cred.dbName)
		db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
	

func ReadFunction (qry string) [] *proto.VacancyStruct {
	var searchQry string = " SELECT vacancies.id, vacancy_name, key_skills, salary, vacancy_desc, job_types.job_type FROM vacancies JOIN job_types ON vacancies.job_type = job_types.id WHERE vacancy_name ILIKE '%"

	db,err:=connectDB()
	if err != nil {
		log.Fatal(err)
		}
	result:=[] *proto.VacancyStruct{}
	rows, err := db.Query(searchQry+qry+"%'")
	if err != nil {
		log.Fatal(err)
		}
	defer rows.Close()
	for i:=0; rows.Next(); i++ {
		result=append(result, &proto.VacancyStruct{})
		err = rows.Scan(&result[i].ID, &result[i].VacName,&result[i].KeySkills, &result[i].Salary, &result[i].VacDesc, &result[i].JobType)
		if err != nil {
			log.Fatal(err)
			}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
			}
	}
	return result
}

func InsertFunction (insertRow *proto.VacancyStruct) string {
	db,err:=connectDB()		
	if err != nil {
		errMsg:=fmt.Sprintf("Ошибка подключения к базе данных: %v", err)
		log.Fatal(err)
		return errMsg
	}
		defer db.Close()
	stmt, err := db.Prepare("INSERT INTO vacancies (vacancy_name,key_skills, vacancy_desc ,  salary, job_type) VALUES($1, $2,$3,$4,$5)")
	if err != nil {
		errMsg:=fmt.Sprintf("Ошибка добавления значений в базу: %v", err)
		log.Fatal(err)
		return errMsg
	}
	_, err = stmt.Exec(insertRow.VacName, insertRow.KeySkills, insertRow.VacDesc, insertRow.Salary, insertRow.JobCode)
	if err != nil {
		errMsg:=fmt.Sprintf("Ошибка добавления значений в базу: %v", err)
		log.Fatal(err)
		return errMsg
	}
	stmt.Close()
	return "Запись добавлена"
}

func DeleteFunction (deleteTarget int32) string {
	db,err:=connectDB()		
	if err != nil {
		log.Fatal(err)
		return "Ошибка подключения к базе данных"
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM products WHERE price = $1;")
	if err != nil {
		errMsg:=fmt.Sprintf("Ошибка добавления значений в базу: %v", err)
		log.Fatal(err)
		return errMsg
	}
	_, err = stmt.Exec(deleteTarget)
	if err != nil {
		errMsg:=fmt.Sprintf("Ошибка добавления значений в базу: %v", err)
		log.Fatal(err)
		return errMsg
	}
	stmt.Close()
	return "Запись удалена"
	

}