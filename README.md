# GO4SQL

<p align="center">
<a href="https://github.com/LissaGreense/GO4SQL/actions">
<img alt="Unit-Tests Status" src="https://github.com/LissaGreense/GO4SQL/workflows/unit-tests/badge.svg?branch=main"/>
</a>

<a href="https://github.com/LissaGreense/GO4SQL/actions">
<img alt="End2End Status" src="https://github.com/LissaGreense/GO4SQL/workflows/end2end-tests/badge.svg?branch=main"/>
</a>

<a href="https://goreportcard.com/report/github.com/LissaGreense/GO4SQL">
<img alt="Report Status" src="https://goreportcard.com/badge/github.com/LissaGreense/GO4SQL"/>
</a>
</p>

GO4SQL is an open source project to write in-memory SQL engine using nothing but Golang.

## HOW TO USE

You can compile the project with ``go build``, this will create ``GO4SQL`` binary.
You can eithier specify file path with ``./GO4SQL -file file_path``, that will read the input data directly into the
program.
Also with ``./GO4SQL -stream`` you can run the program in stream mode, then you provide SQL commands in your console (
from standard input).

### Docker
1. Pull docker image: `docker pull kajedot/go4sql:latest`
2. Run docker container in the interactive mode, remember to provide flag, for example: `docker run -i kajedot/go4sql -stream`
3. You can test this image with `test_file` provided in this repo: `docker run -i kajedot/go4sql -stream < test_file`

## FUNCTIONALITY

* ***CREATE TABLE*** - you can create table with name ``table1`` using
  command: ``CREATE TABLE table1( one TEXT , two INT);``. First column is called ``one`` and it contains strings, second
  one is called ``two`` and it contains integers.
* ***INSERT INTO*** - you can insert values into table called ``table1`` with
  command ``INSERT INTO table1 VALUES( 'hello', 1);``. Please note that the number or arguments and types of the values
  must be the same as you declared.
* ***SELECT FROM*** - you can either select everything from  ``table1`` with ``SELECT * FROM table1;`` command. Or you
  can specify column names that you're intrest in: ``SELECT one, two FROM table1;``, note that culumn names must be the
  same as you declared and duplicated column names will be ignored.

## UNIT TESTS

To run all the tests locally use "go clean -testcache; go test ./..." in root directory.

## DOCKER

To build docker image locally, run this command in the root directory:
```
docker build -t go4sql:test .
```
