# 🦁 Book Forum

This project is a robust web forum designed to tame the "paper pandemonium" of traditional book discussions. It replaces scattered sticky notes and lost comments with a structured, persistent, and engaging online community. Members can spark lively discussions, dissect chapters, share interpretations, and archive insights for future readers.

## ✨ Features

### Core Functionality

* **User Authentication:** Secure registration and login using email/password.
* **Discussions:** Write posts, and engage in comments.
* **Categorization:** Filter discussions by Genre, Book Title, or Author.
* **Reactions:** Like or Dislike posts and comments to show agreement or appreciation.
* **Filtering:** View posts by "My Posts", "My Likes", or specific categories.

### 🚀 Extra Features (Implemented)

* **📸 Image Uploads:** Share photos of book covers or specific passages directly in your posts.
* **🔍 Search Bar:** Instantly find discussions by keywords in titles or content.
* **👤 User Profiles:** View your own activity and history.
* **💬 Real-Time Chat:** A live chat feature for instant communication between online members.
* **🔒 Security:** Passwords are encrypted (hashed) for user safety.

---

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Database:** SQLite3 (Embedded, lightweight relational database)
* **Frontend:** HTML5, CSS
* **Containerization:** Docker

---

## 📋 Prerequisites

Before you begin, ensure you have met the following requirements:

* **Git:** To clone the repository.
* **Docker:** Recommended for running the application in a containerized environment.
* **Go (Golang):** Required only if you intend to run the application locally without Docker.

---

## Database Design

Entity Relationship Diagram (ERD) is available in:
data/ERD.png

---

## 🐳 Installing Docker

If you don't have Docker installed, follow the official Docker installation guide for your operating system:

👉 [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)

After installation, verify Docker is working:

```bash
docker --version
```

---

## 📥 How to Clone the Repository

To get a copy of this project up and running on your local machine:

```bash
git clone https://gitea.kood.tech/aleksanderkruk/forum
cd literary-lions-forum
```

---

## 🐳 Running the Project with Docker (Recommended)

Using Docker is the easiest way to build and run the application without installing Go or other dependencies locally.

### 🔨 Build the Docker Image

From the root of the project (where the `Dockerfile` is located), run:

⚠️ Important for Windows Users: You must run the build command inside your WSL (Windows Subsystem for Linux) terminal (e.g., Ubuntu). Do not use PowerShell or Command Prompt, as this may cause build errors

From the root of the project (where the Dockerfile is located), run:

```bash
docker build -t forum .
```

This command:

* Builds the application inside a Docker image
* Tags the image with the name `forum`

---

### ▶️ Run the Docker Container

Once the image is built, start the application with:

```bash
docker run -p 8080:8080 forum
```

* The application will be available at: **[http://localhost:8080](http://localhost:8080)**
* `-p 8080:8080` maps the container port to your local machine

To run the container in detached (background) mode:

```bash
docker run -d -p 8080:8080 forum
```

---

### 🛑 Stopping the Container

To stop the running container:

```bash
docker ps
docker stop <container_id>
```

---

## 🧑‍💻 Running Without Docker (Optional)

If you prefer to run the project locally without Docker:

1. Ensure **Go** is installed: [https://go.dev/dl/](https://go.dev/dl/)
2. From the project root, run:

```bash
go run ./cmd/web
```

---

Happy reading and discussing! 📚🐾
