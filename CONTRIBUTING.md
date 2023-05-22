# Branching
There are 2 default branches, one for each environment; `development`
for the `development` environment and `main` for the `production` environment.

The development branch is for testing our code. Here is where every newly developed feature is
merged into as per default when creating a pull requests. This is also where we would discover if any
new features were to introduce unforeseen bugs. It is also the branch that every new feature should
branch out of.

The main branch is our production environment. This is our ”live” code and is what the current
state of our product looks like. This is also to define what is ready to use for our users.
Whenever a new feature is to be initiated, this should happen in a new branch. This new branch
should be named after the new feature. Following this naming convention allows other developers to
keep a track of what each branch is for.

# Workflow
A centralized workflow is to be used for this project.

# Reviewing
There must be at least 1 `approved` review for each pull request. The reviewer must not have contributed to the pull request. Once the pull request has been approved, the branch gets merged into the `development` branch