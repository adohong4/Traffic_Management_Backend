package repository

const (
	createUserQuery = `
		INSERT INTO users (
			id, identity_no, password, active, role, version,  creator_id, modifier_id, 
            created_at, updated_at
		) VALUES (
			 $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)RETURNING id, identity_no, password, active, role, version, creator_id, modifier_id, created_at, updated_at`

	updateUserQuery = `
		UPDATE users 
		SET
			identity_no = COALESCE(NULLIF($1, ''), identity_no),
        	password = COALESCE(NULLIF($2, ''), password),
        	active = COALESCE($4, active),
        	role = COALESCE(NULLIF($5, ''), role),
        	creator_id = COALESCE($6, creator_id),
        	modifier_id = COALESCE($7, modifier_id),
        	version = version + 1,
        	updated_at = now()
		WHERE id = $8 AND version = $9
		RETURNING id, identity_no, password, active, role, version,  creator_id, modifier_id, created_at, updated_at`

	deleteUserQuery = `
		UPDATE users
		SET
			active = false
			version = version + 1
			modifier_id = $1
			updated_at = $2
		WHERE id = $3 AND version = $4
		RETURNING id`

	getUserQuery = `
		SELECT 
			id, identity_no, password, active, role, version, creator_id, modifier_id, created_at, updated_at
		FROM users
		WHERE id = $1 AND active = true`

	// getTotalCount counts the total number of active users matching a search term
	getTotalCount = `
        SELECT COUNT(id) 
        FROM users 
        WHERE active = true 
        AND (identity_no ILIKE '%' || $1 || '%' OR role ILIKE '%' || $1 || '%')`

	// findUsers retrieves users matching a search term with pagination
	findUsers = `
        SELECT id, identity_no, password, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users 
        WHERE active = true 
        AND (identity_no ILIKE '%' || $1 || '%' OR role ILIKE '%' || $1 || '%')
        ORDER BY identity_no, role
        OFFSET $2 LIMIT $3`

	// getTotal counts the total number of active users
	getTotal = `
        SELECT COUNT(id) 
        FROM users 
        WHERE active = true`

	// getUsers retrieves all active users with pagination and dynamic sorting
	getUsers = `
        SELECT id, identity_no, password, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users 
        WHERE active = true 
        ORDER BY COALESCE(NULLIF($1, ''), identity_no), role
        OFFSET $2 LIMIT $3`

	findUserByIdentity = `
    SELECT id, identity_no, password, active, role, version, creator_id, modifier_id, created_at, updated_at
    FROM users
    WHERE identity_no = $1 AND active = true`
)
