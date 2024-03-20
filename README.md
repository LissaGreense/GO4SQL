# GO4SQL

<p align="center">
<a href="https://github.com/LissaGreense/GO4SQL/actions/workflows/unit-tests.yml">
<img alt="Unit-Tests Status" src="https://github.com/LissaGreense/GO4SQL/actions/workflows/unit-tests.yml/badge.svg?branch=main"/>
</a>

<a href="https://github.com/LissaGreense/GO4SQL/actions/workflows/end2end-tests.yml">
<img alt="End2End Status" src="https://github.com/LissaGreense/GO4SQL/actions/workflows/end2end-tests.yml/badge.svg?branch=main"/>
</a>

<a href="https://goreportcard.com/report/github.com/LissaGreense/GO4SQL">
<img alt="Report Status" src="https://goreportcard.com/badge/github.com/LissaGreense/GO4SQL"/>
</a>
</p>

GO4SQL is an open source project to write in-memory SQL engine using nothing but Golang.

## HOW TO USE

You can compile the project with ``go build``, this will create ``GO4SQL`` binary.

Currently, there are 3 modes to chose from:

1. `File Mode` - You can specify file path with ``./GO4SQL -file file_path``, that will read the input
   data directly into the program and print the result.

2. `Stream Mode` - With ``./GO4SQL -stream`` you can run the program in stream mode, then you provide SQL commands
   in your console (from standard input).

3. `Socket Mode` - To start Socket Server use `./GO4SQL -socket`, it will be listening on port `1433` by default. To
   choose port different other than that, for example equal to `1444`, go with: `./GO4SQL -socket -port 1444`

### Docker
1. Pull docker image: `docker pull kajedot/go4sql:latest`
2. Run docker container in the interactive mode, remember to provide flag, for example: `docker run -i kajedot/go4sql -stream`
3. You can test this image with `test_file` provided in this repo: `docker run -i kajedot/go4sql -stream < test_file`

## FUNCTIONALITY

* ***CREATE TABLE*** - you can create table with name ``table1`` using
  command:
  ```sql
  CREATE TABLE table1( one TEXT , two INT);
  ```
  First column is called ``one`` and it contains strings (keyword ``TEXT``), second
  one is called ``two`` and it contains integers (keyword ``INT``).


* ***INSERT INTO*** - you can insert values into table called ``table1`` with
  command:
  ```sql
  INSERT INTO table1 VALUES( 'hello', 1);
  ```
  Please note that the number of arguments and types of the values
  must be the same as you declared with ``CREATE``.


* ***SELECT FROM*** - you can either select everything from  ``table1`` with:
  ```SELECT * FROM table1;```
  Or you can specify column names that you're interested in:
  ```sql
  SELECT one, two FROM table1;
  ```
  Note that column names must be the
  same as you declared with ``CREATE`` and also duplicated column names will be ignored.


* ***WHERE*** - is used to filter records. It is used to extract only those records that fulfill a
  specified condition. It can be used with ``SELECT`` like this:
  ```sql
  SELECT column1, column2
  FROM table_name
  WHERE column1 NOT 'goodbye' OR column2 EQUAL 3;
  ```
  Supported logical operations are: ``EQUAL``, ``NOT``, ``OR``, ``AND``, ```FALSE```, ```TRUE```.


* ***DELETE FROM*** is used to delete existing records in a table. It can be used like this:
  ```sql
  DELETE FROM tb1 WHERE two EQUAL 3;
  ```
  ``tb1`` is the name of the table, and ``WHERE`` specify records that fulfill a
  specified condition and afterward will be deleted.


* ***ORDER BY***  is used to sort the result-set in ascending or descending order. It can be used
  with ``SELECT`` like this:
  ```sql
  SELECT column1, column2,
  FROM table_name
  ORDER BY column1 ASC, column2 DESC;
  ```
  In this case, this command will order by ``column1`` in ascending order, but if some rows have the
  same ``column1``, it orders them by column2 in descending order.

## UNIT TESTS

To run all the tests locally run this in root directory:

```shell
go clean -testcache; go test ./...
```

## E2E TEST

In root directory there is **test_file** containing input commands for E2E tests. File
**.github/expected_results/end2end.txt** has expected results for it.
This is integrated into github workflows.

## DOCKER

To build your docker image run this command in root directory:

```shell
docker build -t go4sql:test .
```

### Run docker in interactive stream mode
To run this docker image in interactive stream mode mode use this command:

```shell
docker run -i go4sql:test -stream
```

### Run docker in socket mode
To run this docker image in socket mode use this command:

```shell
docker run go4sql:test -socket
```

### Run docker in file mode
**NOT RECOMMENDED**

Alternatively you can run a docker image in file mode:

```shell
docker run -i go4sql:test -file <PATH_TO_FILE>
```

## HELM

To create a pod deployment using helm chart, there is configuration under `./helm` directory.

Commands:
```shell
cd ./helm
helm install go4sql_pod_name GO4SQL/
```

To check status of pod, use:
```shell
kubectl get pods
```
