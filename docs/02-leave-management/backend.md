# Leave Management Module - Backend Documentation

## Overview

The Leave Management module handles employee leave applications, allocations, policies, encashments, and related workflows. It uses a ledger-based system for accurate leave balance tracking.

---

## Core Doctypes

### 1. Leave Application

**Purpose**: Employee requests for time off

**Auto-naming**: `HR-LAP-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee requesting leave |
| leave_type | Link | Yes | Type of leave |
| from_date | Date | Yes | Leave start date |
| to_date | Date | Yes | Leave end date |
| half_day | Check | No | Half day leave |
| half_day_date | Date | Cond. | Date for half day (if half_day) |
| total_leave_days | Float | - | Calculated leave days |
| description | Small Text | - | Reason for leave |
| leave_approver | Link | - | Approver (from employee/department) |
| leave_balance | Float | - | Balance before application |
| status | Select | Yes | Open/Approved/Rejected/Cancelled |
| salary_slip | Link | - | Linked salary slip (for LWP) |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_active_employee()
    validate_dates()  # to_date >= from_date, half_day_date in range
    validate_balance_leaves()  # Check sufficient balance
    validate_leave_overlap()  # No overlapping leaves
    validate_max_days()  # Check max_continuous_days_allowed
    validate_block_days()  # Check against leave block list
    validate_salary_processed_days()  # For LWP
    validate_attendance()  # No present attendance
    validate_optional_leave()  # Against holiday list
    validate_applicable_after()  # Min working days check
    set_half_day_date()  # Auto-set if from==to
```

**Leave Days Calculation**:
```python
def get_number_of_leave_days(from_date, to_date, half_day, half_day_date, employee, leave_type):
    # Exclude holidays unless leave_type.include_holiday
    # Exclude weekly offs from holiday list
    # Count working days between dates
    # Subtract 0.5 if half_day

    working_days = get_working_days(from_date, to_date, holiday_list)

    if half_day and half_day_date:
        working_days -= 0.5

    return working_days
```

**On Submit**:
```python
def on_submit(self):
    # Only Approved/Rejected can be submitted
    if status not in ["Approved", "Rejected"]:
        frappe.throw("Only Approved/Rejected applications can be submitted")

    # Create leave ledger entry
    if status == "Approved":
        create_leave_ledger_entry(
            employee=employee,
            leave_type=leave_type,
            from_date=from_date,
            to_date=to_date,
            leaves=-total_leave_days,  # Negative for deduction
            transaction_type="Leave Application",
            transaction_name=name
        )

    # Create/update attendance
    update_attendance()

    # Notify employee
    notify_employee()
```

**Attendance Creation**:
```python
def update_attendance(self):
    if status == "Approved":
        for date in get_dates_between(from_date, to_date):
            # Skip holidays unless include_holiday
            if is_holiday(date) and not include_holiday:
                continue

            # Check if attendance exists
            attendance = get_attendance(employee, date)

            if attendance:
                if attendance.docstatus == 1:
                    frappe.throw(f"Attendance already marked for {date}")
                attendance.status = "Half Day" if is_half_day(date) else "On Leave"
                attendance.leave_type = leave_type
                attendance.leave_application = name
                attendance.save()
            else:
                create_attendance(
                    employee=employee,
                    attendance_date=date,
                    status="Half Day" if is_half_day(date) else "On Leave",
                    leave_type=leave_type,
                    leave_application=name
                )
```

---

### 2. Leave Allocation

**Purpose**: Allocates leave balance to employees

**Auto-naming**: `HR-LAL-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee receiving allocation |
| leave_type | Link | Yes | Type of leave |
| from_date | Date | Yes | Allocation period start |
| to_date | Date | Yes | Allocation period end |
| new_leaves_allocated | Float | - | New leaves being allocated |
| carry_forward | Check | No | Add unused from previous |
| unused_leaves | Float | - | Carried forward amount |
| total_leaves_allocated | Float | Yes | Total = new + carry forward |
| leave_period | Link | - | Associated leave period |
| leave_policy_assignment | Link | - | From policy if applicable |
| expired | Check | No | Allocation expired flag |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_period()  # to_date > from_date
    validate_allocation_overlap()  # No overlapping allocations
    validate_lwp()  # Cannot allocate LWP
    set_total_leaves_allocated()
    validate_leave_allocation_days()  # Against max_leaves_allowed
```

**Carry Forward Logic**:
```python
def set_total_leaves_allocated(self):
    if carry_forward:
        # Get previous allocation
        prev = get_previous_allocation(employee, leave_type, from_date)

        if prev:
            # Calculate unused leaves
            allocated = prev.total_leaves_allocated
            used = get_leave_ledger_balance(employee, leave_type, prev.to_date)
            unused_leaves = allocated - abs(used)

            # Apply max carry forward limit
            if leave_type.maximum_carry_forwarded_leaves:
                unused_leaves = min(unused_leaves, leave_type.maximum_carry_forwarded_leaves)

            # Mark carry forward in previous allocation
            prev.carry_forwarded_leaves_count = unused_leaves
            prev.save()

    total_leaves_allocated = new_leaves_allocated + unused_leaves
```

**On Submit**:
```python
def on_submit(self):
    # Create ledger entries
    if carry_forward and unused_leaves:
        # Separate entry for carry forward with expiry
        create_leave_ledger_entry(
            employee=employee,
            leave_type=leave_type,
            from_date=from_date,
            to_date=from_date + expire_days if expire_days else to_date,
            leaves=unused_leaves,
            is_carry_forward=1,
            transaction_type="Leave Allocation"
        )

    # Entry for new allocation
    create_leave_ledger_entry(
        employee=employee,
        leave_type=leave_type,
        from_date=from_date,
        to_date=to_date,
        leaves=new_leaves_allocated,
        transaction_type="Leave Allocation"
    )

    # Expire unused from previous allocation
    if carry_forward:
        expire_previous_allocation_balance()
```

**Earned Leave Scheduling**:
```python
# For earned leave types
def create_earned_leave_schedule(self):
    schedule = []

    frequency_map = {
        "Monthly": 1,
        "Quarterly": 3,
        "Half-Yearly": 6,
        "Yearly": 12
    }

    months = frequency_map[earned_leave_frequency]
    current_date = from_date

    while current_date <= to_date:
        if allocate_on_day == "Last Day":
            allocation_date = get_last_day_of_period(current_date, months)
        else:
            allocation_date = current_date

        schedule.append({
            "date": allocation_date,
            "leaves": get_monthly_earned_leave(annual_allocation, frequency)
        })

        current_date = add_months(current_date, months)

    return schedule
```

---

### 3. Leave Type

**Purpose**: Master data for leave types

**Auto-naming**: By field `leave_type_name`

#### Key Fields

| Field | Type | Description |
|-------|------|-------------|
| leave_type_name | Data | Unique leave type name |
| max_leaves_allowed | Float | Maximum per period |
| applicable_after | Int | Minimum working days |
| max_continuous_days_allowed | Int | Maximum consecutive days |
| is_carry_forward | Check | Can carry forward |
| maximum_carry_forwarded_leaves | Float | Max CF amount |
| expire_carry_forwarded_leaves_after_days | Int | CF expiry days |
| is_lwp | Check | Leave without pay |
| is_ppl | Check | Partially paid leave |
| fraction_of_daily_salary_per_leave | Float | For PPL (0-1) |
| is_optional_leave | Check | Optional holiday |
| allow_negative | Check | Allow negative balance |
| include_holiday | Check | Count holidays in leave |
| is_compensatory | Check | Compensatory leave |
| allow_encashment | Check | Can be encashed |
| earning_component | Link | Salary component for encashment |
| is_earned_leave | Check | Earned periodically |
| earned_leave_frequency | Select | Monthly/Quarterly/etc |
| rounding | Select | 0.25/0.5/1.0 |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_lwp()  # Cannot be LWP if allocations exist
    validate_leave_types()  # Cannot be both compensatory and earned
    # Cannot be both LWP and PPL
    # fraction_of_daily_salary_per_leave must be 0-1
```

---

### 4. Leave Policy & Leave Policy Assignment

**Leave Policy**: Template with annual allocations per leave type

**Leave Policy Assignment**: Assigns policy to employee for period

#### Key Fields (Assignment)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| leave_policy | Link | Yes | Policy to assign |
| assignment_based_on | Select | - | Leave Period/Joining Date |
| leave_period | Link | Cond. | Leave period |
| effective_from | Date | Yes | Start date |
| effective_to | Date | Yes | End date |
| carry_forward | Check | No | Carry forward unused |
| leaves_allocated | Check | - | Allocations created flag |

#### Business Logic

**On Submit**:
```python
def on_submit(self):
    grant_leave_alloc_for_employee()

def grant_leave_alloc_for_employee(self):
    for policy_detail in leave_policy.leave_policy_details:
        # Calculate pro-rated allocation
        if assignment_based_on == "Joining Date":
            prorate_factor = get_prorate_factor(employee.date_of_joining, effective_from, effective_to)
        else:
            prorate_factor = 1.0

        allocated_leaves = policy_detail.annual_allocation * prorate_factor

        # Handle earned leave
        if leave_type.is_earned_leave:
            allocated_leaves = 0  # Start with 0, earn over time

        # Create allocation
        allocation = frappe.new_doc("Leave Allocation")
        allocation.employee = employee
        allocation.leave_type = policy_detail.leave_type
        allocation.from_date = effective_from
        allocation.to_date = effective_to
        allocation.new_leaves_allocated = allocated_leaves
        allocation.leave_policy_assignment = name
        allocation.carry_forward = carry_forward

        # Create earned leave schedule
        if leave_type.is_earned_leave:
            allocation.earned_leave_schedule = get_earned_leave_schedule()

        allocation.insert()
        allocation.submit()

    leaves_allocated = 1
```

---

### 5. Leave Ledger Entry

**Purpose**: Central ledger for all leave transactions

**Naming**: Auto

**Cancellable**: Only expired allocations

#### Key Fields

| Field | Type | Description |
|-------|------|-------------|
| employee | Link | Employee |
| leave_type | Link | Leave type |
| transaction_type | Link | DocType (Allocation/Application/etc) |
| transaction_name | Dynamic Link | Specific document |
| leaves | Float | Positive=allocation, Negative=application |
| from_date | Date | Transaction start |
| to_date | Date | Transaction end |
| is_carry_forward | Check | CF allocation |
| is_expired | Check | Expired entry |

#### Business Logic

**Balance Calculation**:
```python
def get_leave_balance(employee, leave_type, on_date, consider_all_leaves=False):
    filters = {
        "employee": employee,
        "leave_type": leave_type,
        "from_date": ("<=", on_date),
        "docstatus": 1
    }

    if not consider_all_leaves:
        filters["to_date"] = (">=", on_date)

    entries = frappe.get_all("Leave Ledger Entry",
        filters=filters,
        fields=["leaves", "is_expired", "is_carry_forward", "to_date"]
    )

    balance = 0
    for entry in entries:
        # Skip expired CF leaves past their expiry
        if entry.is_carry_forward and entry.is_expired:
            continue

        balance += entry.leaves

    return balance
```

**Expiry Processing**:
```python
def expire_allocation(allocation):
    # Create expiry entry
    create_leave_ledger_entry(
        employee=allocation.employee,
        leave_type=allocation.leave_type,
        from_date=allocation.to_date,
        to_date=allocation.to_date,
        leaves=-allocation.unused_balance,  # Negative to deduct
        is_expired=1,
        transaction_type="Leave Allocation",
        transaction_name=allocation.name
    )
```

---

### 6. Leave Encashment

**Purpose**: Convert unused leave to cash

**Auto-naming**: `HR-LE-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| leave_period | Link | Yes | Leave period |
| leave_type | Link | Yes | Leave type to encash |
| leave_balance | Float | - | Current balance |
| actual_encashable_days | Float | - | After applying limits |
| encashment_days | Float | - | Actual days to encash |
| encashment_amount | Currency | - | Calculated amount |
| encashment_date | Date | Yes | Date of encashment |
| additional_salary | Link | - | Created additional salary |

#### Business Logic

**Calculation**:
```python
def set_encashment_amount(self):
    # Get salary structure
    salary_structure = get_assigned_salary_structure(employee, encashment_date)

    if not salary_structure:
        frappe.throw("No salary structure assigned")

    # Get daily rate
    per_day_encashment = frappe.db.get_value(
        "Salary Detail",
        {
            "parent": salary_structure,
            "salary_component": leave_type.earning_component
        },
        "amount"
    )

    if not per_day_encashment:
        frappe.throw("Earning component not found in salary structure")

    # Calculate
    encashment_amount = encashment_days * per_day_encashment
```

**On Submit**:
```python
def on_submit(self):
    # Create additional salary
    additional_salary = frappe.new_doc("Additional Salary")
    additional_salary.employee = employee
    additional_salary.salary_component = leave_type.earning_component
    additional_salary.amount = encashment_amount
    additional_salary.payroll_date = encashment_date
    additional_salary.insert()
    additional_salary.submit()

    self.additional_salary = additional_salary.name

    # Update leave allocation
    allocation.total_leaves_encashed += encashment_days
    allocation.save()

    # Create ledger entry
    create_leave_ledger_entry(
        employee=employee,
        leave_type=leave_type,
        from_date=encashment_date,
        to_date=encashment_date,
        leaves=-encashment_days,
        transaction_type="Leave Encashment",
        transaction_name=name
    )
```

---

### 7. Leave Block List

**Purpose**: Define dates when leave is restricted

**Auto-naming**: By field `leave_block_list_name`

#### Key Fields

| Field | Type | Description |
|-------|------|-------------|
| leave_block_list_name | Data | Unique name |
| company | Link | Company |
| applies_to_all_departments | Check | Company-wide |
| leave_block_list_dates | Table | Block dates with reasons |
| leave_block_list_allowed | Table | Users exempt from blocks |

#### Business Logic

**Validation in Leave Application**:
```python
def validate_block_days(self):
    block_dates = get_applicable_block_dates(employee, leave_type)

    leave_dates = get_dates_between(from_date, to_date)
    blocked = set(block_dates) & set(leave_dates)

    if blocked:
        # Check if user is in allow list
        if current_user not in get_allowed_users(block_list):
            frappe.throw(f"Leave is blocked for dates: {', '.join(blocked)}")
        else:
            # Show warning only
            frappe.msgprint(f"Note: These dates are blocked: {', '.join(blocked)}")
```

---

## Scheduled Jobs

### Allocate Earned Leaves (Daily Long)

```python
def allocate_earned_leaves():
    # Get all allocations with earned leave schedules
    allocations = frappe.get_all("Leave Allocation",
        filters={
            "docstatus": 1,
            "leave_type.is_earned_leave": 1
        }
    )

    for allocation in allocations:
        # Get pending schedule entries
        for schedule in allocation.earned_leave_schedule:
            if schedule.date == today() and not schedule.allocated:
                # Allocate leaves
                allocation.new_leaves_allocated += schedule.leaves
                allocation.total_leaves_allocated += schedule.leaves

                # Create ledger entry
                create_leave_ledger_entry(
                    employee=allocation.employee,
                    leave_type=allocation.leave_type,
                    from_date=schedule.date,
                    to_date=allocation.to_date,
                    leaves=schedule.leaves,
                    transaction_type="Leave Allocation"
                )

                schedule.allocated = 1
                allocation.save()
```

### Expire Allocations (Daily)

```python
def process_expired_allocation():
    # Get allocations past to_date
    allocations = frappe.get_all("Leave Allocation",
        filters={
            "docstatus": 1,
            "to_date": ("<", today()),
            "expired": 0
        }
    )

    for allocation in allocations:
        # Calculate unused balance
        balance = get_leave_balance(allocation.employee, allocation.leave_type, allocation.to_date)

        if balance > 0 and not allocation.leave_type.is_carry_forward:
            # Expire unused leaves
            expire_allocation(allocation)
            allocation.expired = 1
            allocation.save()
```

---

## Integration Points

### With Attendance

- Approved leaves create/update attendance records
- Status: "On Leave" or "Half Day"
- Linked via leave_application field

### With Salary Processing

- LWP leaves reduce salary
- Tracked in Salary Slip Leave child table
- PPL leaves calculate partial payment

### With Payroll

- Leave encashment creates Additional Salary
- Integrated in salary slip generation

---

## Formula Reference

### Leave Days Calculation
```python
working_days = 0
current = from_date

while current <= to_date:
    # Skip if holiday (unless include_holiday)
    if not is_holiday(current) or include_holiday:
        working_days += 1
    current = add_days(current, 1)

if half_day:
    working_days -= 0.5

return working_days
```

### Pro-rated Allocation
```python
total_period_days = date_diff(period_end, period_start) + 1
actual_period_days = date_diff(period_end, employee_start) + 1
prorate_factor = actual_period_days / total_period_days
allocated = annual_allocation * prorate_factor
```

### Earned Leave Monthly Amount
```python
frequency_months = {"Monthly": 1, "Quarterly": 3, "Half-Yearly": 6, "Yearly": 12}
months = frequency_months[frequency]
monthly_amount = annual_allocation / 12
per_period = monthly_amount * months

if rounding:
    per_period = round(per_period / rounding) * rounding
```

---

This backend documentation provides complete technical implementation details for the Leave Management module.
