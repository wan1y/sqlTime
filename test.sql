select tidb_version();
SHOW VARIABLES LIKE 'aurora\_version';
show variables where variable_name=‘language’ or variable_name=‘net_write_timeout’;
SHOW DRAINER STATUS;
SHOW PUMP STATUS;
SHOW MASTER STATUS;
SHOW BACKUPS;
show analyze status;
SHOW BUILTINS;
SHOW CHARACTER SET;
SHOW COLLATION;
SHOW CONFIG;
SHOW DATABASES;
SHOW SCHEMAS
SHOW CREATE USER 'root';
SHOW DRAINER STATUS;
SHOW ENGINES;
SHOW GRANTS;
SHOW PLUGINS;
show privileges;
SHOW PROCESSLIST;
SHOW STATS_HEALTHY;
show stats_histograms;
show stats_meta;
show status;
SHOW GLOBAL VARIABLES;
SELECT @@sql_mode;
SHOW GLOBAL VARIABLES like '%tidb%'