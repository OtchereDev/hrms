# Appointment Letters Module - Backend Documentation

## Overview

Generate and send appointment letters to new hires.

## Core Doctype: Appointment Letter

### Fields
- `applicant_name`: Job applicant reference
- `date_of_joining`: Start date
- `designation`, `department`, `company`
- `letter_head`: Company letterhead
- `body`: Letter content (Jinja template)

### Business Logic
- Auto-populate from Job Offer
- Template-based generation
- Email delivery
- E-signature support

## Key Rules
- One letter per job offer
- Generated after offer acceptance
- Can be printed or emailed
