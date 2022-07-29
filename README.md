# golang_recruitment_task


As it was written in backend.md:
GET: /random/mean/?requests={r}&length={l}
![Screenshot](api_example.png)


In examples, there was:
{
     "stddev": 1,
     "data": [1, 2, 3, 4, 5]
}
but standard deviation of [1,2,3,4,5] is $\sqrt{2}$, so in my app this will be the result.

It was not specified from what range the numbers should be taken from, so I set it to from 1 to 6 to make the results easy to interpret.

There is sleep(3 seconds) function, for testing concurrency, it can be get rid of later.



### Dockerhub link with app image: https://hub.docker.com/r/ksero/random-api-app
