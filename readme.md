
# Online Banking Web App

## Overview

This is a web-based banking application built using the Echo web framework. The app allows users to manage their finances by providing features such as account creation, deposit, withdrawal, fund transfer, transaction history, and more.

## Features

- **Account Types:** Users can have two types of accounts: Checking and Savings.
- **Deposit and Withdrawal:** Perform transactions to deposit funds into or withdraw funds from your accounts.
- **Fund Transfer:** Easily transfer funds between your checking and savings accounts.
- **Transaction History:** View a detailed history of your transactions over time.
- **Automatic Checking Account Creation:** A checking account is automatically created for users upon signing up to the app.

## Technologies Used

- [Echo](https://github.com/labstack/echo): A high-performance, minimalist Go web framework.
- [SQLC](https://github.com/kyleconroy/sqlc): Generates type-safe Go code from SQL queries.
- [Goose](https://github.com/pressly/goose): A database migration tool for Go.

## Setup

### Prerequisites

- Go installed
- Database (e.g., PostgreSQL) installed and running

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/markmumba/BankApp-.git
    cd online-banking-app
    ```

2. Run database migrations using Goose:

    ```bash
    goose up
    ```

3. Build and run the application:

    ```bash
    go build
    ./online-banking-app
    ```

4. Access the application in your web browser at [http://localhost:4000](http://localhost:4000).

## Usage

1. Sign up for an account to automatically get a checking account.
2. Log in to your account.
3. Explore the various features such as deposit, withdrawal, fund transfer, and transaction history.

## Contributing

If you find any issues or have suggestions for improvement, please feel free to open an issue or submit a pull request. Contributions are welcome!

## License

This project is licensed under the [MIT License](LICENSE).

---

Feel free to customize this template according to your specific application details. Include additional sections or information that you find relevant for your users and contributors.