package controllers

import "pinchito/app"
import "pinchito/app/models"
import "database/sql"
import "github.com/revel/revel"
import "github.com/go-sql-driver/mysql"
import "fmt"
import "strings"
import "errors"

func GetUser(id int) (models.User, error) {
	row := app.DB.QueryRow("SELECT u.id, u.login, u.avatar FROM users u WHERE id = ?", id)
	var autor models.User
	err := row.Scan(&autor.Id, &autor.Login, &autor.Avatar)
	if err != nil {
		return models.User{}, err
	}
	return autor, nil
}

func makePlogFromRows(rows *sql.Rows) (models.Plog, error) {
	var plog models.Plog
	var autor, protagonista int
	var nt mysql.NullTime
	err := rows.Scan(&plog.Id, &plog.Text, &autor, &protagonista, &plog.Titol, &nt, &plog.Nota)
	return makePlog(err, plog, autor, protagonista, nt)
}

func makePlogFromRow(row *sql.Row) (models.Plog, error) {
	var plog models.Plog
	var autor, protagonista int
	var nt mysql.NullTime
	err := row.Scan(&plog.Id, &plog.Text, &autor, &protagonista, &plog.Titol, &nt, &plog.Nota)
	return makePlog(err, plog, autor, protagonista, nt)
}

func makePlog(err error, plog models.Plog, autor int, protagonista int, nt mysql.NullTime) (models.Plog, error) {
	if err != nil {
		revel.ERROR.Println("Error while scanning row to make plog", err)
		return models.Plog{}, err
	}

	if nt.Valid {
		plog.Dia = nt.Time.Format("02/01/2006")
		plog.Hora = nt.Time.Format("15:04")
	} else {
		plog.Dia = "dia desconegut"
		plog.Hora = "tantes"
	}

	// Normalise log title if it does not have any
	if strings.TrimSpace(plog.Titol) == "" {
		plog.Titol = fmt.Sprintf("Log %d", plog.Id)
	}

	if protagonista != 0 {
		plog.Protagonista, err = GetUser(protagonista)
		if err != nil {
			revel.ERROR.Println("Error getting protagonista", err)
			return models.Plog{}, err
		}
	}

	plog.Autor, err = GetUser(autor)
	if err != nil {
		revel.ERROR.Println("Error getting user", err)
		return models.Plog{}, err
	}

	return plog, nil
}

func GetPlog(id int) (models.Plog, error) {
	row := app.DB.QueryRow("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, 0.0 as nota FROM plogs p WHERE id = ?", id)
	return makePlogFromRow(row)
}

func retrievePlogs(rows *sql.Rows) ([]models.Plog, error) {
	var plogs []models.Plog
	var err error
	defer rows.Close()
	for rows.Next() {
		var plog models.Plog
		plog, err = makePlogFromRows(rows)
		if err != nil {
			revel.ERROR.Println("Error making plog", err)
			return []models.Plog{}, err
		}
		plogs = append(plogs, plog)
	}
	err = rows.Err() // get any error encountered during iteration
	if err != nil {
		revel.ERROR.Println("Error while retrieving plogs", err)
		return []models.Plog{}, err
	}
	return plogs, nil
}

func GetPlogBunch(page int, numplogs *int) ([]models.Plog, error) {
	if page <= 0 {
		return []models.Plog{}, errors.New("page cannot be zero or lower than zero")
	}
	var err error

	offset := (page - 1) * app.LogsPerPage

	err = app.DB.QueryRow("SELECT COUNT(*) FROM plogs").Scan(numplogs)
	if err != nil {
		revel.ERROR.Println("Error retrieving number of plogs", err)
		return []models.Plog{}, err
	}

	rows, err := app.DB.Query("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, 0.0 as nota FROM plogs p ORDER BY p.data DESC LIMIT ? OFFSET ?", app.LogsPerPage, offset)
	if err != nil {
		revel.ERROR.Println("Error retrieving rows of plogs", err)
		return []models.Plog{}, err
	}

	return retrievePlogs(rows)
}

func GetTop20Plogs() ([]models.Plog, error) {
	rows, err := app.DB.Query("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, AVG(v.nota) AS rank FROM plogs p INNER JOIN votes v ON p.id = v.log_id GROUP BY p.id ORDER BY rank DESC LIMIT 20")
	if err != nil {
		revel.ERROR.Println("Error retrieving top 20 plogs", err)
		return []models.Plog{}, err
	}

	return retrievePlogs(rows)
}

func GetRandomPlogs() ([]models.Plog, error) {
	rows, err := app.DB.Query("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, 0.0 AS rank FROM plogs p ORDER BY RAND() LIMIT 3")
	if err != nil {
		revel.ERROR.Println("Error retrieving 3 random plogs", err)
		return []models.Plog{}, err
	}

	return retrievePlogs(rows)
}

func GetBlobAvatar(id int) ([]byte, error) {
	user, err := GetUser(id)
	return user.Avatar, err
}

func GetRandomCookie() (models.Cookie, error) {
	row := app.DB.QueryRow("SELECT q.text, q.autor FROM quotes q ORDER BY RAND() LIMIT 1")
	var result models.Cookie
	var err error
	err = row.Scan(&result.Text, &result.Autor)
	return result, err
}

func SearchPlogs(keywords []string, page int, numplogs *int) ([]models.Plog, error) {

	var condText []string = make([]string, len(keywords))
	var condTitol []string = make([]string, len(keywords))
	for i := range keywords {
		condText[i] = "p.text REGEXP ?"
		condTitol[i] = "p.titol REGEXP ?"
	}

	// Sometimes go is weird
	concatRange := append(keywords, keywords...)
	ifaceArray := make([]interface{}, len(concatRange))
	for i, v := range concatRange {
		ifaceArray[i] = fmt.Sprintf("[[:<:]]%v[[:>:]]", v)
	}
	offset := (page - 1) * app.LogsPerPage

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM plogs p WHERE ( %v ) OR ( %v )",
		strings.Join(condText, " AND "),
		strings.Join(condTitol, " AND "))
	countRow := app.DB.QueryRow(countQuery, ifaceArray...)
	err := countRow.Scan(numplogs)
	if err != nil {
		revel.ERROR.Println("Error retrieving rows of plogs after search count", err)
		return []models.Plog{}, err
	}

	ifaceArray = append(ifaceArray, app.LogsPerPage)
	ifaceArray = append(ifaceArray, offset)

	query := fmt.Sprintf("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, 0.0 AS rank FROM plogs p WHERE ( %v ) OR ( %v ) ORDER BY p.data DESC LIMIT ? OFFSET ?",
		strings.Join(condText, " AND "),
		strings.Join(condTitol, " AND "))

	revel.INFO.Println("Error retrieving rows of plogs after search", query)
	rows, err := app.DB.Query(query, ifaceArray...)
	if err != nil {
		revel.ERROR.Println("Error retrieving rows of plogs after search", err)
		return []models.Plog{}, err
	}

	return retrievePlogs(rows)
}
