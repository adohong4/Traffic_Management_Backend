package repository

const (
	createLicenseQuery = `
	INSERT INTO vehicle_registration (
		id, owner_id, brand, type_vehicle, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
		owner_name, seats, issue_date, issuer, registration_date, expiry_date, registration_place, on_blockchain, blockchain_txhash, status, 
		version, creator_id, modifier_id, created_at, updated_at
	)VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
	)RETURNING 
		id, owner_id, brand, type_vehicle, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
		owner_name,seats, issue_date,issuer, registration_date, expiry_date, registration_place, on_blockchain, blockchain_txhash,
		status, version, creator_id, modifier_id, created_at, updated_at
	`

	updateLicenseQuery = `
    UPDATE vehicle_registration
    SET
        owner_id = COALESCE($1, owner_id),
        brand = COALESCE(NULLIF($2, ''), brand),
        type_vehicle = COALESCE(NULLIF($3, ''), type_vehicle),
        vehicle_no = COALESCE(NULLIF($4, ''), vehicle_no),
        color_plate = COALESCE(NULLIF($5, ''), color_plate),
        chassis_no = COALESCE(NULLIF($6, ''), chassis_no),
        engine_no = COALESCE(NULLIF($7, ''), engine_no),
        color_vehicle = COALESCE(NULLIF($8, ''), color_vehicle),
        owner_name = COALESCE(NULLIF($9, ''), owner_name),
        seats = COALESCE($10, seats),
        issue_date = COALESCE($11, issue_date),
		issuer = COALESCE(NULLIF($12, ''), issuer),
		registration_date = COALESCE($13, registration_date),
        expiry_date = COALESCE($14, expiry_date),
		registration_place = COALESCE($15, registration_place),
        status = COALESCE(NULLIF($16, ''), status),
        modifier_id = COALESCE($17, modifier_id),
        version = version + 1,
        updated_at = now(),
        active = COALESCE($18, active)
    WHERE id = $19
    RETURNING *
	`

	updateBlockchainConfirmationQuery = `
    UPDATE vehicle_registration 
    SET
        blockchain_txhash = COALESCE(NULLIF($1, ''), blockchain_txhash),
        on_blockchain = $2,
        modifier_id = COALESCE($3, modifier_id),
        version = version + 1,
        updated_at = $4
    WHERE id = $5
    RETURNING *
    `

	deleteLicenseQuery = `
	UPDATE vehicle_registration
	SET
		active = false,
		version = version + 1,
		modifier_id = $1,
		updated_at = $2
	WHERE id = $3
	RETURNING *
	`

	getLicenseQuery = `
	SELECT *
	FROM vehicle_registration
	WHERE id = $1 AND active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM vehicle_registration
	WHERE active = true
	`

	findByVehiclePlateNOCount = `
		SELECT COUNT(*)
		FROM vehicle_registration
		WHERE active = true
		AND vehicle_no ILIKE '%' || $1 || '%'
	`

	searchByVehiclePlateNO = `
    SELECT * 
    FROM vehicle_registration
    WHERE vehicle_no ILIKE '%' || $1 || '%' AND active = true	
    ORDER BY vehicle_no
    OFFSET $2 LIMIT $3
	`

	getVehicleDocuments = `
	SELECT id, owner_id, brand, type_vehicle, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
    	owner_name, issue_date, expiry_date, issuer, 
    	status, version, creator_id, modifier_id, updated_at, created_at
	FROM vehicle_registration
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`

	findVehiclePlateNO = `
	SELECT *
	FROM vehicle_registration
	WHERE vehicle_no = $1 AND active = true
	`

	// Query for count by type_vehicle
	getCountByType = `
    SELECT type_vehicle, COUNT(*) as count
    FROM vehicle_registration
    WHERE active = true
    GROUP BY type_vehicle
    `

	getRegistrationStatusStats = `
    SELECT 
        COUNT(*) FILTER (WHERE expiry_date >= CURRENT_DATE AND expiry_date IS NOT NULL) AS valid_count,
        COUNT(*) FILTER (WHERE expiry_date < CURRENT_DATE AND expiry_date IS NOT NULL) AS expired_count,
        COUNT(*) FILTER (WHERE expiry_date IS NULL OR registration_date IS NULL) AS pending_count
    FROM vehicle_registration
    WHERE active = true
      AND type_vehicle NOT ILIKE ANY (ARRAY[
        '%xe máy%', '%xe mô tô%', '%xe gắn máy%', 
        '%xe đạp%', '%xe đạp điện%', '%xe máy điện%'
      ])
    `

	// Query for top 5 brands
	getTopBrands = `
    SELECT brand, COUNT(*) as count
    FROM vehicle_registration
    WHERE active = true
    GROUP BY brand
    ORDER BY count DESC
    LIMIT 5
    `

	// Total active vehicles for others calculation
	getTotalActiveVehicles = `
    SELECT COUNT(*)
    FROM vehicle_registration
    WHERE active = true
    `

	// User - Owner ID
	getVehiclesByOwnerID = `
        SELECT *
        FROM vehicle_registration
        WHERE owner_id = $1 AND active = true
        ORDER BY updated_at DESC, created_at DESC
        OFFSET $2 LIMIT $3
    `

	getTotalCountByOwnerID = `
        SELECT COUNT(*)
        FROM vehicle_registration
        WHERE owner_id = $1 AND active = true
    `

	getVehicleByIDAndOwner = `
        SELECT *
        FROM vehicle_registration
        WHERE id = $1 AND owner_id = $2 AND active = true
    `
)

var excludedVehicleTypes = []string{
	"%xe máy%",
	"%xe mô tô%",
	"%xe gắn máy%",
	"%xe đạp%",
	"%xe đạp điện%",
	"%xe máy điện%",
}
