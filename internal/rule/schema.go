package rule

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rule struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
	Matcher  string      		`bson:"matcher"`
	Ineq	 string				`bson:"ineq"`
	Limit    int				`bson:"limit"`
	Obj  	 []primitive.ObjectID `bson:"obj"`
	Notify   string				`bson:"notify"`
	When 	 string				`bson:"when"`
    Owners   []primitive.ObjectID `bson:"owners"`
}
