# A Portfolio Project - Time to take off the training wheels

## Some Helpful Notes

### The modd.conf file

```

**/*.go
**/*.gohtml {
    prep: go build -o ./server ./cmd/server
    daemon +sigterm: ./server
}

```

Config allows for autoreloading of the server, which can be used for development purposes. As it will not be required
in production, if required can be added to the .gitignore file.

Should be thought of as:

- Whenever there are any changes in files _.go or _.gohtml
- rebuild the server
- restart the server as a daemon (in the background) checking for sigterm sent to the process run by the server binary

# Project Structure

### /CMD

the CMD folder should be housing the main applications for this project. This could include

- servers
- command line tools
- other executables, or tools that we would end up building

The directory for each application should match the name of the executable you want to have. For example
**/cmd/myapp/**

It should not hold a lot of code. If the code can be imported and used in other projects, then it should live in the
/pkg directory. If the code is not reusable or rather, you don't want other people to use it, then put it in the
internal directory.

It's common to have a small <mark>main()</mark> function, that invokes all the other code from <mark>internal</mark>
and <mark>pkg</mark> directories, so think of the main function that wires everything up, and gets things ready
to be built.

### /internal

me <mark>/internal</mark> package should be used to house application code that is not intended to be reused by
other projects (both my own and others.) All code that is in the /internal folder **cannot** be imported into other
projects, and so this will make refactoring easier as we would never expose dependencies to other consumers

for now, there is no good reason to have any other main folders besides, these, as all application code should be
within the internal package

### Tasks

- render a form, using the gohtml templating engine, so that we can make a post request
- protect that post request using gorilla/csrf, and validate locally that this is actually working
- try to deploy it to an ec2 instance, on AWS
