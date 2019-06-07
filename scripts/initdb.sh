#!/bin/bash

cd $(dirname $0)
mysql -u root -p < db.sql
