
## redis

### 登录验证相关

predix | description
--- | ---
`username` | 用户名，标记是否已注册
`token` | 用户登录标记

### 业务相关
predix | description
--- | ---
`scene` | 用户信息


## 数据库
### 表
```mysql
CREATE TABLE `scenes` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `created_at` int NOT NULL DEFAULT 0,
    `updated_at` int NOT NULL DEFAULT 0,
    `deleted_at` int DEFAULT 0,
    `scene_id` bigint NOT NULL,
    `parent_scene_id` bigint NOT NULL,
    `choose_content` varchar(255) NOT NULL,
    `creator` varchar(255) NOT NULL,
    `cos_url` varchar(255) NOT NULL,
    `short_desc` varchar(600) NOT NULL,
    `is_init` tinyint NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE INDEX (`scene_id`),
    INDEX (`parent_scene_id`)
);
```

```mysql
CREATE TABLE `background_imgs` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `created_at` int NOT NULL DEFAULT 0,
    `updated_at` int NOT NULL DEFAULT 0,
    `deleted_at` int DEFAULT 0,
    `img_id` BIGINT,
    `cos_url` VARCHAR(255),
    `scene_id` BIGINT,
    `prompt` VARCHAR(255),
    INDEX idx_background_imgs_img_id (`img_id`),
    INDEX idx_background_imgs_scene_id (`scene_id`)
);
```

初始化

```mysql
INSERT INTO `scenes` (scene_id, choose_content, creator, parent_scene_id, cos_url, short_desc) VALUES (0, '-', 'admin123456', -1, 'https://fake-buddha-1300084664.cos.ap-shanghai.myqcloud.com/scene%2F491213694122852611.txt','我是一家三流杂志的记者，一天收到一封匿名信，信中透露本市知名富商周远山在一个偏远村庄隐藏着重大秘密。为了挖掘独家新闻，我带上相机独自启程。');
```

## Docker
### Build
```bash
docker build -t synergy_api_server .
```

### Run
```bash
docker run -d -p 8080:8080 synergy_api_server
```