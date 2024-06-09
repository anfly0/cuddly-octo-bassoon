# Robo server

## Build & test

### Build

To build the robo-server navigate to the root folder off this repo and run `go build ./cmd/robo-server`

### Test

To run all unit tests navigate to the root of this repo and run `go test ./...`
To run benchmarks run `go test ./... -bench=. -benchmem`  

To analyze performance issues is it useful to run the benchmark(s) with the `-memprofile=mem.prof` or `-cpuprofile=cpu.prof`  flag and then use pprof to get a detailed report.

## Running the server

The robo-server executable accepts the following command line flags:

- **-h** prints an overview of the available flags
- **-addr** the address of the interface that the server will bind to
- **-port** the port number the server will bind to

## API Documentation

### Endpoints

#### Create a Robot

**Endpoint:** `POST /robot`

**Description:** This endpoint creates a new robot with the specified direction, room, and starting coordinates.

**Request Body:**

```json
{
  "direction": "N", // Direction the robot is facing ('N', 'E', 'S', 'W')
  "room": {
    "x": 3, 
    "y": 3 
  },
  "start": {
    "x": 0, // X coordinate of the starting position
    "y": 0  // Y coordinate of the starting position
  }
}
```

**Responses:**

- **200 OK:** Robot created successfully.

  ```json
  {
    "direction": "N",
    "x": 0,
    "y": 0,
    "id": "abcd" // ID of the created robot
  }
  ```

- **400 Bad Request:** Invalid request payload.
- **500 Internal Server Error:** Server encountered an error while processing the request.

---

### Get Robot Status

**Endpoint:** `GET /robot/{id}`

**Description:** This endpoint retrieves the status of a robot with the specified ID. The status is guaranteed to be internally consistent i.e combination of x, y, and direction represent the real stat of the robot at the time the request is processed.

**Path Parameters:**

- `id` (string): The ID of the robot.

**Responses:**

- **200 OK:** Robot status retrieved successfully.

  ```json
  {
    "direction": "N",
    "x": 0,
    "y": 0,
    "id": "abcd"
  }
  ```

- **404 Not Found:** Robot with the specified ID not found.

---

### Command a Robot

**Endpoint:** `POST /robot/{id}`

**Description:** This endpoint sends a series of commands to the robot with the specified ID. The robot will process commands up until one of two things occur.

  1. The command string ends
  2. The robot encounters an invalid command i.e a command that is not R, L or F. In this case the robot will be left in the state it was after the last valid command was processed.

If multiple request are made to this endpoint concurrently the robot is guaranteed to process one series of commands at a time. **The order of processing is however not guaranteed**.

**Path Parameters:**

- `id` (string): The ID of the robot.

**Request Body:**

```json
{
  "cmd": "LRF" // Command to be executed by the robot (e.g., 'F', 'L', 'R')
}
```

**Responses:**

- **200 OK:** Command executed successfully.

  ```json
  {
    "direction": "N",
    "x": 1,
    "y": 0,
    "id": "abcd"
  }
  ```

- **400 Bad Request:** Invalid command or request payload.
- **404 Not Found:** Robot with the specified ID not found.
