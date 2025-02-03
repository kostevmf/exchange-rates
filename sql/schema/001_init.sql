-- +goose Up

CREATE TABLE IF NOT EXISTS currency (
    currency_id INT UNSIGNED NOT NULL PRIMARY KEY COMMENT 'Currency ID',
    code VARCHAR(3) NOT NULL COMMENT 'Currency Code',
    name VARCHAR(100) NOT NULL COMMENT 'Currency Name',
    scale INT UNSIGNED NOT NULL COMMENT 'Scale',
    UNIQUE KEY CURRENCY_CODE (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_general_ci COMMENT='Currency';

CREATE TABLE IF NOT EXISTS currency_rate (
    rate_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'Rate ID',
    currency_id INT UNSIGNED NOT NULL COMMENT 'Currency ID',
    rate DECIMAL(8,4) NOT NULL COMMENT 'Currency Conversion Rate',
    date DATE NOT NULL COMMENT 'Date',
    INDEX CURRENCY_RATE_DATE (date),
    UNIQUE KEY CURRENCY_RATE_DATE_CURRENCY_ID (date, currency_id),
    CONSTRAINT CURRENCY_RATE_CURRENCY_ID_CURRENCY_CURRENCY_ID
        FOREIGN KEY (currency_id) REFERENCES currency (currency_id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_general_ci COMMENT='Currency Rate';

-- +goose Down

DROP TABLE currency_rate;
DROP TABLE currency;