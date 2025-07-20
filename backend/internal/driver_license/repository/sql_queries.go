package repository

const (
	createDriverLicenseQuery = `
	INSERT INTO driver_licenses (
		id, full_name, dob, identity_no, owner_address, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, version, creator_id, modifier_id, created_at, updated_at, active
	)VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 
		$15, $16, $17, $18, $19, $20
	)RETURNING id, full_name, dob, identity_no, owner_address, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, version, creator_id, modifier_id, created_at, updated_at, active
	`

	updateDriverLicenseQuery = `
	UPDATE driver_licenses 
	SET
		full_name = COALESCE(NULLIF($1, ''), full_name),
		dob = COALESCE(NULLIF($2, '')::DATE, dob),
		owner_address = COALESCE(NULLIF($3, ''), owner_address),
		expiry_date = COALESCE($4, expiry_date),
		status = COALESCE(NULLIF($5, ''), status),
		nationality = COALESCE(NULLIF($6, ''), nationality),
		point = COALESCE($7, point),
		modifier_id = COALESCE($8, modifier_id),
        version = version + 1,
        updated_at = $9
	WHERE id = $10
	RETURNING *
	`

	deleteDriverLicenseQuery = `
	UPDATE driver_licenses
	SET
		active = false,
		version = version + 1,
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
	SELECT id, full_name, dob, identity_no, owner_address, license_no, 
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
