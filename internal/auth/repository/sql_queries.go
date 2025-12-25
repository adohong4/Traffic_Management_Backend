package repository

const (
	createUserQuery = `
        INSERT INTO users (
            id, user_address, identity_no, active, role, version, creator_id, modifier_id, 
            created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        ) RETURNING id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at`

	updateUserQuery = `
        UPDATE users 
        SET
            user_address = COALESCE(NULLIF($1, ''), user_address),
            identity_no = COALESCE(NULLIF($2, ''), identity_no),
            active = COALESCE($3, active),
            role = COALESCE(NULLIF($4, ''), role),
            creator_id = COALESCE($5, creator_id),
            modifier_id = COALESCE($6, modifier_id),
            version = version + 1,
            updated_at = now()
        WHERE id = $7 AND version = $8
        RETURNING id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at`

	deleteUserQuery = `
        UPDATE users
        SET
            active = false,
            version = version + 1,
            modifier_id = $1,
            updated_at = $2
        WHERE id = $3 AND version = $4
        RETURNING id`

	getUserQuery = `
        SELECT 
            id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users
        WHERE id = $1 AND active = true`

	getTotalCount = `
        SELECT COUNT(id) 
        FROM users 
        WHERE active = true 
        AND (identity_no ILIKE '%' || $1 || '%' OR role ILIKE '%' || $1 || '%')`

	findUsers = `
        SELECT id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users 
        WHERE active = true 
        AND (identity_no ILIKE '%' || $1 || '%' OR role ILIKE '%' || $1 || '%')
        ORDER BY identity_no, role
        OFFSET $2 LIMIT $3`

	getTotal = `
        SELECT COUNT(id) 
        FROM users 
        WHERE active = true`

	getUsers = `
        SELECT id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users 
        WHERE active = true 
        ORDER BY COALESCE(NULLIF($1, ''), identity_no), role
        OFFSET $2 LIMIT $3`

	findUserByIdentity = `
        SELECT id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users
        WHERE identity_no = $1 AND active = true`

	findUserByUserAddress = `
        SELECT id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at
        FROM users
        WHERE user_address = $1 AND active = true`
)
