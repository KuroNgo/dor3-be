## CLEAN ARCHRITECTURE
#### Reference Resource: [github clean-architecture](https://github.com/amitshekhariitbhu/go-backend-clean-architecture/tree/main/mongo)
![Workflow](assets/images/button-view-api-docs.png)

# Overall 
![Workflow](assets/images/CleanArchitecture.jpg)
## Entities: 
    Domain and Repository in Entities
- Directories like domain and repository may contain components related to the core Entities in Clean Architecture. 
- In this case, domain might hold the core objects and business logic of the application, while repository could contain classes or components related to data access.

## UseCase
    Usecase or Interactions in Usecase
- The useCase directory may contain the use cases (or interactions) 
- in Clean Architecture, which are components that contain the business logic of the application.

## Interface Adapters
    API, Infrastructor and internal in Interface Adapter
- Directories such as api, infrastructor, and internal may contain components related to interfaces and infrastructure. 
- In Clean Architecture, these components are often called Interface Adapters, responsible for converting data between internal and external components of the application.
  #### api 
- may contain components related to handling HTTP requests from the client.
  
  #### infrastructor 
- may contain components related to implementing infrastructure such as databases, peripheral services, etc.
  
  #### internal 
- may contain internal application components that are not meant to be accessed from outside, such as authentication mechanisms (in this case OAuth2) or special connections to third-party services (e.g., connecting to Google services).


## Frameworks and Drivers
    Idea, asset and bootstrap
- Directories like .idea, assets, and bootstrap may contain components related to frameworks, tools, and resources such as images, fonts, CSS files, JavaScript files, etc. 
- These components are typically not part of Clean Architecture but are still important for implementing the application.

## File Structure
```
├───.idea
│   └───inspectionProfiles
├───api
│   ├───controller
│   │   ├───activity
│   │   ├───admin
│   │   ├───course
│   │   ├───exam
│   │   ├───exam_answer
│   │   ├───exam_options
│   │   ├───exam_question
│   │   ├───exam_result
│   │   ├───exercise
│   │   ├───exercise_answer
│   │   ├───exercise_options
│   │   ├───exercise_quesiton
│   │   ├───exercise_result
│   │   ├───feedback
│   │   ├───image
│   │   ├───lesson
│   │   ├───mark_list
│   │   ├───mark_vocabulary
│   │   ├───quiz
│   │   ├───quiz_answer
│   │   ├───quiz_options
│   │   ├───quiz_question
│   │   ├───quiz_result
│   │   ├───unit
│   │   ├───user
│   │   ├───user_attempt
│   │   └───vocabulary
│   ├───middleware
│   └───router
│       ├───activity_log
│       ├───admin
│       ├───course
│       ├───exam
│       ├───exam_answer
│       ├───exam_options
│       ├───exam_question
│       ├───exam_result
│       ├───exercise
│       ├───exercise_answer
│       ├───exercise_options
│       ├───exercise_question
│       ├───exercise_result
│       ├───feedback
│       ├───image
│       ├───lesson
│       ├───mark_list
│       ├───mark_vocabulary
│       ├───quiz
│       ├───quiz_answer
│       ├───quiz_options
│       ├───quiz_question
│       ├───quiz_result
│       ├───unit
│       ├───user
│       ├───user_attempt
│       └───vocabulary
├───assets
│   └───images
├───bootstrap
├───config
├───domain
│   ├───activity_log
│   ├───admin
│   ├───box
│   ├───config_system
│   ├───course
│   ├───exam
│   ├───exam_answer
│   ├───exam_options
│   ├───exam_question
│   ├───exam_result
│   ├───exercise
│   ├───exercise_answer
│   ├───exercise_options
│   ├───exercise_questions
│   ├───exercise_result
│   ├───feedback
│   ├───image
│   ├───lesson
│   ├───mark_list
│   ├───mark_vocabulary
│   ├───quiz
│   ├───quiz_answer
│   ├───quiz_options
│   ├───quiz_question
│   ├───quiz_result
│   ├───unit
│   ├───user
│   ├───user_attempt
│   ├───user_detail
│   └───vocabulary
├───infrastructor
│   ├───mongo
│   └───redis
├───internal
│   ├───cloud
│   │   ├───cloudinary
│   │   ├───firebase
│   │   │   └───audio
│   │   └───google
│   │       ├───const
│   │       └───mail
│   ├───file
│   │   └───excel
│   └───oauth2
│       └───google
├───repository
│   ├───activity
│   ├───admin
│   ├───course
│   ├───exam
│   ├───exam_answer
│   ├───exam_options
│   ├───exam_question
│   ├───exam_result
│   ├───exercise
│   ├───exercise_answer
│   ├───exercise_options
│   ├───exercise_question
│   ├───exercise_result
│   ├───feedback
│   ├───image
│   ├───lesson
│   ├───mark_list
│   ├───mark_vacabulary
│   ├───quiz
│   ├───quiz_answer
│   ├───quiz_options
│   ├───quiz_question
│   ├───quiz_result
│   ├───unit
│   ├───user
│   ├───user_attempt
│   ├───user_detail
│   └───vocabulary
├───templates
└───usecase
    ├───activity
    ├───admin
    ├───course
    ├───exam
    ├───exam_answer
    ├───exam_options
    ├───exam_question
    ├───exam_result
    ├───exercise
    ├───exercise_answer
    ├───exercise_options
    ├───exercise_question
    ├───exercise_result
    ├───feedback
    ├───image
    ├───lesson
    ├───mark_list
    ├───mark_vocabulary
    ├───mean
    ├───quiz
    ├───quiz_answer
    ├───quiz_options
    ├───quiz_question
    ├───quiz_result
    ├───system
    ├───unit
    ├───user
    ├───user_attempt
    └───vocabulary
```
## Run Programming
How to run this project?
We can run this Go Backend Clean Architecture project with or without Docker. Here, I am providing both ways to run this project.

#### Clone this project
    cd your-workspace
- Move to your workspace


#### Clone this project into your workspace
- git clone https://github.com/KuroNgo/FEIT.git

#### Move to the project root directory
    cd FEIT

#### Run without Docker
- Create a file .env similar to .env.example at the root directory with your configuration.
- Install go if not installed on your machine.
- Install MongoDB if not installed on your machine.
- Important: Change the DB_HOST to localhost (DB_HOST=localhost) in .env configuration file. DB_HOST=mongodb is needed only when you run with Docker.
- Run go run cmd/main.go.
- Access API using http://localhost:8080
#### Run with Docker
- Create a file .env similar to .env.example at the root directory with your configuration.
- Install Docker and Docker Compose.
- Run docker-compose up -d.
- Access API using http://localhost:8080
#### How to run the test?
#### Run all tests
    go test ./...
#### How to generate the mock code?
- In this project, to test, we need to generate mock code for the use-case, repository, and database.

#### Generate mock code for the usecase and repository
    mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

#### Generate mock code for the database
    mockery --dir=mongo --output=mongo/mocks --outpkg=mocks --all
- Whenever you make changes in the interfaces of these use-cases, repositories, or databases, you need to run the corresponding command to regenerate the mock code for testing.