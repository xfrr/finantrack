# 💸 FinanTrack - Personal Finance Manager

[![License](https://img.shields.io/github/license/yourusername/finantrack)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xfrr/finantrack)](https://goreportcard.com/report/github.com/xfrr/finantrack) [![Release](https://img.shields.io/github/v/release/yourusername/finantrack)](https://github.com/xfrr/finantrack/releases)

## 📊 Overview

**FinanTrack** is an open-source personal finance management tool designed to help you track income, expenses, and savings with ease. The app allows you to categorize your expenses, set financial goals, and generate insightful reports, giving you full control over your finances.

### 🌟 Features
- 🧾 **Track Expenses & Income**: Categorize expenses and gain a clear view of your finances.
- 💡 **Custom Financial Goals**: Set and track your savings goals and budgets.
- 📈 **Generate Reports**: Monthly, quarterly, and yearly reports on your financial health.
- 🔐 **Secure & Private**: Your data is securely stored and only accessible to you.
- 🌍 **Multi-Currency Support**: Easily manage finances across different currencies.

## 🚀 Getting Started

Follow these instructions to set up **FinanTrack** on your local machine.

### Prerequisites
- [Go](https://golang.org/doc/install) (>= 1.23)
- [PostgreSQL](https://www.postgresql.org/download/) for data storage
- [Git](https://git-scm.com/)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/xfrr/finantrack.git
   cd finantrack
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up `.env` file:
   ```bash
   cp .env.example .env
   ```
   Update the `.env` file with your database credentials and other configuration settings.

4. Initialize the database:
   ```bash
   just db:migrate
   ```

5. Start the application using Docker:
   ```bash
   just up
   ```
   The application will be accessible at http://localhost:8080.

## 🛠️ Commands
This repository uses [Just](https://github.com/casey/just) to manage commands. You can view all available commands by running:

```bash
just help
```


## 📜 Usage
Once the application is running, open your browser and navigate to http://localhost:8080. You can start managing your finances by creating accounts, adding expenses, and setting financial goals.

API Documentation
FinanTrack comes with a RESTful API for integration with external systems or mobile apps. You can check the API documentation by visiting http://localhost:8080/docs once the app is up and running.

## 🧪 Testing
Run the following command to execute the test suite:

## 📚 Contributing
We welcome contributions! To contribute:

1. Fork the repository.
2. Create a new branch (git checkout -b feature/my-feature).
3. Make your changes.
4. Commit your changes (git commit -am 'Add my feature').
5. Push to the branch (git push origin feature/my-feature).
6. Create a new Pull Request.

Please make sure to update tests as appropriate and follow the Contributing Guidelines.

## 📝 License
This project is licensed under the MIT License - see the LICENSE file for details.