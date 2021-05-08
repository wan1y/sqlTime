# sqlTime 

Use this tool to compare the execution time of the same SQL on different database clusters.

```bash
# sqlTime --help

sql execution time test tool

Usage:
  sqlTime [flags]
  sqlTime [command]

Available Commands:
  help        Help about any command
  time        Compare SQL execution time

Flags:
      --dsns strings        set DSN of test dbs,take the first one as the standard
      --filenames strings   file name to save the result
  -h, --help                help for sqlTime
      --log-level string    log level: info, warn, error, silent (default "error")
      --sqlfile string      sql file (default "test.sql")
      --threads int         concurrent threads (default 1)
```



for example

~~~bash
# sqlTime time --dsns="root@tcp(192.168.10.1:4000)/test,root:@tcp(192.168.10.2:4000)/test" --filenames="test.txt,test1.txt" --log-level info --threads 100 --sqlfile test.sql
~~~

This command will generate three files in the current directory: text.txt, text1.txt,StandardTime.txt

StandardTime.txt is all the execution time of the first dsn, text.txt and text1.txt save the SQL and its execution time that are too different from the standard execution time.
