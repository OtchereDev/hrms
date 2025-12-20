# Recruitment Module - Frontend Documentation

## Overview

User interactions, form behaviors, and client-side features for the complete recruitment pipeline from job requisitions through candidate selection and employee creation.

---

## Job Requisition

### Form View

**Route**: `/app/job-requisition/{name}`

#### Form Layout

```
┌─── Requisition Details ──────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Naming Series    │ Status *             │ │
│ │ HR-HIREQ-        │ [Pending ▼]          │ │
│ │                  │                      │ │
│ │ Designation *    │ Department           │ │
│ │ [Search...]      │ [Select]             │ │
│ │                  │                      │ │
│ │ Company *        │ No of Positions *    │ │
│ │ [Select]         │ [1]                  │ │
│ │                  │                      │ │
│ │ Expected Comp *  │ Requested By *       │ │
│ │ [Currency]       │ [Employee]           │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Timeline ─────────────────────────────────┐
│ Posting Date *   │ Expected By              │
│ [Today]          │ [📅 Date]                │
│                  │                          │
│ Completed On     │ Time to Fill             │
│ [Auto-filled]    │ [Calculated]             │
└──────────────────────────────────────────────┘

┌─── Details ──────────────────────────────────┐
│ Description *                                │
│ [Rich Text Editor]                           │
│                                              │
│ Reason for Requesting                        │
│ [Text Area]                                  │
└──────────────────────────────────────────────┘
```

#### Interactive Behaviors

**Dashboard Alert for Referrals**:
```javascript
frappe.ui.form.on("Job Requisition", {
    refresh: function(frm) {
        if (!["Filled", "On Hold", "Cancelled"].includes(frm.doc.status)) {
            // Check for pending employee referrals
            frappe.db.get_list("Employee Referral", {
                filters: {
                    for_designation: frm.doc.designation,
                    status: "Pending"
                }
            }).then((data) => {
                if (data && data.length) {
                    // Show clickable headline
                    let link = data.length > 1
                        ? "Employee Referrals"
                        : "Employee Referral";

                    frm.dashboard.set_headline(
                        `${data.length} ${link} open for this position.`,
                        "yellow"
                    );

                    // On click, navigate to referral list
                    $("#referral_links").on("click", () => {
                        frappe.set_route("List", "Employee Referral", {
                            for_designation: frm.doc.designation,
                            status: "Pending"
                        });
                    });
                }
            });
        }
    }
});
```

**Custom Actions (when status = "Open & Approved")**:

1. **Create Job Opening**:
```javascript
frm.add_custom_button("Create Job Opening", () => {
    frappe.model.open_mapped_doc({
        method: "hrms.hr.doctype.job_requisition.job_requisition.make_job_opening",
        frm: frm
    });
}, "Actions");
```
- Opens new Job Opening form
- Pre-fills data from requisition
- Links back to requisition

2. **Associate Job Opening**:
```javascript
frm.add_custom_button("Associate Job Opening", () => {
    frappe.prompt({
        label: "Job Opening",
        fieldname: "job_opening",
        fieldtype: "Link",
        options: "Job Opening",
        reqd: 1,
        get_query: () => {
            return {
                filters: {
                    company: frm.doc.company,
                    status: "Open",
                    designation: frm.doc.designation,
                    department: frm.doc.department || undefined
                }
            };
        }
    }, (values) => {
        frm.call("associate_job_opening", {
            job_opening: values.job_opening
        });
    }, "Associate Job Opening", "Submit");
}, "Actions");
```
- Shows dialog to select existing Job Opening
- Filters only matching, open positions
- Links opening to requisition

---

## Job Opening

### Form View

**Route**: `/app/job-opening/{name}`

#### Form Layout

```
┌─── Basic Information ────────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Job Title *      │ Status               │ │
│ │ [Text]           │ [Open ▼]             │ │
│ │                  │                      │ │
│ │ Designation *    │ Company *            │ │
│ │ [Link]           │ [Link]               │ │
│ │                  │                      │ │
│ │ Department       │ Location             │ │
│ │ [Link - filtered]│ [Link]               │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Posting Details ──────────────────────────┐
│ Posted On        │ Closes On   │ Closed On  │
│ [Datetime]       │ [Date]      │ [Auto]     │
│                                              │
│ Employment Type  │ Vacancies                │
│ [Link]           │ [1]                      │
└──────────────────────────────────────────────┘

┌─── Staffing Plan (auto-populated) ───────────┐
│ Staffing Plan    │ Planned Vacancies        │
│ [Read-only]      │ [Read-only]              │
│                                              │
│ ⓘ Auto-detected based on designation        │
└──────────────────────────────────────────────┘

┌─── Job Details ──────────────────────────────┐
│ Description                                  │
│ [Rich Text Editor]                           │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │ Salary Range                           │  │
│ │ Lower Range      │ Upper Range         │  │
│ │ [Currency]       │ [Currency]          │  │
│ │                  │                     │  │
│ │ ☐ Publish Salary Range                │  │
│ └────────────────────────────────────────┘  │
└──────────────────────────────────────────────┘

┌─── Publishing ───────────────────────────────┐
│ ☐ Publish on Website                         │
│                                              │
│ Route (URL slug)                             │
│ [Auto-generated or custom]                   │
└──────────────────────────────────────────────┘
```

#### Interactive Behaviors

**Department Filtering**:
```javascript
frm.set_query("department", function() {
    return {
        filters: {
            company: frm.doc.company
        }
    };
});
```

**Auto-fetch Staffing Plan**:
```javascript
frappe.ui.form.on("Job Opening", {
    designation: function(frm) {
        if (frm.doc.designation && frm.doc.company) {
            frappe.call({
                method: "hrms.hr.doctype.staffing_plan.staffing_plan.get_active_staffing_plan_details",
                args: {
                    company: frm.doc.company,
                    designation: frm.doc.designation,
                    date: frappe.datetime.now_date()
                },
                callback: function(data) {
                    if (data.message) {
                        frm.set_value("staffing_plan", data.message[0].name);
                        frm.set_value("planned_vacancies", data.message[0].vacancies);
                    } else {
                        frm.set_value("staffing_plan", "");
                        frm.set_value("planned_vacancies", 0);
                        frappe.show_alert({
                            indicator: "orange",
                            message: "No Staffing Plans found for this Designation"
                        });
                    }
                }
            });
        }
    },

    company: function(frm) {
        // Reset designation when company changes
        frm.set_value("designation", "");
    }
});
```

**Dashboard**:
- Shows count of linked Job Applicants
- Quick navigation to applications

---

## Job Applicant

### Form View

**Route**: `/app/job-applicant/{name}`

#### Form Layout

```
┌─── Applicant Information ────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Applicant Name * │ Email ID *           │ │
│ │ [Text]           │ [Email - unique]     │ │
│ │                  │                      │ │
│ │ Phone Number     │ Status *             │ │
│ │ [Phone]          │ [Open ▼]             │ │
│ │                  │                      │ │
│ │ Job Title        │ Designation          │ │
│ │ [Job Opening]    │ [Link]               │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Application Details ──────────────────────┐
│ Source           │ Source Name (if referral)│
│ [Link]           │ [Employee]               │
│                  │                          │
│ Applicant Rating                            │
│ ★★★★★ (interactive)                         │
│                                              │
│ Cover Letter                                 │
│ [Text Area]                                  │
│                                              │
│ Resume Attachment │ Resume Link             │
│ [File Upload]     │ [URL]                   │
│                                              │
│ Salary Expectations                          │
│ Lower Range      │ Upper Range              │
│ [Currency]       │ [Currency]               │
└──────────────────────────────────────────────┘
```

#### Interactive Dashboard

**Interview Summary Section**:
```javascript
frm.events.make_dashboard = function(frm) {
    frappe.call({
        method: "hrms.hr.doctype.job_applicant.job_applicant.get_interview_details",
        args: {
            job_applicant: frm.doc.name
        },
        callback: function(r) {
            if (r.message) {
                // Remove old dashboard
                $("div").remove(".form-dashboard-section.custom");

                // Render interview summary
                frm.dashboard.add_section(
                    frappe.render_template("job_applicant_dashboard", {
                        data: r.message.interviews,
                        number_of_stars: r.message.stars
                    }),
                    "Interview Summary"
                );
            }
        }
    });
};
```

Shows:
- List of all interviews
- Status (Pending/Cleared/Rejected)
- Average rating with star visualization
- Interview dates and rounds

#### Custom Actions

**Create Interview** (when status is not Rejected/Accepted):
```javascript
frm.add_custom_button("Interview", function() {
    // Show dialog to select interview round
    let d = new frappe.ui.Dialog({
        title: "Enter Interview Round",
        fields: [{
            label: "Interview Round",
            fieldname: "interview_round",
            fieldtype: "Link",
            options: "Interview Round"
        }],
        primary_action_label: "Create Interview",
        primary_action(values) {
            frappe.call({
                method: "hrms.hr.doctype.job_applicant.job_applicant.create_interview",
                args: {
                    doc: frm.doc,
                    interview_round: values.interview_round
                },
                callback: function(r) {
                    // Navigate to new interview
                    var doclist = frappe.model.sync(r.message);
                    frappe.set_route("Form", doclist[0].doctype, doclist[0].name);
                }
            });
            d.hide();
        }
    });
    d.show();
}, "Create");
```

**Create Job Offer** (when status = Accepted and no offer exists):
```javascript
frm.add_custom_button("Job Offer", function() {
    // Set route options to pre-fill offer
    frappe.route_options = {
        job_applicant: frm.doc.name,
        applicant_name: frm.doc.applicant_name,
        designation: frm.doc.job_opening || frm.doc.designation
    };
    frappe.new_doc("Job Offer");
}, "Create");
```

**View Job Offer** (when offer exists):
```javascript
frm.add_custom_button("Job Offer", function() {
    frappe.set_route("Form", "Job Offer", frm.doc.__onload.job_offer);
}, "View");
```

#### Email Integration

Form enables email communication directly from applicant record:
```javascript
cur_frm.email_field = "email_id";
```

### List View

**Route**: `/app/job-applicant`

#### Features

**Status Indicators**:
```javascript
frappe.listview_settings["Job Applicant"] = {
    add_fields: ["status"],
    get_indicator: function(doc) {
        if (doc.status == "Accepted") {
            return ["Accepted", "green", "status,=,Accepted"];
        } else if (["Open", "Replied"].includes(doc.status)) {
            return ["In Progress", "orange", "status,=," + doc.status];
        } else if (["Hold", "Rejected"].includes(doc.status)) {
            return [doc.status, "red", "status,=," + doc.status];
        }
    }
};
```

**Quick Filters**:
- Status (Open/Replied/Hold/Accepted/Rejected)
- Job Title
- Source
- Date range

**Bulk Actions**:
- Update status
- Assign to recruiter
- Send bulk emails

---

## Interview

### Form View

**Route**: `/app/interview/{name}`

#### Form Layout

```
┌─── Interview Details ────────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Interview Round* │ Job Applicant *      │ │
│ │ [Link]           │ [Link - filtered]    │ │
│ │                  │                      │ │
│ │ Designation      │ Status               │ │
│ │ [Auto-fetched]   │ [Pending ▼]          │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Schedule ─────────────────────────────────┐
│ ┌──────────────────────────────────────────┐ │
│ │ Scheduled On *   │ From Time   │ To Time│ │
│ │ [📅 Date]        │ [⏰ Time]   │ [Time] │ │
│ └──────────────────────────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Interviewers (Table) ─────────────────────┐
│ ┌──────────────────────────────────────────┐ │
│ │ Interviewer      │ Email               │ │
│ ├──────────────────────────────────────────┤ │
│ │ [Auto-populated from Interview Round]   │ │
│ │ John Doe         │ john@company.com    │ │
│ │ Jane Smith       │ jane@company.com    │ │
│ └──────────────────────────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Skill Assessment (Table) ─────────────────┐
│ ┌──────────────────────────────────────────┐ │
│ │ Skill            │ Rating              │ │
│ ├──────────────────────────────────────────┤ │
│ │ Python           │ ★★★★☆               │ │
│ │ Communication    │ ★★★★★               │ │
│ └──────────────────────────────────────────┘ │
│                                              │
│ Average Rating: 4.2/5.0                      │
└──────────────────────────────────────────────┘

┌─── Feedback Section (HTML) ──────────────────┐
│ [Dynamically rendered feedback from all      │
│  interviewers with ratings and comments]     │
└──────────────────────────────────────────────┘
```

#### Interactive Behaviors

**Auto-populate Interviewers**:
```javascript
frappe.ui.form.on("Interview", {
    interview_round: function(frm) {
        frm.set_value("job_applicant", "");

        frappe.call({
            method: "hrms.hr.doctype.interview.interview.get_interviewers",
            args: {
                interview_round: frm.doc.interview_round || ""
            },
            callback: function(r) {
                frm.clear_table("interview_details");
                r.message.forEach((interviewer) =>
                    frm.add_child("interview_details", interviewer)
                );
                refresh_field("interview_details");
            }
        });
    }
});
```

**Filter Job Applicants**:
```javascript
frm.set_query("job_applicant", function() {
    let filters = {
        status: ["!=", "Rejected"]
    };

    if (frm.doc.designation) {
        filters.designation = frm.doc.designation;
    }

    return { filters: filters };
});
```

#### Primary Action: Submit Feedback

**For Interviewers Only**:
```javascript
frm.page.set_primary_action("Submit Feedback", () => {
    // Get expected skills for this round
    frappe.call({
        method: "hrms.hr.doctype.interview.interview.get_expected_skill_set",
        args: {
            interview_round: frm.doc.interview_round
        },
        callback: function(r) {
            show_feedback_dialog(frm, r.message);
        }
    });
});
```

**Feedback Dialog**:
```javascript
let d = new frappe.ui.Dialog({
    title: "Submit Feedback",
    fields: [
        {
            fieldname: "skill_set",
            fieldtype: "Table",
            label: "Skill Assessment",
            cannot_add_rows: false,
            in_editable_grid: true,
            reqd: 1,
            fields: [
                // Dynamic fields from Skill Assessment doctype
            ],
            data: expected_skills  // Pre-populated
        },
        {
            fieldname: "result",
            fieldtype: "Select",
            options: ["", "Cleared", "Rejected"],
            label: "Result",
            reqd: 1
        },
        {
            fieldname: "feedback",
            fieldtype: "Small Text",
            label: "Feedback"
        }
    ],
    size: "large",
    minimizable: true,
    primary_action: function(values) {
        frappe.call({
            method: "hrms.hr.doctype.interview.interview.create_interview_feedback",
            args: {
                data: values,
                interview_name: frm.doc.name,
                interviewer: frappe.session.user,
                job_applicant: frm.doc.job_applicant
            }
        }).then(() => {
            frm.refresh();
        });
        d.hide();
    }
});
```

#### Custom Actions

**Reschedule Interview** (when status = Pending):
```javascript
frm.add_custom_button("Reschedule Interview", function() {
    let d = new frappe.ui.Dialog({
        title: "Reschedule Interview",
        fields: [
            {
                label: "Schedule On",
                fieldname: "scheduled_on",
                fieldtype: "Date",
                reqd: 1,
                default: frm.doc.scheduled_on
            },
            {
                label: "From Time",
                fieldname: "from_time",
                fieldtype: "Time",
                reqd: 1,
                default: frm.doc.from_time
            },
            {
                label: "To Time",
                fieldname: "to_time",
                fieldtype: "Time",
                reqd: 1,
                default: frm.doc.to_time
            }
        ],
        primary_action_label: "Reschedule",
        primary_action(values) {
            frm.call({
                method: "reschedule_interview",
                doc: frm.doc,
                args: values
            }).then(() => {
                frm.refresh();
                d.hide();
            });
        }
    });
    d.show();
}, "Actions");
```

#### Feedback Visualization

**Skill-wise Average Rating**:
```javascript
frappe.call({
    method: "hrms.hr.doctype.interview.interview.get_skill_wise_average_rating",
    args: { interview: frm.doc.name }
}).then((r) => {
    // Display in HTML section with visual bars
    render_skill_ratings(r.message);
});
```

**Reviews per Rating Distribution**:
```javascript
// Calculate percentage distribution
const reviews_per_rating = [0, 0, 0, 0, 0];  // 1-5 stars
frm.feedback.forEach((x) => {
    reviews_per_rating[Math.floor(x.total_score - 1)] += 1;
});

// Convert to percentages
frm.reviews_per_rating = reviews_per_rating.map((x) =>
    flt((x * 100) / frm.feedback.length, 1)
);

// Render as horizontal bar chart
```

---

## Job Offer

### Form View

**Route**: `/app/job-offer/{name}`

#### Form Layout

```
┌─── Offer Details ────────────────────────────┐
│ ┌──────────────────┬──────────────────────┐ │
│ │ Naming Series    │ Status               │ │
│ │ HR-OFF-.YYYY.-   │ [Awaiting Response▼] │ │
│ │                  │                      │ │
│ │ Job Applicant *  │ Applicant Name *     │ │
│ │ [Link]           │ [Auto-filled]        │ │
│ │                  │                      │ │
│ │ Applicant Email  │ Offer Date *         │ │
│ │ [Auto-filled]    │ [📅 Date]            │ │
│ │                  │                      │ │
│ │ Designation *    │ Company *            │ │
│ │ [Link]           │ [Link]               │ │
│ └──────────────────┴──────────────────────┘ │
└──────────────────────────────────────────────┘

┌─── Offer Terms ──────────────────────────────┐
│ Job Offer Term Template                      │
│ [Select template to auto-populate terms]     │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │ Offer Terms (Table)                    │  │
│ ├────────────────────────────────────────┤  │
│ │ Offer Term       │ Value               │  │
│ ├────────────────────────────────────────┤  │
│ │ Base Salary      │ $80,000/year        │  │
│ │ Joining Bonus    │ $5,000              │  │
│ │ Benefits         │ Health, Dental, 401k│  │
│ │ Vacation Days    │ 20 days             │  │
│ └────────────────────────────────────────┘  │
│                                              │
│ Additional Terms                             │
│ [Rich Text Editor]                           │
└──────────────────────────────────────────────┘

┌─── Terms and Conditions ─────────────────────┐
│ Select Terms (Template)                      │
│ [Link to Terms and Conditions]               │
│                                              │
│ Terms                                        │
│ [Rich Text - auto-populated from template]  │
└──────────────────────────────────────────────┘
```

#### Interactive Behaviors

**Job Offer Term Template Selection**:
```javascript
frappe.ui.form.on("Job Offer", {
    job_offer_term_template: function(frm) {
        if (!frm.doc.job_offer_term_template) return;

        frappe.db.get_doc("Job Offer Term Template", frm.doc.job_offer_term_template)
            .then((doc) => {
                // Clear existing terms
                frm.clear_table("offer_terms");

                // Add terms from template
                doc.offer_terms.forEach((term) => {
                    frm.add_child("offer_terms", term);
                });

                refresh_field("offer_terms");
            });
    }
});
```

**Terms and Conditions Template**:
```javascript
frappe.ui.form.on("Job Offer", {
    select_terms: function(frm) {
        erpnext.utils.get_terms(frm.doc.select_terms, frm.doc, function(r) {
            if (!r.exc) {
                frm.set_value("terms", r.message);
            }
        });
    }
});
```

**Email Field Configuration**:
```javascript
frappe.ui.form.on("Job Offer", {
    setup: function(frm) {
        frm.email_field = "applicant_email";
    }
});
```

#### Custom Actions

**Create Employee** (when status = Accepted and submitted):
```javascript
if (frm.doc.status == "Accepted" &&
    frm.doc.docstatus === 1 &&
    !frm.doc.__onload.employee) {

    frm.add_custom_button("Create Employee", function() {
        frappe.model.open_mapped_doc({
            method: "hrms.hr.doctype.job_offer.job_offer.make_employee",
            frm: frm
        });
    });
}
```

**Show Employee** (when employee created):
```javascript
if (frm.doc.__onload && frm.doc.__onload.employee) {
    frm.add_custom_button("Show Employee", function() {
        frappe.set_route("Form", "Employee", frm.doc.__onload.employee);
    });
}
```

---

## Web Form: Job Application

### Public Form

**Route**: `/job_application?job_opening={job_opening_name}`

#### Form Fields

```
┌─── Apply for Position ───────────────────────┐
│                                              │
│ Position: [Software Engineer]               │
│ [Read-only, from URL parameter]              │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │ Your Information                       │  │
│ ├────────────────────────────────────────┤  │
│ │ Full Name *                            │  │
│ │ [Text Input]                           │  │
│ │                                        │  │
│ │ Email Address *                        │  │
│ │ [Email Input - validated]              │  │
│ │                                        │  │
│ │ Phone Number                           │  │
│ │ [Phone Input]                          │  │
│ │                                        │  │
│ │ Country                                │  │
│ │ [Select]                               │  │
│ └────────────────────────────────────────┘  │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │ Cover Letter                           │  │
│ │ [Large Text Area]                      │  │
│ │                                        │  │
│ │ Resume Link                            │  │
│ │ [URL Input - validated]                │  │
│ │ ⓘ Link to online resume/portfolio     │  │
│ └────────────────────────────────────────┘  │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │ Salary Expectations                    │  │
│ │ Minimum        │ Maximum                │  │
│ │ [Currency]     │ [Currency]             │  │
│ └────────────────────────────────────────┘  │
│                                              │
│ [Submit Application]                         │
└──────────────────────────────────────────────┘
```

#### Client-side Validation

**URL Validation for Resume Link**:
```javascript
frappe.web_form.validate = () => {
    let resume = frappe.web_form.get_value("resume_link");

    if (resume && !isValidURL(resume)) {
        frappe.throw("Please enter a valid URL for resume");
    }
};

function isValidURL(string) {
    try {
        new URL(string);
        return true;
    } catch (_) {
        return false;
    }
}
```

#### Success Message

After submission:
```
┌──────────────────────────────────────────────┐
│          ✓ Application Submitted             │
│                                              │
│  Thank you for applying! We've received      │
│  your application and will review it soon.   │
│                                              │
│  You will receive an email confirmation at:  │
│  your.email@example.com                      │
│                                              │
│  [Return to Job Openings]                    │
└──────────────────────────────────────────────┘
```

---

## Dashboards and Reports

### Recruitment Dashboard

**Charts Available**:

1. **Job Application Status** (Pie Chart)
   - Open
   - Replied
   - Hold
   - Accepted
   - Rejected

2. **Job Applicants by Country** (Bar Chart)
   - Geographic distribution

3. **Job Applicant Source** (Donut Chart)
   - Referral
   - Career Site
   - LinkedIn
   - Other sources

4. **Job Application Frequency** (Line Chart)
   - Applications over time

5. **Job Offer Status** (Pie Chart)
   - Awaiting Response
   - Accepted
   - Rejected
   - Cancelled

6. **Job Applicant Pipeline** (Funnel Chart)
   - Total Applications → Interviews → Offers → Hires

### Number Cards

1. **Job Openings**
   - Count of open positions

2. **Job Offers (This Month)**
   - Monthly offers count

3. **Job Offer Acceptance Rate**
   - Formula: (Accepted / Total Submitted) × 100%

---

## Navigation and Workflows

### Typical User Journeys

**1. Recruiter: Posting a Job**
```
Job Requisition (Create & Submit)
    ↓ Approval
Job Requisition (Status: Open & Approved)
    ↓ [Create Job Opening]
Job Opening (Fill details)
    ↓ Save
☑ Publish on Website
```

**2. Candidate: Applying**
```
Browse Jobs (Website)
    ↓ Click "Apply"
Job Application Form
    ↓ Fill & Submit
Confirmation Message
    ↓
Email Confirmation Sent
```

**3. Recruiter: Screening Applicants**
```
Job Applicant List
    ↓ Click applicant
Review Application
    ↓ Update Status to "Replied"
[Create Interview]
    ↓ Select Round
Interview Form
    ↓ Schedule & Save
Email Sent to Interviewers
```

**4. Interviewer: Submitting Feedback**
```
Interview (from email link)
    ↓ Open form
[Submit Feedback]
    ↓ Fill skill ratings
Rate & Comment
    ↓ Submit
Interview Status Updated
```

**5. Recruiter: Making Offer**
```
Job Applicant (Status: Accepted)
    ↓ [Create Job Offer]
Job Offer Form
    ↓ Select Template
Offer Terms Auto-filled
    ↓ Review & Submit
Email Sent to Candidate
```

**6. HR: Onboarding New Hire**
```
Job Offer (Status: Accepted)
    ↓ [Create Employee]
Employee Form (Pre-filled)
    ↓ Complete details
Save Employee
    ↓
Job Applicant Status → Accepted
Job Offer Status → Accepted
Job Opening → Can be Closed
```

---

## Error Handling

### Validation Messages

**Duplicate Requisition**:
```
┌────────────────────────────────────┐
│ ⚠ Validation Error                │
│                                    │
│ Requisition already exists for     │
│ Software Engineer in Engineering   │
│                                    │
│ [OK]                               │
└────────────────────────────────────┘
```

**No Vacancies Available**:
```
┌────────────────────────────────────┐
│ ⚠ Cannot Create Offer             │
│                                    │
│ No vacancies available for this    │
│ designation according to the       │
│ staffing plan.                     │
│                                    │
│ [OK]                               │
└────────────────────────────────────┘
```

**Closed Job Opening**:
```
┌────────────────────────────────────┐
│ ⚠ Cannot Apply                    │
│                                    │
│ This job opening is closed.        │
│                                    │
│ [OK]                               │
└────────────────────────────────────┘
```

---

## Real-time Features

**Status Updates**:
```javascript
// Listen for offer status changes
frappe.realtime.on("job_offer_accepted", function(data) {
    if (data.job_applicant === current_applicant) {
        frappe.show_alert({
            message: "Candidate accepted the offer!",
            indicator: "green"
        });
        cur_frm.reload_doc();
    }
});
```

---

## Accessibility

- Keyboard navigation for all forms
- Screen reader labels for all fields
- High contrast status indicators
- ARIA labels for custom buttons
- Form validation with clear error messages

---

This frontend documentation covers all user interactions for the Recruitment module.
