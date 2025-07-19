package statusmodel

const (
	//status
	ViolationStatusPending     = "Pending"
	ViolationStatusProcessed   = "Processed"
	ViolationStatusCancelled   = "Cancelled"
	ViolationStatusUnderReview = "UnderReview"
	ViolationStatusOverdue     = "Overdue"

	//type
	ViolationTypeSpeeding              = "Speeding"
	ViolationTypeRedLightViolation     = "RedLightViolation"
	ViolationTypeWrongLane             = "WrongLane"
	ViolationTypeNoHelmet              = "NoHelmet"
	ViolationTypeIllegalParking        = "IllegalParking"
	ViolationTypeDrivingUnderInfluence = "DrivingUnderInfluence"
	ViolationTypeNoLicense             = "NoLicense"
	ViolationTypeOverloading           = "Overloading"
	ViolationTypeSignalViolation       = "SignalViolation"
	ViolationTypeOther                 = "Other"
)
