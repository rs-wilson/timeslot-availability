# Timeslot Mircoservice
This code was written as a technical challenge, and took about 4 hours to complete.

## Instructions
Create an API using the language of your choosing that will be the foundation for an availability microservice that will help us determine if a time slot is available as well as reserve and free up time slots.
The service should be able to perform 3 main functions:

1. Check availability of a `time_slot` (a `time_slot` is a combination of `Start Timestamp` and `Duration`)
   - **Input:** `time_slot`
   - **Output:** `true` if the time slot is available, `false` if it is not available
2. Reserve availability of a `time_slot`
   - **Input:** `time_slot`
   - **Output:** confirmation if successful, user friendly error if not successful
3. Free availability of a previous time slot
   - **Input:** `time_slot`
   - **Output:** confirmation if successful, user friendly error if not successful

### Notes
- Data storage is up to you
- Testing is nice but not required for this
- If you don't have time to fully implement something leave a comment on how you would approach it.

## Development
This project uses a Makefile. You can run `make` to build and test the codebase.

### Building
`make build`

### Testing
`make test`

### Running
`make run`

## Endpoints
All endpoints use the same json request body to determine the timeslot.

Request Body:
```json
{
  "start_timestamp": "{unix_timestamp_seconds}",
  "duration": "{seconds}"
}
```

If an error occurs, the error message will be returned in a response body with this format:

Error Body:
```json
{
  "error": "{message}"
}
```

### POST /v1/timeslot
The POST endpoint checks whether the requested timestamp is available.

Response Body:
```json
{
  "available": "{true/false}"
}
```

### PUT /v1/timeslot
The PUT endpoint attempts to reserve a timeslot.

If the timeslot is available, a 200 with no body is returned and the timeslot is reserved.
If the timeslot is not available, a 409 'conflict' with an error body is returned instead.

### DELETE /v1/timeslot
The DELETE endpoint attempts to free a reserved timeslot.

If the timeslot is not available, a 204 with no body is returned.
If the timeslot is available, a 404 'not found' with an error body is returned instead.

## Roadmap
+ Unit testing
+ Make server configurable through a `config.yaml` file
+ Add user ID to request body for PUT
+ Add "get" endpoint to get timeslots for a given user
+ Add "rescheduled" state to deleted timestamps
+ Replace in-memory store with real DB
+ Make time-checking more efficient
+ Dockerize
