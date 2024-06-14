package rule

import "go.mongodb.org/mongo-driver/bson/primitive"

type MinMax struct {
	Which string				`bson:"which"`
	Obj1  string				`bson:"obj1"`
	Obj2  string				`bson:"obj2"`
	Limit    int				`bson:"limit"`
	Ineq	 bool				`bson:"ineq"`
}

type Rule struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
	Matcher  string      		`bson:"matcher"`
	Ineq	 string				`bson:"ineq"`
	Limit    int				`bson:"limit"`
	Obj  	 []primitive.ObjectID `bson:"obj"`
	//Created  int				`bson:"created"`
	Notify   string				`bson:"notify"`
	When 	 string				`bson:"when"`
    Owners   []primitive.ObjectID `bson:"owners"`
}
