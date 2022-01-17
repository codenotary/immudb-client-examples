# immudb-f1 example
Formula 1 database insertion into immudb

New tables are created in immudb and index set

Index allows for queries like:

```
select drivers.surname, rescount.rcount
from (select driverid, count(*) as rcount from results where statusid = 2 group by driverid order by driverid) as rescount
join drivers on drivers.driverid = rescount.driverid
```
