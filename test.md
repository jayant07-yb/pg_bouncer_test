Deploying PgBouncer
Where do you install and run PgBouncer? There are different answers, with different advantages:

1. On the Postgres server node: You can install it alongside the PostgreSQL server itself, on the same node. The clients connect to the PgBouncer port rather than the Postgres port. This has the effect of an “enhanced” Postgres which does connection pooling internally. You also only have to maintain one copy of the configuration files for PgBouncer. On the other hand, this involves actually running something else also on the PostgreSQL server node, which may not be easy or permitted (firewalls, policies) or even possible (AWS RDS).
2. On client nodes: You can install PgBouncer in each client node, for example each web node runs Apache and PHP, and the PHP scripts connect to the local PgBouncer. This has the advantage of not having to disturb the server setup, and the pool configuration can be used to keep the server load predictable. On the flip side, if the number of client nodes are huge, or can vary a lot depending on the load/traffic, the server can be overloaded quickly.
3. As a standalone cluster: The third option to have a cluster of independent, stateless PgBouncer nodes, fronted by a TCP load balancer like HAProxy. This setup, while being more complicated than the other two options, provides maximum control and configurability.


-   pgbouncer test.ini
-   On client nodes::
```sh
Inserted data: INSERT INTO employee(id, name, age, language) VALUES (1, 'John', 35, 'Go')
Query for id=1 returned: Row[John, 35, Go]
2022/06/13 07:17:29 PGBouncer took 567.967408ms
Created table employee
Inserted data: INSERT INTO employee(id, name, age, language) VALUES (1, 'John', 35, 'Go')
Query for id=1 returned: Row[John, 35, Go]
2022/06/13 07:17:33 PGBouncer took 4.378692815s
Created table employee
Inserted data: INSERT INTO employee(id, name, age, language) VALUES (1, 'John', 35, 'Go')
Query for id=1 returned: Row[John, 35, Go]
2022/06/13 07:17:38 PGBouncer took 4.533311832s
Created table employee
Inserted data: INSERT INTO employee(id, name, age, language) VALUES (1, 'John', 35, 'Go')
Query for id=1 returned: Row[John, 35, Go]
2022/06/13 07:17:43 Without took 4.755971505s
Created table employee
Inserted data: INSERT INTO employee(id, name, age, language) VALUES (1, 'John', 35, 'Go')
Query for id=1 returned: Row[John, 35, Go]
2022/06/13 07:17:47 Without took 4.77600505s
Created table employee
Inserted data: INSERT INTO employee(id, name, age, language) VALUES (1, 'John', 35, 'Go')
Query for id=1 returned: Row[John, 35, Go]
2022/06/13 07:17:52 Without took 4.82102202s
```

-    We will be requiring once pgBouncer for each node.
-   They are mapped with the database name.


Starting the tservers :: 
-------


Launching masters:
ssh -q 10.150.1.135 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-master \
--master_addresses=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100  \
--rpc_bind_addresses=10.150.1.135:7100 --use_private_ip=zone \
--placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.20  </dev/null &

ssh -q 10.150.1.146 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-master \
--master_addresses=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100  \
--rpc_bind_addresses=10.150.1.146:7100 --use_private_ip=zone \
--placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.20  </dev/null &

ssh -q 10.150.1.140 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-master \
--master_addresses=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100  \
--rpc_bind_addresses=10.150.1.140:7100 --use_private_ip=zone \
--placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.20   </dev/null &

--------

ssh -q 10.150.1.135 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-tserver \
--tserver_master_addrs=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100 \
--server_broadcast_addresses=34.82.75.132 --rpc_bind_addresses=10.150.1.135:9100 \
--use_private_ip=zone --placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.30 \
--enable_ysql=true </dev/null &

ssh -q 10.150.1.146 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-tserver \
--tserver_master_addrs=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100 \
--server_broadcast_addresses=34.127.98.189 --rpc_bind_addresses=10.150.1.146:9100 \
--use_private_ip=zone --placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.30 \
--enable_ysql=true </dev/null &

ssh -q 10.150.1.140 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-tserver \
--tserver_master_addrs=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100 \
--server_broadcast_addresses=34.168.86.62 --rpc_bind_addresses=10.150.1.140:9100 \
--use_private_ip=zone --placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.30 \
--enable_ysql=true </dev/null &




--------

postgres=# create user user1 password 'user1pass';
CREATE ROLE
postgres=# create database db1 owner user1;
CREATE DATABASE
postgres=#

mvn -q package exec:java -DskipTests -Dexec.mainClass=com.yugabyte.HelloSqlApp

#   Issues Faced
-   [jayantanand@dev-server-janand-3 MySample]$ mvn -q package exec:java -DskipTests -Dexec.mainClass=com.yugabyte.HelloSqlApp
    Database has been initialized
    FATAL: query_wait_timeout
-   If the number of connection is increased (over 20) then getting this error. 


##   With    min_pool_size = 10 when number of connections were 22 
[jayantanand@dev-server-janand-3 pgbouncer_test]$ grep -c db_name  first.txt && grep -c db_name second.txt && grep -c db_name third.txt 
6
8
8

```sh
curl  10.150.1.135:13000/rpcz  > first.txt  && curl  10.150.1.140:13000/rpcz  > second.txt &&  curl 10.150.1.146:13000/rpcz > third.txt && grep -c db_name  first.txt && grep -c db_name second.txt && grep -c db_name third.txt
```

# Smart Driver settings for Server Broadcast 

ssh -q 10.150.1.135 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-tserver \
--tserver_master_addrs=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100 \
--server_broadcast_addresses=34.82.75.132 --rpc_bind_addresses=10.150.1.135:9100 \
--use_private_ip=zone --placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.30 \
--enable_ysql=true </dev/null &

ssh -q 10.150.1.146 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-tserver \
--tserver_master_addrs=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100 \
--server_broadcast_addresses=34.127.98.189 --rpc_bind_addresses=10.150.1.146:9100 \
--use_private_ip=zone --placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.30 \
--enable_ysql=true </dev/null &

ssh -q 10.150.1.140 /home/jayantanand/yugabyte-2.13.2.0/bin/yb-tserver \
--tserver_master_addrs=10.150.1.135:7100,10.150.1.146:7100,10.150.1.140:7100 \
--server_broadcast_addresses=34.168.86.62 --rpc_bind_addresses=10.150.1.140:9100 \
--use_private_ip=zone --placement_cloud gcp --placement_region us-west1 \
--placement_zone us-west1-b --fs_data_dirs=/home/jayantanand/var/data \
--replication_factor=3 --default_memory_limit_to_ram_ratio=0.30 \
--enable_ysql=true </dev/null &

----------
Setting Up PgBouncer Integration
If you are running PgBouncer instances that pool incoming connections to this PostgreSQL server, you can monitor them using pgDash.

To monitor a PgBouncer instance running on host PGBOUNCER-HOST listening on port PGBOUNCER-PORT, that has an admin user called ADMIN-USER, use the following command:

./pgmetrics_1.13.0_linux_amd64/pgmetrics -h localhost -p 5400 -U ADMIN-USER  -wfjson pgbouncer |
  ./pgdash_1.9.0_linux_amd64/pgdash -a kw9X4F4P6S01QfRIGOex9G report-pgbouncer 10.150.1.146 localhost
If you need to supply a password, you can do:

PGPASSWORD=user1pass \
./pgmetrics_1.13.0_linux_amd64/pgmetrics -h localhost -p 5400 -U user1  -wfjson pgbouncer |
  ./pgdash_1.9.0_linux_amd64/pgdash -a kw9X4F4P6S01QfRIGOex9G report-pgbouncer 10.150.1.146 localhost
Set this command up as a cron job running every few minutes to regularly collect and report your PgBouncer metrics. Repeat for each PgBouncer instance pooling connections to this PostgreSQL server.



-   [jayantanand@dev-server-janand-3 MySample]$ mvn -q package exec:java -DskipTests -Dexec.mainClass=com.yugabyte.pgbouncer_prepared_statement_test


mvn -q package exec:java -DskipTests -Dexec.mainClass=com.yugabyte.TestReturn




// Statement pooling



10.150.1.146,10.150.1.140,10.150.1.135


Transaction without 



mvn -q package exec:java -DskipTests -Dexec.mainClass=com.yugabyte.TestReturn
mvn -q package exec:java -DskipTests -Dexec.mainClass=com.yugabyte.TestPhase1


PREPARE AB2 AS  INSERT INTO employee VALUES ($1, 'John', 35, 'Prepared_Statement')

 http://localhost:13000/rpcz
 