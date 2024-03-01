# Contributing Guidelines
> "Whatsoever you findeth to do, do it with thy might"

To get started, follow the steps below:

1. Clone project and install all required packages/dependencies 
    * see [Project Setup](#project-setup) for instructions
2. Check open issues that interest you, or make an [Issue Request](#github-issues)
3. Assign yourself to the issue
4. Create your branch, and start coding!
    * see [Development Process](#development-process) for some best practices

To Learn more about the technologies used in this project, check [Resources](#resources)

## Project Setup
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

3. Install Go dependencies defined in `go.mod`
    ```
    go mod tidy
    ```


## Github Issues
All planned and ongoing tasks will be listed under the Github Issues of this repository. You may pick any Issues listed to contribute to the project.

If you want to request a new feature or report a bug, you may also create an Issue with your request
* When requesting a feature, explain the what, the why, and show some use cases
* When reporting a bug, the more details the better to help any coders replicate and debug the issue

## Development Process
After finding an open issue to work on, you should:

1. Clone Project
2. Create a [Branch](#how-to-branch) to code on
    * See [How to Branch](#how-to-branch) for branch naming conventions
3. Create a [Pull Request](#pull-requests) when you are ready for a code review
4. Merge to main after approval

### Code Conventions
* Comment your code! Follow the official Go commenting conventions [here](https://go.dev/blog/godoc)
* Check Go code formatting conventions [here](https://go.dev/doc/effective_go)

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
* [Go By Example]()
    * Very quick guide to make you familiar with Go's syntax and features
* [An Introduction to Programming in GO](https://www.golang-book.com/books/intro)
    * A more in depth, yet gentle, introduction to Go

### PostgreSQL
* [PostgreSQL Documentation](https://www.postgresql.org/docs/)
* [PostgreSQL Docs Tutorial](https://www.postgresql.org/docs/current/tutorial-start.html)
    * Quick guide to get started with PostgreSQL
* [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
    * More in depth tutorial with practice examples

