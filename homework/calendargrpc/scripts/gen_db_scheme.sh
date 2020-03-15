#!/bin/bash

sudo -u postgres psql < ../migrations/db_scheme.sql
