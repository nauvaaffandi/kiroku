package backup

import (
	"database/sql"
	"strconv"
	"time"
)

func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullFloat(nf sql.NullFloat64) string {
	if nf.Valid {
		return strconv.FormatFloat(nf.Float64, 'f', -1, 64)
	}
	return ""
}

func nullTimeStr(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}

func scanTransactions(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()
	var out [][]string
	for rows.Next() {
		var (
			id, typ, dateTime, currency          string
			amount                               int64
			fromAcc, toAcc, catID, labels, note  sql.NullString
			payee, paymentType, warranty, status sql.NullString
			placeName, attachments, templateID   sql.NullString
			placeLat, placeLng                   sql.NullFloat64
			createdAt, updatedAt                 sql.NullTime
		)
		if err := rows.Scan(&id, &typ, &dateTime, &amount, &currency,
			&fromAcc, &toAcc, &catID, &labels, &note,
			&payee, &paymentType, &warranty, &status, &placeName,
			&placeLat, &placeLng, &attachments, &templateID, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, []string{
			id, typ, dateTime, strconv.FormatInt(amount, 10), currency,
			nullStr(fromAcc), nullStr(toAcc), nullStr(catID), nullStr(labels), nullStr(note),
			nullStr(payee), nullStr(paymentType), nullStr(warranty), nullStr(status), nullStr(placeName),
			nullFloat(placeLat), nullFloat(placeLng), nullStr(attachments), nullStr(templateID),
			nullTimeStr(createdAt), nullTimeStr(updatedAt),
		})
	}
	return out, rows.Err()
}

func scanAccounts(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()
	var out [][]string
	for rows.Next() {
		var (
			id, name, currency, typ, icon, color string
			archived                             int64
			createdAt, updatedAt                 sql.NullTime
		)
		if err := rows.Scan(&id, &name, &currency, &typ, &icon, &color, &archived, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, []string{id, name, currency, typ, icon, color,
			strconv.FormatInt(archived, 10), nullTimeStr(createdAt), nullTimeStr(updatedAt)})
	}
	return out, rows.Err()
}

func scanCategories(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()
	var out [][]string
	for rows.Next() {
		var (
			id, name             string
			icon, color, kind    sql.NullString
			parentID             sql.NullString
			sortOrder            int64
			createdAt, updatedAt sql.NullTime
		)
		if err := rows.Scan(&id, &parentID, &name, &icon, &color, &kind, &sortOrder, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, []string{id, nullStr(parentID), name, nullStr(icon), nullStr(color), nullStr(kind),
			strconv.FormatInt(sortOrder, 10), nullTimeStr(createdAt), nullTimeStr(updatedAt)})
	}
	return out, rows.Err()
}
