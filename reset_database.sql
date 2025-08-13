-- 创建数据库
CREATE DATABASE IF NOT EXISTS urban_traffic CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE urban_traffic;

/* -- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    user_type VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB;

-- 创建车辆表
CREATE TABLE IF NOT EXISTS vehicles (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    user_id INT UNSIGNED NOT NULL,
    plate_number VARCHAR(30) UNIQUE NOT NULL,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    color VARCHAR(30),
    type VARCHAR(30) NOT NULL DEFAULT '小型汽车',
    reg_date VARCHAR(15),
    is_default BOOLEAN DEFAULT FALSE,
    KEY idx_vehicles_deleted_at (deleted_at),
    KEY idx_vehicles_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 创建停车记录表
CREATE TABLE IF NOT EXISTS parking_records (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    vehicle_id INT UNSIGNED NOT NULL,
    location VARCHAR(100) NOT NULL,
    start_time VARCHAR(25),
    end_time VARCHAR(25),
    fee DECIMAL(10,2) DEFAULT 0.00,
    duration DECIMAL(8,2) DEFAULT 0.00,
    KEY idx_parking_records_deleted_at (deleted_at),
    KEY idx_parking_records_vehicle_id (vehicle_id),
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 插入默认用户数据
INSERT IGNORE INTO users (username, password, email, user_type) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@example.com', 'admin'),
('user', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user@example.com', 'user');

-- 删除所有表并重新创建

-- 先删除有外键约束的表（按依赖关系顺序）
DROP TABLE IF EXISTS parking_records;
DROP TABLE IF EXISTS special_spots;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS parking_lots; */

-- 重新创建表结构（使用BIGINT UNSIGNED匹配GORM默认）
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    user_type VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB;

CREATE TABLE vehicles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    plate_number VARCHAR(30) UNIQUE NOT NULL,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    color VARCHAR(30),
    type VARCHAR(30) NOT NULL DEFAULT '小型汽车',
    reg_date VARCHAR(15),
    is_default BOOLEAN DEFAULT FALSE,
    KEY idx_vehicles_deleted_at (deleted_at),
    KEY idx_vehicles_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE parking_lots (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    total_spots INT NOT NULL DEFAULT 0,
    available_spots INT NOT NULL DEFAULT 0,
    hourly_rate DECIMAL(10,2) DEFAULT 10.00,
    is_active BOOLEAN DEFAULT TRUE,
    operating_hours VARCHAR(50) DEFAULT '24小时',
    payment_methods VARCHAR(100) DEFAULT '微信,支付宝,现金',
    KEY idx_parking_lots_deleted_at (deleted_at),
    KEY idx_parking_lots_location (latitude, longitude)
) ENGINE=InnoDB;

CREATE TABLE special_spots (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    parking_lot_id BIGINT UNSIGNED NOT NULL,
    spot_type ENUM('charging', 'disabled', 'vip') NOT NULL,
    total_count INT NOT NULL DEFAULT 0,
    available_count INT NOT NULL DEFAULT 0,
    additional_fee DECIMAL(10,2) DEFAULT 0.00,
    KEY idx_special_spots_deleted_at (deleted_at),
    KEY idx_special_spots_parking_lot_id (parking_lot_id),
    KEY idx_special_spots_type (spot_type),
    FOREIGN KEY (parking_lot_id) REFERENCES parking_lots(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE parking_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    vehicle_id BIGINT UNSIGNED NOT NULL,
    parking_lot_id BIGINT UNSIGNED,
    location VARCHAR(100) NOT NULL,
    start_time VARCHAR(25),
    end_time VARCHAR(25),
    fee DECIMAL(10,2) DEFAULT 0.00,
    duration DECIMAL(8,2) DEFAULT 0.00,
    spot_type ENUM('normal', 'charging', 'disabled', 'vip') DEFAULT 'normal',
    KEY idx_parking_records_deleted_at (deleted_at),
    KEY idx_parking_records_vehicle_id (vehicle_id),
    KEY idx_parking_records_parking_lot_id (parking_lot_id),
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE,
    FOREIGN KEY (parking_lot_id) REFERENCES parking_lots(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 插入初始停车场数据
INSERT INTO parking_lots (name, address, latitude, longitude, total_spots, available_spots, hourly_rate, is_active, operating_hours, payment_methods) VALUES
('西湖停车场', '杭州市西湖区西湖景区内', 30.2394, 120.1509, 120, 25, 10.00, TRUE, '24小时', '微信,支付宝,现金'),
('中央公园停车场', '杭州市上城区中央公园旁', 30.2635, 120.1709, 200, 42, 10.00, TRUE, '24小时', '微信,支付宝,现金'),
('东方广场停车场', '杭州市江干区东方广场地下', 30.2793, 120.2194, 80, 5, 10.00, TRUE, '24小时', '微信,支付宝,现金'),
('火车站停车场', '杭州市上城区火车站南广场', 30.2431, 120.1816, 150, 60, 10.00, TRUE, '24小时', '微信,支付宝,现金'),
('商业中心停车场', '杭州市西湖区商业中心B1层', 30.2586, 120.1619, 90, 12, 10.00, TRUE, '24小时', '微信,支付宝,现金'),
('医院停车场', '杭州市上城区人民医院停车场', 30.2547, 120.1734, 60, 8, 10.00, TRUE, '24小时', '微信,支付宝,现金');

-- 插入特殊车位数据
INSERT INTO special_spots (parking_lot_id, spot_type, total_count, available_count, additional_fee) VALUES
-- 西湖停车场特殊车位
(1, 'charging', 5, 3, 5.00),
(1, 'disabled', 3, 2, 0.00),
(1, 'vip', 2, 1, 15.00),
-- 中央公园停车场特殊车位
(2, 'charging', 8, 5, 5.00),
(2, 'disabled', 6, 4, 0.00),
(2, 'vip', 4, 2, 15.00),
-- 东方广场停车场特殊车位
(3, 'charging', 3, 1, 5.00),
(3, 'disabled', 2, 1, 0.00),
(3, 'vip', 1, 0, 15.00),
-- 火车站停车场特殊车位
(4, 'charging', 10, 6, 5.00),
(4, 'disabled', 5, 3, 0.00),
(4, 'vip', 3, 2, 15.00),
-- 商业中心停车场特殊车位
(5, 'charging', 4, 2, 5.00),
(5, 'disabled', 2, 1, 0.00),
(5, 'vip', 3, 1, 15.00),
-- 医院停车场特殊车位
(6, 'charging', 2, 1, 5.00),
(6, 'disabled', 4, 3, 0.00),
(6, 'vip', 0, 0, 15.00);

-- 重新插入用户数据（在表重新创建后）
INSERT INTO users (username, password, email, user_type) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@example.com', 'admin'),
('user', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user@example.com', 'user');

-- 插入示例车辆数据
INSERT INTO vehicles (user_id, plate_number, brand, model, color, type) VALUES
(1, '京A12345', '奔驰', 'C200L', '黑色', '小型汽车'),
(2, '京B67890', '宝马', 'X3', '白色', '小型汽车');

-- 插入示例停车记录
INSERT INTO parking_records (vehicle_id, parking_lot_id, location, start_time, end_time, fee, duration, spot_type) VALUES
(1, 1, '西湖停车场', '2025-08-12 09:30:00', NULL, 18.00, 2.75, 'normal'),
(2, 2, '中央公园停车场', '2025-08-10 14:20:00', '2025-08-10 18:05:00', 24.00, 3.75, 'normal'),
(2, 3, '东方广场停车场', '2025-08-03 10:15:00', '2025-08-03 12:30:00', 18.00, 2.25, 'normal');

-- ==================== 交通相关表 ====================

-- 创建交通流量表
CREATE TABLE traffic_flows (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    location VARCHAR(100) NOT NULL,
    flow_count INT NOT NULL,
    speed DECIMAL(5,2),
    direction VARCHAR(20),
    vehicle_type VARCHAR(30),
    timestamp TIMESTAMP NOT NULL,
    KEY idx_traffic_flows_deleted_at (deleted_at),
    KEY idx_traffic_flows_location (location),
    KEY idx_traffic_flows_timestamp (timestamp)
) ENGINE=InnoDB;

-- 创建交通用户统计表
CREATE TABLE traffic_user_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    motor_count INT NOT NULL,
    non_motor_count INT NOT NULL,
    pedestrian_count INT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    location VARCHAR(100),
    KEY idx_traffic_user_stats_deleted_at (deleted_at),
    KEY idx_traffic_user_stats_timestamp (timestamp)
) ENGINE=InnoDB;

-- 创建交通热力图表
CREATE TABLE traffic_heatmaps (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    road_name VARCHAR(100) NOT NULL,
    congestion_level DECIMAL(3,2) NOT NULL,
    grid_x INT NOT NULL,
    grid_y INT NOT NULL,
    time_filter VARCHAR(20),
    timestamp TIMESTAMP NOT NULL,
    KEY idx_traffic_heatmaps_deleted_at (deleted_at),
    KEY idx_traffic_heatmaps_grid (grid_x, grid_y),
    KEY idx_traffic_heatmaps_time_filter (time_filter)
) ENGINE=InnoDB;

-- 创建拥堵播报表
CREATE TABLE congestion_reports (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    location VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    description VARCHAR(500),
    duration INT NOT NULL,
    report_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    KEY idx_congestion_reports_deleted_at (deleted_at),
    KEY idx_congestion_reports_status (status),
    KEY idx_congestion_reports_report_time (report_time)
) ENGINE=InnoDB;

-- 创建出入流量监控表
CREATE TABLE in_out_flow_data (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    location VARCHAR(100) NOT NULL,
    inbound_flow INT NOT NULL,
    outbound_flow INT NOT NULL,
    net_flow INT,
    timestamp TIMESTAMP NOT NULL,
    hour INT NOT NULL,
    KEY idx_in_out_flow_data_deleted_at (deleted_at),
    KEY idx_in_out_flow_data_timestamp (timestamp),
    KEY idx_in_out_flow_data_hour (hour)
) ENGINE=InnoDB;

-- 创建车辆通过率表
CREATE TABLE car_crossing_rates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    location VARCHAR(100) NOT NULL,
    total_count INT NOT NULL,
    passed_count INT NOT NULL,
    crossing_rate DECIMAL(5,2),
    timestamp TIMESTAMP NOT NULL,
    period VARCHAR(20),
    KEY idx_car_crossing_rates_deleted_at (deleted_at),
    KEY idx_car_crossing_rates_timestamp (timestamp)
) ENGINE=InnoDB;

-- ==================== 设备管理相关表 ====================

-- 创建设备表
CREATE TABLE devices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    device_name VARCHAR(100) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    location VARCHAR(100),
    status VARCHAR(20) DEFAULT 'normal',
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    manufacturer VARCHAR(100),
    install_date TIMESTAMP,
    KEY idx_devices_deleted_at (deleted_at),
    KEY idx_devices_type (device_type),
    KEY idx_devices_status (status)
) ENGINE=InnoDB;

-- 创建设备支出表
CREATE TABLE device_expenses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    device_id BIGINT UNSIGNED NOT NULL,
    expense_type VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2),
    description VARCHAR(255),
    expense_date TIMESTAMP,
    period VARCHAR(20),
    KEY idx_device_expenses_deleted_at (deleted_at),
    KEY idx_device_expenses_device_id (device_id),
    KEY idx_device_expenses_type (expense_type),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 创建设备维修记录表
CREATE TABLE device_maintenance_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    device_id BIGINT UNSIGNED NOT NULL,
    maintenance_type VARCHAR(50) NOT NULL,
    description VARCHAR(500),
    status VARCHAR(20) DEFAULT 'pending',
    technician VARCHAR(100),
    cost DECIMAL(10,2),
    start_time TIMESTAMP,
    end_time TIMESTAMP NULL,
    priority VARCHAR(20) DEFAULT 'medium',
    KEY idx_device_maintenance_records_deleted_at (deleted_at),
    KEY idx_device_maintenance_records_device_id (device_id),
    KEY idx_device_maintenance_records_status (status),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 创建设备故障统计表
CREATE TABLE device_fault_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    device_type VARCHAR(50) NOT NULL,
    fault_count INT NOT NULL,
    repair_count INT NOT NULL,
    fault_rate DECIMAL(5,2),
    stat_date TIMESTAMP,
    period VARCHAR(20),
    trend_direction VARCHAR(10),
    percentage_change DECIMAL(5,2),
    KEY idx_device_fault_stats_deleted_at (deleted_at),
    KEY idx_device_fault_stats_type (device_type),
    KEY idx_device_fault_stats_date (stat_date)
) ENGINE=InnoDB;

-- 创建设备报警表
CREATE TABLE device_alarms (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    device_id BIGINT UNSIGNED NOT NULL,
    alarm_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    message VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    alarm_time TIMESTAMP,
    ack_time TIMESTAMP NULL,
    resolve_time TIMESTAMP NULL,
    location VARCHAR(100),
    KEY idx_device_alarms_deleted_at (deleted_at),
    KEY idx_device_alarms_device_id (device_id),
    KEY idx_device_alarms_status (status),
    KEY idx_device_alarms_alarm_time (alarm_time),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 创建设备报警统计表
CREATE TABLE device_alarm_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    alarm_type VARCHAR(50) NOT NULL,
    count INT NOT NULL,
    hour INT NOT NULL,
    stat_date TIMESTAMP,
    severity VARCHAR(20),
    KEY idx_device_alarm_stats_deleted_at (deleted_at),
    KEY idx_device_alarm_stats_type (alarm_type),
    KEY idx_device_alarm_stats_date (stat_date)
) ENGINE=InnoDB;

-- ==================== 停车统计相关表 ====================

-- 创建停车饱和度表
CREATE TABLE parking_saturations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    parking_lot_id BIGINT UNSIGNED,
    saturation_rate DECIMAL(5,2),
    occupied_spots INT NOT NULL,
    total_spots INT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    hour INT NOT NULL,
    KEY idx_parking_saturations_deleted_at (deleted_at),
    KEY idx_parking_saturations_lot_id (parking_lot_id),
    KEY idx_parking_saturations_timestamp (timestamp),
    FOREIGN KEY (parking_lot_id) REFERENCES parking_lots(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 创建停车占用率表
CREATE TABLE parking_occupancy_rates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    parking_lot_id BIGINT UNSIGNED,
    occupancy_rate DECIMAL(5,2),
    period VARCHAR(20),
    timestamp TIMESTAMP NOT NULL,
    KEY idx_parking_occupancy_rates_deleted_at (deleted_at),
    KEY idx_parking_occupancy_rates_lot_id (parking_lot_id),
    KEY idx_parking_occupancy_rates_timestamp (timestamp),
    FOREIGN KEY (parking_lot_id) REFERENCES parking_lots(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 创建总占用率表
CREATE TABLE total_occupancies (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    total_spots INT NOT NULL,
    occupied_spots INT NOT NULL,
    occupancy_rate DECIMAL(5,2),
    available_spots INT,
    timestamp TIMESTAMP NOT NULL,
    period VARCHAR(20),
    KEY idx_total_occupancies_deleted_at (deleted_at),
    KEY idx_total_occupancies_timestamp (timestamp)
) ENGINE=InnoDB;

-- 创建非机动车停车拥堵表
CREATE TABLE motor_parking_congestions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    location VARCHAR(100) NOT NULL,
    congestion_level VARCHAR(20) NOT NULL,
    motor_count INT NOT NULL,
    capacity INT NOT NULL,
    congestion_rate DECIMAL(5,2),
    timestamp TIMESTAMP NOT NULL,
    reported_by VARCHAR(100),
    KEY idx_motor_parking_congestions_deleted_at (deleted_at),
    KEY idx_motor_parking_congestions_location (location),
    KEY idx_motor_parking_congestions_timestamp (timestamp)
) ENGINE=InnoDB;

-- 创建收费记录表
CREATE TABLE toll_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    vehicle_id BIGINT UNSIGNED,
    plate_number VARCHAR(30) NOT NULL,
    toll_station VARCHAR(100) NOT NULL,
    entry_time TIMESTAMP,
    exit_time TIMESTAMP NULL,
    amount DECIMAL(10,2),
    payment_method VARCHAR(50),
    status VARCHAR(20) DEFAULT 'completed',
    distance DECIMAL(10,2),
    KEY idx_toll_records_deleted_at (deleted_at),
    KEY idx_toll_records_vehicle_id (vehicle_id),
    KEY idx_toll_records_plate_number (plate_number),
    KEY idx_toll_records_entry_time (entry_time),
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 创建施工统计表
CREATE TABLE construction_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    project_name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP NULL,
    status VARCHAR(20) NOT NULL,
    progress DECIMAL(5,2),
    impact_level VARCHAR(20),
    traffic_impact VARCHAR(500),
    budget DECIMAL(12,2),
    actual_cost DECIMAL(12,2),
    KEY idx_construction_stats_deleted_at (deleted_at),
    KEY idx_construction_stats_status (status),
    KEY idx_construction_stats_location (location)
) ENGINE=InnoDB;

-- 创建监控摄像头表
CREATE TABLE monitoring_cameras (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    camera_name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    status VARCHAR(20) DEFAULT 'online',
    stream_url VARCHAR(255),
    type VARCHAR(50),
    resolution VARCHAR(20),
    view_angle INT,
    night_vision BOOLEAN DEFAULT FALSE,
    install_date TIMESTAMP,
    KEY idx_monitoring_cameras_deleted_at (deleted_at),
    KEY idx_monitoring_cameras_location (location),
    KEY idx_monitoring_cameras_status (status)
) ENGINE=InnoDB;

-- ==================== 插入模拟数据 ====================

-- 插入交通流量数据
INSERT INTO traffic_flows (location, flow_count, speed, direction, vehicle_type, timestamp) VALUES
('一层北', 1250, 45.5, 'inbound', '小型汽车', NOW() - INTERVAL 1 HOUR),
('一层东', 880, 35.2, 'outbound', '小型汽车', NOW() - INTERVAL 1 HOUR),
('二层南', 950, 42.1, 'inbound', '小型汽车', NOW() - INTERVAL 1 HOUR),
('地下一层', 1100, 25.8, 'inbound', '小型汽车', NOW() - INTERVAL 30 MINUTE),
('东门入口', 1350, 50.2, 'inbound', '混合车型', NOW() - INTERVAL 15 MINUTE),
('西门出口', 920, 38.7, 'outbound', '混合车型', NOW() - INTERVAL 5 MINUTE);

-- 插入交通用户统计数据
INSERT INTO traffic_user_stats (motor_count, non_motor_count, pedestrian_count, timestamp, location) VALUES
(65, 22, 18, NOW() - INTERVAL 30 MINUTE, '市中心区域'),
(80, 30, 25, NOW() - INTERVAL 25 MINUTE, '市中心区域'),
(75, 40, 20, NOW() - INTERVAL 20 MINUTE, '市中心区域'),
(50, 15, 10, NOW() - INTERVAL 15 MINUTE, '市中心区域'),
(60, 25, 15, NOW() - INTERVAL 10 MINUTE, '市中心区域'),
(70, 20, 12, NOW() - INTERVAL 5 MINUTE, '市中心区域'),
(85, 35, 20, NOW(), '市中心区域');

-- 插入交通热力图数据
INSERT INTO traffic_heatmaps (road_name, congestion_level, grid_x, grid_y, time_filter, timestamp) VALUES
('中山北路', 0.75, 0, 0, 'now', NOW()),
('南京路', 0.60, 0, 1, 'now', NOW()),
('延安高架', 0.85, 0, 2, 'now', NOW()),
('内环路', 0.45, 0, 3, 'now', NOW()),
('外环路', 0.30, 0, 4, 'now', NOW()),
('浦东大道', 0.65, 1, 0, 'now', NOW()),
('世纪大道', 0.70, 1, 1, 'now', NOW()),
('陆家嘴环路', 0.55, 1, 2, 'now', NOW()),
('北京路', 0.40, 1, 3, 'now', NOW()),
('四川路', 0.50, 1, 4, 'now', NOW());

-- 插入拥堵播报数据
INSERT INTO congestion_reports (location, severity, description, duration, report_time, status) VALUES
('延安高架路中段', '严重', '因施工作业导致车道缩减，车辆通行缓慢', 120, NOW() - INTERVAL 2 HOUR, 'active'),
('南京路与淮海路交叉口', '轻微', '红绿灯故障，交警现场指挥交通', 45, NOW() - INTERVAL 1 HOUR, 'active'),
('外环路北段', '中等', '多车追尾事故，占用两个车道', 90, NOW() - INTERVAL 30 MINUTE, 'active'),
('人民广场附近', '轻微', '大型活动结束，人流车流量增大', 60, NOW() - INTERVAL 15 MINUTE, 'active'),
('虹桥机场高速', '严重', '恶劣天气影响，能见度低车速缓慢', 180, NOW() - INTERVAL 10 MINUTE, 'active');

-- 插入出入流量数据
INSERT INTO in_out_flow_data (location, inbound_flow, outbound_flow, net_flow, timestamp, hour) VALUES
('主要入口A', 150, 120, 30, NOW() - INTERVAL 1 HOUR, 8),
('主要入口B', 200, 180, 20, NOW() - INTERVAL 1 HOUR, 8),
('次要入口C', 80, 75, 5, NOW() - INTERVAL 1 HOUR, 8),
('主要入口A', 180, 140, 40, NOW() - INTERVAL 30 MINUTE, 9),
('主要入口B', 220, 200, 20, NOW() - INTERVAL 30 MINUTE, 9),
('次要入口C', 100, 90, 10, NOW() - INTERVAL 30 MINUTE, 9);

-- 插入车辆通过率数据
INSERT INTO car_crossing_rates (location, total_count, passed_count, crossing_rate, timestamp, period) VALUES
('十字路口A', 120, 108, 90.00, NOW() - INTERVAL 1 HOUR, 'peak'),
('十字路口B', 150, 135, 90.00, NOW() - INTERVAL 1 HOUR, 'peak'),
('十字路口C', 80, 68, 85.00, NOW() - INTERVAL 1 HOUR, 'peak'),
('十字路口A', 100, 92, 92.00, NOW() - INTERVAL 30 MINUTE, 'normal'),
('十字路口B', 110, 99, 90.00, NOW() - INTERVAL 30 MINUTE, 'normal'),
('十字路口C', 90, 81, 90.00, NOW() - INTERVAL 30 MINUTE, 'normal');

-- 插入设备数据
INSERT INTO devices (device_name, device_type, location, status, serial_number, manufacturer, install_date) VALUES
('摄像头-001', '监控设备', '一层北入口', 'normal', 'CAM001', '海康威视', '2024-01-15'),
('摄像头-002', '监控设备', '一层东出口', 'normal', 'CAM002', '海康威视', '2024-01-20'),
('感应器-001', '车位检测', '二层A区', 'normal', 'SEN001', '博世', '2024-02-01'),
('感应器-002', '车位检测', '二层B区', 'fault', 'SEN002', '博世', '2024-02-05'),
('LED显示屏-001', '信息显示', '入口处', 'normal', 'LED001', '利亚德', '2024-01-10'),
('充电桩-001', '充电设备', '地下一层', 'normal', 'CHG001', '特来电', '2024-03-01'),
('充电桩-002', '充电设备', '地下一层', 'maintenance', 'CHG002', '特来电', '2024-03-05'),
('道闸-001', '出入控制', '主入口', 'normal', 'GATE001', '捷顺', '2024-01-05');

-- 插入设备支出数据
INSERT INTO device_expenses (device_id, expense_type, amount, description, expense_date, period) VALUES
(1, '维护费用', 500.00, '摄像头清洁和调试', '2025-08-01', 'month'),
(2, '维护费用', 450.00, '摄像头镜头更换', '2025-08-03', 'month'),
(3, '电费', 120.00, '车位检测器电费', '2025-08-05', 'month'),
(4, '维修费用', 800.00, '感应器电路板更换', '2025-08-07', 'month'),
(5, '维护费用', 300.00, 'LED屏幕清洁', '2025-08-10', 'month'),
(6, '电费', 2500.00, '充电桩电费', '2025-08-01', 'month'),
(7, '维修费用', 1200.00, '充电桩充电模块维修', '2025-08-08', 'month'),
(8, '维护费用', 200.00, '道闸润滑保养', '2025-08-12', 'month');

-- 插入设备维修记录
INSERT INTO device_maintenance_records (device_id, maintenance_type, description, status, technician, cost, start_time, end_time, priority) VALUES
(4, '故障维修', '感应器无法正常检测车辆，需更换传感器', 'completed', '张工程师', 800.00, '2025-08-07 09:00:00', '2025-08-07 11:30:00', 'high'),
(7, '定期维护', '充电桩定期检查和保养', 'in_progress', '李技师', 300.00, '2025-08-12 08:00:00', NULL, 'medium'),
(2, '预防维护', '摄像头镜头清洁和角度调整', 'completed', '王师傅', 150.00, '2025-08-10 14:00:00', '2025-08-10 15:00:00', 'low'),
(5, '故障维修', 'LED显示屏部分像素不亮', 'pending', '刘技师', 600.00, '2025-08-13 10:00:00', NULL, 'medium'),
(6, '定期维护', '充电桩电气系统检查', 'completed', '陈工程师', 400.00, '2025-08-11 16:00:00', '2025-08-11 18:00:00', 'medium');

-- 插入设备故障统计
INSERT INTO device_fault_stats (device_type, fault_count, repair_count, fault_rate, stat_date, period, trend_direction, percentage_change) VALUES
('监控设备', 2, 2, 12.50, CURDATE(), 'daily', 'down', -5.2),
('车位检测', 3, 2, 25.00, CURDATE(), 'daily', 'up', 8.3),
('信息显示', 1, 0, 8.33, CURDATE(), 'daily', 'stable', 0.0),
('充电设备', 1, 1, 25.00, CURDATE(), 'daily', 'down', -10.5),
('出入控制', 0, 0, 0.00, CURDATE(), 'daily', 'stable', 0.0);

-- 插入设备报警
INSERT INTO device_alarms (device_id, alarm_type, severity, message, status, alarm_time, location) VALUES
(4, '设备故障', '高', '车位检测器响应异常，可能存在硬件故障', 'active', NOW() - INTERVAL 2 HOUR, '二层B区'),
(7, '设备维护', '中', '充电桩需要定期维护检查', 'acknowledged', NOW() - INTERVAL 1 HOUR, '地下一层'),
(5, '通信异常', '中', 'LED显示屏与控制中心通信中断', 'active', NOW() - INTERVAL 30 MINUTE, '入口处'),
(2, '环境异常', '低', '摄像头检测到异常遮挡', 'resolved', NOW() - INTERVAL 3 HOUR, '一层东出口'),
(6, '电源异常', '高', '充电桩电源电压不稳定', 'active', NOW() - INTERVAL 15 MINUTE, '地下一层');

-- 插入设备报警统计
INSERT INTO device_alarm_stats (alarm_type, count, hour, stat_date, severity) VALUES
('设备故障', 3, 8, CURDATE(), '高'),
('设备维护', 2, 9, CURDATE(), '中'),
('通信异常', 1, 10, CURDATE(), '中'),
('环境异常', 1, 11, CURDATE(), '低'),
('电源异常', 2, 12, CURDATE(), '高'),
('设备故障', 1, 13, CURDATE(), '中'),
('设备维护', 1, 14, CURDATE(), '低');

-- 插入停车饱和度数据
INSERT INTO parking_saturations (parking_lot_id, saturation_rate, occupied_spots, total_spots, timestamp, hour) VALUES
(1, 79.17, 95, 120, NOW() - INTERVAL 1 HOUR, 8),
(2, 79.00, 158, 200, NOW() - INTERVAL 1 HOUR, 8),
(3, 93.75, 75, 80, NOW() - INTERVAL 1 HOUR, 8),
(4, 60.00, 90, 150, NOW() - INTERVAL 1 HOUR, 8),
(5, 86.67, 78, 90, NOW() - INTERVAL 1 HOUR, 8),
(6, 86.67, 52, 60, NOW() - INTERVAL 1 HOUR, 8);

-- 插入停车占用率数据
INSERT INTO parking_occupancy_rates (parking_lot_id, occupancy_rate, period, timestamp) VALUES
(1, 75.25, 'morning', NOW() - INTERVAL 2 HOUR),
(2, 82.50, 'morning', NOW() - INTERVAL 2 HOUR),
(3, 90.00, 'morning', NOW() - INTERVAL 2 HOUR),
(4, 65.33, 'morning', NOW() - INTERVAL 2 HOUR),
(5, 85.56, 'morning', NOW() - INTERVAL 2 HOUR),
(6, 80.00, 'morning', NOW() - INTERVAL 2 HOUR);

-- 插入总占用率数据
INSERT INTO total_occupancies (total_spots, occupied_spots, occupancy_rate, available_spots, timestamp, period) VALUES
(700, 548, 78.29, 152, NOW() - INTERVAL 1 HOUR, 'hourly'),
(700, 562, 80.29, 138, NOW() - INTERVAL 30 MINUTE, 'hourly'),
(700, 575, 82.14, 125, NOW(), 'hourly');

-- 插入非机动车停车拥堵数据
INSERT INTO motor_parking_congestions (location, congestion_level, motor_count, capacity, congestion_rate, timestamp, reported_by) VALUES
('地铁站A出口', '中等', 45, 60, 75.00, NOW() - INTERVAL 1 HOUR, '巡查员001'),
('商业街B段', '严重', 80, 90, 88.89, NOW() - INTERVAL 30 MINUTE, '巡查员002'),
('学校门口C区', '轻微', 25, 50, 50.00, NOW() - INTERVAL 15 MINUTE, '巡查员003'),
('公园入口D处', '中等', 35, 45, 77.78, NOW() - INTERVAL 45 MINUTE, '巡查员001');

-- 插入收费记录数据
INSERT INTO toll_records (vehicle_id, plate_number, toll_station, entry_time, exit_time, amount, payment_method, status, distance) VALUES
(1, '京A12345', '杭州西收费站', '2025-08-12 08:30:00', '2025-08-12 09:15:00', 15.00, '微信支付', 'completed', 35.5),
(2, '京B67890', '杭州东收费站', '2025-08-12 07:45:00', '2025-08-12 08:30:00', 12.00, '支付宝', 'completed', 28.2),
(1, '京A12345', '萧山收费站', '2025-08-11 16:20:00', '2025-08-11 17:05:00', 18.00, 'ETC', 'completed', 42.8),
(2, '京B67890', '余杭收费站', '2025-08-11 14:10:00', '2025-08-11 15:00:00', 20.00, '现金', 'completed', 48.3);

-- 插入施工统计数据
INSERT INTO construction_stats (project_name, location, start_date, end_date, status, progress, impact_level, traffic_impact, budget, actual_cost) VALUES
('延安高架路面修复工程', '延安高架中段', '2025-08-01', '2025-08-20', 'in_progress', 65.00, '高', '车道缩减，通行速度降低50%', 500000.00, 320000.00),
('地铁5号线建设', '人民广场地下', '2025-07-15', '2025-12-30', 'in_progress', 25.00, '中', '部分路段封闭，绕行增加10分钟', 2000000.00, 1200000.00),
('智能信号灯升级', '主要路口', '2025-08-10', '2025-08-25', 'in_progress', 40.00, '低', '夜间施工，白天正常通行', 150000.00, 95000.00),
('桥梁加固工程', '外滩大桥', '2025-07-01', '2025-08-15', 'completed', 100.00, '高', '单向通行，通行能力减半', 800000.00, 750000.00);

-- 插入监控摄像头数据
INSERT INTO monitoring_cameras (camera_name, location, latitude, longitude, status, stream_url, type, resolution, view_angle, night_vision, install_date) VALUES
('入口监控-001', '主入口A', 30.2594, 120.1644, 'online', 'rtmp://192.168.1.101/live/cam001', '球机', '1080P', 360, TRUE, '2024-01-15'),
('出口监控-002', '主出口B', 30.2595, 120.1645, 'online', 'rtmp://192.168.1.102/live/cam002', '枪机', '4K', 90, TRUE, '2024-01-20'),
('区域监控-003', '二层停车区', 30.2596, 120.1646, 'online', 'rtmp://192.168.1.103/live/cam003', '半球', '1080P', 180, FALSE, '2024-02-01'),
('通道监控-004', '地下通道', 30.2597, 120.1647, 'maintenance', 'rtmp://192.168.1.104/live/cam004', '筒机', '720P', 60, TRUE, '2024-02-10'),
('周界监控-005', '停车场周界', 30.2598, 120.1648, 'online', 'rtmp://192.168.1.105/live/cam005', '球机', '4K', 360, TRUE, '2024-03-01'),
('应急监控-006', '应急通道', 30.2599, 120.1649, 'offline', 'rtmp://192.168.1.106/live/cam006', '枪机', '1080P', 120, FALSE, '2024-03-15');

-- 停车会话表建表语句
CREATE TABLE IF NOT EXISTS `parking_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `vehicle_id` bigint unsigned NOT NULL,
  `parking_lot_id` bigint unsigned NOT NULL,
  `spot_code` varchar(20) NOT NULL COMMENT '车位编号',
  `spot_type` varchar(20) DEFAULT 'normal' COMMENT '车位类型：normal,charging,disabled,vip',
  `start_time` datetime(3) NOT NULL COMMENT '停车开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '停车结束时间',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '状态：active,ended,paid',
  `fee_rate` decimal(10,2) NOT NULL COMMENT '每小时费率',
  `fee_current` decimal(10,2) DEFAULT '0.00' COMMENT '当前费用',
  `next_billing_time` datetime(3) DEFAULT NULL COMMENT '下次计费时间',
  `next_fee_amount` decimal(10,2) DEFAULT NULL COMMENT '下次计费金额',
  `current_billing_cycle` int DEFAULT '0' COMMENT '当前计费周期',
  `pricing_rule` varchar(100) DEFAULT '首小时10元，后续每小时5元' COMMENT '计费规则描述',
  `navigation_status` varchar(20) DEFAULT 'en_route' COMMENT '导航状态：en_route,in_garage,parked',
  `remaining_distance_m` int DEFAULT '0' COMMENT '剩余距离（米）',
  `estimated_minutes` int DEFAULT '0' COMMENT '预计到达时间（分钟）',
  `destination_lat` decimal(10,8) DEFAULT NULL COMMENT '目的地纬度',
  `destination_lon` decimal(11,8) DEFAULT NULL COMMENT '目的地经度',
  `user_position_lat` decimal(10,8) DEFAULT NULL COMMENT '用户当前纬度',
  `user_position_lon` decimal(11,8) DEFAULT NULL COMMENT '用户当前经度',
  `progress_to_destination_percent` int DEFAULT '0' COMMENT '到达目的地进度百分比',
  PRIMARY KEY (`id`),
  KEY `idx_parking_sessions_deleted_at` (`deleted_at`),
  KEY `idx_parking_sessions_user_id` (`user_id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_parking_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_parking_sessions_vehicle` FOREIGN KEY (`vehicle_id`) REFERENCES `vehicles` (`id`),
  CONSTRAINT `fk_parking_sessions_parking_lot` FOREIGN KEY (`parking_lot_id`) REFERENCES `parking_lots` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='停车会话表';

-- 插入测试数据
INSERT INTO `parking_sessions` (`created_at`, `updated_at`, `user_id`, `vehicle_id`, `parking_lot_id`, `spot_code`, `spot_type`, `start_time`, `end_time`, `status`, `fee_rate`, `fee_current`, `next_billing_time`, `next_fee_amount`, `current_billing_cycle`, `pricing_rule`, `navigation_status`, `remaining_distance_m`, `estimated_minutes`, `destination_lat`, `destination_lon`, `user_position_lat`, `user_position_lon`, `progress_to_destination_percent`) VALUES
(NOW(), NOW(), 2, 1, 1, 'A-101', 'normal', DATE_SUB(NOW(), INTERVAL 2 HOUR), NULL, 'active', 10.00, 10.00, DATE_ADD(NOW(), INTERVAL 58 MINUTE), 5.00, 1, '首小时10元，后续每小时5元', 'parked', 0, 0, 30.23940000, 120.15090000, 30.23940000, 120.15090000, 100);

-- 插入一些历史停车记录
INSERT INTO `parking_sessions` (`created_at`, `updated_at`, `user_id`, `vehicle_id`, `parking_lot_id`, `spot_code`, `spot_type`, `start_time`, `end_time`, `status`, `fee_rate`, `fee_current`, `next_billing_time`, `next_fee_amount`, `current_billing_cycle`, `pricing_rule`, `navigation_status`, `remaining_distance_m`, `estimated_minutes`, `destination_lat`, `destination_lon`, `user_position_lat`, `user_position_lon`, `progress_to_destination_percent`) VALUES
(DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 2, 1, 1, 'B-205', 'charging', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(DATE_SUB(NOW(), INTERVAL 1 DAY), INTERVAL -3 HOUR), 'ended', 12.00, 25.00, NULL, NULL, 3, '充电车位首小时12元，后续每小时8元', 'parked', 0, 0, 30.23940000, 120.15090000, 30.23940000, 120.15090000, 100),
(DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 2, 1, 2, 'C-108', 'normal', DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(DATE_SUB(NOW(), INTERVAL 3 DAY), INTERVAL -2 HOUR), 'ended', 8.00, 16.00, NULL, NULL, 2, '普通车位首小时8元，后续每小时6元', 'parked', 0, 0, 30.26400000, 120.15500000, 30.26400000, 120.15500000, 100),
(DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), 2, 1, 3, 'A-301', 'vip', DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(DATE_SUB(NOW(), INTERVAL 5 DAY), INTERVAL -4 HOUR), 'ended', 15.00, 45.00, NULL, NULL, 4, 'VIP车位首小时15元，后续每小时10元', 'parked', 0, 0, 30.24000000, 120.13000000, 30.24000000, 120.13000000, 100),
(DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), 2, 1, 1, 'A-105', 'normal', DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(DATE_SUB(NOW(), INTERVAL 7 DAY), INTERVAL -1 HOUR), 'ended', 10.00, 10.00, NULL, NULL, 1, '首小时10元，后续每小时5元', 'parked', 0, 0, 30.23940000, 120.15090000, 30.23940000, 120.15090000, 100);

-- 为演示导航功能，插入一个即将到达的停车会话（可选）
INSERT INTO `parking_sessions` (`created_at`, `updated_at`, `user_id`, `vehicle_id`, `parking_lot_id`, `spot_code`, `spot_type`, `start_time`, `end_time`, `status`, `fee_rate`, `fee_current`, `next_billing_time`, `next_fee_amount`, `current_billing_cycle`, `pricing_rule`, `navigation_status`, `remaining_distance_m`, `estimated_minutes`, `destination_lat`, `destination_lon`, `user_position_lat`, `user_position_lon`, `progress_to_destination_percent`) VALUES
(NOW(), NOW(), 2, 1, 1, 'A-102', 'normal', NOW(), NULL, 'active', 10.00, 0.00, DATE_ADD(NOW(), INTERVAL 1 HOUR), 10.00, 0, '首小时10元，后续每小时5元', 'en_route', 500, 8, 30.23940000, 120.15090000, 30.23500000, 120.14800000, 25);

-- 如果需要，可以删除重复的活跃会话，只保留一个
DELETE FROM `parking_sessions` 
WHERE `user_id` = 2 AND `status` = 'active' AND `id` NOT IN (
    SELECT * FROM (
        SELECT MIN(`id`) FROM `parking_sessions` 
        WHERE `user_id` = 2 AND `status` = 'active'
    ) AS temp
);

-- 创建交通流量数据表
CREATE TABLE `traffic_flows` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `location` VARCHAR(100) NOT NULL COMMENT '路段位置',
  `flow_count` INT NOT NULL COMMENT '车流量',
  `speed` DECIMAL(5,2) DEFAULT '0.00' COMMENT '平均速度 km/h',
  `direction` VARCHAR(20) DEFAULT NULL COMMENT '方向 inbound/outbound',
  `vehicle_type` VARCHAR(30) DEFAULT NULL COMMENT '车辆类型',
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '数据时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_traffic_flows_deleted_at` (`deleted_at`),
  KEY `idx_traffic_flows_location` (`location`),
  KEY `idx_traffic_flows_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='交通流量数据表';

-- 创建交通用户统计表
CREATE TABLE `traffic_user_stats` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `motor_count` INT NOT NULL DEFAULT '0' COMMENT '机动车数量',
  `non_motor_count` INT NOT NULL DEFAULT '0' COMMENT '非机动车数量',
  `pedestrian_count` INT NOT NULL DEFAULT '0' COMMENT '行人数量',
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '统计时间',
  `location` VARCHAR(100) DEFAULT NULL COMMENT '统计位置',
  PRIMARY KEY (`id`),
  KEY `idx_traffic_user_stats_deleted_at` (`deleted_at`),
  KEY `idx_traffic_user_stats_timestamp` (`timestamp`),
  KEY `idx_traffic_user_stats_location` (`location`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='交通用户统计表';

-- 创建交通热力图数据表
CREATE TABLE `traffic_heatmaps` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `road_name` VARCHAR(100) NOT NULL COMMENT '道路名称',
  `congestion_level` DECIMAL(3,2) DEFAULT '0.00' COMMENT '拥堵程度 0-1',
  `grid_x` INT NOT NULL COMMENT '网格X坐标',
  `grid_y` INT NOT NULL COMMENT '网格Y坐标',
  `time_filter` VARCHAR(20) DEFAULT NULL COMMENT '时间过滤器',
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '数据时间',
  PRIMARY KEY (`id`),
  KEY `idx_traffic_heatmaps_deleted_at` (`deleted_at`),
  KEY `idx_traffic_heatmaps_grid` (`grid_x`, `grid_y`),
  KEY `idx_traffic_heatmaps_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='交通热力图数据表';

-- 创建车辆通过率数据表
CREATE TABLE `car_crossing_rates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `location` VARCHAR(100) NOT NULL COMMENT '检测位置',
  `crossing_rate` DECIMAL(5,2) DEFAULT '0.00' COMMENT '通过率',
  `total_vehicles` INT DEFAULT '0' COMMENT '总车辆数',
  `passed_vehicles` INT DEFAULT '0' COMMENT '通过车辆数',
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '统计时间',
  PRIMARY KEY (`id`),
  KEY `idx_car_crossing_rates_deleted_at` (`deleted_at`),
  KEY `idx_car_crossing_rates_location` (`location`),
  KEY `idx_car_crossing_rates_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='车辆通过率数据表';

-- 创建拥堵播报表
CREATE TABLE `congestion_reports` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `road_name` VARCHAR(100) NOT NULL COMMENT '道路名称',
  `report_content` TEXT COMMENT '播报内容',
  `severity_level` VARCHAR(20) DEFAULT 'medium' COMMENT '严重程度: low, medium, high',
  `status` VARCHAR(20) DEFAULT 'active' COMMENT '状态: active, resolved',
  `report_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '播报时间',
  `location_lat` DECIMAL(10,8) DEFAULT NULL COMMENT '位置纬度',
  `location_lon` DECIMAL(11,8) DEFAULT NULL COMMENT '位置经度',
  PRIMARY KEY (`id`),
  KEY `idx_congestion_reports_deleted_at` (`deleted_at`),
  KEY `idx_congestion_reports_status` (`status`),
  KEY `idx_congestion_reports_time` (`report_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='拥堵播报表';

-- 插入交通流量测试数据
INSERT INTO `traffic_flows` (`location`, `flow_count`, `speed`, `direction`, `vehicle_type`, `timestamp`) VALUES
('一层北入口', 1250, 35.5, 'inbound', '小型汽车', NOW()),
('一层东入口', 880, 42.1, 'inbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 10 MINUTE)),
('一层南入口', 950, 28.3, 'inbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 20 MINUTE)),
('负一层北入口', 1100, 31.2, 'inbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 30 MINUTE)),
('负一层东入口', 720, 45.8, 'inbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 40 MINUTE)),
('一层北出口', 890, 38.7, 'outbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 5 MINUTE)),
('一层东出口', 650, 41.3, 'outbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 15 MINUTE)),
('一层南出口', 780, 33.9, 'outbound', '小型汽车', DATE_SUB(NOW(), INTERVAL 25 MINUTE));

-- 插入交通用户统计测试数据
INSERT INTO `traffic_user_stats` (`motor_count`, `non_motor_count`, `pedestrian_count`, `location`, `timestamp`) VALUES
(65, 22, 18, '主要路口', NOW()),
(50, 15, 10, '主要路口', DATE_SUB(NOW(), INTERVAL 5 MINUTE)),
(75, 40, 25, '主要路口', DATE_SUB(NOW(), INTERVAL 10 MINUTE)),
(80, 30, 20, '主要路口', DATE_SUB(NOW(), INTERVAL 15 MINUTE)),
(60, 25, 15, '主要路口', DATE_SUB(NOW(), INTERVAL 20 MINUTE)),
(15, 8, 5, '主要路口', DATE_SUB(NOW(), INTERVAL 25 MINUTE)),
(20, 5, 3, '主要路口', DATE_SUB(NOW(), INTERVAL 30 MINUTE));

-- 插入交通热力图测试数据
INSERT INTO `traffic_heatmaps` (`road_name`, `congestion_level`, `grid_x`, `grid_y`, `time_filter`, `timestamp`) VALUES
('主干道A段', 0.8, 0, 0, 'now', NOW()),
('主干道B段', 0.6, 0, 1, 'now', NOW()),
('支路C段', 0.3, 1, 0, 'now', NOW()),
('支路D段', 0.9, 1, 1, 'now', NOW()),
('环路E段', 0.5, 2, 0, 'now', NOW()),
('环路F段', 0.7, 2, 1, 'now', NOW());

-- 插入车辆通过率测试数据
INSERT INTO `car_crossing_rates` (`location`, `crossing_rate`, `total_vehicles`, `passed_vehicles`, `timestamp`) VALUES
('路口A', 85.50, 200, 171, NOW()),
('路口B', 92.30, 150, 138, DATE_SUB(NOW(), INTERVAL 1 HOUR)),
('路口C', 78.20, 180, 141, DATE_SUB(NOW(), INTERVAL 2 HOUR)),
('路口D', 88.70, 220, 195, DATE_SUB(NOW(), INTERVAL 3 HOUR));

-- 插入拥堵播报测试数据
INSERT INTO `congestion_reports` (`road_name`, `report_content`, `severity_level`, `status`, `report_time`, `location_lat`, `location_lon`) VALUES
('西湖大道', '西湖大道主干道出现轻微拥堵，预计15分钟后缓解', 'medium', 'active', NOW(), 30.2394, 120.1509),
('中山路', '中山路与建国路交叉口信号灯故障，交通缓行', 'high', 'active', DATE_SUB(NOW(), INTERVAL 30 MINUTE), 30.2635, 120.1709),
('延安路', '延安路施工路段车辆缓行，建议绕行', 'low', 'resolved', DATE_SUB(NOW(), INTERVAL 1 HOUR), 30.2793, 120.2194);

-- 创建空气质量数据表
CREATE TABLE `air_qualities` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `location` VARCHAR(100) NOT NULL COMMENT '监测位置',
  `aqi` INT NOT NULL COMMENT '空气质量指数',
  `level` VARCHAR(20) NOT NULL COMMENT '空气质量等级',
  `pm25` DECIMAL(5,2) NOT NULL COMMENT 'PM2.5浓度',
  `pm10` DECIMAL(5,2) NOT NULL COMMENT 'PM10浓度',
  `o3` DECIMAL(5,2) NOT NULL COMMENT '臭氧浓度',
  `no2` DECIMAL(5,2) NOT NULL COMMENT '二氧化氮浓度',
  `so2` DECIMAL(5,2) NOT NULL COMMENT '二氧化硫浓度',
  `co` DECIMAL(5,2) NOT NULL COMMENT '一氧化碳浓度',
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '数据时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_air_qualities_deleted_at` (`deleted_at`),
  KEY `idx_air_qualities_location` (`location`),
  KEY `idx_air_qualities_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='空气质量数据表';

-- 创建空气质量统计表
CREATE TABLE `air_quality_stats` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `date` DATE NOT NULL COMMENT '统计日期',
  `avg_aqi` DECIMAL(5,2) NOT NULL COMMENT '平均AQI',
  `max_aqi` INT NOT NULL COMMENT '最高AQI',
  `min_aqi` INT NOT NULL COMMENT '最低AQI',
  `primary_pollutant` VARCHAR(20) DEFAULT NULL COMMENT '主要污染物',
  `good_hours` INT DEFAULT '0' COMMENT '优良小时数',
  `polluted_hours` INT DEFAULT '0' COMMENT '污染小时数',
  PRIMARY KEY (`id`),
  KEY `idx_air_quality_stats_deleted_at` (`deleted_at`),
  KEY `idx_air_quality_stats_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='空气质量统计表';

-- 插入空气质量测试数据
INSERT INTO `air_qualities` (`location`, `aqi`, `level`, `pm25`, `pm10`, `o3`, `no2`, `so2`, `co`, `timestamp`) VALUES
('停车场周边', 75, '良', 28.0, 55.0, 90.0, 30.0, 10.0, 0.8, NOW()),
('停车场周边', 68, '良', 25.0, 52.0, 85.0, 28.0, 8.0, 0.7, DATE_SUB(NOW(), INTERVAL 1 HOUR)),
('停车场周边', 82, '良', 32.0, 58.0, 95.0, 35.0, 12.0, 0.9, DATE_SUB(NOW(), INTERVAL 2 HOUR)),
('停车场周边', 71, '良', 26.0, 48.0, 88.0, 25.0, 9.0, 0.6, DATE_SUB(NOW(), INTERVAL 3 HOUR)),
('停车场周边', 89, '良', 35.0, 62.0, 98.0, 38.0, 15.0, 1.0, DATE_SUB(NOW(), INTERVAL 4 HOUR)),
('停车场周边', 65, '良', 22.0, 45.0, 82.0, 22.0, 7.0, 0.5, DATE_SUB(NOW(), INTERVAL 5 HOUR)),
('停车场周边', 78, '良', 29.0, 56.0, 92.0, 32.0, 11.0, 0.8, DATE_SUB(NOW(), INTERVAL 6 HOUR));

-- 插入空气质量统计测试数据
INSERT INTO `air_quality_stats` (`date`, `avg_aqi`, `max_aqi`, `min_aqi`, `primary_pollutant`, `good_hours`, `polluted_hours`) VALUES
(CURDATE(), 75.5, 98, 52, 'PM2.5', 18, 6),
(DATE_SUB(CURDATE(), INTERVAL 1 DAY), 82.3, 105, 65, 'PM10', 16, 8),
(DATE_SUB(CURDATE(), INTERVAL 2 DAY), 69.8, 88, 45, 'O3', 20, 4),
(DATE_SUB(CURDATE(), INTERVAL 3 DAY), 91.2, 125, 72, 'PM2.5', 14, 10),
(DATE_SUB(CURDATE(), INTERVAL 4 DAY), 67.5, 85, 48, 'NO2', 22, 2),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 73.8, 92, 58, 'PM10', 19, 5),
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 88.4, 118, 69, 'PM2.5', 15, 9);

