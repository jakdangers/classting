package user

const createUserQuery = `INSERT INTO users (user_name, password, user_type) VALUES (?, ?, ?)`

const findUserByUserNameQuery = `SELECT id, user_name, password, user_type FROM users WHERE user_name = ?`
