# deduper

## APIs

**get(key string)**

**set(key string, value string, time time)**

**expire(key string, time time)** `## Set a timeout on key. After the timeout has expired, the key will automatically be deleted`

**incr(key string)** `## Increments the number stored at key by one. If the key does not exist, it is set to 0 before performing the operation. An error is returned if the key contains a value of the wrong type or contains a string that can not be represented as integer.`

**decr(key string)** `## Decrements the number stored at key by one. If the key does not exist, it is set to 0 before performing the operation. An error is returned if the key contains a value of the wrong type or contains a string that can not be represented as integer.`

**multi()** `## transaction starts. TODO may not be needed.`