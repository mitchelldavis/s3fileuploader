S3FileUploader
==============

[![Build Status](https://travis-ci.org/mitchelldavis/s3fileuploader.svg?branch=master)](https://travis-ci.org/mitchelldavis/s3fileuploader)

This application can be used to give users the ability to upload files to non-public buckets in AWS S3 without provisioning IAM credentials.

The application leverages AWS Cognito User and Identiy pools to do so.  This means a Cognito User Pool associated with an Identity pool is required and the Ids of said resources should be passed to the parameters of the application.

USAGE
---

```
./s3fileupload -h

Usage of ./s3fileupload:
    -bucket string
        The Bucket to upload the file to.
    -clientId string
        The Client Id to pass to the identity provider.
    -filename string
        The File to upload.
    -identityPoolId string
        The Identity Pool Id that we're authenticating against.
    -region string
        The AWS region the user and identity pools exist in.
    -userPoolId string
        The User Pool Id that has the user information we're authenticating against.
    -username string
        The Username to sign in with
```

BUILD
-----

```
make
./.bin/s3fileupload -h
```

LICENSE
-------

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org>
