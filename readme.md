## Problem

We've added our password hashing functionality, and we set it as our password before saving
a new user. * You should never, ever store plain-text passwords *

Also, on authentication, we check against the hashed password.

Now, traditionally, we can securely authenticate a user against the database, however, we need
a mechanism in which we can do this across our user interfaces and distributed services. There
are many ways in which to do this, but the simplest solution I've come across, which we can use
across our services and web, is JWT.

- JWT vs OAuth

### JWT
JWT stands for JSON web tokens and is a distributed security protocol. Similar to OAuth.
The concept is simple, you use an algorithm to generate a unique hash for a user, 
which can be compared and validated against.

But not only that, the token itself can contain and be made up of our users' metadata.
In other words, their data can itself become a part of the token. So let's look at an example
of a JWT:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

The token is separated into three by .'s. Each segment has a significance. 

- The first segment is made up of some metadata about the token itself. Such as the type of token and the algorithm used to create the token. This allows clients to understand how to decode the token. 
- The second segment is made up of user-defined metadata. This can be your user's details, an expiration time, anything you wish. 
- The final segment is the verification signature, which is information on how to hash the token and what data to use.

