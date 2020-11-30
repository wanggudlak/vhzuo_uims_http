/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : localhost:3306
 Source Schema         : uims

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 27/04/2020 18:12:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for uims_access_resource
-- ----------------------------
DROP TABLE IF EXISTS `uims_access_resource`;
CREATE TABLE `uims_access_resource` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '资源ID',
  `client_id` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '客户端业务系统ID',
  `org_id` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '客户端组织ID',
  `res_code` char(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源编码',
  `res_front_code` varchar(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '和前端约定的资源编码',
  `res_type` char(1) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源类型，A：逻辑资源；B：实体资源',
  `res_sub_type` char(3) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源',
  `platform` varchar(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '所用平台，all：所有，back_desk：后台；front_desk：前台；vzhuo_back：结算后台；vzhuo_front：结算前台',
  `res_name_en` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源的英文名称',
  `res_name_cn` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源的中文名称',
  `res_endp_route` varchar(128) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源的后端路由URI',
  `res_data_location` json DEFAULT NULL COMMENT '资源所在的位置，主要用于数据权限，json存储，包含以下属性：客户端id、数据库名、表名、行记录属性名、行记录属性值',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，默认N：未软删除；Y：已软删除',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限资源表';

-- ----------------------------
-- Table structure for uims_client
-- ----------------------------
DROP TABLE IF EXISTS `uims_client`;
CREATE TABLE `uims_client` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '客户端业务系统ID',
  `app_id` char(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端系统APPID，用来唯一标识客户端系统',
  `app_secret` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端系统与UIMS系统之间用来对称加解密的秘钥，base64编码后存储',
  `client_type` char(3) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端类型，VDK：微桌',
  `client_flag_code` varchar(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端业务系统标识，VDK_CASS：微桌结算系统；VDK_MP：微桌任务系统平台；VDK_CRM：微桌CRM系统；VDK_INVO：微桌代开发票系统；VDK_ESIGN：微桌电签系统；VDK_ES_SAPP：微桌电签小程序；',
  `client_spm1_code` char(4) CHARACTER SET utf8 NOT NULL DEFAULT '1024' COMMENT '客户端业务系统的SPM编码中的第一部分（外站类型ID），微桌内部系统用1024；外部系统用2048',
  `client_spm2_code` char(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端业务系统的SPM编码中的第二部分（外站APP ID），和APPID一致',
  `client_name` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端业务系统名称',
  `status` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '客户端业务系统使用UIMS的状态，默认N：未授权不可用；Y：已授权可用；F-被禁用',
  `client_host_ip` json NOT NULL COMMENT '客户端使用的IP，可用于白名单',
  `client_host_url` json NOT NULL COMMENT '客户端业务系统当前使用的域名，例如微桌结算系统是https://fuwu.skysharing.cn',
  `client_pub_key_path` varchar(128) CHARACTER SET utf8 NOT NULL COMMENT '客户端业务系统的RSA公钥key文件路径，以appid作为各自的目录',
  `uims_pub_key_path` varchar(128) CHARACTER SET utf8 NOT NULL COMMENT 'UIMS系统的RSA公钥文件路径',
  `uims_pri_key_path` varchar(128) CHARACTER SET utf8 NOT NULL COMMENT 'UIMS系统的RSA私钥文件路径',
  `in_at` datetime(6) NOT NULL COMMENT '入驻可以使用的开始时间点',
  `forget_at` datetime(6) NOT NULL COMMENT '在什么时间点，客户端系统不能使用UIMS，默认是空字符串',
  `created_at` datetime(6) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户端即使用uims系统的业务系统';

-- ----------------------------
-- Table structure for uims_client_settings
-- ----------------------------
DROP TABLE IF EXISTS `uims_client_settings`;
CREATE TABLE `uims_client_settings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT ,
  `client_id` int(11) unsigned NOT NULL COMMENT '客户端业务系统ID',
  `type` char(3) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '类型：LGN-用于登录的设置；REG-用于注册的设置；',
  `bus_channel_id` char(3) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '频道ID，对于登录业务，频道ID为100；注册业务频道ID为200',
  `page_id` char(3) NOT NULL DEFAULT '' COMMENT '页面ID，对于登录业务的登录页面ID为101；注册业务的注册页面ID为201',
  `spm_full_code` char(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT 'spm编码，由以下组成：client_spm1_code.client_spm2_code.频道ID.页面ID',
  `form_fields` json DEFAULT NULL COMMENT '表单域属性数据{"a": [{"attr_id": "account", "attr_cn": "账号"},{"attr_id": "passwd", "attr_cn": "密码"}, {"attr_id": "sms_code", "attr_cn": "验证码", "type": "phone_sms"}]}',
  `page_template_file` json DEFAULT NULL COMMENT '登录页或注册页html模板文件(.tmpl后缀)所在的位置(UIMS系统的根目录的相对路径，并以appid作为子目录名)，因为涉及多阶段登录，采用json存储，{"a":"/downloads/app_id/login_a.tmpl", "b": "/sownloads/app_id/login_b.tmpl"}',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，默认N：未删除；Y：已软删除',
  `created_at` datetime(6) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户端业务系统设置';

-- ----------------------------
-- Table structure for uims_organization
-- ----------------------------
DROP TABLE IF EXISTS `uims_organization`;
CREATE TABLE `uims_organization` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '组织ID',
  `parent_org_id` int(11) unsigned NOT NULL COMMENT '直接父级组织ID',
  `client_id` int(11) unsigned NOT NULL COMMENT '客户端业务系统id',
  `client_app_id` char(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '客户端APPID',
  `org_name_cn` varchar(255) NOT NULL DEFAULT '' COMMENT '组织中文名',
  `org_name_en` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '组织英文名',
  `org_code` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '组织代码',
  `org_level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '组织层级，0：顶级；1：第1级，以此类推',
  `org_full_pinyin` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '组织全拼音',
  `org_first_pinyin` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '组织拼音搜字母',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组织表';

-- ----------------------------
-- Table structure for uims_res_group
-- ----------------------------
DROP TABLE IF EXISTS `uims_res_group`;
CREATE TABLE `uims_res_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '资源组ID',
  `res_group_code` char(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '权限资源策略组编码',
  `platform` varchar(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '所用平台，all：所有，back_desk：后台；front_desk：前台；vzhuo_back：结算后台；vzhuo_front：结算前台',
  `res_group_en` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '权限策略组英文名称',
  `res_group_cn` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '权限策略组中文名称',
  `res_group_type` char(10) CHARACTER SET utf8 NOT NULL DEFAULT 'DEFAULT' COMMENT '权限策略组类型：DEFAULT-默认策略组；SELF-自定义配置的策略组',
  `res_of_curr` json DEFAULT NULL COMMENT '属于当前策略组的资源id list',
  `client_id` int(11) unsigned NOT NULL COMMENT '客户端或业务系统ID，默认是0，即不区分客户端业务系统，属于跨业务系统通用类型策略组',
  `org_id` int(11) unsigned NOT NULL COMMENT '组织id，默认是0，标识不区分组织，即是跨组织型的策略组',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，默认N：未软删除；Y：已软删除',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源策略组';

-- ----------------------------
-- Table structure for uims_role
-- ----------------------------
DROP TABLE IF EXISTS `uims_role`;
CREATE TABLE `uims_role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT  COMMENT '角色ID',
  `client_id` int(11) unsigned NOT NULL DEFAULT  '0' COMMENT '客户端ID',
  `org_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属的组织ID',
  `role_type` char(1) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '角色类型：A：UIMS的角色；F：通过页面增加的角色',
  `role_code` varchar(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '资源编码，UIMS系统的角色固定用UIMS.SUPERADMIN.001',
  `role_name_en` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '角色英文名称',
  `role_name_cn` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '角色中文名称',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，默认N：未软删除；Y：已软删除',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ----------------------------
-- Table structure for uims_role_res_map
-- ----------------------------
DROP TABLE IF EXISTS `uims_role_res_map`;
CREATE TABLE `uims_role_res_map` (
  `id` bigint(11) unsigned AUTO_INCREMENT NOT NULL,
  `role_id` int(11) unsigned NOT NULL COMMENT '角色ID',
  `res_grp_id` int(11) unsigned NOT NULL COMMENT '资源组ID',
  `client_id` int(11) unsigned NOT NULL COMMENT '客户端id，默认是0，即不区分客户端业务系统',
  `org_id` int(11) unsigned NOT NULL COMMENT '组织id，默认是0，即不区分组织',
  `start_valid_at` datetime NOT NULL COMMENT '开始生效时间，默认值是created_at',
  `forget_at` datetime NOT NULL COMMENT '忘记时间，即这个角色与这个资源的关系将在什么时刻自动过期，默认值是6666-01-01 00:00:00',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，N：没有软删除；Y：已经软删除',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色与资源组关系';

-- ----------------------------
-- Table structure for uims_user_auth
-- ----------------------------
DROP TABLE IF EXISTS `uims_user_auth`;
CREATE TABLE `uims_user_auth` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户或会员ID',
  `user_type` char(3) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户类型，取client_type字段的值，VDK：微桌全平台下面的用户',
  `account` char(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '普通登录账号，允许最大长度16',
  `user_code` char(11) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户编码，全数字最多11位，不同的组织下可以重复',
  `na_code` char(5) CHARACTER SET utf8 NOT NULL COMMENT '国家代码，中国：+86',
  `phone` char(12) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '手机号',
  `wechat` varchar(24) NOT NULL DEFAULT '' COMMENT '微信号',
  `wechat_id` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户的微信ID',
  `email` varchar(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '电子邮箱地址',
  `salt` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '盐值',
  `encrypt_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '结算传0, 微桌系统传1',
  `passwd` varchar(512) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '密码密文',
  `status` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'Y' COMMENT '账号的状态，默认Y：正常；N：已冻结，禁止登录；',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，默认N：未软删除；Y：已软删除',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户登录鉴权';

-- ----------------------------
-- Table structure for uims_user_info
-- ----------------------------
DROP TABLE IF EXISTS `uims_user_info`;
CREATE TABLE `uims_user_info` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(11) unsigned NOT NULL COMMENT '用户ID',
  `is_identify` char(1) DEFAULT NULL COMMENT '是否实名认证，默认N：没有；Y：已经实名认证',
  `user_code` char(11) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户编码，全数字最多11位，不同的组织下可以重复',
  `user_type` char(3) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户类型，取client_type字段的值，VDK：微桌全平台下面的用户',
  `name_en` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户英文姓名',
  `name_cn` varchar(64) NOT NULL DEFAULT '' COMMENT '用户中文名',
  `name_cn_alias` varchar(16) NOT NULL DEFAULT '' COMMENT '用户别名',
  `name_abbr_py` varchar(16) CHARACTER SET utf8 NOT NULL COMMENT '用户姓名拼音首字母',
  `name_full_py` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户姓名全拼音，英文空格分隔',
  `wechat` varchar(24) NOT NULL DEFAULT '' COMMENT '用户的微信号',
  `identity_card_no` varchar(20) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '身份证号码',
  `na_code` char(5) NOT NULL COMMENT '国家代码，中国：+86',
  `phone` char(12) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '手机号',
  `landline_phone` varchar(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '座机号',
  `sex` char(1) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '性别，M：男；F：女',
  `birthday` date COMMENT '出生日期，年月日',
  `nickname` varchar(32) NOT NULL DEFAULT '' COMMENT '昵称',
  `taxer_type` char(1) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '纳税人类型，A：一般纳税人',
  `taxer_no` varchar(16) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '纳税人识别号',
  `header_img_url` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户头像图片相对地址',
  `identity_card_person_img` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户身份证人像面图片相对地址',
  `identity_card_emblem_img` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户身份证国徽面图片相对地址',
  `isdel` char(1) CHARACTER SET utf8 NOT NULL DEFAULT 'N' COMMENT '是否软删除，默认N：未软删除；Y：已软删除',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户资料库';

-- ----------------------------
-- Table structure for uims_user_org
-- ----------------------------
DROP TABLE IF EXISTS `uims_user_org`;
CREATE TABLE `uims_user_org` (
  `id` bigint(11) unsigned AUTO_INCREMENT NOT NULL,
  `user_id` bigint(11) unsigned NOT NULL,
  `org_id` int(11) unsigned NOT NULL,
  `client_id` int(11) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户与组织关系表';

-- ----------------------------
-- Table structure for uims_user_role
-- ----------------------------
DROP TABLE IF EXISTS `uims_user_role`;
CREATE TABLE `uims_user_role` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(11) unsigned NOT NULL COMMENT '用户ID',
  `role_id` int(11) unsigned NOT NULL COMMENT '角色ID',
  `client_id` int(11) unsigned NOT NULL COMMENT '客户端业务系统ID',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户与角色关系表';

SET FOREIGN_KEY_CHECKS = 1;
