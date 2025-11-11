package ical

func Create(name, description string) *Calendar {
	if name == "" || description == "" {
		return nil
	}

	return &Calendar{
		Name:        name,
		Description: description,
	}
}
