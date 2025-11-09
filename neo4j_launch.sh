#!/usr/bin/env bash
set -e

# --- CONFIG ---
CONTAINER_NAME="neo4j"
NEO4J_PASS="wikiscrolls"
NEO4J_IMAGE="neo4j:5"
DATA_DIR="$HOME/neo4j/data"
IMPORT_DIR="$HOME/neo4j/import"
PLUGINS_DIR="$HOME/neo4j/plugins"

# --- CREATE DIRECTORIES ---
mkdir -p "$DATA_DIR" "$IMPORT_DIR"

# --- STOP OLD CONTAINER IF RUNNING ---
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
  echo "Stopping old Neo4j container..."
  docker stop $CONTAINER_NAME >/dev/null 2>&1 || true
  docker rm $CONTAINER_NAME >/dev/null 2>&1 || true
fi

# --- RUN NEO4J CONTAINER ---
echo "Starting Neo4j container..."
docker run -d \
  --name $CONTAINER_NAME \
  -p 7474:7474 -p 7687:7687 \
  -e NEO4J_AUTH=neo4j/$NEO4J_PASS \
  -e NEO4J_dbms_memory_heap_max__size=6G \
  -e NEO4J_dbms_memory_pagecache_size=8G \
  -e NEO4J_PLUGINS='["graph-data-science"]' \
  -v "$DATA_DIR":/data \
  -v "$IMPORT_DIR":/import \
  -v "$PLUGINS_DIR":/plugins\
  $NEO4J_IMAGE

# --- INFO ---
echo
echo "Neo4j is running!"
echo "→ Browser:  http://localhost:7474"
echo "→ Username: neo4j"
echo "→ Password: $NEO4J_PASS"
echo
echo "TSV files available at /import inside container."
echo "Example import command:"
echo "  docker exec -it $CONTAINER_NAME neo4j-admin database import full --delimiter='\t' --nodes=/import/nodes.tsv --relationships=/import/relationships.tsv --database=wikipedia --overwrite-destination=true"

