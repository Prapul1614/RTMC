package rule

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rule struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
	MinMax   string 			`bson:"min_max"`
    Name     string             `bson:"name"`
	Word     string             `bson:"word"`
	Char	 rune				`bson:"char"`
	IsTheir	 bool				`bson:"is_their"`
	Limit    int				`bson:"limit"`
	Ineq	 bool				`bson:"ineq"`
    Owners   []primitive.ObjectID `bson:"owners"`
}
