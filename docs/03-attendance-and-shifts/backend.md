# Attendance & Shifts Module - Backend Documentation

## Overview

This module manages employee attendance tracking, shift scheduling, check-ins, and automated attendance marking. It supports flexible shift patterns, overnight shifts, geolocation validation, and auto-attendance from check-in/check-out logs.

---

## Core Doctypes

### 1. Attendance

**Purpose**: Records employee attendance for a specific date

**Auto-naming**: `HR-ATT-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| attendance_date | Date | Yes | Attendance date |
| status | Select | Yes | Present/Absent/On Leave/Half Day/Work From Home |
| shift | Link | No | Assigned shift |
| in_time | Datetime | No | Check-in time (from employee checkin) |
| out_time | Datetime | No | Check-out time (from employee checkin) |
| working_hours | Float | No | Calculated hours worked |
| late_entry | Check | No | Marked late |
| early_exit | Check | No | Left early |
| leave_type | Link | Cond. | Leave type (if On Leave/Half Day) |
| leave_application | Link | No | Linked leave application |
| attendance_request | Link | No | Created from attendance request |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_attendance_date()  # Cannot be before DOJ
    validate_duplicate_record()  # No duplicate for employee+date+shift
    validate_overlapping_shift_attendance()  # Check shift time overlaps
    validate_employee_status()  # Employee must be active
    check_leave_record()  # Auto-link leave applications
```

**Duplicate Prevention**:
```python
def validate_duplicate_record(self):
    # Check for existing attendance
    duplicate = frappe.db.exists("Attendance", {
        "employee": self.employee,
        "attendance_date": self.attendance_date,
        "shift": self.shift or "",
        "name": ("!=", self.name),
        "docstatus": ("<", 2)
    })

    if duplicate:
        frappe.throw(_("Attendance for {0} on {1} already exists").format(
            self.employee, self.attendance_date
        ))
```

**Leave Integration**:
```python
def check_leave_record(self):
    # Check for approved leave
    leave_record = frappe.db.get_value("Leave Application", {
        "employee": self.employee,
        "from_date": ("<=", self.attendance_date),
        "to_date": (">=", self.attendance_date),
        "status": "Approved",
        "docstatus": 1
    }, ["name", "leave_type", "half_day"], as_dict=True)

    if leave_record:
        self.leave_application = leave_record.name
        self.leave_type = leave_record.leave_type
        self.status = "Half Day" if leave_record.half_day else "On Leave"
```

**On Cancel**:
```python
def on_cancel(self):
    # Unlink from employee checkins
    frappe.db.sql("""
        UPDATE `tabEmployee Checkin`
        SET attendance = NULL
        WHERE attendance = %s
    """, self.name)
```

---

### 2. Employee Checkin

**Purpose**: Records employee check-in/check-out events

**Auto-naming**: Auto-generated

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| log_type | Select | No | IN/OUT |
| time | Datetime | Yes | Check-in/out timestamp |
| device_id | Data | No | Location/device identifier |
| shift | Link | No | Auto-detected shift |
| shift_start | Datetime | No | Shift start time |
| shift_end | Datetime | No | Shift end time |
| shift_actual_start | Datetime | No | With begin_check_in grace |
| shift_actual_end | Datetime | No | With allow_check_out grace |
| attendance | Link | No | Created attendance |
| skip_auto_attendance | Check | No | Skip auto-processing |
| latitude | Float | No | Geolocation lat |
| longitude | Float | No | Geolocation long |

#### Business Logic

**Shift Detection**:
```python
def fetch_shift(self):
    # Get shift for employee at check-in time
    shift_details = get_employee_shift(
        employee=self.employee,
        for_timestamp=self.time,
        consider_default_shift=True
    )

    if shift_details:
        self.shift = shift_details.shift_type.name
        self.shift_start = shift_details.start_datetime
        self.shift_end = shift_details.end_datetime
        self.shift_actual_start = shift_details.actual_start
        self.shift_actual_end = shift_details.actual_end
```

**Geolocation Validation**:
```python
def validate_distance_from_shift_location(self):
    if not self.shift:
        return

    shift_type = frappe.get_doc("Shift Type", self.shift)

    # Check if shift has assigned location
    shift_location = get_shift_location(self.employee, self.time)

    if shift_location and shift_location.checkin_radius:
        if self.latitude and self.longitude:
            distance = get_distance_between_coordinates(
                (self.latitude, self.longitude),
                (shift_location.latitude, shift_location.longitude)
            )

            if distance > shift_location.checkin_radius:
                frappe.throw(
                    _("You are {0}m away from shift location. Maximum allowed: {1}m").format(
                        int(distance), shift_location.checkin_radius
                    ),
                    exc=CheckinRadiusExceededError
                )
```

**Duplicate Prevention**:
```python
def validate_duplicate_log(self):
    # Prevent duplicate at exact same timestamp
    duplicate = frappe.db.exists("Employee Checkin", {
        "employee": self.employee,
        "time": self.time,
        "name": ("!=", self.name)
    })

    if duplicate:
        frappe.throw(_("Duplicate log entry found"))
```

---

### 3. Shift Type

**Purpose**: Defines shift timings and auto-attendance rules

**Auto-naming**: By field `name`

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | Data | Yes | Shift name |
| start_time | Time | Yes | Shift start (e.g., 09:00) |
| end_time | Time | Yes | Shift end (e.g., 18:00) |
| holiday_list | Link | No | Holiday calendar |
| enable_auto_attendance | Check | No | Enable auto-marking |
| begin_check_in_before_shift_start_time | Int | No | Grace before (minutes, default 60) |
| allow_check_out_after_shift_end_time | Int | No | Grace after (minutes, default 60) |
| working_hours_calculation_based_on | Select | No | First+Last / Every Valid |
| determine_check_in_and_check_out | Select | No | Alternating / Strictly Log Type |
| working_hours_threshold_for_half_day | Float | No | Hours for half day |
| working_hours_threshold_for_absent | Float | No | Hours for absent |
| process_attendance_after | Date | No | Don't process before this date |
| last_sync_of_checkin | Datetime | No | Last processed timestamp |
| enable_late_entry_marking | Check | No | Mark late entries |
| late_entry_grace_period | Int | No | Minutes after start |
| enable_early_exit_marking | Check | No | Mark early exits |
| early_exit_grace_period | Int | No | Minutes before end |
| allow_overtime | Check | No | Track overtime |
| overtime_type | Link | No | Overtime calculation type |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_same_start_and_end()  # Start != end
    validate_circular_shift()  # Not >= 24 hours
    validate_unlinked_logs()  # No pending checkins if changing start_time
```

**Auto-Attendance Processing**:
```python
def process_auto_attendance(self):
    # Get unprocessed checkins
    if auto_update_last_sync:
        from_time = last_sync_of_checkin
    else:
        from_time = process_attendance_after

    to_time = now_datetime()

    # Get all checkins for this shift
    logs = frappe.get_all("Employee Checkin",
        filters={
            "skip_auto_attendance": 0,
            "attendance": ("is", "not set"),
            "time": ("between", [from_time, to_time]),
            "shift": self.name
        },
        fields=["*"],
        order_by="time"
    )

    # Group by employee and shift start
    logs_by_employee = {}
    for log in logs:
        key = (log.employee, log.shift_start)
        if key not in logs_by_employee:
            logs_by_employee[key] = []
        logs_by_employee[key].append(log)

    # Process each group
    for (employee, shift_start), shift_logs in logs_by_employee.items():
        attendance = get_attendance(shift_logs)

        if attendance.working_hours > 0:
            # Create attendance
            att_doc = frappe.new_doc("Attendance")
            att_doc.employee = employee
            att_doc.attendance_date = getdate(shift_start)
            att_doc.shift = self.name
            att_doc.status = attendance.status
            att_doc.working_hours = attendance.working_hours
            att_doc.in_time = attendance.in_time
            att_doc.out_time = attendance.out_time
            att_doc.late_entry = attendance.late_entry
            att_doc.early_exit = attendance.early_exit
            att_doc.insert(ignore_permissions=True)
            att_doc.submit()

            # Link checkins to attendance
            mark_attendance_and_link_log(shift_logs, att_doc.name)

    # Mark absent for dates with no checkins
    mark_absent_for_dates_with_no_attendance()
```

**Attendance Determination**:
```python
def get_attendance(self, logs):
    # Determine IN/OUT pairs
    if determine_check_in_and_check_out == "Alternating entries as IN and OUT":
        in_out_pairs = []
        for i in range(0, len(logs), 2):
            if i + 1 < len(logs):
                in_out_pairs.append((logs[i], logs[i+1]))
    else:
        # Use log_type field
        in_logs = [l for l in logs if l.log_type == "IN"]
        out_logs = [l for l in logs if l.log_type == "OUT"]
        in_out_pairs = list(zip(in_logs, out_logs))

    # Calculate working hours
    if working_hours_calculation_based_on == "First Check-in and Last Check-out":
        if logs:
            in_time = min(l.time for l in logs)
            out_time = max(l.time for l in logs)
            working_hours = time_diff_in_hours(out_time, in_time)
    else:
        # Sum all valid pairs
        working_hours = 0
        for in_log, out_log in in_out_pairs:
            working_hours += time_diff_in_hours(out_log.time, in_log.time)

    # Determine status
    if working_hours_threshold_for_absent and working_hours < working_hours_threshold_for_absent:
        status = "Absent"
    elif working_hours_threshold_for_half_day and working_hours < working_hours_threshold_for_half_day:
        status = "Half Day"
    else:
        status = "Present"

    # Check late/early
    late_entry = in_time > (shift_start + timedelta(minutes=late_entry_grace_period))
    early_exit = out_time < (shift_end - timedelta(minutes=early_exit_grace_period))

    return AttendanceResult(
        status=status,
        working_hours=working_hours,
        in_time=in_time,
        out_time=out_time,
        late_entry=late_entry,
        early_exit=early_exit
    )
```

---

### 4. Shift Assignment

**Purpose**: Assigns shifts to employees for specific periods

**Auto-naming**: Auto-generated

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| shift_type | Link | Yes | Shift to assign |
| start_date | Date | Yes | Assignment start |
| end_date | Date | No | Assignment end (null = indefinite) |
| status | Select | No | Active/Inactive |
| shift_location | Link | No | Shift location for geofencing |
| shift_request | Link | No | Created from shift request |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_from_to_dates()  # start_date <= end_date
    validate_overlapping_shifts()  # Check for conflicts

def validate_overlapping_shifts(self):
    # Get overlapping assignments
    overlaps = frappe.db.sql("""
        SELECT name, shift_type, start_date, end_date
        FROM `tabShift Assignment`
        WHERE employee = %s
        AND name != %s
        AND status = 'Active'
        AND docstatus < 2
        AND (
            (start_date BETWEEN %s AND %s)
            OR (end_date BETWEEN %s AND %s)
            OR (start_date <= %s AND (end_date >= %s OR end_date IS NULL))
        )
    """, (self.employee, self.name, self.start_date, self.end_date or "2099-12-31",
          self.start_date, self.end_date or "2099-12-31",
          self.start_date, self.end_date or "2099-12-31"))

    if overlaps:
        # Check if shifts have overlapping time windows
        for overlap in overlaps:
            if has_overlapping_timings(self.shift_type, overlap[1]):
                frappe.throw(_("Overlapping shift assignment found"))
```

**Shift Timing Overlap**:
```python
def has_overlapping_timings(shift1_name, shift2_name):
    shift1 = frappe.get_cached_doc("Shift Type", shift1_name)
    shift2 = frappe.get_cached_doc("Shift Type", shift2_name)

    # Convert to minutes from midnight
    s1_start = time_to_minutes(shift1.start_time)
    s1_end = time_to_minutes(shift1.end_time)
    s2_start = time_to_minutes(shift2.start_time)
    s2_end = time_to_minutes(shift2.end_time)

    # Handle overnight shifts
    if s1_end < s1_start:
        s1_end += 1440  # Add 24 hours
    if s2_end < s2_start:
        s2_end += 1440

    # Check overlap
    return not (s1_end <= s2_start or s2_end <= s1_start)
```

---

### 5. Shift Request

**Purpose**: Employee requests for shift changes

**Auto-naming**: Auto-generated

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee requesting |
| shift_type | Link | Yes | Requested shift |
| from_date | Date | Yes | Request start |
| to_date | Date | No | Request end |
| approver | Link | Yes | Shift request approver |
| status | Select | Yes | Draft/Approved/Rejected |
| company | Link | Yes | Company |

#### Business Logic

**On Submit (if Approved)**:
```python
def on_submit(self):
    if self.status == "Approved":
        # Create shift assignment
        assignment = frappe.new_doc("Shift Assignment")
        assignment.employee = self.employee
        assignment.shift_type = self.shift_type
        assignment.start_date = self.from_date
        assignment.end_date = self.to_date
        assignment.shift_request = self.name
        assignment.status = "Active"
        assignment.insert()
        assignment.submit()

**Notifications**:
```python
# After insert
def after_insert(self):
    notify_approver()

# On update
def on_update(self):
    if self.status in ["Approved", "Rejected"]:
        notify_employee()
```

---

### 6. Shift Schedule & Shift Schedule Assignment

**Shift Schedule**: Template defining recurring shift patterns

**Shift Schedule Assignment**: Applies schedule to employee with auto-creation

#### Key Fields (Shift Schedule)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | Data | Yes | Schedule name |
| shift_type | Link | Yes | Shift to apply |
| frequency | Select | Yes | Every Week/2 Weeks/3 Weeks/4 Weeks |
| repeat_on_days | Table | Yes | Days of week to apply |

#### Key Fields (Shift Schedule Assignment)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| shift_schedule | Link | Yes | Schedule template |
| shift_location | Link | No | Location for geofencing |
| shift_status | Select | No | Active/Inactive |
| enabled | Check | No | Auto-create shifts |
| create_shifts_after | Date | Cond. | Next creation date |

#### Business Logic

**Automatic Shift Creation**:
```python
def create_shifts(self, start_date, end_date=None):
    if not end_date:
        end_date = add_days(start_date, 90)  # Create 90 days ahead

    schedule = frappe.get_doc("Shift Schedule", self.shift_schedule)

    # Get frequency in weeks
    frequency_weeks = int(schedule.frequency.split()[1]) if "Week" in schedule.frequency else 1

    # Get days to apply
    days_to_apply = [d.day for d in schedule.repeat_on_days]

    current_date = start_date
    assignments = []

    while current_date <= end_date:
        day_name = current_date.strftime("%A")

        # Check if this day matches schedule
        week_offset = (current_date - start_date).days // 7
        if week_offset % frequency_weeks == 0 and day_name in days_to_apply:
            # Check for existing assignment
            exists = frappe.db.exists("Shift Assignment", {
                "employee": self.employee,
                "shift_type": schedule.shift_type,
                "start_date": current_date,
                "docstatus": ("<", 2)
            })

            if not exists:
                # Create assignment
                assignment = create_shift_assignment(
                    employee=self.employee,
                    shift_type=schedule.shift_type,
                    start_date=current_date,
                    shift_schedule_assignment=self.name
                )
                assignments.append(assignment)

        current_date = add_days(current_date, 1)

    # Update next creation date
    self.create_shifts_after = end_date

    return assignments
```

**Scheduled Job**:
```python
def process_auto_shift_creation():
    # Run hourly_long
    # Get all enabled assignments
    assignments = frappe.get_all("Shift Schedule Assignment",
        filters={
            "enabled": 1,
            "create_shifts_after": ("<=", today())
        }
    )

    for assignment_name in assignments:
        assignment = frappe.get_doc("Shift Schedule Assignment", assignment_name)
        try:
            assignment.create_shifts(
                start_date=assignment.create_shifts_after,
                end_date=add_days(today(), 90)
            )
            assignment.save()
        except Exception as e:
            frappe.log_error(f"Shift creation failed for {assignment_name}", e)
```

---

### 7. Attendance Request

**Purpose**: Employees request attendance marking for past dates

**Auto-naming**: Auto-generated

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee |
| from_date | Date | Yes | Period start |
| to_date | Date | Yes | Period end |
| reason | Select | Yes | Work From Home/On Duty |
| half_day | Check | No | Half day request |
| half_day_date | Date | Cond. | Half day date |
| shift | Link | No | Shift (won't overwrite existing) |
| include_holidays | Check | No | Mark on holidays too |

#### Business Logic

**On Submit**:
```python
def on_submit(self):
    create_attendance_records()

def create_attendance_records(self):
    dates = get_dates_between(self.from_date, self.to_date)

    for date in dates:
        # Skip if holiday (unless include_holidays)
        if is_holiday(date, self.employee) and not self.include_holidays:
            continue

        # Skip if on leave
        if has_leave_application(self.employee, date):
            continue

        # Check existing attendance
        existing = frappe.db.get_value("Attendance", {
            "employee": self.employee,
            "attendance_date": date,
            "docstatus": ("<", 2)
        }, "name")

        status = get_attendance_status(self.reason, self.half_day, date, self.half_day_date)

        if existing:
            # Update existing
            att = frappe.get_doc("Attendance", existing)
            if att.docstatus == 0:
                att.status = status
                att.attendance_request = self.name
                if self.shift:
                    att.shift = self.shift
                att.save()
                att.submit()
        else:
            # Create new
            att = frappe.new_doc("Attendance")
            att.employee = self.employee
            att.attendance_date = date
            att.status = status
            att.shift = self.shift
            att.attendance_request = self.name
            att.insert()
            att.submit()
```

---

## Scheduled Jobs

### 1. Update Last Sync of Checkin (Hourly Long)

```python
def update_last_sync_of_checkin():
    # Get shifts with auto_update_last_sync enabled
    shifts = frappe.get_all("Shift Type",
        filters={
            "enable_auto_attendance": 1,
            "auto_update_last_sync": 1
        }
    )

    for shift_name in shifts:
        shift = frappe.get_doc("Shift Type", shift_name)

        # Calculate actual shift end (with grace period)
        shift_end = combine_datetime(today(), shift.end_time)
        shift_end = shift_end + timedelta(minutes=shift.allow_check_out_after_shift_end_time or 60)

        # Only update if shift has ended
        if now_datetime() > shift_end:
            shift.last_sync_of_checkin = shift_end + timedelta(minutes=1)
            shift.save()
```

### 2. Process Auto Attendance (Hourly Long)

```python
def process_auto_attendance_for_all_shifts():
    shifts = frappe.get_all("Shift Type",
        filters={"enable_auto_attendance": 1}
    )

    for shift_name in shifts:
        shift = frappe.get_doc("Shift Type", shift_name)
        try:
            shift.process_auto_attendance()
        except Exception as e:
            frappe.log_error(f"Auto attendance failed for {shift_name}", e)
```

### 3. Process Auto Shift Creation (Hourly Long)

Covered in Shift Schedule Assignment section above.

---

## Formulas & Calculations

### Working Hours

```python
def time_diff_in_hours(end_time, start_time):
    diff = end_time - start_time
    return diff.total_seconds() / 3600
```

### Shift Time Detection

```python
def get_employee_shift(employee, for_timestamp, consider_default_shift=False):
    # Get active shift assignment for timestamp
    assignment = frappe.db.get_value("Shift Assignment",
        filters={
            "employee": employee,
            "start_date": ("<=", getdate(for_timestamp)),
            "docstatus": 1,
            "status": "Active",
            "OR": [
                ["end_date", ">=", getdate(for_timestamp)],
                ["end_date", "is", "not set"]
            ]
        },
        fieldname=["shift_type", "start_date", "end_date"],
        as_dict=True
    )

    if not assignment and consider_default_shift:
        # Check employee default shift
        assignment = get_employee_default_shift(employee)

    if assignment:
        shift_type = frappe.get_cached_doc("Shift Type", assignment.shift_type)

        # Calculate actual start/end with grace periods
        shift_date = getdate(for_timestamp)
        shift_start = combine_datetime(shift_date, shift_type.start_time)
        shift_end = combine_datetime(shift_date, shift_type.end_time)

        # Handle overnight shifts
        if shift_type.end_time < shift_type.start_time:
            # Determine if timestamp is for current day or next day
            if for_timestamp.time() < shift_type.end_time:
                shift_start = shift_start - timedelta(days=1)
            else:
                shift_end = shift_end + timedelta(days=1)

        actual_start = shift_start - timedelta(minutes=shift_type.begin_check_in_before_shift_start_time or 60)
        actual_end = shift_end + timedelta(minutes=shift_type.allow_check_out_after_shift_end_time or 60)

        return ShiftDetails(
            shift_type=shift_type,
            start_datetime=shift_start,
            end_datetime=shift_end,
            actual_start=actual_start,
            actual_end=actual_end
        )

    return None
```

---

This backend documentation provides complete technical implementation details for the Attendance & Shifts module.
