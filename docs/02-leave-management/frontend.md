# Leave Management Module - Frontend Documentation

## Overview

This document describes user interactions, form behaviors, and client-side logic for the Leave Management module, including both web interface and mobile app.

---

## Leave Application

### Web Form View

**Route**: `/app/leave-application/{name}`

**Access**: All employees (own applications), HR User, HR Manager, Leave Approver

#### Form Sections

**Section 1: Basic Information**
```
┌───────────────────────┬───────────────────────┐
│ Employee *            │ Leave Type *          │
│ [Auto: Current User]  │ [Dropdown]            │
│                       │                       │
│ Employee Name         │ Company               │
│ [Read-only]           │ [Read-only]           │
│                       │                       │
│ Department            │ Leave Balance         │
│ [Read-only]           │ [Read-only, Bold]     │
└───────────────────────┴───────────────────────┘
```

**Section 2: Dates & Duration**
```
┌───────────────────────┬───────────────────────┐
│ From Date *           │ To Date *             │
│ [Date picker]         │ [Date picker]         │
│                       │                       │
│ ☐ Half Day            │ Half Day Date         │
│                       │ [Date picker]         │
│                       │ (if half_day checked) │
│                       │                       │
│ Total Leave Days      │                       │
│ [Read-only, Bold]     │                       │
│ 5.0 days              │                       │
└───────────────────────┴───────────────────────┘
```

**Section 3: Details**
```
Description (Reason for Leave)
[Text area]
```

**Section 4: Approval**
```
┌───────────────────────┬───────────────────────┐
│ Leave Approver        │ Status                │
│ [Auto-filled]         │ [Dropdown]            │
│                       │ Open/Approved/Rejected│
│                       │                       │
│ Posting Date          │ ☐ Follow via Email    │
│ [Date: Today]         │                       │
└───────────────────────┴───────────────────────┘
```

#### Interactive Behaviors

**Auto-fill on Load**:
```javascript
frappe.ui.form.on("Leave Application", {
    onload: function(frm) {
        if (frm.is_new() && !frm.doc.employee) {
            // Set current user's employee
            frm.set_value("employee", frappe.defaults.get_user_default("Employee"));
        }
    },

    employee: function(frm) {
        // Fetch leave approver
        frappe.db.get_value("Employee", frm.doc.employee, "leave_approver", (r) => {
            if (r && r.leave_approver) {
                frm.set_value("leave_approver", r.leave_approver);
            } else {
                // Get from department
                get_department_approver(frm);
            }
        });
    }
});
```

**Leave Type Selection**:
```javascript
// Filter to show only allocated leave types
frm.set_query("leave_type", function() {
    return {
        query: "hrms.hr.doctype.leave_application.leave_application.get_leave_types",
        filters: {
            employee: frm.doc.employee,
            date: frm.doc.from_date || frappe.datetime.get_today()
        }
    };
});
```

**Leave Balance Display**:
```javascript
frappe.ui.form.on("Leave Application", {
    leave_type: function(frm) {
        if (frm.doc.employee && frm.doc.leave_type && frm.doc.from_date) {
            frappe.call({
                method: "hrms.hr.doctype.leave_application.leave_application.get_leave_balance_on",
                args: {
                    employee: frm.doc.employee,
                    date: frm.doc.from_date,
                    leave_type: frm.doc.leave_type
                },
                callback: function(r) {
                    frm.set_value("leave_balance", r.message);

                    // Show balance indicator
                    if (r.message < 1) {
                        frm.dashboard.set_headline_alert(
                            __("Insufficient leave balance: {0}", [r.message]),
                            "red"
                        );
                    }
                }
            });
        }
    }
});
```

**Real-time Leave Days Calculation**:
```javascript
function calculate_total_days(frm) {
    if (frm.doc.from_date && frm.doc.to_date) {
        frappe.call({
            method: "hrms.hr.doctype.leave_application.leave_application.get_number_of_leave_days",
            args: {
                employee: frm.doc.employee,
                leave_type: frm.doc.leave_type,
                from_date: frm.doc.from_date,
                to_date: frm.doc.to_date,
                half_day: frm.doc.half_day,
                half_day_date: frm.doc.half_day_date
            },
            callback: function(r) {
                if (r.message) {
                    frm.set_value("total_leave_days", r.message.leave_days);

                    // Show warnings if any
                    if (r.message.warnings) {
                        show_alert(r.message.warnings, "orange");
                    }
                }
            }
        });
    }
}

// Trigger on date changes
frm.fields_dict.from_date.$input.on("change", () => calculate_total_days(frm));
frm.fields_dict.to_date.$input.on("change", () => calculate_total_days(frm));
```

**Half Day Handling**:
```javascript
frappe.ui.form.on("Leave Application", {
    half_day: function(frm) {
        if (frm.doc.half_day) {
            // Show half day date field
            frm.set_df_property("half_day_date", "reqd", 1);

            // Auto-set if single day
            if (frm.doc.from_date === frm.doc.to_date) {
                frm.set_value("half_day_date", frm.doc.from_date);
            }
        } else {
            frm.set_df_property("half_day_date", "reqd", 0);
            frm.set_value("half_day_date", null);
        }
        calculate_total_days(frm);
    }
});
```

**Validation Messages**:
```javascript
// Before save
frappe.ui.form.on("Leave Application", {
    validate: function(frm) {
        // Date validation
        if (frm.doc.from_date > frm.doc.to_date) {
            frappe.msgprint(__("To Date cannot be before From Date"));
            frappe.validated = false;
        }

        // Half day validation
        if (frm.doc.half_day && !frm.doc.half_day_date) {
            frappe.msgprint(__("Half Day Date is required"));
            frappe.validated = false;
        }

        // Balance check (warning only)
        if (frm.doc.total_leave_days > frm.doc.leave_balance) {
            frappe.confirm(
                __("Leave balance is {0}. You are applying for {1} days. Do you want to continue?",
                    [frm.doc.leave_balance, frm.doc.total_leave_days]),
                () => {},
                () => { frappe.validated = false; }
            );
        }
    }
});
```

#### Dashboard

**Leave Allocation Dashboard** (shown in sidebar):
```javascript
// Display leave balances
frm.dashboard.add_section(
    frappe.render_template("leave_balance_dashboard", {
        allocations: get_leave_allocations(frm.doc.employee)
    })
);
```

**Dashboard Template**:
```html
<div class="leave-balance-dashboard">
    {% for alloc in allocations %}
    <div class="leave-balance-row">
        <span class="leave-type">{{ alloc.leave_type }}</span>
        <span class="leave-balance {{ 'low' if alloc.balance < 2 else '' }}">
            {{ alloc.balance }} / {{ alloc.total }}
        </span>
        <div class="leave-balance-bar">
            <div class="progress">
                <div class="progress-bar" style="width: {{ alloc.used_percent }}%"></div>
            </div>
        </div>
    </div>
    {% endfor %}
</div>
```

#### List View

**Route**: `/app/leave-application`

**Filters**: My Leaves (default), Team Leaves, All Leaves

**List Columns**:
- Employee Name
- Leave Type
- From Date → To Date
- Total Days
- Status (colored badge)

**Indicator Colors**:
```javascript
get_indicator: function(doc) {
    let colors = {
        "Open": "orange",
        "Approved": "green",
        "Rejected": "red",
        "Cancelled": "gray"
    };
    return [__(doc.status), colors[doc.status]];
}
```

**Quick Filters**:
```javascript
frappe.listview_settings["Leave Application"] = {
    onload: function(listview) {
        listview.page.add_inner_button(__("My Leaves"), function() {
            listview.filter_area.add([
                ["Leave Application", "employee", "=", frappe.defaults.get_user_default("Employee")]
            ]);
        });

        listview.page.add_inner_button(__("Team Leaves"), function() {
            // Filter by reporting manager
            get_team_members(listview);
        });
    }
};
```

---

## Mobile App (Vue.js)

### Leave Dashboard

**File**: `/frontend/src/views/leave/Dashboard.vue`

**Components**:
1. Leave Balance Cards
2. Request Leave Button
3. Recent Applications List
4. Upcoming Holidays

**Template Structure**:
```vue
<template>
  <base-layout page-title="Leave">
    <!-- Leave Balances -->
    <ion-grid>
      <ion-row>
        <ion-col v-for="balance in leaveBalances" :key="balance.leave_type">
          <leave-balance-card
            :leave-type="balance.leave_type"
            :available="balance.available"
            :total="balance.total"
          />
        </ion-col>
      </ion-row>
    </ion-grid>

    <!-- Request Leave Button -->
    <ion-button expand="block" @click="requestLeave">
      Request a Leave
    </ion-button>

    <!-- Recent Leaves -->
    <ion-list>
      <ion-list-header>Recent Leaves</ion-list-header>
      <leave-item
        v-for="leave in recentLeaves"
        :key="leave.name"
        :leave="leave"
      />
    </ion-list>

    <!-- Holidays -->
    <holidays-component :employee="employee" />
  </base-layout>
</template>
```

**Data Fetching**:
```javascript
export default {
  data() {
    return {
      leaveBalances: [],
      recentLeaves: []
    };
  },
  async mounted() {
    await this.fetchLeaveBalances();
    await this.fetchRecentLeaves();
  },
  methods: {
    async fetchLeaveBalances() {
      const { data } = await call("hrms.api.get_leave_balance_for_dashboard", {
        employee: this.employee
      });
      this.leaveBalances = data.message;
    },
    async fetchRecentLeaves() {
      const { data } = await this.$resources.leaveApplications.fetch({
        filters: {
          employee: this.employee,
          docstatus: 1
        },
        limit: 5,
        order_by: "modified desc"
      });
      this.recentLeaves = data;
    }
  }
};
```

### Leave Application Form

**File**: `/frontend/src/views/leave/Form.vue`

**Features**:
- Auto-fetch leave types with balance
- Date range picker
- Real-time leave days calculation
- Balance validation
- Offline support

**Template**:
```vue
<template>
  <form-view
    doctype="Leave Application"
    :doc="doc"
    @save="handleSave"
  >
    <!-- Leave Type with Balance -->
    <ion-item>
      <ion-label position="stacked">Leave Type *</ion-label>
      <ion-select v-model="doc.leave_type" @ionChange="onLeaveTypeChange">
        <ion-select-option
          v-for="type in leaveTypes"
          :key="type.name"
          :value="type.name"
        >
          {{ type.leave_type_name }} ({{ type.balance }} available)
        </ion-select-option>
      </ion-select>
    </ion-item>

    <!-- Date Range -->
    <ion-item>
      <ion-label position="stacked">From Date *</ion-label>
      <ion-datetime
        v-model="doc.from_date"
        presentation="date"
        @ionChange="calculateDays"
      />
    </ion-item>

    <ion-item>
      <ion-label position="stacked">To Date *</ion-label>
      <ion-datetime
        v-model="doc.to_date"
        presentation="date"
        :min="doc.from_date"
        @ionChange="calculateDays"
      />
    </ion-item>

    <!-- Half Day -->
    <ion-item>
      <ion-label>Half Day</ion-label>
      <ion-checkbox v-model="doc.half_day" @ionChange="calculateDays" />
    </ion-item>

    <ion-item v-if="doc.half_day">
      <ion-label position="stacked">Half Day Date *</ion-label>
      <ion-datetime
        v-model="doc.half_day_date"
        presentation="date"
        :min="doc.from_date"
        :max="doc.to_date"
      />
    </ion-item>

    <!-- Total Days Display -->
    <ion-item lines="none">
      <ion-label>
        <h2>Total Leave Days</h2>
        <p class="total-days">{{ totalLeaveDays }} days</p>
      </ion-label>
    </ion-item>

    <!-- Reason -->
    <ion-item>
      <ion-label position="stacked">Reason</ion-label>
      <ion-textarea v-model="doc.description" rows="4" />
    </ion-item>

    <!-- Balance Warning -->
    <ion-item v-if="showBalanceWarning" lines="none" color="warning">
      <ion-icon :icon="warningOutline" slot="start" />
      <ion-label class="ion-text-wrap">
        Insufficient balance. You have {{ leaveBalance }} days available.
      </ion-label>
    </ion-item>
  </form-view>
</template>
```

**Calculation Logic**:
```javascript
methods: {
  async calculateDays() {
    if (!this.doc.from_date || !this.doc.to_date) return;

    try {
      const { data } = await call(
        "hrms.hr.doctype.leave_application.leave_application.get_number_of_leave_days",
        {
          employee: this.doc.employee,
          leave_type: this.doc.leave_type,
          from_date: this.doc.from_date,
          to_date: this.doc.to_date,
          half_day: this.doc.half_day,
          half_day_date: this.doc.half_day_date
        }
      );

      this.totalLeaveDays = data.message.leave_days;

      // Show warning if exceeds balance
      this.showBalanceWarning = this.totalLeaveDays > this.leaveBalance;
    } catch (error) {
      console.error("Error calculating leave days:", error);
    }
  }
}
```

### Leave List View

**File**: `/frontend/src/views/leave/List.vue`

**Features**:
- Tabbed view: My Leaves / Team Leaves
- Filter by status, leave type, date range
- Pull to refresh
- Infinite scroll

**Template**:
```vue
<template>
  <base-layout page-title="Leave Applications">
    <!-- Tabs -->
    <ion-segment v-model="activeTab">
      <ion-segment-button value="my">My Leaves</ion-segment-button>
      <ion-segment-button value="team">Team Leaves</ion-segment-button>
    </ion-segment>

    <!-- Filters -->
    <ion-toolbar>
      <ion-buttons slot="start">
        <ion-button @click="openFilters">
          <ion-icon :icon="filterOutline" />
          Filters
        </ion-button>
      </ion-buttons>
    </ion-toolbar>

    <!-- List -->
    <ion-content>
      <ion-refresher slot="fixed" @ionRefresh="refresh">
        <ion-refresher-content />
      </ion-refresher>

      <ion-list>
        <leave-list-item
          v-for="leave in leaves"
          :key="leave.name"
          :leave="leave"
          @click="openLeave(leave)"
        />
      </ion-list>

      <ion-infinite-scroll @ionInfinite="loadMore">
        <ion-infinite-scroll-content />
      </ion-infinite-scroll>
    </ion-content>
  </base-layout>
</template>
```

---

## Common UI Patterns

### Date Validation
- From date cannot be in past (configurable)
- To date must be >= from date
- Half day date must be between from and to
- Visual date range selector

### Balance Display
- Color-coded (green: >5 days, yellow: 2-5 days, red: <2 days)
- Progress bar showing used/total
- Real-time updates

### Status Badges
```html
<span class="badge badge-{{ status.toLowerCase() }}">
  {{ status }}
</span>
```

### Error Handling
```javascript
try {
  await saveLeaveApplication();
  showToast("Leave application submitted successfully", "success");
  router.push("/leave");
} catch (error) {
  if (error.exc_type === "InsufficientLeaveBalanceError") {
    showAlert("Insufficient Leave Balance", error.message);
  } else {
    showAlert("Error", error.message);
  }
}
```

---

## Real-time Features

**Leave Balance Updates**:
```javascript
// Listen for allocation updates
frappe.realtime.on("leave_allocation_updated", (data) => {
  if (data.employee === current_employee) {
    refresh_leave_balances();
  }
});
```

**Approval Notifications**:
```javascript
// Listen for status changes
frappe.realtime.on("leave_application_update", (data) => {
  if (data.name === current_doc) {
    frappe.show_alert({
      message: `Leave ${data.status}`,
      indicator: data.status === "Approved" ? "green" : "red"
    });
    cur_frm.reload_doc();
  }
});
```

---

## Accessibility

- ARIA labels on all form fields
- Keyboard navigation support
- Screen reader announcements for status changes
- High contrast color scheme option
- Touch target sizes ≥ 44px

---

This frontend documentation provides complete details for implementing all user-facing features of the Leave Management module.
