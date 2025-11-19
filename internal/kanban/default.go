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
						"tasks": []interface{}{
							map[string]interface{}{
								"id":    "t1",
								"title": "Progress one task",
							},
							map[string]interface{}{
								"id":    "t2",
								"title": "Create your first list",
							},
						},
					},
					map[string]interface{}{
						"id":    "inprogress",
						"title": "In Progress",
						"tasks": []interface{}{
							map[string]interface{}{
								"id":    "t3",
								"title": "Learn how the kanban works",
							},
						},
					},
					map[string]interface{}{
						"id":    "done",
						"title": "Done",
						"tasks": []interface{}{
							map[string]interface{}{
								"id":    "t4",
								"title": "Created Account",
							},
						},
					},
				},
			},
		},
	}
}
