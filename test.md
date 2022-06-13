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

