## CLEAN ARCHRITECTURE
#### Reference Resource: [github clean-architecture](https://github.com/amitshekhariitbhu/go-backend-clean-architecture/tree/main/mongo)
![Workflow](./assets/button-view-api-docs.png)
### API
### Assets
### Bootstrap
### CMD
### Router
First of all, the request comes to the Router.

Further, this router gets divided into two routers as follows:

Public Router: All the public APIs should go through this router.
Protected Router: All the private APIs should go through this router.
The Public API request flow:
![Workflow](./assets/go-arch-public-api-request-flow.png)
The Private API request flow:
![Workflow](./assets/go-arch-private-api-request-flow.png)
JWT Authentication Middleware for Access Token Validation.

In between both routers, a middleware gets added to check the validity of the access token. So, the private request with the invalid access token should not reach the protected router at all.

### Controller
- So now, the request is with the controller. First, it will validate the data present inside the request. If anything is invalid, it returns a "400 Bad Request" as the error response.
- If everything is valid inside the request, it will call the usecase layer to perform an operation.

### Usecase
- The usecase layer is dependent on the repository layer.
- This layer uses the repository layer to perform an operation. 
It is completely up to the repository how it is going to perform an operation.

### Repository
- The repository is the dependency of the usecase. The Usecase layer asks the repository to perform an operation.
- The repository layer is free to choose any database, in fact, it can call any other independent services based on the requirement.
- In the project, the repository layer makes the database query for performing operations asked by the Usecase layer.

### Domain

### File Structure
```
│   .gitignore
│   app.env
│   go.mod
│   go.sum
│   go.work
│   go.work.sum
│   main.go
│   README.md
│
├───.idea
│       .gitignore
│       clean-architecture.iml
│       modules.xml
│       workspace.xml
│
├───api
│   ├───controller
│   │       user.controller.go
│   │
│   ├───middleware
│   │       cors.go
│   │       deserialization.go
│   │       rate_limit.go
│   │
│   └───route
│           route.go
│
├───assets
│       button-view-api-docs.png
│       go-arch-private-api-request-flow.png
│       go-arch-public-api-request-flow.png
│       go-backend-arch-diagram.png
│
├───bootstrap
│       app.go
│       db.go
│       env.go
│
├───domain
│   ├───request
│   │       course.go
│   │       course.impl.go
│   │       lesson.go
│   │       lesson.impl.go
│   │       permission.go
│   │       permission.impl.go
│   │       quiz.go
│   │       store.go
│   │       store.impl.go
│   │       user.go
│   │       user.impl.go
│   │       user_process.go
│   │       user_process.impl.go
│   │       user_quiz_attempt.go
│   │       user_quiz_attempt.impl.go
│   │
│   └───response
│           error_response.go
│           success_response.go
│
├───infrastructor
│   └───mongo
│       │   mongo.go
│       │
│       └───mocks
│               client.go
│               collection.go
│               cursor.go
│               database.go
│               single_result.go
│
├───internal
│       token.go
│
├───repository
│       permission.repo.go
│       user.repo.go
│
└───usecase
    ├───admin
    │       statistical.go
    │
    ├───system
    │       role.usecase.go
    │
    └───user
            user.usecase.go


```
### Run Programming