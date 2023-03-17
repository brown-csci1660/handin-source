package main

// grade assigns a grade to a student by modifying the in-memory
// database. grade returns true if successful.
func grade(problem string, student string, points float64, comment string, files []string, database *DB) bool {

	if database == nil {
		return false
	}

	// Search for the correct assignment in the database.
	for i, asgn := range database.Asgn {
		if asgn.Name == problem {
			// Found the correct assignment. Search for our user & update / insert.
			for i, report := range asgn.Reports {
				if report.User == student {
					report.Grade = points
					report.Comment = comment
					report.Files = files
					asgn.Reports[i] = report
					database.Asgn[i] = asgn
					return true
				}
			}

			// If we get here, the user hadn't already inserted.
			var report Report
			report.User = student
			report.Grade = points
			report.Comment = comment
			report.Files = files

			// Insert into our list of reports
			asgn.Reports = append(asgn.Reports, report)
			database.Asgn[i] = asgn
			return true
		}
	}

	// Assignment not found. Make a new assignment / add a report.
	var report Report
	report.User = student
	report.Grade = points
	report.Comment = comment
	report.Files = files

	var asgn Assignment
	asgn.Name = problem
	asgn.Reports = []Report{report}

	database.Asgn = append(database.Asgn, asgn)
	return true
}
