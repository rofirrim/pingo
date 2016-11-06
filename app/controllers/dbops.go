package controllers

import "pinchito/app"
import "pinchito/app/models"
import "pinchito/app/helpers"
import "database/sql"
import "github.com/revel/revel"
import "github.com/go-sql-driver/mysql"
import "errors"

func GetUser(id int) (models.User, error) {
    row := app.DB.QueryRow("SELECT u.id, u.login, u.avatar FROM users u WHERE id = ?", id);
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
        plog.Data = nt.Time
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

    plog.Text, err = helpers.ProcessLogText(plog.Text)
    if err != nil {
        return models.Plog{}, err
    }
    return plog, nil
}

func GetPlog(id int) (models.Plog, error) {
    row := app.DB.QueryRow("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, 0.0 as nota FROM plogs p WHERE id = ?", id);
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

    rows, err := app.DB.Query("SELECT p.id, p.text, p.autor, p.protagonista, p.titol, p.data, 0.0 as nota FROM plogs p ORDER BY data DESC LIMIT ? OFFSET ?", app.LogsPerPage, offset);
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
