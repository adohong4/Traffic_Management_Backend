package repository

const (
	createTrafficViolationQuery = `
    INSERT INTO traffic_violations (
        id, vehicle_no, date, type, address, description, points, fine_amount, expiry_date, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
    ) RETURNING id, vehicle_no, date, type, address, description, points, fine_amount, status, 
        version, creator_id, modifier_id, created_at, updated_at, active
    `

	updateTrafficViolationQuery = `
    UPDATE traffic_violations 
    SET
        vehicle_no = COALESCE(NULLIF($1, ''), vehicle_no),
        date = COALESCE($2, date),
        type = COALESCE(NULLIF($3, ''), type),
        address = COALESCE($4, address),
        description = COALESCE(NULLIF($5, ''), description),
        points = COALESCE($6, points),
        fine_amount = COALESCE($7, fine_amount),
        expiry_date = COALESCE($8, expiry_date),
        status = COALESCE(NULLIF($9, ''), status),
        modifier_id = COALESCE($10, modifier_id),
        version = version + 1,
        updated_at = $11
    WHERE id = $12
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

	getTrafficViolationStatsQuery = `
        SELECT 
            COUNT(*) AS total_violations,
            COALESCE(SUM(fine_amount), 0) AS total_fine_amount,
            COALESCE(SUM(CASE WHEN status = 'completed' THEN fine_amount ELSE 0 END), 0) AS total_paid_fine_amount,
            COALESCE(SUM(CASE WHEN status != 'cancelled' THEN fine_amount ELSE 0 END), 0) 
            - COALESCE(SUM(CASE WHEN status = 'completed' THEN fine_amount ELSE 0 END), 0) AS total_unpaid_fine_amount
        FROM traffic_violations
        WHERE active = true
    `

	getTrafficViolationStatusStatsQuery = `
        SELECT 
            status,
            COUNT(*) AS total_count,
            COALESCE(SUM(fine_amount), 0) AS total_fine_amount,
            COUNT(*) FILTER (WHERE expiry_date < CURRENT_DATE) AS overdue_count,
            COALESCE(SUM(fine_amount) FILTER (WHERE expiry_date < CURRENT_DATE), 0) AS overdue_fine_amount,
            COUNT(*) FILTER (WHERE expiry_date IS NULL OR expiry_date >= CURRENT_DATE) AS not_overdue_count,
            COALESCE(SUM(fine_amount) FILTER (WHERE expiry_date IS NULL OR expiry_date >= CURRENT_DATE), 0) AS not_overdue_amount
        FROM traffic_violations
        WHERE active = true
        GROUP BY status
        ORDER BY status
    `

	//-----USER-------------
	getViolationsByPlateNo = `
        SELECT *
        FROM traffic_violations
        WHERE vehicle_no = $1 AND active = true
        ORDER BY date DESC
        OFFSET $2 LIMIT $3
    `

	getTotalByPlateNo = `
        SELECT COUNT(*)
        FROM traffic_violations
        WHERE vehicle_no = $1 AND active = true
    `

	getViolationsByOwnerID = `
        SELECT tv.*
        FROM traffic_violations tv
        JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
        WHERE vr.owner_id = $1 AND tv.active = true AND vr.active = true
        ORDER BY tv.date DESC
        OFFSET $2 LIMIT $3
    `

	getTotalViolationsByOwnerID = `
        SELECT COUNT(*)
        FROM traffic_violations tv
        JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
        WHERE vr.owner_id = $1 AND tv.active = true AND vr.active = true
    `

	getViolationsByWalletAddress = `
        SELECT tv.*
        FROM traffic_violations tv
        JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
        JOIN driver_licenses dl ON vr.owner_id = dl.creator_id OR dl.wallet_address = $1
        WHERE (dl.wallet_address = $1 OR vr.owner_id IN (
            SELECT id FROM users WHERE wallet_address = $1
        )) AND tv.active = true AND vr.active = true AND dl.active = true
        ORDER BY tv.date DESC
        OFFSET $2 LIMIT $3
    `

	getViolationsByWallet = `
        SELECT tv.*
        FROM traffic_violations tv
        JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
        JOIN driver_licenses dl ON vr.owner_id = dl.creator_id
        WHERE dl.wallet_address = $1 AND tv.active = true AND vr.active = true
        ORDER BY tv.date DESC
        OFFSET $2 LIMIT $3
    `

	getTotalViolationsByWallet = `
        SELECT COUNT(*)
        FROM traffic_violations tv
        JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
        JOIN driver_licenses dl ON vr.owner_id = dl.creator_id
        WHERE dl.wallet_address = $1 AND tv.active = true AND vr.active = true
    `

	getTrafficViolationByIDAndOwner = `
        SELECT tv.*
        FROM traffic_violations tv
        JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
        WHERE tv.id = $1 
          AND vr.owner_id = $2 
          AND tv.active = true 
          AND vr.active = true
    `

	getViolationsByLicenseWallet = `
    SELECT tv.*
    FROM traffic_violations tv
    JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
    JOIN driver_licenses dl ON vr.owner_id = dl.creator_id
    WHERE dl.wallet_address = $1 
      AND tv.active = true 
      AND vr.active = true 
      AND dl.active = true
    ORDER BY tv.date DESC
    OFFSET $2 LIMIT $3
`

	getTotalViolationsByLicenseWallet = `
    SELECT COUNT(*)
    FROM traffic_violations tv
    JOIN vehicle_registration vr ON tv.vehicle_no = vr.vehicle_no
    JOIN driver_licenses dl ON vr.owner_id = dl.creator_id
    WHERE dl.wallet_address = $1 
      AND tv.active = true 
      AND vr.active = true 
      AND dl.active = true
`
)
