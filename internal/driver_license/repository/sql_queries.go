package repository

const (
	createDriverLicenseQuery = `
	INSERT INTO driver_licenses (
		id, full_name, avatar, dob, identity_no, owner_address, owner_city, license_no,  
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority, 
		nationality, point, wallet_address, on_blockchain, blockchain_txhash,
		version, creator_id, modifier_id, created_at, updated_at, active
	)VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 
		$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
	)RETURNING id, full_name, avatar, dob, identity_no, owner_address, owner_city, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, wallet_address, on_blockchain, blockchain_txhash, 
		version, creator_id, modifier_id, created_at, updated_at, active
	`

	updateDriverLicenseQuery = `
	UPDATE driver_licenses 
	SET
		full_name = COALESCE(NULLIF($1, ''), full_name),
		dob = COALESCE(NULLIF($2, '')::DATE, dob),
		owner_address = COALESCE(NULLIF($3, ''), owner_address),
		owner_city = COALESCE(NULLIF($4, ''), owner_city),
		expiry_date = COALESCE($5, expiry_date),
		status = COALESCE(NULLIF($6, ''), status),
		nationality = COALESCE(NULLIF($7, ''), nationality),
		point = COALESCE($8, point),
		modifier_id = COALESCE($9, modifier_id),
        version = version + 1,
        updated_at = $10
	WHERE id = $11
	RETURNING *
	`
	updateBlockchainConfirmationQuery = `
    UPDATE driver_licenses 
    SET
        blockchain_txhash = COALESCE(NULLIF($1, ''), blockchain_txhash),
        on_blockchain = $2,
        modifier_id = COALESCE($3, modifier_id),
        version = version + 1,
        updated_at = $4
    WHERE id = $5
    RETURNING *
    `

	updateWalletAddressQuery = `
    UPDATE driver_licenses 
    SET
        wallet_address = COALESCE(NULLIF($1, ''), wallet_address),
        modifier_id = COALESCE($2, modifier_id),
        version = version + 1,
        updated_at = $3
    WHERE id = $4
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

	getDriverLicenseByWalletAddressQuery = `
	SELECT *
	FROM driver_licenses
	WHERE wallet_address = $1 AND active = true
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
	SELECT id, full_name, dob, identity_no, owner_address, owner_city, license_no, 
		issue_date, expiry_date, status, license_type, authority_id, issuing_authority,
		nationality, point, wallet_address, on_blockchain, blockchain_txhash, 
		version, creator_id, modifier_id, created_at, updated_at, active
	FROM driver_licenses
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`

	findLicenseNO = `
	SELECT license_no
	FROM driver_licenses
	WHERE license_no = $1 AND active = true
	`

	// Statistic
	getStatusDistributionQuery = `
        SELECT status, COUNT(*) as count
        FROM driver_licenses
        WHERE active = true
        GROUP BY status
        ORDER BY count DESC
    `

	getLicenseTypeDistributionQuery = `
        SELECT license_type, COUNT(*) as count
        FROM driver_licenses
        WHERE active = true AND license_type IS NOT NULL AND license_type != ''
        GROUP BY license_type
        ORDER BY count DESC
    `

	getLicenseTypeStatusDistributionQuery = `
        SELECT 
            license_type,
            status,
            COUNT(*) as count
        FROM driver_licenses
        WHERE active = true 
          AND license_type IS NOT NULL 
          AND license_type != ''
        GROUP BY license_type, status
        ORDER BY license_type, 
                 count DESC,
                 status
    `

	getCityStatusDistributionQuery = `
        SELECT 
            COALESCE(owner_city, 'Không xác định') as owner_city,
            status,
            COUNT(*) as count
        FROM driver_licenses
        WHERE active = true
        GROUP BY owner_city, status
        ORDER BY count DESC, owner_city, status
    `
)
