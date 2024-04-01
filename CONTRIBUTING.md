# Contributing Guidelines
> "Whatsoever thy hand findeth to do, do it with thy might"

To get started, follow the steps below:

1. Check open issues that interest you, or make an [Issue Request](#github-issues)
2. Assign yourself to the issue
3. Setup the project locally, create your branch, and start coding!
    * see [Project Setup](#project-setup) for setup instructions
    * see [Development Process](#development-process) for best practices 

To Learn more about the technologies used in this project, check [Resources](#resources)

## Github Issues
All planned and ongoing tasks will be listed under the Github Issues of this repository. You may pick any Issues listed to contribute to the project.

If you want to request a new feature or report a bug, you may also create an Issue with your request
* When requesting a feature, explain the what, the why, and show some use cases
* When reporting a bug, the more details the better to help any coders replicate and debug the issue

## Project Setup
Instructions for setting up repo for local development:
1. clone repository
    - cloning via HTTPS:
        ```
        git clone https://github.com/pacific-theta-tau/tt-db.git
        ```

    - cloning via SSH:
        ```
        git clone git@github.com:pacific-theta-tau/tt-db.git
        ```

2. Install Go. You can see instructions [here](https://go.dev/doc/install)

3. Install Docker. Instructions [here](https://docs.docker.com/engine/install/)

3. Install dependencies defined in `go.mod`
    ```
    go mod tidy
    ```
5. run server using Docker Compose
    ```
    # Run dev containers
    docker compose --profile dev up 

    # Add `--build` after changes in code to rebuild services
    docker compose --profile dev up --build

    # Shut down dev containers
    docker compose --profile dev down
    ```


## Development Process
When working on an issue, you should:

1. Create a [Branch](#how-to-branch) to code on
    * See [How to Branch](#how-to-branch) for branch naming conventions
2. Create a [Pull Request](#pull-requests) when you are ready for a code review
3. Merge to main after approval

### Code Conventions
* Comment your code! Follow the official Go commenting conventions [here](https://go.dev/blog/godoc)
* Check Go code formatting conventions [here](https://go.dev/doc/effective_go)
* For API conventions, check [this](https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design)

### Branching 
#### How to Branch
`git checkout -b "name-of-the-branch"`
- you will be automatically moved to your new branch after hitting enter

#### **Branch Naming and Organization**
Branch names should be separated by dashes and should be organized in a hierarchical structure like folders:

```git checkout -b "<your-username>/<short-description>"```

#### **Example names:**

`NickAhn/add-main-page`

`thlau21/fix-loading-screen`

### Pull Requests
When creating a Pull Request, describe what you did and assign reviewers to approve your request
* Reviewers can be other members working in your area (Frontend, Backend, etc...), or experienced/older members
* If you are unsure who to assign as a reviewer, contact the current project maintainer.

## Resources
### Learn Go
* [Tour of Go](https://go.dev/tour/welcome/1)
    * A very quick guide of Go's syntax and features, including some practice problems
* [Go By Example](https://gobyexample.com/)
    * Guide to make you familiar with Go's syntax and features
* [An Introduction to Programming in GO](https://www.golang-book.com/books/intro)
    * A more in depth, yet gentle, introduction to Go

### PostgreSQL
* [PostgreSQL Documentation](https://www.postgresql.org/docs/)
* [PostgreSQL Docs Tutorial](https://www.postgresql.org/docs/current/tutorial-start.html)
    * Quick guide to get started with PostgreSQL
* [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
    * More in depth tutorial with practice examples

### Docker
* [Dockerfile Reference](https://docs.docker.com/reference/dockerfile/)
* [Docker Compose File Reference](https://docs.docker.com/compose/compose-file/compose-file-v3/)