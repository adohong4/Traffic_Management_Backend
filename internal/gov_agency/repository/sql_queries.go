package repository

const (
	createGovAgencyQuery = `
	INSERT INTO gov_agency (
		id, name, address, city, type, phone, email, status, 
		version, created_at, updated_at, active
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, 
		$9, $10, $11, $12
	) RETURNING 
		id, name, address, city, type, phone, email, status, 
		version, created_at, updated_at, active
	`

	updateGovAgencyQuery = `
	UPDATE gov_agency
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

	deleteGovAgencyQuery = `
	UPDATE gov_agency
	SET
		active = false,
		version = version + 1,
		updated_at = $2
	WHERE id = $3
	RETURNING *
	`

	getGovAgencyQuery = `
	SELECT *
	FROM gov_agency
	WHERE id = $1 AND active = true
	`

	getTotalGovAgencyCount = `
	SELECT COUNT(id)
	FROM gov_agency
	WHERE active = true
	`

	searchGovAgencyByNameCount = `
	SELECT COUNT(*)
	FROM gov_agency
	WHERE active = true
	AND name ILIKE '%' || $1 || '%'
	`

	searchGovAgencyByName = `
	SELECT * 
	FROM gov_agency
	WHERE name ILIKE '%' || $1 || '%' AND active = true	
	ORDER BY name
	OFFSET $2 LIMIT $3
	`

	getAllGovAgency = `
	SELECT id, name, address, city, type, phone, email, status, 
		version, updated_at, created_at, active
	FROM gov_agency
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`

	findGovAgencyByName = `
	SELECT *
	FROM gov_agency
	WHERE name = $1 AND active = true
	`
)
