package subscription

const createSubscriptionQuery = `INSERT INTO subscriptions (school_id, user_id) VALUES (?, ?)`

const listSubscriptionSchoolsQuery = `SELECT subscriptions.id, subscriptions.create_date, subscriptions.update_date, subscriptions.school_id, schools.name, schools.region FROM schools JOIN subscriptions ON schools.id = subscriptions.school_id  WHERE subscriptions.user_id = ? %s ORDER BY id DESC LIMIT 20`

const findSubscriptionByUserIDAndSchoolIDQuery = `SELECT id, create_date, update_date, school_id, user_id FROM subscriptions WHERE user_id = ? AND school_id = ?`

const deleteSubscriptionQuery = `DELETE FROM subscriptions WHERE id = ?`
