package repository

const (
	createNotificationQuery = `
	INSERT INTO notifications (
		id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	)VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)
	RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`

	updateNotificationQuery = `
		UPADATE notifications
		SET 
			code COALESCE(NULLIF($1, ''), code),
			title = COALESCE(NULLIF($2, ''), title),
			content = COALESCE(NULLIF($3, ''), content),
			type = COALESCE(NULLIF($4, ''), type),
			target = COALESCE(NULLIF($5, ''), target),
			target_user = COALESCE(NULLIF($6, ''), target_user),
			status = COALESCE(NULLIF($7, ''), status),
			modifier_id = COALESCE($8, modifier_id),
			version = version + 1,
			updated_at = $9
		WHERE id = $10
		RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`

	deleteNotificationQuery = `
	UPDATE notifications
	SET
		active = false,
		version = version + 1,
		modifier_id = $1,
		updated_at = $2
	WHERE id = $1 
	RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`

	getNotificationByIdQuery = `
	SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE id = $1 AND active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM notifications
	WHERE active = true
	`

	searchByTitleQuery = `
	SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
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
	SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`
)
