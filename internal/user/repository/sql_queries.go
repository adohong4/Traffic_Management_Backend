package repository

const (
	getUser = `
	INSERT INTO notifications (
		id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	)VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)
	RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`
)
