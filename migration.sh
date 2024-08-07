#!/bin/bash

name=$1

last_migration=$(ls migrations | grep -oE '^[0-9]+' | sort -n | tail -n 1)

if [ -z "$last_migration" ]; then
  new_migration=1
else
  new_migration=$((last_migration + 1))
fi

touch "migrations/${new_migration}_${name}.up.sql"
touch "migrations/${new_migration}_${name}.down.sql"

echo "Созданы новые файлы миграции: ${new_migration}_${name}.up.sql и ${new_migration}_${name}.down.sql"
