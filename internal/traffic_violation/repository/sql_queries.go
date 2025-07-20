package repository

const (
	createTrafficViolationQuery = `
    INSERT INTO traffic_violations (
        id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
    ) RETURNING id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    `

	updateTrafficViolationQuery = `
    UPDATE traffic_violations 
    SET
        vehicle_no = COALESCE(NULLIF($1, ''), vehicle_no),
        date = COALESCE($2, date),
        type = COALESCE(NULLIF($3, ''), type),
        description = COALESCE(NULLIF($4, ''), description),
        points = COALESCE($5, points),
        fine_amount = COALESCE($6, fine_amount),
        status = COALESCE(NULLIF($7, ''), status),
        modifier_id = COALESCE($8, modifier_id),
        version = version + 1,
        updated_at = $9
    WHERE id = $10
    RETURNING id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    `

	deleteTrafficViolationQuery = `
    UPDATE traffic_violations
    SET
        active = false,
        version = version + 1,
        modifier_id = $1,
        updated_at = $2
    WHERE id = $3
    RETURNING id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    `

	getTrafficViolationByIdQuery = `
    SELECT id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    FROM traffic_violations
    WHERE id = $1 AND active = true
    `

	getTrafficViolationTotalCount = `
    SELECT COUNT(id)
    FROM traffic_violations
    WHERE active = true
    `
	getTrafficViolationQuery = `
    SELECT id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    FROM traffic_violations
    WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
    `

	findVehiclePlateNoCount = `
    SELECT COUNT(*)
    FROM traffic_violations
    WHERE active = true
    AND vehicle_no ILIKE '%' || $1 || '%'
    `

	searchByVehicleNo = `
    SELECT id, vehicle_no, date, type, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    FROM traffic_violations
    WHERE vehicle_no ILIKE '%' || $1 || '%' AND active = true
    ORDER BY vehicle_no
    OFFSET $2 LIMIT $3
    `

	findVehicleNo = `
    SELECT vehicle_no
    FROM traffic_violations
    WHERE vehicle_no = $1 AND active = true
    `
)
