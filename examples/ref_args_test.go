package examples_test

import (
	"encoding/json"
	"fmt"

	"github.com/sr9000/go-ptr-tools/ref"
)

type User struct {
	FirstName    string `json:"first-name"`
	LastName     string `json:"last-name"`
	IsRegistered bool   `json:"is-registered"`
	TotalVisits  *int   `json:"total-visits"`
}

func (u User) String() string {
	var totalVisits any

	if u.TotalVisits != nil {
		totalVisits = *u.TotalVisits
	}

	return fmt.Sprintf("{%s %s %t %v}", u.FirstName, u.LastName, u.IsRegistered, totalVisits)
}

func PlusOne(r ref.Ref[int]) {
	*r.Ptr()++ // dereference and increment without checking for nil
}

func VisitsIncrement(u ref.Ref[User]) error {
	if u.Ptr().IsRegistered {
		// new reference from the TotalVisits pointer field
		r, err := ref.New(u.Ptr().TotalVisits)
		if err != nil {
			return fmt.Errorf("visits increment: %w", err)
		}

		PlusOne(r)
	}

	return nil
}

func Example_guaranteedReference() {
	var n int

	PlusOne(ref.Guaranteed(&n))
	PlusOne(ref.Guaranteed(&n))
	PlusOne(ref.Guaranteed(&n))

	fmt.Println(n)
	// Output: 3
}

func Example_usersVisits() {
	// simulating data from http request
	users := []struct {
		tag, data string
	}{
		{"registered", `{"first-name":"John","last-name":"Doe","is-registered":true,"total-visits":10}`},
		{"unregistered", `{"first-name":"Jane","last-name":"Doe","is-registered":false}`},
		{"invalid", `{"first-name":"Invalid","last-name":"User","is-registered":true}`},
	}

	for _, packet := range users {
		var user User

		err := json.Unmarshal([]byte(packet.data), &user)
		if err != nil {
			panic(err)
		}

		err = VisitsIncrement(ref.Guaranteed(&user))
		fmt.Println(packet.tag, user, err)
	}

	// Output:
	// registered {John Doe true 11} <nil>
	// unregistered {Jane Doe false <nil>} <nil>
	// invalid {Invalid User true <nil>} visits increment: ptr must be not nil
}
