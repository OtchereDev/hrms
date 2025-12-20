# Travel Management Module - Backend Documentation

## Overview

Business logic and data processing for employee travel requests, itinerary planning, and travel cost estimation.

---

## Core Doctypes

### 1. Travel Request

**Purpose**: Employee travel request with itinerary and costing details.

#### Fields

**Travel Details**:
- `travel_type` (Select, required): Domestic, International
- `purpose_of_travel` (Link to Purpose of Travel, required)
- `travel_funding` (Select): Require Full Funding, Fully Sponsored, Partially Sponsored
- `travel_proof` (Attach): Copy of invitation/announcement
- `details_of_sponsor` (Data): Sponsor name and location
- `description` (Small Text): Additional details

**Employee Details**:
- `employee` (Link, required)
- `employee_name` (Data, read-only)
- `cell_number` (Data)
- `prefered_email` (Data)
- `company` (Link, required)
- `date_of_birth` (Date)
- `personal_id_type` (Data)
- `personal_id_number` (Data)
- `passport_number` (Data)

**Itinerary (Child Table - Travel Itinerary)**:
- `travel_from` (Data): Origin
- `travel_to` (Data): Destination
- `travel_date` (Date)
- `expected_arrival_date` (Date)
- `lodging_required` (Check)
- `preferred_hotel_type` (Link to Hotel Type)
- `mode_of_travel` (Link to Mode of Travel)

**Costing (Child Table - Travel Request Costing)**:
- `item` (Data): Cost item description
- `expense_type` (Link to Expense Claim Type)
- `estimated_cost` (Currency)

**Event Details**:
- `name_of_organizer` (Data)
- `address_of_organizer` (Data)

**Accounting**:
- `cost_center` (Link)
- Accounting dimensions (Project, Department, etc.)

#### Business Logic

**Validation**:
```python
class TravelRequest(Document):
    def validate(self):
        validate_active_employee(self.employee)
```

**Auto-naming**: HR-TRQ-.YYYY.-.##### (e.g., HR-TRQ-2024-00001)

**Status Workflow**: Can implement custom workflow for approval

---

### 2. Purpose of Travel (Master)

**Purpose**: Predefined travel purposes.

#### Fields

- `purpose` (Data, required): e.g., "Conference", "Client Meeting", "Training"

---

### 3. Hotel Type (Master)

**Purpose**: Hotel category preferences.

#### Fields

- `hotel_type` (Data, required): e.g., "5 Star", "4 Star", "Budget"

---

### 4. Mode of Travel (Master)

**Purpose**: Transportation modes.

#### Fields

- `mode_of_travel` (Data, required): e.g., "Air", "Train", "Car", "Bus"

---

## Key Business Rules

1. **Active Employee**: Only active employees can create travel requests
2. **Travel Type**: Must specify Domestic or International
3. **Itinerary Required**: At least one travel leg should be specified
4. **Costing Optional**: Estimated costs can be added for budgeting
5. **Supporting Documents**: Can attach invitation/announcement proof
6. **Multi-leg Travel**: Supports multiple destinations in itinerary
7. **Lodging Tracking**: Can specify hotel requirements per leg
8. **Accounting Integration**: Links to cost center and dimensions

---

This backend documentation provides complete business logic for the Travel Management module.
