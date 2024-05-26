package types

type HotUpdate struct {
	Name string
	App  *App
}

type HotUpdateSetting struct {
	MsgID string
	Name  string
	Error error
	Data  []byte
}

type App struct {
	Title      string `mapstructure:"title" json:"title" yaml:"title"`
	Listen     string `mapstructure:"listen" json:"listen" yaml:"listen"`
	Port       string `mapstructure:"port" json:"port" yaml:"port"`
	Mode       string `mapstructure:"mode" json:"mode" yaml:"mode"`
	RootRouter string `mapstructure:"root_router" json:"root_router" yaml:"root_router"`
	*Nacos     `mapstructure:"nacos" json:"nacos" yaml:"nacos"`
	*Log       `mapstructure:"log" json:"log" yaml:"log"`
	*Sql       `mapstructure:"sql" json:"sql" yaml:"sql"`
	*Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`
	*Mail      `mapstructure:"mail" json:"mail" yaml:"mail"`
	*Jwt       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	*Oss       `mapstructure:"oss" json:"oss" yaml:"oss"`
	*Mq        `mapstructure:"mq" json:"mq" yaml:"mq"`
	*WarmUp    `mapstructure:"warm_up" json:"warm_up" yaml:"warm_up"`
	*Qps       `mapstructure:"qps" json:"qps" yaml:"qps"`
	*Prefix    `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	*CoolDown  `mapstructure:"cool_down" json:"cool_down" yaml:"cool_down"`
	*DubboGo   `mapstructure:"dubbo_go" json:"dubbo_go" yaml:"dubbo_go"`
}

type Log struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	MaxSizeMB    int    `mapstructure:"max_size_mb" json:"max_size_mb" yaml:"max_size_mb"`
	MaxAgeDay    int    `mapstructure:"max_age_day" json:"max_age_day" yaml:"max_age_day"`
	MaxBackupDay int    `mapstructure:"max_backup_day" json:"max_backup_day" yaml:"max_backup_day"`
	Compress     bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type Sql struct {
	SqlAddr         string `mapstructure:"sql_addr" json:"sql_addr" yaml:"sql_addr"`
	SqlPort         string `mapstructure:"sql_port" json:"sql_port" yaml:"sql_port"`
	SqlUser         string `mapstructure:"sql_user" json:"sql_user" yaml:"sql_user"`
	SqlPass         string `mapstructure:"sql_pass" json:"sql_pass" yaml:"sql_pass"`
	DBName          string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	StdPrint        bool   `mapstructure:"std_print" json:"std_print" yaml:"std_print"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn" json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxOpenConn     int    `mapstructure:"max_open_conn" json:"max_open_conn" yaml:"max_open_conn"`
	ConnMaxLifeHour int    `mapstructure:"conn_max_life_hour" json:"conn_max_life_hour" yaml:"conn_max_life_hour"`
	IdleMaxLifeHour int    `mapstructure:"idle_max_life_hour" json:"idle_max_life_hour" yaml:"idle_max_life_hour"`
}

func (c *Sql) MysqlDSN() string {
	return c.SqlUser + ":" + c.SqlPass + "@tcp(" + c.SqlAddr + ":" + c.SqlPort + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

type Redis struct {
	RedisAddr          string `mapstructure:"redis_addr" json:"redis_addr" yaml:"redis_addr"`
	RedisUser          string `mapstructure:"redis_user" json:"redis_user" yaml:"redis_user"`
	RedisPort          int    `mapstructure:"redis_port" json:"redis_port" yaml:"redis_port"`
	RedisPass          string `mapstructure:"redis_pass" json:"redis_pass" yaml:"redis_pass"`
	RedisDB            int    `mapstructure:"redis_db" json:"redis_db" yaml:"redis_db"`
	WithExpirySecond   int    `mapstructure:"with_expiry_second" json:"with_expiry_second" yaml:"with_expiry_second"`
	WithTriesTimes     int    `mapstructure:"with_tries_times" json:"with_tries_times" yaml:"with_tries_times"`
	ReadTimeoutSecond  int    `mapstructure:"read_timeout_second" json:"read_timeout_second" yaml:"read_timeout_second"`
	WriteTimeoutSecond int    `mapstructure:"write_timeout_second" json:"write_timeout_second" yaml:"write_timeout_second"`
	PoolSize           int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	MinIdleConn        int    `mapstructure:"min_idle_conn" json:"min_idle_conn" yaml:"min_idle_conn"`
	MaxConnAgeMinute   int    `mapstructure:"max_conn_age_minute" json:"max_conn_age_minute" yaml:"max_conn_age_minute"`
	PoolTimeoutMinute  int    `mapstructure:"pool_timeout_minute" json:"pool_timeout_minute" yaml:"pool_timeout_minute"`
	IdleTimeoutMinute  int    `mapstructure:"idle_timeout_minute" json:"idle_timeout_minute" yaml:"idle_timeout_minute"`
}

type Mail struct {
	MailHost string `mapstructure:"mail_host" json:"mail_host" yaml:"mail_host"`
	MailPort int    `mapstructure:"mail_port" json:"mail_port" yaml:"mail_port"`
	MailUser string `mapstructure:"mail_user" json:"user" yaml:"mail_user"`
	MailPass string `mapstructure:"mail_pass" json:"mail_pass" yaml:"mail_pass"`
}

type Jwt struct {
	Key          string `mapstructure:"key" json:"key" yaml:"key"`
	Issuer       string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	Subject      string `mapstructure:"subject" json:"subject" yaml:"subject"`
	ExpireMinute int    `mapstructure:"expire_minute" json:"expire_minute" yaml:"expire_minute"`
	ExpireDay    int    `mapstructure:"expire_day" json:"expire_day" yaml:"expire_day"`
}

type Oss struct {
	StrIpPort       string `mapstructure:"str_ip_port" json:"str_ip_port" yaml:"str_ip_port"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name" json:"bucket_name" yaml:"bucket_name"`
	ObjectName      string `mapstructure:"object_name" json:"object_name" yaml:"object_name"`
	HostURL         string `mapstructure:"host_url" json:"host_url" yaml:"host_url"`
	CallbackURL     string `mapstructure:"callback_url" json:"callback_url" yaml:"callback_url"`
	Prefix          string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	ExpireTimeSec   int64  `mapstructure:"expires_time_sec" json:"expire_time_sec" yaml:"expires_time_sec"`
}

type Mq struct {
	Host       []string `mapstructure:"host" json:"host" yaml:"host"`
	User       string   `mapstructure:"user" json:"user" yaml:"user"`
	Pass       string   `mapstructure:"pass" json:"pass" yaml:"pass"`
	Broker     string   `mapstructure:"broker" json:"broker" yaml:"broker"`
	BatchCount int      `mapstructure:"batch_count" json:"batch_count" yaml:"batch_count"`
	RetryTimes int      `mapstructure:"retry_times" json:"retry_times" yaml:"retry_times"`
}

type WarmUp struct {
	Resource               string  `mapstructure:"resource" json:"resource" yaml:"resource"`
	TokenCalculateStrategy string  `mapstructure:"token_calculate_strategy" json:"token_calculate_strategy" yaml:"token_calculate_strategy"`
	ControlBehavior        string  `mapstructure:"control_behavior" json:"control_behavior" yaml:"control_behavior"`
	Threshold              float64 `mapstructure:"threshold" json:"threshold" yaml:"threshold"`
	WarmUpPeriodSec        uint32  `mapstructure:"warm_up_period_sec" json:"warm_up_period_sec" yaml:"warm_up_period_sec"`
	WarmUpColdFactor       uint32  `mapstructure:"warm_up_cold_factor" json:"warm_up_cold_factor" yaml:"warm_up_cold_factor"`
}

type Qps struct {
	Resource               string  `mapstructure:"resource" json:"resource" yaml:"resource"`
	WarmOrQps              string  `mapstructure:"warm_or_qps" json:"warm_or_qps" yaml:"warm_or_qps"`
	TokenCalculateStrategy string  `mapstructure:"token_calculate_strategy" json:"token_calculate_strategy" yaml:"token_calculate_strategy"`
	ControlBehavior        string  `mapstructure:"control_behavior" json:"control_behavior" yaml:"control_behavior"`
	WarmUpPeriodSec        uint32  `mapstructure:"warm_up_period_sec" json:"warm_up_period_sec" yaml:"warm_up_period_sec"`
	WarmUpColdFactor       uint32  `mapstructure:"warm_up_cold_factor" json:"warm_up_cold_factor" yaml:"warm_up_cold_factor"`
	Threshold              float64 `mapstructure:"threshold" json:"threshold" yaml:"threshold"`
	StatIntervalMs         uint32  `mapstructure:"stat_interval_ms" json:"stat_interval_ms" yaml:"stat_interval_ms"`
}

type Nacos struct {
	NacosAddr   string `mapstructure:"nacos_addr" json:"nacos_addr" yaml:"nacos_addr"`
	NacosPort   uint64 `mapstructure:"nacos_port" json:"nacos_port" yaml:"nacos_port"`
	Username    string `mapstructure:"username" json:"username" yaml:"username"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	NamespaceId string `mapstructure:"namespace_id" json:"namespace_id" yaml:"namespace_id"`
	TimeoutMs   uint64 `mapstructure:"timeout_ms" json:"timeout_ms" yaml:"timeout_ms"`
	LoadCache   bool   `mapstructure:"load_cache" json:"load_cache" yaml:"load_cache"`
	LogDir      string `mapstructure:"log_dir" json:"log_dir" yaml:"log_dir"`
	CacheDir    string `mapstructure:"cache_dir" json:"cache_dir" yaml:"cache_dir"`
	RotateTime  string `mapstructure:"rotate_time" json:"rotate_time" yaml:"rotate_time"`
	MaxAge      int64  `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	LogLevel    string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	DataId      string `mapstructure:"data_id" json:"data_id" yaml:"data_id"`
	Group       string `mapstructure:"group" json:"group" yaml:"group"`
}

type Prefix struct {
	DownloadPrefix string `mapstructure:"download_prefix" json:"download_prefix" yaml:"download_prefix"`
	EmailPrefix    string `mapstructure:"email_prefix" json:"email_prefix" yaml:"email_prefix"`
}

type CoolDown struct {
	CoolDownMs int `mapstructure:"cool_down_ms" json:"cool_down_ms" yaml:"cool_down_ms"`
	CoolDownS  int `mapstructure:"cool_down_s" json:"cool_down_s" yaml:"cool_down_s"`
	CoolDownM  int `mapstructure:"cool_down_m" json:"cool_down_m" yaml:"cool_down_m"`
	CoolDownH  int `mapstructure:"cool_down_h" json:"cool_down_h" yaml:"cool_down_h"`
}

type DubboGo struct {
	ServerName    string `mapstructure:"server_name" json:"server_name" yaml:"server_name"`
	RegistryTo    string `mapstructure:"registry_to" json:"registry_to" yaml:"registry_to"` // nacos or zookeeper ...
	RegisterAddr  string `mapstructure:"register_addr" json:"register_addr" yaml:"register_addr"`
	RegisterGroup string `mapstructure:"register_group" json:"register_group" yaml:"register_group"`
}
