# 🚀 TDD Workshop Repository

**Learn and practice Test-Driven Development (TDD) with Golang through hands-on exercises.**


## 📌 Overview
This repository is dedicated to learning **Test-Driven Development (TDD) in Golang** through practical exercises.  
The goal is to understand how TDD helps in writing better code, improving maintainability, and ensuring business logic correctness.


## 📚 Topics Covered

### 1️⃣ Understanding TDD: What, How, and Why?
- What is **Test-Driven Development (TDD)**?
- The **Red-Green-Refactor** cycle
- **Benefits** and best practices

### 2️⃣ Managing Dependencies with Test Doubles
- Using **Dummy, Stub, Spy, Fake, and Mock**
- Decoupling external services and dependencies
- Writing effective **unit tests**

### 3️⃣ Using Test Containers for Fake Databases
- Setting up **Testcontainer** for database testing
- Simulating **real database behavior** in tests
- Ensuring **migration and ORM behavior correctness**

### 4️⃣ Writing Meaningful Test Case Names
- Structuring test names to reflect **business cases**
- The **Given-When-Then** pattern
- Improving test **readability and maintainability**

## 🛠 Prerequisites
- Basic knowledge of **Golang**
- Understanding of **unit testing**
- Familiarity with **Go testing frameworks** (e.g., `testing`, `testify`)


## 🚀 Getting Started

### 🔹 Clone the repository
```sh
git clone https://github.com/kaweel/workshop-tdd.git
```

### 🚀 TDD Workshop - Makefile Commands
This project uses **Makefile** to streamline development and testing processes. Below are the available commands categorized by their functionality.

```
cd payment
```
🧪 Testing
```
✅ Run unit tests only
make unit-test

✅ Run integration tests only
make integration-test

✅ Run all tests and generate a coverage report
make test-gen-cov
```
📊 Code Coverage
```
✅ Show function-wise coverage report
make watch-cov-func

✅ Open HTML coverage report
make watch-cov
```
🐳 Database Management
```
✅ Start the database using Docker Compose
make db-up

✅ Stop and remove the database container
make db-down
```

🔗 Full Workflow (Start DB, Run app, Stop DB)
```
make start
```

### 📜 License
---
This project is licensed under the MIT License.

### 🚀 Happy Coding & TDD! 😊
