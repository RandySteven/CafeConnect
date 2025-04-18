package queries

const (
	CreateUserTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS users (
		    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
		    name VARCHAR(64) NOT NULL,
		    username VARCHAR(64) NOT NULL,
		    email VARCHAR(64) NOT NULL,
		    password VARCHAR(244) NOT NULL,
		    phone_number VARCHAR(16) NOT NULL,
		    profile_picture VARCHAR(244) NOT NULL,
		    dob DATE NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreatePointTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS points (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    point INT NOT NULL DEFAULT 0,
			user_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP DEFAULT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`

	CreateAddressTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS addresses (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    address TEXT NOT NULL,
		    coordinate POINT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreateAddressUserTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS address_users (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    address_id BIGINT NOT NULL,
		    user_id BIGINT NOT NULL,
		    is_default BOOLEAN DEFAULT FALSE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
		    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`

	CreateReferralTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS referrals (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    code CHAR(16),
		    user_id BIGINT NOT NULL,
		    expired_time TIMESTAMP NOT NULL,
		    status VARCHAR(24) NOT NULL,
		    numb_of_usage INT DEFAULT 0,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`

	CreateRoleTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS roles (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    role VARCHAR(24) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreateRoleUserTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS role_users (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    role_id BIGINT NOT NULL,
		    user_id BIGINT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreateCafeFranchiseTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS cafe_franchises (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    name VARCHAR(64) NOT NULL,
		    logo_url VARCHAR(244) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreateCafeTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS cafes (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    cafe_franchise_id BIGINT NOT NULL,
		    cafe_type VARCHAR(16) NOT NULL,
		    photo_url VARCHAR(244) NOT NULL,
		    open_hour TIME NOT NULL,
		    close_hour TIME NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    cafe_franchise_id 
		)
	`
)
