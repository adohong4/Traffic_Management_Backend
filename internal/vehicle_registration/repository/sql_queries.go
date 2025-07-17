package repository

const (
	createLicenseQuery = `
	INSERT INTO vehicle_documents (
		id, owner_id, brand, type_vehicle, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
		owner_name, document_type, document_no, issue_date, expiry_date, issuer, status, 
		version, creator_id, modifier_id, created_at, updated_at
	)VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $13, $14, $15, $16, $17, $18, $19, $20, $21
	)RETURNING 
		id, owner_id, brand, type_vehicle vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
		owner_name, document_type, document_no, issue_date, expiry_date, issuer, 
		status, version, creator_id, modifier_id, created_at, updated_at
	`

	updateLicenseQuery = `
    UPDATE vehicle_documents
    SET
        owner_id = COALESCE(NULLIF($1, ''), owner_id),
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
        expiry_date = COALESCE($12, expiry_date),
        issuer = COALESCE(NULLIF($13, ''), issuer),
        status = COALESCE(NULLIF($14, ''), status),
        modifier_id = COALESCE($15, modifier_id),
        version = version + 1,
        updated_at = now(),
        active = COALESCE($16, active)
    WHERE id = $17
    RETURNING *
`

	deleteLicenseQuery = `
	UPDATE vehicle_documents
	SET
		active = false
		version = version + 1
		modifier_id = $1,
		updated_at = $2
	WHERE id = $3
	RETURNING *
	`

	getLicenseQuery = `
	SELECT *
	FROM vehicle_documents
	WHERE id = $1 AND active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM vehicle_documents
	WHERE active = true
	`

	findByVehiclePlateNOCount = `
	SELECT COUNT(*)
	FROM  vehicle_documents
	WHERE active = true
	AND vehicle_no ILIKE '%' || $1 || '%'
	ORDER BY vehicle_no
	OFFSET $2 LIMIT $3
	`

	findByVehiclePlateNO = `
    SELECT * 
    FROM vehicle_documents
    WHERE vehicle_no ILIKE '%' || $1 || '%' AND active = true
    ORDER BY vehicle_no
    OFFSET $2 LIMIT $3
`

	getVehicleDocuments = `
	SELECT id, owner_id, brand, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,owner_name,
		document_type, document_no, issue_date, expiry_date, issuer, 
		status, version, creator_id, modifier_id, created_at, updated_at,
	FROM vehicle_documents
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`
)
