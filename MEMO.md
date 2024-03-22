# Memo for oracle test

## Host info
* docker image: truevoly/oracle-12c
* localhost:1521/xe
* id/pw: sys as sysdba/oracle

## Command

1. By sys(sysdba)
```sql
CREATE USER ninem IDENTIFIED BY 1111;
GRANT CREATE TABLE TO ninem;
GRANT CREATE SESSION TO ninem;
```

2. By ninem schema
```sql
CREATE TABLE example_table (
    id NUMBER,
    name VARCHAR2(50)
);
```
