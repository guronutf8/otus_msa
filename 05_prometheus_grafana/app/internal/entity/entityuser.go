package entity

type User struct {
	Id        string `json:"Id,omitempty" bson:"_id,omitempty"`
	Username  string `json:"username" jsonschema:"title=The user name,required" jsonschema_extras:"example=Jane"`
	FirstName string `json:"firstName" jsonschema:"title=The user firstName,required" jsonschema_extras:"example=Jane"`
	LastName  string `json:"lastName" jsonschema:"title=The user name,required" jsonschema_extras:"example=Jane"`
	Email     string `json:"email" jsonschema:"title=Email of the user"`
	Phone     string `json:"phone" jsonschema:"title=mobile number of user" jsonschema_extras:"example=+71002003040"`
}

type UserShort struct {
	FirstName string `json:"firstName" jsonschema:"title=The user firstName,required" jsonschema_extras:"example=Jane"`
	LastName  string `json:"lastName" jsonschema:"title=The user name,required" jsonschema_extras:"example=Jane"`
	Email     string `json:"email" jsonschema:"title=Email of the user"`
	Phone     string `json:"phone" jsonschema:"title=mobile number of user" jsonschema_extras:"example=+71002003040"`
}
