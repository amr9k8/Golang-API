## Golang-Authentication-Authorization-System

This project is an organization-user JWT authentication system implemented in Golang using the Gin framework. It integrates MongoDB for data storage, Redis for caching, and Docker for containerization


#### INSTALLATION

#### Using Docker 
 - adjust env variable in yaml file
 - Run `docker-compose up`
#### Without Docker
- Clone this repo
- Change directory to the cloned repo
- Ensure you have `go` installed on your machine.
- Run `go mod download`
- Ensure you have `mongodb` installed on your machine OR you can create a cloud `monog` database.
- Run `go run main.go`


#### FRAMEWORKS AND LIBRARIES USED

- `gin` - It is most popular go framework that is designed for building APIs. Hence, I have used gin to complete the task.
- `go-jwt` - It is the most popular package for implementing JWTs in Go. It has many pre-defined functions for signing and validating JWTs. Hence, used I have used it here, to sign and validate access_token and refresh_token.
- `redis` - A key-value store that acts as a caching layer in our system, enhancing performance and speeding up data retrieval. Redis is instrumental in optimizing the overall responsiveness of the authentication system.
- `mongodb` - A NoSQL database utilized for its flexibility and scalability. MongoDB stores hashed passwords securely using the Provos and Mazières's bcrypt adaptive hashing algorithm, ensuring robust password management within the application.



#### APIs&Features

- **LOGIN** POST -  User authentication with JWT.
- **SIGNUP** POST - User registration functionality.
- **RefreshJWTToken** POST - Refresh JWT token for prolonged sessions.
- **GetAllUsers** GET -  Retrieve all users (for Admin or testing purposes).
- **GetUserbyID** GET - Retrieve user details by ID.
- **DeleteUserbyID** DELETE - Delete user by ID.
- **UpdateUser** PUT - Update user information.
- **GetAllOrganizations** GET - Retrieve details of all organizations.
- **GetOrganizationbyID** GET - Retrieve organization details by ID.
- **AddOrganization** POST - Add a new organization.
- **InviteMemberIntoOrganization** POST -  Invite members with readonly access.
- **UpdateOrganization** PUT - Update organization details (Full Access required).
- **DeleteOrganization** DELETE - Delete organization by id (Full Access required).
- **TokenRefresh** POST - Refresh authentication token.
- **Bearer Authorization** -  Secure API access using Bearer token.
- **JWT Authentication** - JSON Web Token-based user authentication.


#### Project Structure
<pre>
|-- pkg/
|   |-- controllers/
|   |   - Controllers manage the application's flow and business logic. They receive input from the handlers, process it using the models, and return results to be presented by the views.
|   |-- db/models/
|   |   - Represents the data layer of the application. It typically includes data models and the logic for interacting with the database or other data sources.
|   |-- db/repository/
|   |   - Contain Database operations.
|   |-- utils/
|       - Contains utility functions or modules that can be used across different parts of the application. Utilities might include helper functions, generic components, or modules that provide common functionalities.
|-- Api/
|   |-- routes/
|   |   - Contains the definitions of routes and their corresponding handlers, responsible for routing incoming requests to the appropriate controllers.
|   |-- middlewares/
|   |   - Includes middleware components that can be executed before or after the request reaches the controller. Middlewares often handle tasks like authentication, logging, etc.
|   |-- handlers/
|       - This folder might contain modules or classes that handle specific types of requests. Handlers are often responsible for interacting with the request and response objects, processing data, and calling the appropriate controller methods.
|-- cmd/
|   - main.go: entry point of app
|-- go.mod
|-- go.sum
|-- Dockerfile : Instructions for building the application image.
|-- docker-compose.yaml: Configuration for Docker Compose.
</pre>
#### Technologies Used

- **Golang**
- **Gin (Web framework)**
- **MongoDB (Database)**
- **Redis (Caching)**
- **Docker (Containerization)**
- **docker-compose (Container orchestration)**
