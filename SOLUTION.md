# Solution Steps

1. 1. Define the APIRequest struct to hold request metadata and results, including fields for ID, Endpoint, Timestamp, Response, Duration, and Timeout flag.

2. 2. Write a utility function randomEndpoint() to pick a random API endpoint from a small list.

3. 3. Write the apiCall(req *APIRequest) function to simulate a network call, sleeping for a random duration between 400ms and 1200ms, updating the request object on completion.

4. 4. Implement processRequest: each request must acquire a slot from a semaphore (implemented by a buffered channel), run the apiCall in a goroutine, and use a select statement to enforce a 1.5s timeout; mark Timeout if appropriate, always release the slot, and send the result into the results channel.

5. 5. In main, create the semaphore buffered channel (size 3), a results channel, and a WaitGroup.

6. 6. Loop to create and launch 10 *concurrent* requests, each as a goroutine running processRequest, passing all required arguments; increment the WaitGroup for every launch.

7. 7. Wait for all requests to finish using wg.Wait(), then close the results channel.

8. 8. Finally, consume and print results from the results channel, showing each request's ID, endpoint, elapsed time, and noting 'timeout' when appropriate.

