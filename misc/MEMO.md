# Memo for oracle test

## Host info
* docker image: truevoly/oracle-12c
* localhost:1521/xe
* id/pw: sys as sysdba/oracle

## Command

1. By sys(sysdba)
```sql
CREATE USER ninem IDENTIFIED BY 1111;
GRANT CONNECT, RESOURCE TO ninem;
GRANT CREATE TABLE TO ninem;
GRANT CREATE SESSION TO ninem;
GRANT UNLIMITED TABLESPACE TO ninem;
```

2. By ninem schema
```sql
CREATE TABLE example_table (
    id NUMBER,
    name VARCHAR2(50)
);
```

3. Drop schema
```sql
DROP USER ninem CASCADE;
```