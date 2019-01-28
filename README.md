RUN:

```
SELECT pg_locks.*, pg_stat_activity.* FROM pg_stat_activity left join pg_locks on pg_stat_activity.pid = pg_locks.pid;

```

To check locks

