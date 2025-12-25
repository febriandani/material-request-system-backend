
---

## 1. Database & Naming Convention

**Database name**

```sql
managed_material_db
```

**Naming rules (recommended)**

* Database: `snake_case`
* Table: plural (`material_requests`)
* PK: `id`
* FK: `{table}_id`
* Timestamps: `created_at`, `updated_at`

---

## 2. Master Tables (Dropdown Data)

### 2.1 Departments

```sql
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### 2.2 Request Types

```sql
CREATE TABLE request_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### 2.3 Users (Requester / Approver / Admin)

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(50), -- requester, approver, admin
    department_id INT REFERENCES departments(id),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### 2.4 Units of Measure (UOM)

```sql
CREATE TABLE uoms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL -- pcs, kg, liter, mtr, set
);
```

---

### 2.5 Categories (for parts)

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL -- Bearing, Oil & Lubricant, etc
);
```

---

## 3. Material Request (Header)

### 3.1 Material Requests

```sql
CREATE TABLE material_requests (
    id SERIAL PRIMARY KEY,
    mr_number VARCHAR(50) NOT NULL UNIQUE,
    department_id INT REFERENCES departments(id),
    requested_by VARCHAR(100) NOT NULL,
    request_date DATE NOT NULL,
    request_type_id INT REFERENCES request_types(id),
    approver_id INT REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Status values**

* OPEN
* HOLD
* IN_PROGRESS
* DONE
* CLOSED

(You can later convert this to ENUM if you want)

---

## 4. Material Request Items (Detail)

### 4.1 Material Request Items

```sql
CREATE TABLE material_request_items (
    id SERIAL PRIMARY KEY,
    material_request_id INT REFERENCES material_requests(id) ON DELETE CASCADE,

    part_no VARCHAR(50) NOT NULL,
    part_name VARCHAR(100) NOT NULL,
    description VARCHAR(100),

    qty_requested NUMERIC(10,2) NOT NULL,
    uom_id INT REFERENCES uoms(id),

    qty_issued NUMERIC(10,2) DEFAULT 0,
    on_hand NUMERIC(10,2) DEFAULT 0,

    location VARCHAR(100),
    balance NUMERIC(10,2) GENERATED ALWAYS AS (qty_issued - qty_requested) STORED,

    remarks VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);
```

✔ Supports **1–50 items per MR**
✔ Balance auto-calculated
✔ Cascade delete when MR deleted

---

## 5. Inventory / Availability Parts

### 5.1 Parts Master (Inventory)

```sql
CREATE TABLE parts (
    id SERIAL PRIMARY KEY,
    part_no VARCHAR(50) NOT NULL UNIQUE,
    part_name VARCHAR(100) NOT NULL,
    description VARCHAR(100),

    category_id INT REFERENCES categories(id),

    quantity NUMERIC(10,2) NOT NULL,
    min_stock NUMERIC(10,2),
    max_stock NUMERIC(10,2),

    uom_id INT REFERENCES uoms(id),
    location VARCHAR(100),

    price NUMERIC(15,2),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

---

### 5.2 Stock Status (Derived View – Recommended)

Instead of storing status manually, **derive it automatically** 👇

```sql
CREATE VIEW parts_availability AS
SELECT
    p.*,
    CASE
        WHEN quantity <= 0 THEN 'OUT_OF_STOCK'
        WHEN quantity < min_stock THEN 'LOW_STOCK'
        ELSE 'AVAILABLE'
    END AS stock_status
FROM parts p;
```

✔ No data inconsistency
✔ Always real-time

---

## 6. Status History (Audit Trail – IMPORTANT)

### 6.1 Material Request Status History

```sql
CREATE TABLE material_request_status_histories (
    id SERIAL PRIMARY KEY,
    material_request_id INT REFERENCES material_requests(id) ON DELETE CASCADE,
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    changed_by INT REFERENCES users(id),
    changed_at TIMESTAMP DEFAULT NOW(),
    remarks VARCHAR(200)
);
```

✔ Tracks **who changed what & when**
✔ Mandatory for approval systems

---

## 7. Create Database (Run First)

```sql
CREATE DATABASE managed_material_db;
```

Then connect:

```sql
\c managed_material_db
```

---

## 8. Indexes (Performance)

```sql
CREATE INDEX idx_mr_number ON material_requests(mr_number);
CREATE INDEX idx_part_no ON parts(part_no);
CREATE INDEX idx_mr_status ON material_requests(status);
```

---

## 9. Summary ERD (Simple)

```
users ───────┐
             ├── material_requests ─── material_request_items
departments ─┘

parts ── categories
      └─ uoms

material_requests ── status_histories
```

---# material-request-system-backend
