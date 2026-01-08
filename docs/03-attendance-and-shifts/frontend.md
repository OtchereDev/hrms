# Attendance & Shifts Module - Frontend Documentation

## Overview

User interactions, form behaviors, and client-side features for attendance tracking, shift management, and employee check-ins across web and mobile interfaces.

---

## Attendance

### Form View

**Route**: `/app/attendance/{name}`

#### Form Layout

```
┌─── Attendance Details ───────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Naming Series    │ Status *             │ │
│ │ HR-ATT-.YYYY.-   │ [Present ▼]          │ │
│ │                  │                      │ │
│ │ Employee *       │ Company              │ │
│ │ [Search...]      │ [Read-only]          │ │
│ │                  │                      │ │
│ │ Employee Name    │ Department           │ │
│ │ [Read-only]      │ [Read-only]          │ │
│ │                  │                      │ │
│ │ Attendance Date* │ Shift                │ │
│ │ [📅 Date]        │ [Select shift]       │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Details (if shift selected) ──────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ In Time          │ Out Time             │ │
│ │ [Datetime]       │ [Datetime]           │ │
│ │                  │                      │ │
│ │ Working Hours    │                      │ │
│ │ [Calculated]     │                      │ │
│ │                  │                      │ │
│ │ ☐ Late Entry     │ ☐ Early Exit         │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Leave (if On Leave/Half Day) ─────────────┐
│ Leave Type * │ Leave Application             │
│ [Select]     │ [Link]                        │
└──────────────────────────────────────────────┘
```

#### Interactive Behaviors

**Status Change**:
```javascript
frappe.ui.form.on("Attendance", {
    status: function(frm) {
        // Show/hide leave fields
        let show_leave = ["On Leave", "Half Day"].includes(frm.doc.status);
        frm.toggle_reqd("leave_type", show_leave);
        frm.toggle_display("leave_type", show_leave);
        frm.toggle_display("leave_application", show_leave);

        // Auto-link leave application
        if (show_leave && frm.doc.employee && frm.doc.attendance_date) {
            frappe.call({
                method: "hrms.hr.doctype.attendance.attendance.get_leave_for_date",
                args: {
                    employee: frm.doc.employee,
                    date: frm.doc.attendance_date
                },
                callback: function(r) {
                    if (r.message) {
                        frm.set_value("leave_type", r.message.leave_type);
                        frm.set_value("leave_application", r.message.name);
                    }
                }
            });
        }
    }
});
```

**Working Hours Calculation**:
```javascript
function calculate_working_hours(frm) {
    if (frm.doc.in_time && frm.doc.out_time) {
        let hours = frappe.datetime.get_hour_diff(frm.doc.out_time, frm.doc.in_time);
        frm.set_value("working_hours", hours.toFixed(2));
    }
}

frm.fields_dict.in_time.$input.on("change", () => calculate_working_hours(frm));
frm.fields_dict.out_time.$input.on("change", () => calculate_working_hours(frm));
```

**Custom Buttons**:
```javascript
// After submit, if absent
if (frm.doc.status === "Absent" && frm.doc.docstatus === 1) {
    frm.add_custom_button(__("Create Attendance Request"), function() {
        frappe.model.open_mapped_doc({
            method: "hrms.hr.doctype.attendance.attendance.make_attendance_request",
            frm: frm
        });
    });
}
```

### Calendar View

**Route**: `/app/attendance/view/calendar`

**Features**:
- Monthly/weekly calendar views
- Color-coded attendance status
- Quick view on hover
- Click to edit

**Calendar Configuration**:
```javascript
frappe.views.calendar["Attendance"] = {
    field_map: {
        start: "attendance_date",
        end: "attendance_date",
        id: "name",
        title: "employee_name",
        status: "status"
    },
    get_events_method: "hrms.hr.doctype.attendance.attendance.get_events",
    get_css_class: function(data) {
        let color_map = {
            "Present": "success",
            "Absent": "danger",
            "On Leave": "info",
            "Half Day": "warning",
            "Work From Home": "primary"
        };
        return color_map[data.status] || "default";
    }
};
```

---

## Employee Checkin

### Form View

**Route**: `/app/employee-checkin/{name}`

#### Form Layout

```
┌─── Check-in Details ─────────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Employee *       │ Time *               │ │
│ │ [Auto: Current]  │ [Now]                │ │
│ │                  │                      │ │
│ │ Employee Name    │ Log Type             │ │
│ │ [Read-only]      │ [IN ▼]               │ │
│ │                  │                      │ │
│ │ Device/Location  │ Shift                │ │
│ │ [Text]           │ [Auto-detected]      │ │
│ │                  │                      │ │
│ │ ☐ Skip Auto Attendance                 │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Location (if geolocation enabled) ────────┐
│ Latitude  │ Longitude  │ [Fetch Geolocation] │
│ 28.6139   │ 77.2090    │                     │
│                                              │
│ [Map showing current location]               │
└──────────────────────────────────────────────┘
```

#### Mobile Quick Check-in

**Component**: Floating action button

```vue
<ion-fab vertical="bottom" horizontal="end">
  <ion-fab-button @click="checkIn" :color="isCheckedIn ? 'danger' : 'success'">
    <ion-icon :icon="isCheckedIn ? logOutOutline : logInOutline" />
  </ion-fab-button>
</ion-fab>
```

**Check-in Logic**:
```javascript
async checkIn() {
    try {
        // Get current location
        const position = await Geolocation.getCurrentPosition();

        // Create checkin
        const checkin = await this.$resources.employeeCheckin.insert({
            employee: this.employee,
            time: new Date().toISOString(),
            log_type: this.isCheckedIn ? "OUT" : "IN",
            latitude: position.coords.latitude,
            longitude: position.coords.longitude
        });

        // Show success
        this.$toast.success(`Checked ${this.isCheckedIn ? 'out' : 'in'} successfully`);

        // Toggle state
        this.isCheckedIn = !this.isCheckedIn;

    } catch (error) {
        if (error.exc_type === "CheckinRadiusExceededError") {
            this.$alert.show("Location Error", error.message);
        } else {
            this.$alert.show("Error", "Failed to check in");
        }
    }
}
```

---

## Shift Type

### Form View

**Route**: `/app/shift-type/{name}`

#### Form Sections

**Basic Information**:
```
Name *          │ Start Time *  │ End Time *
[8-Hour Shift]  │ [09:00]       │ [18:00]

Holiday List    │ Color
[Company Default] │ [Blue ▼]
```

**Auto Attendance Settings** (collapsible):
```
☑ Enable Auto Attendance

Check-in/Check-out Method:
○ Alternating entries as IN and OUT during the same shift
● Strictly based on Log Type in Employee Checkin

Working Hours Calculation:
○ First Check-in and Last Check-out
● Every Valid Check-in and Check-out

┌────────────────────────────────────────────┐
│ Begin Check-in Before: 60 minutes         │
│ Allow Check-out After:  60 minutes        │
│                                            │
│ Half Day Threshold:     4.0 hours         │
│ Absent Threshold:       1.0 hours         │
│                                            │
│ Process Attendance After: [📅 Date]       │
│ Last Sync:               [Datetime]       │
└────────────────────────────────────────────┘

☑ Auto Update Last Sync
```

#### Custom Buttons

**Mark Attendance** (manually trigger processing):
```javascript
frm.add_custom_button(__("Mark Attendance"), function() {
    frappe.call({
        method: "hrms.hr.doctype.shift_type.shift_type.process_auto_attendance",
        args: {
            shift_type: frm.doc.name
        },
        callback: function(r) {
            frappe.show_alert({
                message: __("Attendance marked for {0} employees", [r.message.count]),
                indicator: "green"
            });
        }
    });
});
```

---

## Shift Assignment

### Bulk Assignment Tool

**Route**: `/app/shift-assignment-tool`

**Purpose**: Assign shifts to multiple employees

#### Layout

```
┌─── Shift Assignment Tool ────────────────────┐
│                                              │
│ Action *                                     │
│ ● Assign Shift                               │
│ ○ Assign Shift Schedule                     │
│ ○ Process Shift Requests                    │
│                                              │
│ Company *  │ Shift Type *                    │
│ [Select]   │ [Select]                        │
│                                              │
│ Start Date *  │ End Date                     │
│ [📅 Date]     │ [📅 Optional]                │
│                                              │
│ Shift Location  │ Status                     │
│ [Optional]      │ [Active ▼]                 │
│                                              │
├─── Quick Filters ─────────────────────────────┤
│ Branch │ Department │ Designation │ Grade   │
│                                              │
├─── Select Employees ──────────────────────────┤
│ [Get Employees]                              │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │ ☑ Employee Name    Department    Grade │  │
│ ├────────────────────────────────────────┤  │
│ │ ☑ John Doe        Engineering    L3    │  │
│ │ ☑ Jane Smith      Sales          L4    │  │
│ │ ☐ Bob Johnson     Marketing      L2    │  │
│ └────────────────────────────────────────┘  │
│                                              │
│ [Assign Shifts to Selected Employees]        │
└──────────────────────────────────────────────┘
```

#### Interactive Features

**Get Employees**:
```javascript
frm.add_custom_button(__("Get Employees"), function() {
    frappe.call({
        method: "hrms.hr.doctype.shift_assignment_tool.shift_assignment_tool.get_employees",
        args: {
            company: frm.doc.company,
            filters: get_filters(frm)
        },
        callback: function(r) {
            if (r.message) {
                render_employee_grid(r.message);
            }
        }
    });
});
```

**Bulk Assign with Progress**:
```javascript
frm.add_custom_button(__("Assign Shifts"), function() {
    let selected = get_selected_employees();

    if (selected.length === 0) {
        frappe.msgprint(__("Please select employees"));
        return;
    }

    // Show progress dialog
    let progress = new frappe.ui.Progress({
        title: __("Assigning Shifts"),
        total: selected.length
    });

    frappe.call({
        method: "hrms.hr.doctype.shift_assignment_tool.shift_assignment_tool.bulk_assign",
        args: {
            employees: selected,
            shift_type: frm.doc.shift_type,
            start_date: frm.doc.start_date,
            end_date: frm.doc.end_date
        },
        callback: function(r) {
            progress.hide();
            if (r.message) {
                frappe.show_alert({
                    message: __("Shifts assigned to {0} employees", [r.message.success]),
                    indicator: "green"
                });
            }
        }
    });

    // Update progress (via realtime events)
    frappe.realtime.on("bulk_assign_progress", function(data) {
        progress.update_progress(data.completed, data.total);
    });
});
```

---

## Mobile App Features

### Attendance Dashboard

**File**: `/frontend/src/views/attendance/Dashboard.vue`

```vue
<template>
  <base-layout page-title="Attendance">
    <!-- Today's Status -->
    <ion-card>
      <ion-card-header>
        <ion-card-title>Today's Attendance</ion-card-title>
      </ion-card-header>
      <ion-card-content>
        <div class="attendance-status" :class="todayStatus">
          <ion-icon :icon="getStatusIcon(todayStatus)" />
          <span>{{ todayStatus }}</span>
        </div>

        <div class="checkin-info" v-if="lastCheckin">
          <p>Last check-in: {{ formatTime(lastCheckin.time) }}</p>
          <p>Type: {{ lastCheckin.log_type }}</p>
        </div>

        <!-- Quick Check-in Button -->
        <ion-button expand="block" @click="quickCheckin" :color="checkInButtonColor">
          <ion-icon :icon="isCheckedIn ? logOutOutline : logInOutline" slot="start" />
          {{ isCheckedIn ? 'Check Out' : 'Check In' }}
        </ion-button>
      </ion-card-content>
    </ion-card>

    <!-- Attendance Calendar -->
    <attendance-calendar :employee="employee" />

    <!-- Recent Attendance -->
    <ion-list>
      <ion-list-header>Recent Attendance</ion-list-header>
      <attendance-item
        v-for="att in recentAttendance"
        :key="att.name"
        :attendance="att"
      />
    </ion-list>
  </base-layout>
</template>
```

### Check-in List

**File**: `/frontend/src/views/attendance/EmployeeCheckinList.vue`

**Features**:
- Today's check-ins
- Check-in/out history
- Location map view
- Filter by date range

---

## Real-time Features

**Attendance Updates**:
```javascript
// Listen for auto-attendance creation
frappe.realtime.on("attendance_marked", function(data) {
    if (data.employee === current_employee) {
        frappe.show_alert({
            message: __("Attendance marked: {0}", [data.status]),
            indicator: "green"
        });

        // Refresh calendar if visible
        if (cur_list && cur_list.doctype === "Attendance") {
            cur_list.refresh();
        }
    }
});
```

**Shift Assignment Updates**:
```javascript
frappe.realtime.on("shift_assigned", function(data) {
    if (data.employee === current_employee) {
        frappe.show_alert({
            message: __("Shift {0} assigned", [data.shift_type]),
            indicator: "blue"
        });
    }
});
```

---

## List Views

### Attendance List

**Filters**:
- My Attendance (default)
- Team Attendance
- Department Attendance
- Status filter
- Date range

**Quick Actions**:
- Mark Present
- Mark Absent
- Bulk status update

### Shift Assignment Calendar

**Route**: `/app/shift-assignment/view/calendar`

**Features**:
- Color-coded by shift type
- Weekly/monthly view
- Drag & drop to reschedule
- Click to edit

---

## Error Handling

**Geolocation Errors**:
```javascript
try {
    const position = await getGeolocation();
} catch (error) {
    if (error.code === PositionError.PERMISSION_DENIED) {
        showAlert("Please enable location permissions");
    } else if (error.code === PositionError.POSITION_UNAVAILABLE) {
        // Allow check-in without location
        proceedWithoutLocation();
    }
}
```

**Shift Overlap Errors**:
```javascript
if (error.exc_type === "OverlappingShiftError") {
    frappe.confirm(
        __("Employee has overlapping shift. Do you want to proceed?"),
        () => {
            // Proceed with force flag
            save_with_override();
        }
    );
}
```

---

## Accessibility

- High contrast status indicators
- Touch-friendly check-in buttons (min 48px)
- Screen reader support for calendar
- Keyboard shortcuts for common actions

---

This frontend documentation covers all user interactions for the Attendance & Shifts module.
