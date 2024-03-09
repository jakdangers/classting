package news

const createNewsQuery = `INSERT INTO news (school_id, user_id, title) VALUES (?, ?, ?)`

const listNewsQuery = `SELECT id, create_date, update_date, delete_date, school_id, user_id, title FROM news WHERE delete_date IS NULL %s %s %s ORDER BY id DESC LIMIT 10`

const findNewsByIDQuery = `SELECT id, create_date, update_date, delete_date, school_id, user_id, title FROM news WHERE id = ?`

const deleteNewsQuery = `UPDATE news SET delete_date = ? WHERE id = ?`
