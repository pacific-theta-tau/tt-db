# Contributing Guidelines
> "Whatsoever you findeth to do, do it with thy might"

To get started, follow the 

1. Check open issues that interest you, or make a [Issue Request](#github-issues)
2. Assign yourself to the issue
3. Create your branch, and start coding!
    * see [Development Process](#development-process) for some best practices

## Github Issues
All planned and ongoing tasks will be listed under the Github Issues of this repository. You may pick any Issues listed to contribute to the project.

If you want to request a new feature or report a bug, you may also create an Issue with your request
* When requesting a feature, explain the what, the why, and show some use cases
* When reporting a bug, the more details the better to help any coders replicate and debug the issue

## Development Process
After finding an open issue to work on, you should:

1. Clone Project
2. Create a [Branch](#how-to-branch) to code on
3. Create a [Pull Request](#pull-requests) when you are ready for a code review
4. Merge to main after approval

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

### Code Conventions
* Comment functions appropriately
* Check Go code formatting conventions [here](https://go.dev/doc/effective_go)

## Pull Requests
When creating a Pull Request, describe what you did and assign reviewers to approve your request
* Reviewers can be other members working in your area (Frontend, Backend, etc...), or experienced/older members
* If you are unsure who to assign as a reviewer, contact the current project maintainer.