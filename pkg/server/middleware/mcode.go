package middleware

import "strings"

// MCode represents a message code with predefined messages
// Naming convention: "{PackageInitial}{CamelCaseFunction}{Number}"
// Example: CLC1 = Config package, LoadConfig function, number 1
type MCode struct {
	Code    string
	Message string
}

// maxCodeLength is the maximum length of all mcode codes, calculated at initialization
var maxCodeLength int

// PaddedCode returns the code padded to maxCodeLength for aligned log output
func (rcvr MCode) PaddedCode() string {
	if len(rcvr.Code) >= maxCodeLength {
		return rcvr.Code
	}
	return rcvr.Code + strings.Repeat(" ", maxCodeLength-len(rcvr.Code))
}

// GetMaxCodeLength returns the current maximum code length
func GetMaxCodeLength() int {
	return maxCodeLength
}

// calculateMaxCodeLength calculates the maximum length from all defined mcodes
func calculateMaxCodeLength() int {
	maxLen := 0
	allMCodes := []MCode{
		// Config package codes
		CLC1, CLC2, CLC3, CLC4,
		CNBC1, CNBC2,
		CNBCWC1, CNBCWC2, CNBCWC3, CNBCWC4, CNBCWC5, CNBCWC6,
		CNBCFS1, CNBCFS2, CNBCFS3, CNBCFS4, CNBCFS5,
		CCDB1, CCDB2, CCDB3, CCDB4, CCDB5,
		CNDBC1, CNDBC2, CNDBC3,
		CSLF1,
		CNSMC1, CNSMC2, CNSMC3, CNSMC4,
		CGCE1, CLCSM1, CLCSM2, CLCSM3, CLCSM4, CLCSM5,
		// Server package codes
		SM1, SM2, SM3, SM4, SM5,
		SIR1, SIR2, SIR3, SIR4, SIR5, SIR6, SIR7,
		// Repository codes
		RNRC1, RNRC2, RNRC3, RNRC4, RNRC5,
		RNCR1, RNCR2,
		RNUR1, RNUR2,
		RNGR1, RNGR2,
		RNMR1, RNMR2,
		RNRR1, RNRR2,
		// Middleware codes
		MNL1, MNL2, MNL3, MNL4, MNL5,
		MLWC1, MLWC2, MLWC3, MLWC4, MLWC5,
		MFP1, MFP2,
		MFI1, MFI2, MFI3, MFI4, MFI5,
		MFPR1, MFPR2, MFPR3, MFPR4, MFPR5,
		MCA1, MCA2, MCA3, MCA4, MCA5, MCA6,
		MVJT1, MVJT2, MVJT3, MVJT4, MVJT5, MVJT6,
		MSUC1, MSUC2,
		MGUFC1, MGUFC2, MGUFC3,
		MGUI1, MGUI2, MGUI3,
		MGUU1, MGUU2, MGUU3,
		MGUE1, MGUE2, MGUE3,
		MGUN1, MGUN2, MGUN3,
		MGUR1, MGUR2, MGUR3,
		MGUC1, MGUC2, MGUC3,
		// System codes
		SYSE1, SYSE2, SYSE3,
		// Controller codes
		UCPCU0, UCPCU1, UCPCU2, UCPCU3, UCPCU4, UCPCU5, UCPCU6,
		UCPGU0, UCPGU1, UCPGU2, UCPGU3,
	}

	for _, mcode := range allMCodes {
		if len(mcode.Code) > maxLen {
			maxLen = len(mcode.Code)
		}
	}
	return maxLen
}

// init initializes the maxCodeLength by calculating from all mcodes
func init() {
	maxCodeLength = calculateMaxCodeLength()
}

// Message codes for locky
var (
	// Config package (C prefix)
	// LoadConfig related
	CLC1 = MCode{"CLC1", "Configuration load start"}
	CLC2 = MCode{"CLC2", "Configuration load success"}
	CLC3 = MCode{"CLC3", "Configuration load failed"}
	CLC4 = MCode{"CLC4", "Configuration parse error"}

	// NewBaseConfig related
	CNBC1 = MCode{"CNBC1", "NewBaseConfig start"}
	CNBC2 = MCode{"CNBC2", "NewBaseConfig success"}

	// NewBaseConfigWithContext related
	CNBCWC1 = MCode{"CNBCWC1", "NewBaseConfigWithContext start"}
	CNBCWC2 = MCode{"CNBCWC2", "Using Secrets Manager"}
	CNBCWC3 = MCode{"CNBCWC3", "SECRET_ID not set, falling back to file"}
	CNBCWC4 = MCode{"CNBCWC4", "Secrets Manager load failed, falling back"}
	CNBCWC5 = MCode{"CNBCWC5", "Secrets Manager load success"}
	CNBCWC6 = MCode{"CNBCWC6", "File-based config success"}

	// NewBaseConfigFromSource related
	CNBCFS1 = MCode{"CNBCFS1", "NewBaseConfigFromSource start"}
	CNBCFS2 = MCode{"CNBCFS2", "Using secretsmanager source"}
	CNBCFS3 = MCode{"CNBCFS3", "Using localfile source"}
	CNBCFS4 = MCode{"CNBCFS4", "Using default source"}
	CNBCFS5 = MCode{"CNBCFS5", "Invalid CONFIG_SOURCE"}

	// ConnectDB related
	CCDB1 = MCode{"CCDB1", "ConnectDB start"}
	CCDB2 = MCode{"CCDB2", "Database already connected"}
	CCDB3 = MCode{"CCDB3", "Database connection attempt"}
	CCDB4 = MCode{"CCDB4", "Database connected successfully"}
	CCDB5 = MCode{"CCDB5", "Database connection failed"}

	// NewDBConnection related
	CNDBC1 = MCode{"CNDBC1", "NewDBConnection start"}
	CNDBC2 = MCode{"CNDBC2", "Database connection success"}
	CNDBC3 = MCode{"CNDBC3", "Database connection error"}

	// SetLoggerFactory related
	CSLF1 = MCode{"CSLF1", "Logger factory set"}

	// Secrets Manager related
	CNSMC1 = MCode{"CNSMC1", "NewSecretsManagerClient start"}
	CNSMC2 = MCode{"CNSMC2", "Using LocalStack"}
	CNSMC3 = MCode{"CNSMC3", "Using AWS production"}
	CNSMC4 = MCode{"CNSMC4", "Client created successfully"}
	CGCE1  = MCode{"CGCE1", "GetConfigFromEnv start"}
	CLCSM1 = MCode{"CLCSM1", "LoadConfigFromSecretsManager start"}
	CLCSM2 = MCode{"CLCSM2", "Secret retrieved successfully"}
	CLCSM3 = MCode{"CLCSM3", "Secret unmarshal success"}
	CLCSM4 = MCode{"CLCSM4", "Secret retrieval failed"}
	CLCSM5 = MCode{"CLCSM5", "Secret unmarshal failed"}

	// Server package (S prefix)
	// Main function related
	SM1 = MCode{"SM1", "Server Main start"}
	SM2 = MCode{"SM2", "Server starting on port"}
	SM3 = MCode{"SM3", "Server ready"}
	SM4 = MCode{"SM4", "Server Run start"}
	SM5 = MCode{"SM5", "Server stopping"}

	// InitRouter related
	SIR1 = MCode{"SIR1", "InitRouter start"}
	SIR2 = MCode{"SIR2", "Redis client created"}
	SIR3 = MCode{"SIR3", "Redis client creation failed"}
	SIR4 = MCode{"SIR4", "Casbin enforcers initialized"}
	SIR5 = MCode{"SIR5", "Controllers initialized"}
	SIR6 = MCode{"SIR6", "Routes registered"}
	SIR7 = MCode{"SIR7", "Router initialization complete"}

	// Repository - NewRedisClient related
	RNRC1 = MCode{"RNRC1", "NewRedisClient start"}
	RNRC2 = MCode{"RNRC2", "Redis connection success"}
	RNRC3 = MCode{"RNRC3", "Redis connection failed"}
	RNRC4 = MCode{"RNRC4", "Redis ping success"}
	RNRC5 = MCode{"RNRC5", "Redis ping failed"}

	// Repository - NewCommonRepository related
	RNCR1 = MCode{"RNCR1", "NewCommonRepository start"}
	RNCR2 = MCode{"RNCR2", "CommonRepository created"}

	// Repository - NewUserRepository related
	RNUR1 = MCode{"RNUR1", "NewUserRepository start"}
	RNUR2 = MCode{"RNUR2", "UserRepository created"}

	// Repository - NewGroupRepository related
	RNGR1 = MCode{"RNGR1", "NewGroupRepository start"}
	RNGR2 = MCode{"RNGR2", "GroupRepository created"}

	// Repository - NewMemberRepository related
	RNMR1 = MCode{"RNMR1", "NewMemberRepository start"}
	RNMR2 = MCode{"RNMR2", "MemberRepository created"}

	// Repository - NewRoleRepository related
	RNRR1 = MCode{"RNRR1", "NewRoleRepository start"}
	RNRR2 = MCode{"RNRR2", "RoleRepository created"}

	// Middleware package (M prefix)
	// NewLogger related
	MNL1 = MCode{"MNL1", "NewLogger start"}
	MNL2 = MCode{"MNL2", "Log level set"}
	MNL3 = MCode{"MNL3", "Output destination set"}
	MNL4 = MCode{"MNL4", "Logger created"}
	MNL5 = MCode{"MNL5", "Failed to open log file"}

	// LoggerWithConfig (request logging)
	MLWC1 = MCode{"MLWC1", "Request start"}
	MLWC2 = MCode{"MLWC2", "Request processed"}
	MLWC3 = MCode{"MLWC3", "Request success"}
	MLWC4 = MCode{"MLWC4", "Client error"}
	MLWC5 = MCode{"MLWC5", "Server error"}

	// ForPublic middleware related
	MFP1 = MCode{"MFP1", "ForPublic middleware start"}
	MFP2 = MCode{"MFP2", "ForPublic middleware end"}

	// ForInternal middleware related
	MFI1 = MCode{"MFI1", "ForInternal middleware start"}
	MFI2 = MCode{"MFI2", "JWT validation start"}
	MFI3 = MCode{"MFI3", "JWT validation success"}
	MFI4 = MCode{"MFI4", "JWT validation failed"}
	MFI5 = MCode{"MFI5", "ForInternal middleware end"}

	// ForPrivate middleware related
	MFPR1 = MCode{"MFPR1", "ForPrivate middleware start"}
	MFPR2 = MCode{"MFPR2", "JWT validation start"}
	MFPR3 = MCode{"MFPR3", "JWT validation success"}
	MFPR4 = MCode{"MFPR4", "JWT validation failed"}
	MFPR5 = MCode{"MFPR5", "ForPrivate middleware end"}

	// CasbinAuthorization related
	MCA1 = MCode{"MCA1", "CasbinAuthorization start"}
	MCA2 = MCode{"MCA2", "User extracted from context"}
	MCA3 = MCode{"MCA3", "User not found in context"}
	MCA4 = MCode{"MCA4", "Authorization check start"}
	MCA5 = MCode{"MCA5", "Authorization granted"}
	MCA6 = MCode{"MCA6", "Authorization denied"}

	// validateJWTToken related
	MVJT1 = MCode{"MVJT1", "validateJWTToken start"}
	MVJT2 = MCode{"MVJT2", "Token extracted from header"}
	MVJT3 = MCode{"MVJT3", "Token missing"}
	MVJT4 = MCode{"MVJT4", "Token validation start"}
	MVJT5 = MCode{"MVJT5", "Token validation success"}
	MVJT6 = MCode{"MVJT6", "Token validation failed"}

	// setUserContext related
	MSUC1 = MCode{"MSUC1", "setUserContext start"}
	MSUC2 = MCode{"MSUC2", "User context set"}

	// getUserFromContext related
	MGUFC1 = MCode{"MGUFC1", "getUserFromContext start"}
	MGUFC2 = MCode{"MGUFC2", "User found in context"}
	MGUFC3 = MCode{"MGUFC3", "User not found in context"}

	// GetUserID related
	MGUI1 = MCode{"MGUI1", "GetUserID start"}
	MGUI2 = MCode{"MGUI2", "UserID retrieved"}
	MGUI3 = MCode{"MGUI3", "UserID not found"}

	// GetUserUUID related
	MGUU1 = MCode{"MGUU1", "GetUserUUID start"}
	MGUU2 = MCode{"MGUU2", "UserUUID retrieved"}
	MGUU3 = MCode{"MGUU3", "UserUUID not found"}

	// GetUserEmail related
	MGUE1 = MCode{"MGUE1", "GetUserEmail start"}
	MGUE2 = MCode{"MGUE2", "UserEmail retrieved"}
	MGUE3 = MCode{"MGUE3", "UserEmail not found"}

	// GetUserName related
	MGUN1 = MCode{"MGUN1", "GetUserName start"}
	MGUN2 = MCode{"MGUN2", "UserName retrieved"}
	MGUN3 = MCode{"MGUN3", "UserName not found"}

	// GetUserRole related
	MGUR1 = MCode{"MGUR1", "GetUserRole start"}
	MGUR2 = MCode{"MGUR2", "UserRole retrieved"}
	MGUR3 = MCode{"MGUR3", "UserRole not found"}

	// GetUserClaims related
	MGUC1 = MCode{"MGUC1", "GetUserClaims start"}
	MGUC2 = MCode{"MGUC2", "UserClaims retrieved"}
	MGUC3 = MCode{"MGUC3", "UserClaims not found"}

	// System Error codes (SYS prefix)
	SYSE1 = MCode{"SYSE1", "System error"}
	SYSE2 = MCode{"SYSE2", "Unexpected error"}
	SYSE3 = MCode{"SYSE3", "Fatal error"}

	// Controller package codes (UC = User Controller, GC = Group Controller, etc.)
	// User Controller Public (UCPCU = User Controller Public Create User)
	UCPCU0 = MCode{"UCPCU0", "Request received"}
	UCPCU1 = MCode{"UCPCU1", "Bind request failed"}
	UCPCU2 = MCode{"UCPCU2", "Required fields missing"}
	UCPCU3 = MCode{"UCPCU3", "Email already exists"}
	UCPCU4 = MCode{"UCPCU4", "Password hash failed"}
	UCPCU5 = MCode{"UCPCU5", "User created"}
	UCPCU6 = MCode{"UCPCU6", "Response sent"}

	// User Controller Public Get Users
	UCPGU0 = MCode{"UCPGU0", "Request received"}
	UCPGU1 = MCode{"UCPGU1", "Bind request failed"}
	UCPGU2 = MCode{"UCPGU2", "Users retrieved"}
	UCPGU3 = MCode{"UCPGU3", "Response sent"}
)
