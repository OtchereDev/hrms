# Daily Work Summary Module - Backend Documentation

## Overview

Daily work logging and productivity tracking for employees.

## Core Doctype: Daily Work Summary

### Fields
- `employee`: Employee reference
- `posting_date`: Summary date
- `work_summary`: Text description
- `total_hours`: Hours worked

### Business Logic
- Email reminders for pending summaries
- Weekly/monthly rollups
- Manager visibility

## Key Rules
- One summary per employee per day
- Auto-email reminders
- Manager can view team summaries
