package repository

const (
	createNewsQuery = `
        INSERT INTO news (
            id, code, image, title, content, category, author, type, target,
            target_user, tag, view, status, version, creator_id, modifier_id,
            created_at, updated_at, active
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
        ) RETURNING *
    `

	updateNewsQuery = `
        UPDATE news 
        SET
            code = COALESCE(NULLIF($1, ''), code),
            image = COALESCE(NULLIF($2, ''), image),
            title = COALESCE(NULLIF($3, ''), title),
            content = COALESCE(NULLIF($4, ''), content),
            category = COALESCE(NULLIF($5, ''), category),
            author = COALESCE(NULLIF($6, ''), author),
            type = COALESCE(NULLIF($7, ''), type),
            target = COALESCE(NULLIF($8, ''), target),
            target_user = COALESCE(NULLIF($9, ''), target_user),
            tag = COALESCE($10, tag),
            status = COALESCE(NULLIF($11, ''), status),
            modifier_id = $12,
            version = version + 1,
            updated_at = $13
        WHERE id = $14 AND active = true
        RETURNING *
    `

	deleteNewsQuery = `
        UPDATE news
        SET active = false,
            modifier_id = $1,
            version = version + 1,
            updated_at = $2
        WHERE id = $3 AND active = true
        RETURNING *
    `

	getNewsByIdQuery = `
        SELECT * FROM news WHERE id = $1 AND active = true
    `

	incrementViewQuery = `
        UPDATE news SET view = view + 1, updated_at = NOW()
        WHERE id = $1 AND active = true
        RETURNING view
    `

	getAllNewsQuery = `
        SELECT * FROM news
        WHERE active = true
        ORDER BY created_at DESC
        OFFSET $1 LIMIT $2
    `

	getTotalCountQuery = `
        SELECT COUNT(*) FROM news WHERE active = true
    `
)
