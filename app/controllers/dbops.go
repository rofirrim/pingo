package controllers

import "pinchito/app"
import "pinchito/app/models"
import "database/sql"
import "github.com/revel/revel"
import "github.com/go-sql-driver/mysql"
import "errors"

func GetUser(id int) (models.User, error) {
    row := app.DB.QueryRow("SELECT login, avatar FROM users WHERE id = ?", id);
    var autor models.User
    err := row.Scan(&autor.Login, &autor.Avatar)
    if err != nil {
        return models.User{}, nil
    }
    return autor, nil
}

func makePlogFromRows(rows *sql.Rows) (models.Plog, error) {
    var plog models.Plog
    var autor, protagonista int
    var nt mysql.NullTime
    err := rows.Scan(&plog.Text, &autor, &protagonista, &plog.Titol, &nt)
    return makePlog(err, plog, autor, protagonista, nt)
}

func makePlogFromRow(row *sql.Row) (models.Plog, error) {
    var plog models.Plog
    var autor, protagonista int
    var nt mysql.NullTime
    err := row.Scan(&plog.Text, &autor, &protagonista, &plog.Titol, &nt)
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
    plog.Protagonista, err = GetUser(protagonista)
    if err != nil {
		revel.ERROR.Println("Error getting protagonista", err)
        return models.Plog{}, err
    }
    plog.Autor, err = GetUser(autor)
    if err != nil {
		revel.ERROR.Println("Error getting user", err)
        return models.Plog{}, err
    }
    return plog, nil
}

func GetPlog(id int) (models.Plog, error) {
    row := app.DB.QueryRow("SELECT text, autor, protagonista, titol, data FROM plogs WHERE id = ?", id);
    return makePlogFromRow(row)
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

    rows, err := app.DB.Query("SELECT text, autor, protagonista, titol, data FROM plogs ORDER BY data DESC LIMIT ? OFFSET ?", app.LogsPerPage, offset);
    if err != nil {
		revel.ERROR.Println("Error retrieving rows of plogs", err)
        return []models.Plog{}, err
    }

    var plogs []models.Plog
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
