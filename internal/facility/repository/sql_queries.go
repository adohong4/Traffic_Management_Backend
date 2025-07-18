package repository

const (
	createFacilityQuery = `
	INSERT INTO facilities (
		id, name, address, city, type, phone, email, status, 
		version, creator_id, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, 
		$9, $10, $11, $12
	) RETURNING 
		id, name, address, city, type, phone, email, status, 
		version, creator_id, created_at, updated_at
	`

	updateFacilityQuery = `
	UPDATE facilities
	SET
		name = COALESCE(NULLIF($1, ''), name),
		address = COALESCE(NULLIF($2, ''), address),
		city = COALESCE(NULLIF($3, ''), city),
		type = COALESCE(NULLIF($4, ''), type),
		phone = COALESCE(NULLIF($5, ''), phone),
		email = COALESCE(NULLIF($6, ''), email),
		status = COALESCE(NULLIF($7, ''), status),
		version = version + 1,
		updated_at = NOW(),
	WHERE id = $8
	RETURNING *
	`

	deleteFacilityQuery = `
	UPDATE facilities
	SET
		active = false,
		version = version + 1,
		updated_at = $2
	WHERE id = $3
	RETURNING *
	`

	getFacilityQuery = `
	SELECT *
	FROM facilities
	WHERE id = $1 AND active = true
	`

	getTotalFacilityCount = `
	SELECT COUNT(id)
	FROM facilities
	WHERE active = true
	`

	searchFacilityByNameCount = `
	SELECT COUNT(*)
	FROM facilities
	WHERE active = true
	AND name ILIKE '%' || $1 || '%'
	`

	searchFacilityByName = `
	SELECT * 
	FROM facilities
	WHERE name ILIKE '%' || $1 || '%' AND active = true	
	ORDER BY name
	OFFSET $2 LIMIT $3
	`

	getAllFacilities = `
	SELECT id, name, address, city, type, phone, email, status, 
		version, creator_id, updated_at, created_at
	FROM facilities
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`

	findFacilityByName = `
	SELECT *
	FROM facilities
	WHERE name = $1 AND active = true
	`
)
