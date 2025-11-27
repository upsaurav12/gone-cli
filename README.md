# 🐹⚡ Go Bootstrapper

**Go Bootstrapper** is a CLI tool that scaffolds production-ready Golang projects — no dependency headaches, no manual setup.  
Just run a command and get a fully configured project with linters, routers, and structure ready to code.

* * *

## ✨ Features

*   🏗 **Create new Golang projects instantly** — skip the boilerplate setup.
    
*   ⚡ **Framework-ready templates** — built-in support for `Gin`, `Chi`, and more.
    
*   📂 **Standardized structure** — organized directories: `cmd/`, `internal/`, `router/`, etc.
    
*   🔮 **Extensible design** — bring your own templates or modify existing ones.
    
*   🧱 **Preconfigured tooling** — includes Makefile, linters, and testing setup (coming soon).
    

* * *

## 📦 Installation

Install globally using `go install`:

`go install github.com/upsaurav12/bootstrap@latest`

Once installed, confirm the installation:

`bootstrap --help`

* * *

## 🚀 Quick Start

Create a REST API project using **Gin**:

```
bootstrap new myapp --type=rest --router=gin --port=8080
```

Create a project with **PostgreSQL** integration:

```
bootstrap new myapp --type=rest --router=gin --db=postgres
```

* * *

## 📁 Example Project Structure

```
myapp/ ├── Makefile ├── README.md ├── cmd/ │   └── main.go ├── internal/ │   ├── config/ │   │   └── config.go │   ├── handler/ │   │   └── user_handler.go │   ├── router/ │   │   └── routes.go │   └── db/               ← created only if --db flag is passed │       └── db.go └── go.mod
```

* * *

## ⚙️ CLI Options

| Flag | Description | Example |
| --- | --- | --- |
| --type | Type of project (rest, grpc, etc.) | --type=rest |
| --router | Router framework (gin, chi, echo) | --router=gin |
| --port | Application port | --port=8080 |
| --db | Database integration | --db=postgres |

* * *

## 💡 Why Go Bootstrapper?

Developers often waste time repeating setup tasks — creating folders, configuring routers, writing Makefiles, adding linters, etc.

**Go Bootstrapper** automates all that.  
You focus on business logic — it handles the rest.

It’s like:

> `create-react-app`, but for Golang 🐹

* * *

## 🛣️ Roadmap

*    Add `--with-auth` flag for JWT + middleware setup
    
*    Add Docker & Docker Compose templates
    
*    Support for Fiber, Echo, and gRPC
    
*    Generate Swagger / OpenAPI docs
    
*    Add custom template registry (`bootstrap add template`)
    

* * *

## 🤝 Contributing

Contributions, feedback, and ideas are welcome!  
Feel free to open an issue or PR on [GitHub](https://github.com/upsaurav12/bootstrap).

* * *

## 📄 License

Licensed under the **MIT License** © 2025 [Saurav Upadhyay](https://github.com/upsaurav12)

* * *