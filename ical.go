package ical

// Create initializes a new Calendar with the given name and description.
//
// If either the name or description is an empty string, it returns nil.
// Otherwise, it returns a pointer to the newly created Calendar.
func Create(name, description string) *Calendar {
	if name == "" || description == "" {
		return nil
	}

	return &Calendar{
		Name:        name,
		Description: description,
	}
}
