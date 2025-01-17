-- 启用 TimescaleDB 扩展
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- 创建 events 表
CREATE TABLE IF NOT EXISTS events
(
    id                      SERIAL         NOT NULL,
    client_id               VARCHAR(64)    NOT NULL,
    event_type              TEXT           NOT NULL,
    timestamp               TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    user_id                 VARCHAR(64)    NOT NULL,
    home_world              INT            NOT NULL,
    cheat_banned_hash_valid BOOLEAN        NOT NULL,
    client                  VARCHAR(64)    NOT NULL, -- aid
    os                      TEXT           NOT NULL,
    dalamud_version         TEXT           NOT NULL,
    is_testing              BOOLEAN        NOT NULL DEFAULT FALSE,
    plugin3rd_count         INT            NOT NULL DEFAULT 0,
    machine_id              VARCHAR(255)   NOT NULL,
    PRIMARY KEY (id, timestamp)                   -- 将 timestamp 包含在主键中
);

-- 创建单独的索引，用于通过 client_id 全匹配检索
CREATE INDEX IF NOT EXISTS idx_events_client_id
    ON events (client_id);

-- 创建单独的索引，用于通过 client 全匹配检索
CREATE INDEX IF NOT EXISTS idx_events_client
    ON events (client);

CREATE INDEX IF NOT EXISTS idx_events_plugin3rd_count
    ON events (plugin3rd_count);

CREATE INDEX IF NOT EXISTS idx_events_machine_id
    ON events (machine_id);

-- 将 events 表转换为 hypertable，以利用 TimescaleDB 的时间序列功能
SELECT create_hypertable('events', 'timestamp', if_not_exists => TRUE);


-- 创建 machine_id_plugins 表
CREATE TABLE IF NOT EXISTS machine_id_plugins
(
    machine_id              VARCHAR(255)   NOT NULL PRIMARY KEY,
    plugin3rd_list          JSONB          NOT NULL DEFAULT [],
    plugin3rd_count         INT            NOT NULL DEFAULT 0,
    last_seen               TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_machine_id_plugins_plugin3rd_count
    ON machine_id_plugins (plugin3rd_count);

CREATE INDEX IF NOT EXISTS idx_machine_id_plugins_machine_id
    ON machine_id_plugins (machine_id);