package repository

const (
	createNotificationQuery = `
	INSERT INTO notifications (
		id, title, content, target, creator_id, created_at, updated_at, active
	)VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	)
	RETURNING id, title, content, target, creator_id, created_at, updated_at, active
	`

	updateNotificationQuery = `
		UPADATE notifications
		SET 
			title = COALESCE(NULLIF($1, ''), title),
			content = COALESCE(NULLIF($2, ''), content),
			target = COALESCE(NULLIF($3, ''), target),
			modifier_id = COALESCE($4, modifier_id),
			version = version + 1,
			updated_at = $5
		WHERE id = $6
		RETURNING id, title, content, target, creator_id, created_at, updated_at, active
	`

	deleteNotificationQuery = `
	UPDATE notifications
	SET
		active = false,
		version = version + 1,
		modifier_id = $1,
		updated_at = $2
	WHERE id = $1 
	RETURNING id, title, content, target, creator_id, created_at, updated_at, active
	`

	getNotificationByIdQuery = `
	SELECT id, title, content, target, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE id = $1 AND active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM notifications
	WHERE active = true
	`

	searchByTitleQuery = `
	SELECT id, title, content, target, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE active = true AND title ILIKE '%' || $1 || '%'
	ORDER BY created_at DESC`

	findByTitleCount = `
	SELECT COUNT(*)
	FROM notifications
	WHERE title ILIKE '%' || $1 || '%' AND active = true
	`

	findByTitle = `
	SELECT *
	FROM notifications
	WHERE title = $1 AND active = true
	`

	getNotification = `
	SELECT id, title, content, target, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`
)
