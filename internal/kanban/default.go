package kanban

func DefaultKanban() map[string]interface{} {
	return map[string]interface{}{
		"boards": []interface{}{
			map[string]interface{}{
				"id":    "default",
				"title": "My First Board",
				"lists": []interface{}{
					map[string]interface{}{
						"id":    "todo",
						"title": "To Do",
						"tasks": []interface{}{},
					},
					map[string]interface{}{
						"id":    "inprogress",
						"title": "In Progress",
						"tasks": []interface{}{},
					},
					map[string]interface{}{
						"id":    "done",
						"title": "Done",
						"tasks": []interface{}{},
					},
				},
			},
		},
	}
}
