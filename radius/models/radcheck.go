package models

type RadCheck struct {
	id        int32
	Username  string `json:"username"`
	attribute string
	op        string
	value     string
}

func (r *RadCheck) SetId(id int32){
	r.id = id
}

func (r *RadCheck) SetAttribute(attr string) {
	r.attribute = attr
}

func (r *RadCheck) SetOp(op string) {
	r.op = op
}

func (r *RadCheck) SetValue(value string){
	r.value = value
}