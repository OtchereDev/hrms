# Employee Lifecycle Module - Frontend Documentation

## Overview

This document describes all user interactions, form behaviors, buttons, client-side validations, and data flows for the Employee Lifecycle module. It covers what users see and interact with in the web interface.

---

## Employee Onboarding

### Form View

**Route**: `/app/employee-onboarding/{name}`

**Access**: HR Manager, System Manager

#### Form Sections

**Section 1: Basic Information**
- **Left Column**:
  - Job Applicant (link field with search)
  - Job Offer (link field with search)
  - Employee Onboarding Template (link field)

- **Right Column**:
  - Company (read-only, auto-filled)
  - Boarding Status (read-only, badge)
  - Project (read-only, link appears after submit)

**Section 2: Employee Details**
- **Left Column**:
  - Employee (read-only, link appears when created)
  - Employee Name (read-only, auto-filled)
  - Department (auto-filled from template)
  - Designation (auto-filled from template)
  - Employee Grade (auto-filled from template)
  - Holiday List (optional)

- **Right Column**:
  - Date of Joining (date picker, required)
  - Boarding Begins On (date picker, required)

**Section 3: Onboarding Activities**
- Activities Table (editable grid)
  - Activity Name
  - User (link to assign specific user)
  - Role (link to assign all users with role)
  - Begin On (number of days)
  - Duration (number of days)
  - Required for Employee Creation (checkbox)
  - Description (text editor)

- Notify Users by Email (checkbox at bottom)

#### Interactive Behaviors

**1. Job Applicant Selection**:
```javascript
// Filters only "Accepted" applicants
frm.set_query("job_applicant", function() {
    return {
        filters: {
            "status": "Accepted"
        }
    };
});

// Auto-sets employee if already created
frappe.db.get_value("Employee", {
    "job_applicant": frm.doc.job_applicant
}, "name", (r) => {
    if (r && r.name) {
        frm.set_value("employee", r.name);
    }
});
```

**2. Job Offer Selection**:
```javascript
// Filters by selected job applicant and submitted offers
frm.set_query("job_offer", function() {
    return {
        filters: {
            "job_applicant": frm.doc.job_applicant,
            "docstatus": 1
        }
    };
});
```

**3. Template Selection**:
```javascript
// When template is selected, loads activities
frm.fields_dict.employee_onboarding_template.get_query = function() {
    return {
        filters: {
            "company": frm.doc.company
        }
    };
};

// On template change, populate activities table
frappe.model.with_doc("Employee Onboarding Template", template_name, function() {
    let template = frappe.model.get_doc("Employee Onboarding Template", template_name);

    frm.clear_table("activities");

    template.activities.forEach(activity => {
        let row = frm.add_child("activities");
        row.activity_name = activity.activity_name;
        row.user = activity.user;
        row.role = activity.role;
        row.begin_on = activity.begin_on;
        row.duration = activity.duration;
        row.description = activity.description;
        row.required_for_employee_creation = activity.required_for_employee_creation;
    });

    frm.refresh_field("activities");
});
```

#### Custom Buttons

**After Submit**:

1. **View → Employee** (if employee is created)
   ```javascript
   if (frm.doc.employee) {
       frm.add_custom_button(__("Employee"), function() {
           frappe.set_route("Form", "Employee", frm.doc.employee);
       }, __("View"));
   }
   ```

2. **View → Project**
   ```javascript
   if (frm.doc.project) {
       frm.add_custom_button(__("Project"), function() {
           frappe.set_route("Form", "Project", frm.doc.project);
       }, __("View"));
   }
   ```

3. **View → Task**
   ```javascript
   if (frm.doc.project) {
       frm.add_custom_button(__("Task"), function() {
           frappe.set_route("List", "Task", {
               "project": frm.doc.project
           });
       }, __("View"));
   }
   ```

4. **Create → Employee** (if not created yet)
   ```javascript
   if (!frm.doc.employee) {
       frm.add_custom_button(__("Employee"), function() {
           frappe.model.open_mapped_doc({
               method: "hrms.hr.doctype.employee_onboarding.employee_onboarding.make_employee",
               frm: frm
           });
       }, __("Create"));
       frm.page.set_inner_btn_group_as_primary(__("Create"));
   }
   ```

5. **Mark as Completed** (if status is Pending or In Process)
   ```javascript
   if (frm.doc.boarding_status !== "Completed") {
       frm.add_custom_button(__("Mark as Completed"), function() {
           frappe.confirm(
               __("Are you sure you want to mark all activities as completed?"),
               function() {
                   frappe.call({
                       method: "hrms.hr.doctype.employee_onboarding.employee_onboarding.mark_onboarding_as_completed",
                       args: {
                           onboarding_name: frm.doc.name
                       },
                       callback: function() {
                           frm.reload_doc();
                           frappe.show_alert({
                               message: __("Onboarding marked as completed"),
                               indicator: "green"
                           });
                       }
                   });
               }
           );
       });
   }
   ```

#### List View

**Route**: `/app/employee-onboarding`

**Default Filters**: `boarding_status = "Pending"`

**List Columns**:
- Employee Name
- Date of Joining
- Department
- Boarding Status (with color indicators)

**Indicator Colors**:
```javascript
get_indicator: function(doc) {
    let colors = {
        "Pending": "orange",
        "In Process": "yellow",
        "Completed": "green"
    };
    return [__(doc.boarding_status), colors[doc.boarding_status], "boarding_status,=," + doc.boarding_status];
}
```

**Bulk Actions**: None

---

## Employee Separation

### Form View

**Route**: `/app/employee-separation/{name}`

**Access**: System Manager

#### Form Sections

**Section 1: Basic Information**
- **Left Column**:
  - Employee (link field, required)
  - Employee Name (read-only)
  - Department (read-only)
  - Designation (read-only)
  - Employee Grade (read-only)

- **Right Column**:
  - Company (read-only)
  - Boarding Status (read-only, badge)
  - Resignation Letter Date (read-only, from employee)
  - Boarding Begins On (date picker, required)
  - Project (read-only, link after submit)

**Section 2: Separation Activities**
- Employee Separation Template (link field)
- Activities Table (same as onboarding)
- Notify Users by Email (checkbox)

**Section 3: Exit Interview**
- Exit Interview (text editor) - Exit interview summary

#### Interactive Behaviors

**Employee Selection**:
```javascript
// Auto-fetches employee details
frm.set_value("employee_name", employee.employee_name);
frm.set_value("company", employee.company);
frm.set_value("department", employee.department);
frm.set_value("designation", employee.designation);
frm.set_value("employee_grade", employee.grade);
frm.set_value("resignation_letter_date", employee.resignation_letter_date);
```

**Template Selection**:
Same behavior as onboarding - loads activities from template.

#### Custom Buttons

**After Submit**:

1. **View → Employee**
2. **View → Project**
3. **View → Task**

Same implementation as onboarding.

#### List View

**Route**: `/app/employee-separation`

**Default Filters**: `boarding_status = "Pending"`

**List Columns**:
- Employee Name
- Department
- Boarding Status

---

## Employee Transfer

### Form View

**Route**: `/app/employee-transfer/{name}`

**Access**: HR Manager, HR User, Employee (read-only)

#### Form Sections

**Section 1: Basic Information**
- **Left Column**:
  - Employee (link, required, filtered to Active only)
  - Employee Name (read-only)
  - Transfer Date (date picker, required)

- **Right Column**:
  - Company (read-only)
  - New Company (link, for inter-company transfers)
  - Department (read-only, bold)

**Section 2: Employee Transfer Details**
- Transfer Details Table (read-only grid)
  - Property
  - Current
  - New
- Reallocate Leaves (checkbox, hidden)
- Create New Employee ID (checkbox)
- New Employee ID (read-only, shows after submit if created)

#### Interactive Behaviors

**Add Employee Property Button**:
```javascript
frm.add_custom_button(__("Add Employee Property"), function() {
    // Opens dialog
    let dialog = new frappe.ui.Dialog({
        title: __("Add Property"),
        fields: [
            {
                fieldname: "property",
                fieldtype: "Autocomplete",
                label: __("Property"),
                options: get_employee_properties(),  // List of all employee fields
                reqd: 1
            },
            {
                fieldname: "current",
                fieldtype: "Data",
                label: __("Current Value"),
                read_only: 1
            },
            {
                fieldname: "new",
                fieldtype: "Dynamic",  // Type changes based on selected property
                label: __("New Value"),
                reqd: 1
            }
        ],
        primary_action_label: __("Add"),
        primary_action: function(values) {
            // Validate
            if (values.current === values.new) {
                frappe.throw(__("New value cannot be same as current value"));
            }

            // Check duplicate
            let exists = frm.doc.transfer_details.find(d => d.property === values.property);
            if (exists) {
                frappe.throw(__("Property {0} already added", [values.property]));
            }

            // Add row
            let row = frm.add_child("transfer_details");
            row.property = values.property;
            row.fieldname = get_fieldname(values.property);
            row.current = values.current;
            row.new = values.new;

            frm.refresh_field("transfer_details");
            dialog.hide();
        }
    });

    // When property is selected, fetch current value
    dialog.fields_dict.property.$input.on("change", function() {
        let property = dialog.get_value("property");
        let fieldname = get_fieldname(property);

        frappe.db.get_value("Employee", frm.doc.employee, fieldname, (r) => {
            dialog.set_value("current", r[fieldname]);

            // Change field type based on property type
            let field_type = get_field_type(fieldname);
            dialog.fields_dict.new.df.fieldtype = field_type;
            dialog.fields_dict.new.refresh();
        });
    });

    dialog.show();
});
```

**Employee Query**:
```javascript
frm.set_query("employee", function() {
    return {
        filters: {
            "status": "Active"
        }
    };
});
```

**Grid Behavior**:
```javascript
// Hide "Add Row" button
frm.fields_dict.transfer_details.grid.cannot_add_rows = true;
```

#### Validation Messages

```javascript
// Before submit
if (frm.doc.transfer_date > frappe.datetime.get_today()) {
    frappe.msgprint(__("Transfer cannot be done for future dates."));
    frappe.validated = false;
}
```

#### List View

**Route**: `/app/employee-transfer`

**List Columns**:
- Employee
- Transfer Date

---

## Employee Promotion

### Form View

**Route**: `/app/employee-promotion/{name}`

**Access**: HR Manager, HR User, Employee (read-only)

#### Form Sections

**Section 1: Basic Information**
- **Left Column**:
  - Employee (link, filtered to Active)
  - Employee Name (read-only)
  - Department (read-only)
  - Salary Currency (read-only)

- **Right Column**:
  - Promotion Date (date picker, required)
  - Company (read-only)

**Section 2: Employee Promotion Details**
- Description: "Set the properties that should be updated in the Employee master on promotion submission"
- Promotion Details Table (read-only grid, same structure as transfer)

**Section 3: Salary Details**
- **Left Column**: Current CTC (currency)
- **Right Column**: Revised CTC (currency)

#### Interactive Behaviors

**Add Employee Property**: Same as transfer

**CTC Dependency**:
```javascript
// Revised CTC depends on Current CTC
frm.set_df_property("revised_ctc", "reqd", frm.doc.current_ctc ? 1 : 0);
```

#### List View

**Route**: `/app/employee-promotion`

**List Columns**:
- Employee
- Promotion Date

---

## Exit Interview

### Form View

**Route**: `/app/exit-interview/{name}`

**Access**: System Manager

#### Form Sections

**Section 1: Basic Information**
- **Left Column**:
  - Naming Series
  - Employee (link, required)
  - Employee Name (read-only)
  - Email (read-only)

- **Right Column**:
  - Company (link, required)
  - Status (select: Pending/Scheduled/Completed/Cancelled)
  - Date (date picker, required if Scheduled)

**Section 2: Employee Details** (collapsible)
- Department, Designation, Reports To
- Date of Joining, Relieving Date

**Section 3: Exit Questionnaire**
- Reference Doctype (link)
- Questionnaire Email Sent (checkbox, read-only)
- Reference Document Name (dynamic link)

**Section 4: Interview Details**
- Interviewers (table multiselect)
  - User
- Interview Summary (text editor)

**Section 5: Final Decision**
- Employee Status (select: Employee Retained/Exit Confirmed, required if Completed)

#### Interactive Behaviors

**Employee Selection Validation**:
```javascript
frm.set_query("employee", function() {
    return {
        query: "hrms.hr.doctype.exit_interview.exit_interview.get_employees_for_exit",
        filters: {
            "company": frm.doc.company
        }
    };
});

// Validates relieving date exists
frappe.db.get_value("Employee", frm.doc.employee, "relieving_date", (r) => {
    if (!r.relieving_date) {
        frappe.msgprint(__("Please set Relieving Date for employee {0}", [frm.doc.employee]));
        frm.set_value("employee", "");
    }
});
```

**Custom Buttons**:

**Send Exit Questionnaire** (if not sent):
```javascript
if (!frm.doc.questionnaire_email_sent && !frm.is_new()) {
    frm.add_custom_button(__("Send Exit Questionnaire"), function() {
        frappe.call({
            method: "hrms.hr.doctype.exit_interview.exit_interview.send_exit_questionnaire",
            args: {
                interviews: [frm.doc.name]
            },
            callback: function(r) {
                if (r.message) {
                    frappe.msgprint(__("Exit Questionnaire sent successfully"));
                    frm.reload_doc();
                }
            }
        });
    });
}
```

#### List View

**Route**: `/app/exit-interview`

**List Columns**:
- Employee Name
- Date
- Relieving Date
- Employee Status
- Status

**List Actions**:

**Send Exit Questionnaires** (bulk action):
```javascript
list_view.page.add_actions_menu_item(__("Send Exit Questionnaires"), function() {
    let selected = list_view.get_checked_items();
    frappe.call({
        method: "hrms.hr.doctype.exit_interview.exit_interview.send_exit_questionnaire",
        args: {
            interviews: selected.map(d => d.name)
        },
        callback: function(r) {
            list_view.clear_checked_items();
            list_view.refresh();
        }
    });
});
```

**Indicator Colors**:
```javascript
get_indicator: function(doc) {
    let colors = {
        "Pending": "orange",
        "Scheduled": "yellow",
        "Completed": "green",
        "Cancelled": "red"
    };
    return [__(doc.status), colors[doc.status], "status,=," + doc.status];
}
```

---

## Appointment Letter

### Form View

**Route**: `/app/appointment-letter/{name}`

**Access**: HR Manager, System Manager

#### Form Sections

**Section 1: Basic Information**
- **Left Column**:
  - Job Applicant (link, required)
  - Applicant Name (read-only)

- **Right Column**:
  - Company (link, required)
  - Appointment Date (date picker, required)
  - Appointment Letter Template (link, required)

**Section 2: Body**
- Introduction (long text, required)
- Terms (table)
  - Offer Term
  - Value
- Closing Notes (text)

**Section 3: Printing Details** (collapsible)
- Letter Head (link)
- Print Heading (link)

#### Interactive Behaviors

**Template Selection**:
```javascript
frm.add_fetch("appointment_letter_template", "introduction", "introduction");
frm.add_fetch("appointment_letter_template", "closing_notes", "closing_notes");

frappe.ui.form.on("Appointment Letter", "appointment_letter_template", function(frm) {
    if (frm.doc.appointment_letter_template) {
        frappe.call({
            method: "hrms.hr.doctype.appointment_letter.appointment_letter.get_appointment_letter_details",
            args: {
                template: frm.doc.appointment_letter_template
            },
            callback: function(r) {
                if (r.message) {
                    frm.set_value("introduction", r.message.introduction);
                    frm.set_value("closing_notes", r.message.closing_notes);

                    frm.clear_table("terms");
                    r.message.terms.forEach(term => {
                        let row = frm.add_child("terms");
                        row.offer_term = term.title;
                        row.value = term.description;
                    });
                    frm.refresh_field("terms");
                }
            }
        });
    }
});
```

#### List View

**Route**: `/app/appointment-letter`

**List Columns**:
- Applicant Name
- Company

---

## Common UI Patterns

### Date Pickers
- Min/Max date validation
- Date format: System default
- Calendar popup

### Link Fields
- Autocomplete search
- Quick view on hover
- Click to navigate

### Tables
- Editable grids with inline editing
- Add/remove rows
- Column sorting
- Bulk operations

### Status Badges
- Color-coded indicators
- Real-time status updates
- Filterable in list views

### Notifications
- Success: Green toast
- Error: Red dialog
- Warning: Orange banner
- Info: Blue alert

---

## Real-time Updates

**Boarding Status**:
```javascript
// Listens for project updates
frappe.realtime.on("update_boarding_status", function(data) {
    if (cur_frm && cur_frm.doc.name === data.onboarding_name) {
        cur_frm.set_value("boarding_status", data.status);
        cur_frm.refresh();
    }
});
```

---

## Error Handling

### Client-side Validation
- Required field checks
- Date range validation
- Duplicate prevention
- Format validation

### Server Error Display
```javascript
callback: function(r) {
    if (r.exc) {
        frappe.msgprint({
            title: __("Error"),
            indicator: "red",
            message: r.exc
        });
    }
}
```

---

This frontend documentation provides complete details of all user interactions and client-side behaviors for the Employee Lifecycle module.
