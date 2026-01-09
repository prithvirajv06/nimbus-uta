package models

type Orcestartion struct {
	NIMB_ID     string   `bson:"nimb_id" json:"nimb_id"`
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description" json:"description"`
	Type        string   `bson:"type" json:"type"`
	Category    string   `bson:"category" json:"category"`
	Workflow    Workflow `bson:"workflow" json:"workflow"`
	Audit       Audit    `bson:"audit" json:"audit"`
}

type Workflow struct {
	Nodes []Node `json:"nodes"  gorm:"type:jsonb" bson:"nodes"`
	Links []Link `json:"links" gorm:"type:jsonb" bson:"links"`
}

type Node struct {
	ID       int      `json:"id" gorm:"primaryKey;autoIncrement" bson:"id"`
	Type     string   `json:"type" gorm:"type:varchar(100)" bson:"type"`
	Title    string   `json:"title" gorm:"type:varchar(255)" bson:"title"`
	X        float64  `json:"x" gorm:"type:float" bson:"x"`
	Y        float64  `json:"y" gorm:"type:float" bson:"y"`
	Width    float64  `json:"width" gorm:"type:float" bson:"width"`
	Height   float64  `json:"height" gorm:"type:float" bson:"height"`
	Metadata Metadata `json:"metadata" gorm:"type:jsonb" bson:"metadata"`
}

type Metadata struct {
	ReferenceId string `json:"referenceId" gorm:"type:varchar(255)" bson:"reference_id"`
	Version     int32  `json:"version" gorm:"type:integer" bson:"version"`
}

type Link struct {
	ID   int `json:"id" gorm:"primaryKey;autoIncrement" bson:"id"`
	From int `json:"from" gorm:"type:integer" bson:"from"`
	To   int `json:"to" gorm:"type:integer" bson:"to"`
}

func (w Workflow) GetNodeByID(nextTask int) Node {
	for _, node := range w.Nodes {
		if node.ID == nextTask {
			return node
		}
	}
	return Node{}
}
