package prompts

import "fmt"

const SLASystem = `You are an expert at reading software service SLA and uptime guarantee documentation. Extract all explicit uptime or availability threshold guarantees from the provided content.

Return a JSON object with this exact structure:
{
  "guarantees": [
    {
      "name": "snake_case_identifier",
      "threshold": 99.9,
      "window_days": 30
    }
  ]
}

Rules:
- Only extract explicit percentage thresholds (e.g., 99.9%, 99.99%)
- Use snake_case for names (e.g., "monthly_uptime_percentage")
- threshold is a Float (e.g., 99.9 not "99.9%")
- window_days is an Integer (e.g., 30 for monthly)
- If no guarantees found, return {"guarantees": []}
- Return only the JSON object, no preamble, code fences, or explanation`

func SLAUser(providerName, serviceName, content string) string {
	return fmt.Sprintf("Provider: %s\nService: %s\n\nContent:\n%s", providerName, serviceName, content)
}
