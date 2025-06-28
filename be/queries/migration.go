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
		    address_id BIGINT NOT NULL,
		    cafe_franchise_id BIGINT NOT NULL,
		    cafe_type VARCHAR(16) NOT NULL,
		    photo_urls TEXT NOT NULL,
		    open_hour TIME NOT NULL,
		    close_hour TIME NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (cafe_franchise_id) REFERENCES cafe_franchises(id) ON DELETE CASCADE,
		    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE
		)
	`

	CreateProductCategoryTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS product_categories (
		 	id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		 	category VARCHAR(24) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreateProductTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS products (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    name VARCHAR(244) NOT NULL,
		    photo_url VARCHAR(244) NOT NULL,
		    product_category_id BIGINT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (product_category_id) REFERENCES product_categories(id) ON DELETE CASCADE
		)
	`

	CreateCafeProductTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS cafe_products (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    cafe_id BIGINT NOT NULL,
		    product_id BIGINT NOT NULL,
		    price BIGINT NOT NULL,
		    stock INT DEFAULT 0,
		    status VARCHAR(24) DEFAULT "AVAILABLE",
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (cafe_id) REFERENCES cafes (id) ON DELETE CASCADE,
		    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
		)
	`

	CreateReviewTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS reviews (
			id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			cafe_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			score FLOAT NOT NULL,
			comment VARCHAR(244) DEFAULT "",
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (cafe_id) REFERENCES cafes (id) ON DELETE CASCADE,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`

	CreateCartTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS carts (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    user_id BIGINT NOT NULL,
		    cafe_product_id BIGINT NOT NULL,
		    qty INT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		    FOREIGN KEY (cafe_product_id) REFERENCES cafe_products (id) ON DELETE CASCADE
		)
	`

	CreateTransactionHeaderTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS transaction_headers (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    user_id BIGINT NOT NULL,
		    cafe_id BIGINT NOT NULL,
		    transaction_code CHAR(24) NOT NULL,
		    status VARCHAR(16) NOT NULL,
		    transaction_at TIMESTAMP NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		    FOREIGN KEY (cafe_id) REFERENCES cafes(id) ON DELETE CASCADE
		)	
	`

	CreateTransactionDetailTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS transaction_details (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    transaction_id BIGINT NOT NULL,
		    cafe_product_id BIGINT NOT NULL,
		    qty BIGINT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL,
		    FOREIGN KEY (transaction_id) REFERENCES transaction_headers (id) ON DELETE CASCADE,
		    FOREIGN KEY (cafe_product_id) REFERENCES cafe_products(id) ON DELETE CASCADE
		)
	`

	CreateMidtransTransactionTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS midtrans_transactions (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    transaction_code CHAR(24) NOT NULL,
		    token VARCHAR(64) NOT NULL,
		    total_amt BIGINT NOT NULL,
		    redirect_url VARCHAR(144) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP DEFAULT NULL
		)
	`

	CreateVerifyTokenTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS verify_tokens (
			id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT, 
			token VARCHAR(144) NOT NULL,
			user_id BIGINT NOT NULL,
			is_clicked BOOLEAN DEFAULT FALSE,
			expired_time TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP DEFAULT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE    
		)
	`

	CreateSettlementTable MigrationQuery = `
		CREATE TABLE IF NOT EXISTS settlements (
		    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		    midtrans_transaction_id VARCHAR(144) NOT NULL,
		    transaction_code CHAR(24) NOT NULL,
		    franchise_id BIGINT NOT NULL,
		    payed_amount BIGINT NOT NULL,
		    point BIGINT NOT NULL,
		    gross_amount BIGINT NOT NULL,
		    metadata TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAUlT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP DEFAULT NULL,
		)
	`
)
