package api

type Journal struct {
	Members []Student `json:"members"`
	Marks []Mark `json:"marks"`
}

func (j Journal) GetStudentsNames() []string {
	var studentsNames []string
	for _, student := range j.Members {
		studentsNames = append(studentsNames, student.Alias)
	}
	return studentsNames
}

type Student struct {
	Id string `json:"id"`
	TypeId string `json:"type_id"`
	Alias string `json:"alias"`
}

type Grade struct {
	Name string `json:"name"`
	Id string `json:"id"`
	TypeId string `json:"type_id"`
	Subjects[] Subject `json:"items"`
}

func (g Grade) GetSubjectsNames() []string {
	var subjectsNames []string
	for _, subj := range g.Subjects {
		subjectsNames = append(subjectsNames, subj.Name)
	}
	return subjectsNames
}

type Subject struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Classes[] Class `json:"items"`
}

func (s Subject) GetClassesNames() []string {
	var classesNames []string
	for _, class := range s.Classes {
		classesNames = append(classesNames, class.Name)
	}
	return classesNames
}

type Class struct {
	Id string `json:"id"`
	TypeId string `json:"type_id"`
	GradeId string `json:"grade_id"` 
	Name string `json:"name"`
	  
}

type Mark struct {
	StudentId string `json:"student_id"`
	Id string `json:"id"`
	Text string `json:"text"`
}

type MarkDetails struct {
	Id string `json:"id"`
	Value string `json:"val"`
	DateAdd string `json:"date_add"`
	DateModify string `json:"date_modify"`
	UserId string `json:"user_id"`
	TeacherName string `json:"teacher_name"`
	Modified string `json:"modified"`
	OldId string `json:"old_id"`
	OldVal string `json:"old_val"`
	OldAdd string `json:"old_add"`
}