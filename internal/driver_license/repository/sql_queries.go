package repository

const (
	createDriverLicenseQuery = `
	INSERT INTO driver_licenses (
		id, owner_id, full_name, dob, identity_no, owner_address, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, version, creator_id, modifier_id, created_at, updated_at, active
	)VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 
		$15, $16, $17, $18, $19, $20, $21
	)RETURNING id, owner_id, full_name, dob, identity_no, owner_address, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, version, creator_id, modifier_id, created_at, updated_at, active
	`

	updateDriverLicenseQuery = `
	UPDATE driver_licenses 
	SET
		owner_id = COALESCE($1, owner_id),
		full_name = COALESCE(NULLIF($2, ''), full_name),
		dob = COALESCE(NULLIF($3, ''), dob),
		owner_address = COALESCE(NULLIF($4, ''), owner_address),
		license_no = COALESCE(NULLIF($5, ''), license_no),
		expiry_date = COALESCE($6, expiry_date),
		status = COALESCE(NULLIF($7, ''), status),
		nationality = COALESCE(NULLIF($8, ''), nationality),
		point = COALESCE($9, point),
		modifier_id = COALESCE($10, modifier_id),
        version = version + 1,
        updated_at = now(),
	WHERE id = $11
	RETURNING *
	`

	deleteDriverLicenseQuery = `
	UPDATE driver_licenses
	SET
		active = false,
		version = version + 1
		modifier_id = $1,
		updated_at = $2
	WHERE id = $3
	RETURNING *
	`

	getDriverLicenseByIdQuery = `
	SELECT *
	FROM driver_licenses
	WHERE id = $1 AND active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM driver_licenses
	WHERE active = true
	`

	findLicenseNOCount = `
		SELECT COUNT(*)
		FROM driver_licenses
		WHERE active = true
		AND license_no ILIKE '%' || $1 || '%'
	`

	searchByLicenseNo = `
    SELECT * 
    FROM driver_licenses
    WHERE license_no ILIKE '%' || $1 || '%' AND active = true	
    ORDER BY license_no
    OFFSET $2 LIMIT $3
`

	getDriverLicense = `
	SELECT id, owner_id, full_name, dob, identity_no, owner_address, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, version, creator_id, modifier_id, created_at, updated_at, active
	FROM driver_licenses
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`

	findLicenseNO = `
	SELECT license_no
	FROM driver_licenses
	WHERE license_no = $1 AND active = true
	`
)
