package serveDB

import (
	"database/sql"
	"fmt"
	"log"
)

type vacQuery struct {
	ID int
	vacName   string
	keySkills string
	vacDesc   string
	salary    int
	jobType string
}

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
	fmt.Println(psqlInfo)

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
	

func GetFunction (qry string) string {
	scannedRow:=vacQuery{}
	var result string
	// Connect to DB
		db,err:=connectDB()		
		if err != nil {
			log.Fatal(err)
			return "Ошибка подключения к базе данных"
		}
		defer db.Close()
			rows, err := db.Query(qry)
		if err != nil {
			log.Fatal(err)
			return "Ошибка чтения базы данных"
		}
		defer rows.Close()
		
		// make queries to DB
		for i:=0; rows.Next(); i++ {
			err = rows.Scan(&scannedRow.ID, &scannedRow.vacName,&scannedRow.keySkills, &scannedRow.salary, &scannedRow.vacDesc, &scannedRow.jobType)
			if err != nil {
				log.Fatal(err)
				return "Ошибка чтения базы данных"
			}
			err = rows.Err()
		if err != nil {
			log.Fatal(err)
			return "Ошибка чтения базы данных"
		}
		result+=fmt.Sprintf("%v||%v||%v||%v||%v||%v\n",scannedRow.ID,scannedRow.vacName,scannedRow.keySkills,scannedRow.salary,scannedRow.vacDesc,scannedRow.jobType)
	}
	return  result
}

func PutFunction (vacName string, keySkills string, vacDesc string, salary int, jobCode int ) string {
	db,err:=connectDB()		
	if err != nil {
		log.Fatal(err)
		return "Ошибка подключения к базе данных"
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO vacancies (vacancy_name,key_skills, vacancy_desc ,  salary, job_type) VALUES($1, $2,$3,$4,$5)")
	if err != nil {
		log.Fatal(err)
		return "Ошибка добавления значений в базу"
	}
	_, err = stmt.Exec(vacName, keySkills, vacDesc, salary, jobCode)
	if err != nil {
		log.Fatal(err)
		return "Ошибка добавления значений в базу"
		}
	stmt.Close()
	return "Запись добавлена"
}