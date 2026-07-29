\echo '=== 1. Connection information ==='

SELECT
    current_database() AS database_name,
    current_user AS connected_user,
    inet_server_addr() AS PostgreSQL_server_ip,
    pg_is_in_recovery() AS is_standby;

\echo '=== 2. Create the demonstration table ==='

CREATE TABLE IF NOT EXISTS paas_demo (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

\echo '=== 3. Insert test data ==='

INSERT INTO paas_demo (message)
VALUES ('CloudNativePG is running successfully on STACKIT SKE');

\echo '=== 4. Read the stored data ==='

SELECT
    id,
    message,
    created_at
FROM paas_demo
ORDER BY id DESC
LIMIT 5;
