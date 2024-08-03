# tinymapper
> A small library to handle type mapping intended for HTTP APIs inspired by the AutoMapper library from the C# world.

## 👩‍💻 Usage
To map between types, you define how to map between types, and then use either `tinymapper.To`, or `tinymapper.ArrayTo`.

### Single object mappings
```go
package main

import "fmt"
import "github.com/imthatgin/tinymapper"

type User struct {
    Id       uint
    Username string
    
    // When mapping, these should get turned into a display name property
    Firstname string
    Lastname  string
	
    PasswordHash string
}

type UserDTO struct {
    Id          uint
    Username    string
    DisplayName string
}

func main() {
    m := tinymapper.New()
    
    // Define how to map between the type here.
    // Fields with the same name and type will map automatically.
    tinymapper.Register(m, func(user User, dto *UserDTO) {
        dto.DisplayName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
    })
    
    // An example object
    user := User{
        Id:           121,
        Username:     "test_user_121",
        Firstname:    "Test",
        Lastname:     "User",
        PasswordHash: "aaabbbccc",
    }
    
    dto, err := tinymapper.To[UserDTO](m, user)
}
```

### Array mappings
```go
// Same setup as the single object mapping.

users := []User{
    {
        Id:           121,
        Username:     "test_user_121",
        Firstname:    "Test",
        Lastname:     "User",
        PasswordHash: "aaabbbccc",
    },
}

mappedUsers, err := tinymapper.ArrayTo[UserDTO](m, users)
```

### Nested mappings
If you need a nested mapping (ie. Map to a DTO inside another DTO), you can map like this:
```go
tinymapper.Register(m, func(user User, dto *UserDTO) {
    dto.DisplayName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
})

tinymapper.Register(m, func(o ObjectWithAnUser, dto *ObjectWithAnUserDTO) {
    dto.UserId = o.Id
    dto.Owner, err = tinymapper.To[UserDTO](m, o.Owner)
})
```

## 🏗 Contributions
* All contributions are very much welcome, and I hope that this library is useful for your use cases.

## 💡  Motivation
*I like the AutoMapper feature of .NET, and usually face similiar requirements when writing Go APIs.*

