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
			FOREIGN KEY user_id REFERENCES users(id) ON DELETE CASCADE
		)
	`

	CreateAddressTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS addresses (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    address TEXT NOT NULL,
		    coordinate POINT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY user_id REFERENCES users(id) ON DELETE CASCADE
		)
	`

	CreateAddressOwnerTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS address_owner (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    address_id BIGINT NOT NULL,
		    owner_id BIGINT NOT NULL,
		    owner_type VARCHAR(12) NOT NULL,
		    is_default BOOLEAN DEFAULT VALUE FALSE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTMAP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY address_id REFERENCES addresses(id) ON DELETE CASCADE
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
		    created_at TIMESTAMP NOT NULL DEFAUT CURRENT_TIMESTAMP,
		    updated_at TIMESTMAP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY user_id REFERENCES users(id) ON DELETE CASCADE
		)
	`
)
