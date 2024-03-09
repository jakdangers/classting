package school

const createSchoolQuery = `INSERT INTO schools (user_id, name, region) values (?, ?, ?)`

const listSchoolQuery = `SELECT id, user_id, name, region FROM schools WHERE TRUE %s %s ORDER BY id ASC LIMIT 20`

const findSchoolByID = `SELECT id, user_id, name, region FROM schools WHERE id = ?`
