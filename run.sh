#!/bin/bash
go run . network | sed 's/Network running on \[::\]//'
